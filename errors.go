package oakhttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"net/http"

	"golang.org/x/exp/slog"
)

func LogError(l *slog.Logger, err error, r *http.Request) {
	code, unwrapped := UnwrapError(err)
	logError(l, unwrapped, code, r)
}

func logError(l *slog.Logger, err error, code int, r *http.Request) {
	l.Log(
		r.Context(),
		slog.LevelError,
		"HTTP request failed",
		slog.Int("code", code),
		slog.Any("error", err),
		slog.Any("address", r.RemoteAddr),
		slog.Group("request",
			slog.String("host", r.URL.Hostname()),
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
		),
	)
}

type Error interface {
	error
	HyperTextStatusCode() int
}

func UnwrapError(err error) (code int, unwrapped error) {
	var httpError Error
	if errors.As(err, &httpError) {
		return httpError.HyperTextStatusCode(), httpError
	}
	return http.StatusInternalServerError, err
}

type errorHandler struct {
	next   Handler
	logger *slog.Logger
}

func NewErrorHandler(logger *slog.Logger) Middleware {
	return func(next Handler) Handler {
		if next == nil {
			panic("error handler given a <nil> handler")
		}
		if logger == nil {
			logger = slog.Default()
		}
		return &errorHandler{next: next, logger: logger}
	}
}

func (e *errorHandler) ServeHyperText(
	w http.ResponseWriter,
	r *http.Request,
) error {
	nerr := e.next.ServeHyperText(w, r)
	if nerr == nil {
		return nil
	}
	code, err := UnwrapError(nerr)
	logError(e.logger, err, code, r)
	writeError(w, err, code)
	return nil
}

type errorHandlerWithTemplate struct {
	next     Handler
	logger   *slog.Logger
	template *template.Template
}

func NewErrorHandlerWithTemplate(logger *slog.Logger, t *template.Template) Middleware {
	return func(next Handler) Handler {
		if next == nil {
			panic("error handler with template given a <nil> handler")
		}
		if t == nil {
			panic("error handler with template given a <nil> template")
		}
		if logger == nil {
			logger = slog.Default()
		}
		return &errorHandlerWithTemplate{next: next, logger: logger, template: t}
	}
}

func (e *errorHandlerWithTemplate) ServeHyperText(
	w http.ResponseWriter,
	r *http.Request,
) error {
	nerr := e.next.ServeHyperText(w, r)
	if nerr == nil {
		return nil
	}
	code, err := UnwrapError(nerr)
	logError(e.logger, err, code, r)
	var b = &bytes.Buffer{}
	if err = e.template.Execute(b, struct {
		Code  int
		Error string
	}{
		Code:  code,
		Error: err.Error(),
	}); err != nil {
		return err
	}
	w.Header().Set("Context-Type", "text/html")
	w.WriteHeader(code)
	_, err = io.Copy(w, b)
	return err
}

type errorHandlerJSON struct {
	next   Handler
	logger *slog.Logger
}

func NewJSONErrorHandler(logger *slog.Logger) Middleware {
	return func(next Handler) Handler {
		if next == nil {
			panic("JSON error handler given a <nil> handler")
		}
		if logger == nil {
			logger = slog.Default()
		}
		return &errorHandlerJSON{next: next, logger: logger}
	}
}

func (e *errorHandlerJSON) ServeHyperText(
	w http.ResponseWriter,
	r *http.Request,
) error {
	nerr := e.next.ServeHyperText(w, r)
	if nerr == nil {
		return nil
	}
	code, err := UnwrapError(nerr)
	logError(e.logger, err, code, r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(struct {
		Error string
	}{err.Error()})
}

type notFoundError struct {
	cause error
}

func NewNotFoundError(cause error) error {
	return &notFoundError{cause: cause}
}

func (e *notFoundError) Error() string {
	return http.StatusText(http.StatusNotFound)
}

func (e *notFoundError) HyperTextStatusCode() int {
	return http.StatusNotFound
}

func (e *notFoundError) LogValue() slog.Value {
	if e.cause == nil {
		slog.StringValue(e.Error())
	}
	return slog.StringValue("not found: " + e.cause.Error())
}

type accessDeniedError struct {
	cause error
}

func NewAccessDeniedError(cause error) error {
	return &accessDeniedError{cause: cause}
}

func (e *accessDeniedError) Error() string {
	return http.StatusText(http.StatusForbidden)
}

func (e *accessDeniedError) HyperTextStatusCode() int {
	return http.StatusForbidden
}

func (e *accessDeniedError) LogValue() slog.Value {
	if e.cause == nil {
		slog.StringValue(e.Error())
	}
	return slog.StringValue("access denied: " + e.cause.Error())
}

func IsAccessDeniedError(err error) bool {
	var httpError Error
	return errors.As(err, &httpError) && httpError.HyperTextStatusCode() == http.StatusForbidden
}

type InvalidRequestError struct {
	error
}

func NewInvalidRequestError(fromError error) *InvalidRequestError {
	return &InvalidRequestError{fromError}
}

func (e *InvalidRequestError) Error() string {
	return "invalid request: " + e.error.Error()
}

func (e *InvalidRequestError) Unwrap() error {
	return e.error
}

func (e *InvalidRequestError) HyperTextStatusCode() int {
	return http.StatusUnprocessableEntity
}

type methodNotAllowedError struct {
	method string
}

func NewMethodNotAllowedError(method string) Error {
	return &methodNotAllowedError{method: method}
}

func (e *methodNotAllowedError) Error() string {
	if e.method == "" {
		return "unspecified method is not allowed"
	}
	return "method not allowed: " + e.method
}

func (e *methodNotAllowedError) HyperTextStatusCode() int {
	return http.StatusMethodNotAllowed
}

// ObfuscatedError presents itself as a missing resource. It is useful for hiding certain API endpoints, like the health check, from un-authenticated requests. The real error is always logged using the [slog.Value] interface.
type ObfuscatedError struct {
	real error
}

func NewObfuscatedError(real error) error {
	return &ObfuscatedError{real: real}
}

func (e *ObfuscatedError) Unwrap() error {
	return e.real
}

func (e *ObfuscatedError) Error() string {
	return http.StatusText(http.StatusNotFound)
}

func (e *ObfuscatedError) HyperTextStatusCode() int {
	return http.StatusNotFound
}

func (e *ObfuscatedError) LogValue() slog.Value {
	return slog.StringValue("obfuscated error: " + e.real.Error())
}
