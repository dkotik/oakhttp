package oakhttp

import (
	"context"
	"net/http"
)

// NewContextInjector replaces the [context.Context] for [http.Request]s. Use it to enrich context with authenticated session information.
func NewContextInjector(
	injector func(*http.Request) (context.Context, error),
) Middleware {
	return func(next Handler) Handler {
		return HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) (err error) {
			ctx, err := injector(r)
			if err != nil {
				return err
			}
			return next.ServeHyperText(w, r.WithContext(ctx))
		})
	}
}
