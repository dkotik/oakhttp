package oakhttp

import (
	"errors"
	"fmt"
)

type adaptorOptions struct {
	readLimit   int64
	middlewares []Middleware
}

type AdaptorOption func(*adaptorOptions) error

func newAdaptorOptions(from []AdaptorOption) adaptorOptions {
	o := adaptorOptions{}
	var err error
	for _, option := range append(from, WithDefaultReadLimitOf1MB()) {
		if err = option(&o); err != nil {
			panic(fmt.Errorf("cannot initiate an OakHTTP adaptor: %w", err))
		}
	}
	return o
}

// WithDefaultReadLimitOf1MB sets default limit to one megabyte, if it was not yet set by [WithReadLimit] option. `parsePostForm` in the standard library has the default limit of ten megabytes, but most JSON requests are around 2500 bytes in size. One megabyte was chosen as a sweet spot with plenty of room to grow.
//
// Read limit is not an effective mitigation for denial of service attacks. Time and rate limits are much better. Read limit, however, is an important part of defensive coding that help detect bugs and intrusion attempts, so it remains with generous headroom.
func WithDefaultReadLimitOf1MB() AdaptorOption {
	return func(a *adaptorOptions) error {
		if a.readLimit == 0 {
			return WithReadLimit(1 << 20)(a)
		}
		return nil
	}
}

func WithReadLimit(limit int64) AdaptorOption {
	return func(a *adaptorOptions) error {
		if a.readLimit != 0 {
			return errors.New("read limit is already set")
		}
		if limit <= 0 {
			return errors.New("read limit must be greater than 0 bytes")
		}
		a.readLimit = limit
		return nil
	}
}

func WithMiddleware(ms ...Middleware) AdaptorOption {
	return func(a *adaptorOptions) error {
		for _, m := range ms {
			if m == nil {
				return errors.New("cannot use uninitialized <nil> middleware")
			}
			a.middlewares = append(a.middlewares, m)
		}
		return nil
	}
}
