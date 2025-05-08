package oakhttp

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

// ApplyMiddleware applies [Middleware] in reverse to preserve logical order.
func ApplyMiddleware(h http.Handler, mws []Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
