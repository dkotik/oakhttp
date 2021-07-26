package oakhttp

import (
	"context"
	"net/http"
)

const sessionContextKey = sessionContextKeyType("session")

type sessionContextKeyType string

// Protect wraps a handler and injects session into its context after checking throttling and access.
func (acs *AccessControlSystem) Protect(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create session
		// inject new context
		h.ServeHTTP(w, r.WithContext(
			context.WithValue(r.Context(), acs.sessionContextKey, "xid.ID")))
	})
}
