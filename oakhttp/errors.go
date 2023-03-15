package oakhttp

import (
	"encoding/json"
	"errors"
	"net/http"
)

func DefaultErrorHandler(h Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

// var (
// // not needed, just print
// 	ErrNotFound = errors.New("requested resource was not found")
// )

// Unwrap asserts that some error wraps [Error] and returns it. Use the result together with [http.Error]. Log the original error to prevent leaking potentially sensitive information to the client while recording it on the server.
func Unwrap(wrapped error) (err Error, ok bool) {
	if errors.As(wrapped, &err) {
		return err, true
	}
	return nil, false
}

type Error interface {
	error
	HTTPStatusCode() int
}

type NotFoundError struct {
	resource string
}

func NewNotFoundError(resource string) *NotFoundError {
	return &NotFoundError{resource: resource}
}

func (e *NotFoundError) HTTPStatusCode() int {
	return http.StatusNotFound
}

func (e *NotFoundError) Error() string {
	return "resource \"" + e.resource + "\" was not found"
}

type ValidationError struct {
	errors []error
}

func NewValidationError(fieldError ...error) *ValidationError {
	return &ValidationError{
		errors: fieldError,
	}
}

func (e *ValidationError) MarshallJSON() []byte {
	return []byte(`[]`)
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

func (e *ValidationError) Unwrap() []error {
	return e.errors
}

func (e *ValidationError) HTTPStatusCode() int {
	return http.StatusUnprocessableEntity
}

// type EncodingError struct{} // not needed, just print

// type DecodingError struct{} // not needed, just print

// type RequestFactoryError struct{} // not needed, just print

func MarshalErrorToJSON(w http.ResponseWriter, err error) error {
	if httpError, ok := err.(Error); ok {
		w.WriteHeader(httpError.HTTPStatusCode())
	}
	// stringer, ok := str.(fmt.Stringer)
	return json.NewEncoder(w).Encode(map[string]string{
		"Error": err.Error(),
	})
}
