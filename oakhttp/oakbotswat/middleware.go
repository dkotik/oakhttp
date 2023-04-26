package oakbotswat

import (
	"fmt"
	"net/http"

	"github.com/dkotik/oakacs/oakhttp"
)

func New(withOptions ...Option) (oakhttp.Middleware, error) {
	o := &options{}

	var err error
	for _, option := range append(withOptions, WithDefaultOptions()) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("failed to initialize an OakBotSWAT middleware: %w", err)
		}
	}

	verify := func(r *http.Request) (string, error) {
		token, err := o.ResponseExtractor(r)
		if err != nil {
			return "", err
		}
		return o.Verifier(r.Context(), token, r.RemoteAddr)
	}

	return func(next oakhttp.Handler) oakhttp.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			_, err := verify(r)
			if err != nil {
				return o.Encoder(w, err)
			}
			return next(w, r)
		}
	}, nil
}
