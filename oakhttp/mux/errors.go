package mux

import (
	"errors"
	"net/http"
	// "golang.org/x/slog"
)

var (
	ErrDoubleSlash        = NewRoutingError(errors.New("path contains double slash"))
	ErrPathNotFound       = NewRoutingError(errors.New("path not found"))
	ErrNotInteger         = NewRoutingError(errors.New("field is not an integer"))
	ErrNotUnsignedInteger = NewRoutingError(errors.New("field is not an unsigned integer"))
	ErrNotPageNumber      = NewRoutingError(errors.New("field is not an page number"))
)

type RoutingError struct {
	cause error
}

func NewRoutingError(cause error) *RoutingError {
	return &RoutingError{cause: cause}
}

func (r *RoutingError) Unwrap() error {
	return r.cause
}

func (r *RoutingError) HyperTextStatusCode() int {
	return http.StatusNotFound
}

func (r *RoutingError) Error() string {
	return http.StatusText(http.StatusNotFound)
}

// TODO: re-enable after slog is merged into standard library
// func (r *RoutingError) LogValue() slog.Value {
// 	return slog.StringValue("routing error: " + e.cause.Error())
// }
