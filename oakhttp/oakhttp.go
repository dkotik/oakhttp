package oakhttp

import (
	"context"
	"net/http"

	"golang.org/x/exp/slog"
)

type Handler interface {
	ServeHyperText(http.ResponseWriter, *http.Request) error
}

type HandlerFunc func(http.ResponseWriter, *http.Request) error

func (f HandlerFunc) ServeHyperText(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if herr := f(w, r); herr != nil {
		code, err := UnwrapError(herr)
		logError(slog.Default(), err, code, r)
		writeError(w, err, code)
	}
}

type Middleware func(Handler) Handler

// ApplyMiddleware applies [Middleware] in reverse to preserve logical order.
func ApplyMiddleware(h Handler, mws []Middleware) Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

type DomainRequest[T any, V ValidatableNormalizable[T]] func(context.Context, V) error

type DomainRequestResponse[T any, V ValidatableNormalizable[T], O any] func(context.Context, V) (O, error)

// ValidatableNormalizable constrains a domain request. Validation errors will be wrapped as InvalidRequestError by the adapter.
type ValidatableNormalizable[T any] interface {
	*T
	Validate() error
}

func NewRedirect(URL string, statusCode int) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		http.Redirect(w, r, URL, statusCode)
		return nil
	}
}
