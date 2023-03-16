package oakhttp

import (
	"net/http"
)

type ValidatableNormalizable[T comparable] interface {
	*T
	Validate() error
	Normalize() error
}

type InvalidRequestError struct {
	error
}

func NewInvalidRequestError(fromError error) *InvalidRequestError {
	return &InvalidRequestError{fromError}
}

// func (e *InvalidRequestError) Error() string {
// 	return "invalid request"
// }

func (e *InvalidRequestError) Unwrap() error {
	return e.error
}

func (e *InvalidRequestError) HTTPStatusCode() int {
	return http.StatusUnprocessableEntity
}
