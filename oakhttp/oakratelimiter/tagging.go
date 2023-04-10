package oakratelimiter

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dkotik/oakacs/oakhttp"
)

type SingleTagging struct {
	basic
	taggedBucketMap
}

func NewSingleDiscriminating(withOptions ...Option) (*SingleTagging, error) {
	o, err := newOptions(append(
		withOptions,
		func(o *options) error { // validate
			if len(o.Tagging) != 1 {
				return errors.New("single-tagged rate limiter must be initiated with exactly one tagger")
			}
			return nil
		},
	)...)
	if err != nil {
		return nil, fmt.Errorf("cannot create single-tagged rate limiter: %w", err)
	}

	s := &SingleTagging{
		basic:           *o.Basic,
		taggedBucketMap: o.Tagging[0],
	}

	if o.CleanUpContext == nil {
		o.CleanUpContext = context.Background()
	}
	go s.purgeLoop(o.CleanUpContext, o.CleanUpPeriod)
	return s, nil
}

// Rate returns discriminating [Rate] or global [Rate], whichever is slower.
func (d *SingleTagging) Rate() Rate {
	if d.taggedBucketMap.rate < d.basic.rate {
		return d.taggedBucketMap.rate
	}
	return d.basic.rate
}

func (d *SingleTagging) Take(r *http.Request) (err error) {
	from := time.Now()
	d.basic.mu.Lock()
	defer d.basic.mu.Unlock()

	if !d.basic.bucket.Take(
		d.basic.limit,
		d.basic.rate,
		from,
		from.Add(d.basic.interval),
	) {
		err = ErrTooManyRequests
	}
	return errors.Join(err, d.taggedBucketMap.Take(r, from))
}

func (d *SingleTagging) purgeLoop(ctx context.Context, interval time.Duration) {
	var t time.Time
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case t = <-ticker.C:
		}

		d.basic.mu.Lock()
		d.bucketMap.Purge(t)
		// for _, bm := range r.bucketMaps {
		// 	bm.Purge(t)
		// }
		d.basic.mu.Unlock()
	}
}

func (d *SingleTagging) Middleware() oakhttp.Middleware {
	return NewMiddleware(d, d.Rate())
}

func (d *SingleTagging) ObfuscatedMiddleware(displayRate Rate) oakhttp.Middleware {
	return NewMiddleware(d, displayRate)
}

type MultiTagging struct {
	basic
	taggedBucketMap []taggedBucketMap
}

func NewMultiDiscriminating(withOptions ...Option) (*MultiTagging, error) {
	o, err := newOptions(append(
		withOptions,
		func(o *options) error { // validate
			if len(o.Tagging) < 2 {
				return errors.New("tagged rate limiter must be initiated with more than one tagger")
			}
			return nil
		},
	)...)
	if err != nil {
		return nil, fmt.Errorf("cannot create tagged rate limiter: %w", err)
	}

	m := &MultiTagging{
		basic:           *o.Basic,
		taggedBucketMap: o.Tagging,
	}

	if o.CleanUpContext == nil {
		o.CleanUpContext = context.Background()
	}
	go m.purgeLoop(o.CleanUpContext, o.CleanUpPeriod)
	return m, nil
}

// Rate returns discriminating [Rate] or global [Rate], whichever is slower.
func (d *MultiTagging) Rate() (r Rate) {
	r = d.basic.rate
	for _, child := range d.taggedBucketMap {
		if child.rate < r {
			r = child.rate
		}
	}
	return
}

func (d *MultiTagging) Take(r *http.Request) (err error) {
	from := time.Now()
	d.basic.mu.Lock()
	defer d.basic.mu.Unlock()

	if !d.basic.bucket.Take(
		d.basic.limit,
		d.basic.rate,
		from,
		from.Add(d.basic.interval),
	) {
		err = ErrTooManyRequests
	}

	l := len(d.taggedBucketMap)
	cerr := make([]error, l+1)
	cerr[l] = err // last cell is basic error
	for i, child := range d.taggedBucketMap {
		cerr[i] = child.Take(r, from)
	}

	return errors.Join(cerr...)
}

func (d *MultiTagging) purgeLoop(ctx context.Context, interval time.Duration) {
	var t time.Time
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case t = <-ticker.C:
		}

		d.basic.mu.Lock()
		for _, child := range d.taggedBucketMap {
			child.bucketMap.Purge(t)
		}
		d.basic.mu.Unlock()
	}
}

func (d *MultiTagging) Middleware() oakhttp.Middleware {
	return NewMiddleware(d, d.Rate())
}

func (d *MultiTagging) ObfuscatedMiddleware(displayRate Rate) oakhttp.Middleware {
	return NewMiddleware(d, displayRate)
}

// func (t *Discriminating) Take(r *http.Request) (Rate, error) {
// 	t := time.Now()
// 	b.mu.Lock()
// 	defer b.mu.Unlock()
//
// 	remaining := b.bucket.Take(b.limit, b.rate, t, t.Add(b.interval))
// 	// log.Println("remaining", remaining)
// 	if remaining < 0 {
// 		return b.rate, ErrTooManyRequests
// 	}
// 	return b.rate, nil
// }

// type RateLimiter struct {
// 	tokenizers []Tokenizer
// 	mu         sync.Mutex
// 	global     *leakybucketmap.LeakyBucket
// 	bucketMaps []*leakybucketmap.Map
// }
