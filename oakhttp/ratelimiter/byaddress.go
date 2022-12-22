package ratelimiter

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type byClientAddress struct {
	http.Handler
	limiters        map[string]*rate.Limiter
	rate            rate.Limit
	detectionLimit  int
	maxAddressCount int
	minAddressCount int
	mu              *sync.Mutex
}

func (b *byClientAddress) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !b.Allow(r.RemoteAddr) {
		writeError(w, r)
		return
	}
	b.Handler.ServeHTTP(w, r)
}

func (b *byClientAddress) Allow(IP string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.limiters) > b.maxAddressCount {
		return false // tracking too many IP addresses
	}

	limit, ok := b.limiters[IP]
	if !ok {
		limit = rate.NewLimiter(b.rate, b.detectionLimit)
		b.limiters[IP] = limit
	}
	return limit.Allow()
}

func (b *byClientAddress) Clean(limit float64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.limiters) <= b.minAddressCount {
		// keep some minimum of the addresses in memory
		return
	}

	var eliminationQueue []string
	for key, limiter := range b.limiters {
		if limiter.Tokens() == limit {
			eliminationQueue = append(eliminationQueue, key)
			if len(eliminationQueue) == b.minAddressCount {
				break // do not evict more than minimum at a time
			}
		}
	}

	for _, key := range eliminationQueue {
		delete(b.limiters, key)
	}
}

func New(h http.Handler, withOptions ...Option) (http.Handler, error) {
	o := &options{}
	var err error
	for _, option := range append(
		withOptions,
		WithDefaultLimit(),
		WithDefaultCleanUpPeriod(),
		WithDefaultMaximumAddressCount(),
		WithDefaultMinimumAddressCount(),
		func(o *options) error {
			// validate options as the last step
			if h == nil {
				return errors.New("http.Handler is required")
			}
			if o.detectionPeriod == 0 || o.detectionLimit == 0 {
				return errors.New("WithLimit is a required option")
			}
			if o.cleanUpPeriod == 0 {
				return errors.New("WithCleanUpPeriod is a required option")
			}
			if o.minAddressCount == 0 {
				return errors.New("WithMinimumAddressCount is a required option")
			}
			if o.maxAddressCount == 0 {
				return errors.New("WithMaximumAddressCount is a required option")
			}
			return nil
		},
	) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("IP-based rate limiter initialization failed: %w", err)
		}
	}

	rateLimiter := &byClientAddress{
		Handler:         h,
		rate:            rate.Every(o.detectionPeriod),
		limiters:        make(map[string]*rate.Limiter),
		detectionLimit:  o.detectionLimit,
		maxAddressCount: o.maxAddressCount,
		minAddressCount: o.maxAddressCount / 5,
		mu:              &sync.Mutex{},
	}
	go func(d time.Duration, limit float64) { // clean up in parallel
		t := time.NewTicker(d)
		for {
			select {
			case <-t.C:
				rateLimiter.Clean(limit)
			}
		}
	}(o.cleanUpPeriod, float64(o.detectionLimit))
	return rateLimiter, nil
}
