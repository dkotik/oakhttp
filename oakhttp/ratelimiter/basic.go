package ratelimiter

import (
	"net/http"
	"sync"
	"time"

	"github.com/dkotik/oakacs/oakhttp"
)

type Basic struct {
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
		return ErrTooManyRequests
	}
	return nil
}

func (b *Basic) Middleware() oakhttp.Middleware {
	return NewMiddleware(b, b.rate)
}

func NewBasic(limit float64, interval time.Duration) *Basic {
	return &Basic{
		rate:     NewRate(limit, interval),
		limit:    limit,
		interval: interval,
		mu:       sync.Mutex{},
		bucket:   bucket{},
	}
}
