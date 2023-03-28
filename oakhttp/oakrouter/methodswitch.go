package oakrouter

import (
	"net/http"

	"github.com/dkotik/oakacs/oakhttp"
)

type Methods struct {
	Get    oakhttp.Handler
	Post   oakhttp.Handler
	Put    oakhttp.Handler
	Delete oakhttp.Handler
	Patch  oakhttp.Handler
}

func NewMethodSwitch(m Methods) oakhttp.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		switch r.Method {
		case http.MethodGet, http.MethodHead:
			if m.Get != nil {
				return m.Get(w, r)
			}
		case http.MethodPost:
			if m.Post != nil {
				return m.Post(w, r)
			}
		case http.MethodPut:
			if m.Put != nil {
				return m.Put(w, r)
			}
		case http.MethodDelete:
			if m.Delete != nil {
				return m.Delete(w, r)
			}
		case http.MethodPatch:
			if m.Patch != nil {
				return m.Patch(w, r)
			}
		}
		return NewMethodNotAllowedError(r.Method)
	}
}

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
