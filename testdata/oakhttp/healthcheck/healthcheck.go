package healthcheck

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/tigerperformanceinstitute/tprograms/service/oakhttp"
	"golang.org/x/exp/slog"
)

type HealthCheck interface {
	ReportStatus(context.Context, *slog.Logger) (any, error)
}

type HealthCheckFunc func(context.Context, *slog.Logger) (any, error)

func (f HealthCheckFunc) ReportStatus(ctx context.Context, l *slog.Logger) (any, error) {
	return f(ctx, l)
}

type healthCheckHandler struct {
	authenticator Authenticator
	names         []string
	checks        []HealthCheck
	logger        *slog.Logger

	mu         sync.Mutex
	statusCode int
	report     map[string]any
}

func New(withOptions ...Option) oakhttp.Handler {
	o := &options{}
	var err error
	for _, option := range append(
		withOptions,
		WithDefaultFrequencyEveryFiveMinutes(),
		WithDefaultLimitHalfOfFrequency(),
		WithDefaultLogger(),
		func(o *options) error { // validate
			if len(o.checks) == 0 {
				return errors.New("provide at least one health check")
			}
			if o.frequency < o.limit {
				return errors.New("frequency must not be greater than the limit")
			}
			if o.authenticator == nil {
				o.authenticator = func(r *http.Request) error {
					return nil // no operation
				}
				o.logger.Error("health check created without any authentication")
			}
			return nil
		},
	) {
		if err = option(o); err != nil {
			panic(fmt.Errorf("failed to create health check handler: %w", err))
		}
	}

	h := &healthCheckHandler{
		authenticator: o.authenticator,
		mu:            sync.Mutex{},
		statusCode:    http.StatusProcessing,
		names:         o.names,
		checks:        o.checks,
		logger:        o.logger,
	}
	go h.loop(context.Background(), o.frequency)
	return h
}

func (h *healthCheckHandler) loop(parent context.Context, frequency time.Duration) {
	report := make(map[string]any)
	ctx, cancel := context.WithTimeout(parent, frequency)
	err := h.runAllChecks(ctx, report)
	h.updateReport(ctx, report, err)
	cancel()

	t := time.NewTicker(frequency)
	contextTimeout := frequency - time.Second
	for {
		select {
		case <-parent.Done():
			t.Stop()
			return
		case <-t.C:
			report = make(map[string]any)
			ctx, cancel = context.WithTimeout(parent, contextTimeout)
			err = h.runAllChecks(ctx, report)
			h.updateReport(ctx, report, err)
			cancel()
		}
	}
}

func (h *healthCheckHandler) updateReport(
	ctx context.Context,
	report map[string]any,
	err error,
) {
	var statusCode int
	if err != nil {
		statusCode = http.StatusInternalServerError
	} else {
		statusCode = http.StatusOK
		h.logger.Log(
			ctx,
			slog.LevelDebug,
			"all health checks passed",
		)
	}
	h.mu.Lock()
	h.report = report
	h.statusCode = statusCode
	h.mu.Unlock()
}

func (h *healthCheckHandler) runAllChecks(ctx context.Context, report map[string]any) error {
	var cumulative []error
	for i, check := range h.checks {
		name := h.names[i]
		status, err := check.ReportStatus(ctx, h.logger)
		report[name] = status
		if err != nil {
			cumulative = append(
				cumulative,
				fmt.Errorf("health check %q failed: %w", name, err),
			)
		}
	}
	if len(cumulative) > 0 {
		return &HealthCheckError{
			causes: cumulative,
			report: report,
		}
	}
	return nil
}

func (h *healthCheckHandler) ServeHyperText(
	w http.ResponseWriter,
	r *http.Request,
) (err error) {
	if err = h.authenticator(r); err != nil {
		return oakhttp.NewObfuscatedError(err)
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}
	codec, err := oakhttp.GetCodec(contentType)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(h.statusCode)
	return codec.Encode(w, h.report)
}

func (h *healthCheckHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	err := h.ServeHyperText(w, r)
	if err != nil {
		code, unwrapped := oakhttp.UnwrapError(err)
		http.Error(w, unwrapped.Error(), code)
		h.logger.Error(
			"health checks are failing",
			// slog.Any("report", report),
			slog.Any("error", err),
		)
	}
}
