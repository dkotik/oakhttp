package oakratelimiter

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/dkotik/oakacs/oakhttp"
)

// Basic rate limiter enforces the limit using one leaky token bucket.
type Basic struct {
	failure  error
	interval time.Duration
	rate     Rate
	limit    float64

	mu sync.Mutex
	bucket
}

// Rate returns the rate limiter [Rate].
func (b *Basic) Rate() Rate {
	return b.rate
}

// Take consumes one token per request. Returns a [TooManyRequestsError] if the leaky bucket is drained.
func (b *Basic) Take(r *http.Request) error {
	t := time.Now()
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.bucket.Take(b.limit, b.rate, t, t.Add(b.interval)) {
		return b.failure
	}
	return nil
}

// Middleware creates an [oakhttp.Middleware] from the [Basic] rate limiter.
func (b *Basic) Middleware() oakhttp.Middleware {
	return NewMiddleware(b, b.rate)
}

// ObfuscatedMiddleware creates an [oakhttp.Middleware] from the [Basic] rate limiter with a display [Rate] different from the actual.
func (b *Basic) ObfuscatedMiddleware(displayRate Rate) oakhttp.Middleware {
	return NewMiddleware(b, displayRate)
}

// NewBasic initializes a [Basic] rate limiter.
func NewBasic(withOptions ...LimitOption) (*Basic, error) {
	o, err := newSupervisingLimitOptions(withOptions...)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize basic rate limiter: %w", err)
	}

	return &Basic{
		failure: NewTooManyRequestsError(
			fmt.Errorf("rate limiter %q ran out of tokens", o.Name)),
		rate:     NewRate(o.Limit, o.Interval),
		limit:    o.Limit,
		interval: o.Interval,
		mu:       sync.Mutex{},
		bucket:   bucket{},
	}, nil
}
