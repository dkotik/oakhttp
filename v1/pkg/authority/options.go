package cueroles

import "errors"

type Option func(a *Authority) error

func WithOptions(options ...Option) Option {
	return func(a *Authority) (err error) {
		for _, option := range options {
			if err = option(a); err != nil {
				return err
			}
		}
		return nil
	}
}

func WithMethod(name string, m Method) Option {
	return func(a *Authority) error {
		if a.methodSet == nil {
			a.methodSet = make(map[string]Method)
		}
		if m == nil {
			return errors.New("cannot use an empty authority method")
		}
		a.methodSet[name] = m
		return nil
	}
}

// WithCache

// WithBackend(r Repository)

func WithDefaults() Option {
	return func(a *Authority) error {
		if a.methodSet == nil {
			a.methodSet = make(map[string]Method)
		}
		if _, ok := a.methodSet["exact"]; !ok {
			WithMethod("exact", MethodExact)(a)
		}
		return nil
	}
}
