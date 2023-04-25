package oakswitch

import (
	"net/http"

	"github.com/dkotik/oakacs/oakhttp"
)

type methodNotAllowedError struct {
	method string
}

func NewMethodNotAllowedError(method string) oakhttp.Error {
	return &methodNotAllowedError{method: method}
}

func (e *methodNotAllowedError) Error() string {
	if e.method == "" {
		return "unspecified method is not allowed"
	}
	return "method not allowed: " + e.method
}

func (e *methodNotAllowedError) HTTPStatusCode() int {
	return http.StatusMethodNotAllowed
}
