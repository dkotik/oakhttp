package oakratelimiter

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/dkotik/oakacs/oakhttp"
	"golang.org/x/exp/slog"
)

type RateLimiter interface {
	Take(*http.Request) error
}

type TooManyRequestsError struct {
	cause error
}

func (e *TooManyRequestsError) Unwrap() error {
	return e.cause
}

func (e *TooManyRequestsError) Error() string {
	return http.StatusText(http.StatusTooManyRequests)
}

func (e *TooManyRequestsError) HTTPStatusCode() int {
	return http.StatusTooManyRequests
}

func (e *TooManyRequestsError) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("message", "too many requests"),
		slog.Any("cause", e.cause),
	)
}

var ErrTooManyRequests = &TooManyRequestsError{
	cause: errors.New("no more tokens available"),
}

func New(withOptions ...Option) (RateLimiter, error) {
	o, err := newOptions(append(
		withOptions,
		func(o *options) error { // validate
			return nil
		},
	)...)
	if err != nil {
		return nil, fmt.Errorf("cannot create the rate limiter: %w", err)
	}

	if len(o.Tagging) == 0 {
		return o.Basic, nil
	}

	if o.CleanUpContext == nil {
		o.CleanUpContext = context.Background()
	}

	if len(o.Tagging) == 1 {
		s := &SingleTagging{
			basic:           *o.Basic,
			taggedBucketMap: o.Tagging[0],
		}
		go s.purgeLoop(o.CleanUpContext, o.CleanUpPeriod)
		return s, nil
	}

	m := &MultiTagging{
		basic:           *o.Basic,
		taggedBucketMap: o.Tagging,
	}
	go m.purgeLoop(o.CleanUpContext, o.CleanUpPeriod)
	return m, nil
}

// NewMiddleware protects an [oakhttp.Handler] using a [RateLimiter]. The display [Rate] can be used to obfuscate the true [RateLimiter] throughput. HTTP headers are set to promise availability of no more than one call. This is done to conceal the performance capacity of the system, while giving some useful information to API callers regarding service availability. "X-RateLimit-*" headers are experimental, inconsistent in implementation, and meant to be approximate. If display [Rate] is 0, the headers are ommitted.
func NewMiddleware(l RateLimiter, displayRate Rate) oakhttp.Middleware {
	if l == nil {
		panic("<nil> rate limiter")
	}

	if displayRate == Rate(0) {
		return func(next oakhttp.Handler) oakhttp.Handler {
			return func(w http.ResponseWriter, r *http.Request) error {
				if err := l.Take(r); err != nil {
					return err
				}
				return next(w, r)
			}
		}
	}

	limit := uint(1)
	oneTokenWindow := time.Nanosecond * time.Duration(1.05/displayRate)
	if oneTokenWindow < time.Second {
		limit = uint(math.Min(
			math.Floor(float64(time.Second.Nanoseconds())*float64(displayRate*0.95)),
			1,
		))
		oneTokenWindow = time.Second
	}
	displayLimit := fmt.Sprintf("%d", limit)
	return func(next oakhttp.Handler) oakhttp.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			t := time.Now().
				Add(oneTokenWindow).
				UTC().
				Format(time.RFC1123)

			header := w.Header()
			header.Set("X-RateLimit-Limit", displayLimit)
			header.Set("X-RateLimit-Reset", t)

			if err := l.Take(r); err != nil {
				header.Set("X-RateLimit-Remaining", "0")
				header.Set("Retry-After", t)
				return err
			}
			header.Set("X-RateLimit-Remaining", "1")
			return next(w, r)
		}
	}
}
