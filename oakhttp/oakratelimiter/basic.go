package oakratelimiter

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/dkotik/oakacs/oakhttp"
)

type basic struct {
	failure  error
	interval time.Duration
	rate     Rate
	limit    float64

	mu sync.Mutex
	bucket
}

func (b *basic) Rate() Rate {
	return b.rate
}

func (b *basic) Take(r *http.Request) error {
	t := time.Now()
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.bucket.Take(b.limit, b.rate, t, t.Add(b.interval)) {
		return b.failure
	}
	return nil
}

func (b *basic) Middleware() oakhttp.Middleware {
	return NewMiddleware(b, b.rate)
}

func (b *basic) ObfuscatedMiddleware(displayRate Rate) oakhttp.Middleware {
	return NewMiddleware(b, displayRate)
}

func newBasic(withOptions ...LimitOption) (*basic, error) {
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

	return &basic{
		failure: NewTooManyRequestsError(
			fmt.Errorf("rate limiter %q ran out of tokens", o.Name)),
		rate:     NewRate(o.Limit, o.Interval),
		limit:    o.Limit,
		interval: o.Interval,
		mu:       sync.Mutex{},
		bucket:   bucket{},
	}, nil
}

func NewBasic(withOptions ...LimitOption) (RateLimiter, error) {
	return newBasic(withOptions...)
}
