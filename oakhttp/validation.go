package oakhttp

import (
	"net/http"
)

// ValidatableNormalizable constrains a domain request. Validation and normalization errors will be wrapped as InvalidRequestError by the adapter.
type ValidatableNormalizable[T comparable] interface {
	*T
	Validate() error
	Normalize() error
}

// InvalidRequestError represents
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

func (e *InvalidRequestError) HTTPStatusCode() int {
	return http.StatusUnprocessableEntity
}
