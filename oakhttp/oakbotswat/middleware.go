package oakbotswat

import (
	"fmt"
	"net"
	"net/http"

	"github.com/dkotik/oakacs/oakhttp"
	"golang.org/x/exp/slog"
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
		cached, err := o.Cache.GetToken(r.Context(), token)
		if err != nil {
			return "", err
		}
		if cached {
			return token, nil // cache hit
		}

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		// first returned value is data
		_, err = o.Verifier(r.Context(), token, ip)
		if err != nil {
			return "", err
		}
		o.Cache.SetToken(r.Context(), token)
		return token, nil
	}

	return func(next oakhttp.Handler) oakhttp.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			_, err := verify(r)
			if err != nil {
				slog.Error("failed botswat verification", slog.Any("error", err))
				// panic(err)
				return o.Encoder(w, err)
			}
			return next(w, r)
		}
	}, nil
}
