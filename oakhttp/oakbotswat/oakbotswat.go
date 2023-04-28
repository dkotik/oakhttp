package oakbotswat

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/dkotik/oakacs/oakhttp"
)

type OakBotSWAT struct {
	verifier               Verifier
	humanityTokenExtractor HumanityTokenExtractor
}

// IsHuman returns `nil` for human agents, an [Error] if humanity [Verifier] was not passed, or an [error] for any other condition. In rare cases when you need access to user data proven by the UI library, such as the Turnstile implementation, use that [Verifier] directly.
func (b *OakBotSWAT) IsHuman(r *http.Request) error {
	token, err := b.humanityTokenExtractor(r)
	if err != nil {
		return fmt.Errorf("cannot recover request token proving humanity: %w", err)
	}
	if token == "" {
		return ErrTokenEmpty
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	// first returned value is data
	_, err = b.verifier.VerifyHumanityToken(r.Context(), token, ip)
	return err
}

func (b *OakBotSWAT) Middleware() oakhttp.Middleware {
	return func(next oakhttp.Handler) oakhttp.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			if err := b.IsHuman(r); err != nil {
				return err
			}
			return next(w, r)
		}
	}
}

func (b *OakBotSWAT) Gate(onError oakhttp.Encoder) oakhttp.Middleware {
	if onError == nil {
		panic("got <nil> onError encoder")
	}
	return func(next oakhttp.Handler) oakhttp.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			if err := b.IsHuman(r); err != nil {
				return onError(w, err)
			}
			return next(w, r)
		}
	}
}

func New(withOptions ...Option) (b *OakBotSWAT, err error) {
	o := &options{}
	for _, option := range append(withOptions, func(o *options) error {
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
			return nil, fmt.Errorf("failed to initialize OakBotSWAT team: %w", err)
		}
	}

	if o.Cache != nil {
		o.Verifier = &cachedVerifier{
			cache:   o.Cache,
			backend: o.Verifier,
		}
	}

	return &OakBotSWAT{
		verifier:               o.Verifier,
		humanityTokenExtractor: o.HumanityTokenExtractor,
	}, nil
}
