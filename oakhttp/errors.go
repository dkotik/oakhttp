package oakhttp

import (
	"encoding/json"
	"net/http"
)

// var (
// // not needed, just print
// 	ErrNotFound = errors.New("requested resource was not found")
// )

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
