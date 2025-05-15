package botswat

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/dkotik/oakhttp"
)

type botswat struct {
	ErrorHandler           oakhttp.ErrorHandler
	Verifier               Verifier
	HumanityTokenExtractor HumanityTokenExtractor
}

// IsHuman returns `nil` for human agents, an [Error] if humanity [Verifier] was not passed, or an [error] for any other condition. In rare cases when you need access to user data proven by the UI library, such as the Turnstile implementation, use that [Verifier] directly.
func (b *botswat) IsHuman(r *http.Request) error {
	token, err := b.HumanityTokenExtractor(r)
	if err != nil {
		return fmt.Errorf("cannot recover request token proving humanity: %w", err)
	}
	if token == "" {
		return ErrTokenEmpty
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	// first returned value is data
	_, err = b.Verifier.VerifyHumanityToken(r.Context(), token, ip)
	return err
}

func (b *botswat) Middleware() oakhttp.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := b.IsHuman(r); err != nil {
				b.ErrorHandler.HandleError(w, r, err)
			}
			next.ServeHTTP(w, r)
		})
	}
}

// func (b *botswat) Gate(onError oakhttp.Encoder) oakhttp.Middleware {
// 	if onError == nil {
// 		panic("got <nil> onError encoder")
// 	}
// 	return func(next oakhttp.Handler) oakhttp.Handler {
// 		return func(w http.ResponseWriter, r *http.Request) error {
// 			if err := b.IsHuman(r); err != nil {
// 				return onError(w, err)
// 			}
// 			return next(w, r)
// 		}
// 	}
// }

func New(withOptions ...Option) (b *botswat, err error) {
	o := &options{}
	for _, option := range append(withOptions, func(o *options) error {
		if o.ErrorHandler == nil {
			if err := WithErrorHandler(oakhttp.NewErrorHandler(nil, nil, nil))(o); err != nil {
				return fmt.Errorf("unable to set up default error handler: %w", err)
			}
		}
		if o.Verifier == nil {
			return errors.New("WithVerifier option is required")
		}
		if o.HumanityTokenExtractor == nil {
			if err := WithDefaultCookieHumanityTokenExtractor()(o); err != nil {
				return err
			}
		}
		return nil
	}) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("failed to initialize botswat team: %w", err)
		}
	}

	if o.Cache != nil {
		o.Verifier = &cachedVerifier{
			cache:   o.Cache,
			backend: o.Verifier,
		}
	}

	return &botswat{
		Verifier:               o.Verifier,
		HumanityTokenExtractor: o.HumanityTokenExtractor,
	}, nil
}
