package oakhttp

import (
	"net/http"
)

// ValidatableNormalizable constrains a domain request. Validation and normalization errors will be wrapped as InvalidRequestError by the adapter.
type ValidatableNormalizable[T any] interface {
	*T
	Validate() error
	Normalize() error
}

// valid must wrap the [RequestFactory] call within every adapter.
func valid[T any, P ValidatableNormalizable[T]](request P, err error) (P, error) {
	if err != nil {
		return nil, NewInvalidRequestError(err)
	}

	if err = request.Validate(); err != nil {
		return nil, NewInvalidRequestError(err)
	}
	if err = request.Normalize(); err != nil {
		return nil, NewInvalidRequestError(err)
	}
	return request, nil
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
