package turnstile

import "errors"

type middlewareOptions struct {
	Verifier Verifier
	// Cache func()
}

type MiddlewareOption func(*middlwareOptions) error

func WithMiddlewareVerifier(v Verifier) MiddlewareOption {
	return func(o *middlwareOptions) error {
		if o.Verifier != nil {
			return errors.New("verifier is already set")
		}
		if v == nil {
			return errors.New("cannot use a <nil> verifier")
		}
		o.Verifier = v
		return nil
	}
}
