package healthcheck

import (
	"errors"
	"net/http"

	"golang.org/x/exp/slog"
)

var (
	ErrTokenRejected = errors.New("invalid token")
)

type HealthCheckError struct {
	causes []error
	report map[string]any
}

func (e *HealthCheckError) Unwrap() []error {
	return e.causes
}

func (e *HealthCheckError) Error() string {
	return "health checks are failing"
}

func (e *HealthCheckError) HyperTextStatusCode() int {
	return http.StatusInternalServerError
}

func (e *HealthCheckError) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Any("errors", e.causes),
		slog.Any("report", e.report),
	)
}
