package botswat

import "errors"

type options struct {
	Verifier          Verifier
	ResponseExtractor ResponseExtractor
	// Cache func()
}

type Option func(*options) error

func WithVerifier(v Verifier) Option {
	return func(o *options) error {
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

func WithResponseExtractor(e ResponseExtractor) Option {
	return func(o *options) error {
		if o.ResponseExtractor != nil {
			return errors.New("response extractor is already set")
		}
		if e == nil {
			return errors.New("cannot use a <nil> response extractor")
		}
		o.ResponseExtractor = e
		return nil
	}
}
