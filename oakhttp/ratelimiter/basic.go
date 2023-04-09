package ratelimiter

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/dkotik/oakacs/oakhttp"
)

type Basic struct {
	failure  error
	interval time.Duration
	rate     Rate
	limit    float64

	mu sync.Mutex
	bucket
}

func (b *Basic) Rate() Rate {
	return b.rate
}

func (b *Basic) Take(r *http.Request) error {
	t := time.Now()
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.bucket.Take(b.limit, b.rate, t, t.Add(b.interval)) {
		return b.failure
	}
	return nil
}

func (b *Basic) Middleware() oakhttp.Middleware {
	return NewMiddleware(b, b.rate)
}

func (b *Basic) ObfuscatedMiddleware(displayRate Rate) oakhttp.Middleware {
	return NewMiddleware(b, displayRate)
}

func NewBasic(withOptions ...LimitOption) (*Basic, error) {
	o, err := newLimitOptions(append(
		withOptions,
		WithDefaultName(),
		func(o *limitOptions) error {
			// if o.Name == "" {
			// 	return errors.New("name option is required")
			// }
			if o.InitialAllocationSize != 0 {
				return errors.New("initial allocation size option cannot be applied to a basic rate limiter")
			}
			return nil
		},
	)...)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize basic rate limiter: %w", err)
	}

	return &Basic{
		failure: &TooManyRequestsError{
			cause: fmt.Errorf("rate limiter %q ran out of tokens", o.Name),
		},
		rate:     NewRate(o.Limit, o.Interval),
		limit:    o.Limit,
		interval: o.Interval,
		mu:       sync.Mutex{},
		bucket:   bucket{},
	}, nil
}
