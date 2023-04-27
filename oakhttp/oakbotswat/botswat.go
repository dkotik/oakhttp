package oakbotswat

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/dkotik/oakacs/oakhttp"
)

// IsHuman returns `nil` for human agents, an [Error] if humanity [Verifier] was not passed, or an [error] for any other condition. In rare cases when you need access to user data proven by the UI library, such as the Turnstile implementation, use that [Verifier] directly.
type IsHuman func(*http.Request) error

// Verifier returns [Error] if client response was not recognized as valid.
type Verifier func(
	ctx context.Context,
	clientResponseToken string,
	clientIPAddress string,
) (
	userData string,
	err error,
)

type ResponseExtractor func(
	r *http.Request,
) (
	clientResponseToken string,
	err error,
)

func New(withOptions ...Option) (IsHuman, error) {
	o := &options{}

	var err error
	for _, option := range append(withOptions, WithDefaultOptions()) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("failed to initialize OakBotSWAT team: %w", err)
		}
	}

	if o.Cache == nil {
		return func(r *http.Request) error {
			token, err := o.ResponseExtractor(r)
			if err != nil {
				return fmt.Errorf("cannot recover request token proving humanity: %w", err)
			}
			if token == "" {
				return ErrTokenEmpty
			}
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			// first returned value is data
			_, err = o.Verifier(r.Context(), token, ip)
			if err != nil {
				return err
			}
			return nil
		}, nil
	}

	return func(r *http.Request) error {
		token, err := o.ResponseExtractor(r)
		if err != nil {
			return fmt.Errorf("cannot recover request token proving humanity: %w", err)
		}
		if token == "" {
			return ErrTokenEmpty
		}
		cached, err := o.Cache.GetToken(r.Context(), token)
		if err != nil {
			return fmt.Errorf("cannot access humanity token cache: %w", err)
		}
		if cached {
			return nil // cache hit
		}

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		// first returned value is data
		_, err = o.Verifier(r.Context(), token, ip)
		if err != nil {
			return err
		}
		o.Cache.SetToken(r.Context(), token)
		return nil
	}, nil
}

func NewMiddleware(withOptions ...Option) (oakhttp.Middleware, error) {
	isHuman, err := New(withOptions...)
	if err != nil {
		return nil, err
	}

	return func(next oakhttp.Handler) oakhttp.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			if err := isHuman(r); err != nil {
				return err
			}
			return next(w, r)
		}
	}, nil
}

func NewGate(gate oakhttp.Encoder, withOptions ...Option) (oakhttp.Middleware, error) {
	isHuman, err := New(withOptions...)
	if err != nil {
		return nil, err
	}

	return func(next oakhttp.Handler) oakhttp.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			if err := isHuman(r); err != nil {
				return gate(w, err)
			}
			return next(w, r)
		}
	}, nil
}
