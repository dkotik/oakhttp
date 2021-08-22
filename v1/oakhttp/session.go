package oakhttp

import (
	"net/http"

	"github.com/dkotik/oakacs/v1"
	"github.com/rs/xid"
)

const sessionContextKey = sessionContextKeyType("session")

type sessionContextKeyType string

// Protect wraps a handler and injects session into its context after checking throttling and access.
func Protect(acs *oakacs.AccessControlSystem, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// retrieve session
		// inject new context
		id := xid.New()
		ctx, err := acs.SessionBind(r.Context(), id, "differentiator")
		if err != nil {
			panic(err) // TODO:
		}
		h.ServeHTTP(w, r.WithContext(ctx))
		// context.WithValue(r.Context(), acs.sessionContextKey, "xid.ID")))
	})
}
