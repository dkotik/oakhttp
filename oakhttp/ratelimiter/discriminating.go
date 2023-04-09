package ratelimiter

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dkotik/oakacs/oakhttp"
)

type discriminating struct {
	name          string
	interval      time.Duration
	rate          Rate
	limit         float64
	bucketMap     bucketMap
	discriminator Discriminator
}

// Must run inside mutex lock.
func (d *discriminating) Take(r *http.Request, from time.Time) (err error) {
	to := from.Add(d.interval)

	tag, err := d.discriminator(r)
	if err != nil {
		if errors.Is(err, SkipDiscriminator) {
			return nil
		}
		return fmt.Errorf("discriminator %q failed to execute: %w", d.name, err)
	}
	if !d.bucketMap.Take(tag, d.limit, d.rate, from, to) {
		return fmt.Errorf("discriminator %q maxed out on tag: %s", d.name, tag)
	}
	return nil
}

type SingleDiscriminating struct {
	Basic
	discriminating
}

func NewSingleDiscriminating(withOptions ...Option) (*SingleDiscriminating, error) {
	o, err := newOptions(append(
		withOptions,
		func(o *options) error { // validate
			if len(o.Discriminating) != 1 {
				return errors.New("single-discriminating rate limiter must be initiated with exactly one tagger")
			}
			return nil
		},
	)...)
	if err != nil {
		return nil, fmt.Errorf("cannot create single-discriminating rate limiter: %w", err)
	}

	s := &SingleDiscriminating{
		Basic:          *o.Basic,
		discriminating: *o.Discriminating[0],
	}

	if o.CleanUpContext == nil {
		o.CleanUpContext = context.Background()
	}
	go s.purgeLoop(o.CleanUpContext, o.CleanUpPeriod)
	return s, nil
}

// Rate returns discriminating [Rate] or global [Rate], whichever is slower.
func (d *SingleDiscriminating) Rate() Rate {
	if d.discriminating.rate < d.Basic.rate {
		return d.discriminating.rate
	}
	return d.Basic.rate
}

func (d *SingleDiscriminating) Take(r *http.Request) (err error) {
	from := time.Now()
	d.Basic.mu.Lock()
	defer d.Basic.mu.Unlock()

	if !d.Basic.bucket.Take(
		d.Basic.limit,
		d.Basic.rate,
		from,
		from.Add(d.Basic.interval),
	) {
		err = ErrTooManyRequests
	}
	return errors.Join(err, d.discriminating.Take(r, from))
}

func (d *SingleDiscriminating) purgeLoop(ctx context.Context, interval time.Duration) {
	var t time.Time
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case t = <-ticker.C:
		}

		d.Basic.mu.Lock()
		d.bucketMap.Purge(t)
		// for _, bm := range r.bucketMaps {
		// 	bm.Purge(t)
		// }
		d.Basic.mu.Unlock()
	}
}

func (d *SingleDiscriminating) Middleware() oakhttp.Middleware {
	return NewMiddleware(d, d.Rate())
}

func (d *SingleDiscriminating) ObfuscatedMiddleware(displayRate Rate) oakhttp.Middleware {
	return NewMiddleware(d, displayRate)
}

type MultiDiscriminating struct {
	Basic
	discriminating []*discriminating
}

func NewMultiDiscriminating(withOptions ...Option) (*MultiDiscriminating, error) {
	o, err := newOptions(append(
		withOptions,
		func(o *options) error { // validate
			if len(o.Discriminating) < 2 {
				return errors.New("multi-discriminating rate limiter must be initiated with more than one tagger")
			}
			return nil
		},
	)...)
	if err != nil {
		return nil, fmt.Errorf("cannot create multi-discriminating rate limiter: %w", err)
	}

	m := &MultiDiscriminating{
		Basic:          *o.Basic,
		discriminating: o.Discriminating,
	}

	if o.CleanUpContext == nil {
		o.CleanUpContext = context.Background()
	}
	go m.purgeLoop(o.CleanUpContext, o.CleanUpPeriod)
	return m, nil
}

// Rate returns discriminating [Rate] or global [Rate], whichever is slower.
func (d *MultiDiscriminating) Rate() (r Rate) {
	r = d.Basic.rate
	for _, child := range d.discriminating {
		if child.rate < r {
			r = child.rate
		}
	}
	return
}

func (d *MultiDiscriminating) Take(r *http.Request) (err error) {
	from := time.Now()
	d.Basic.mu.Lock()
	defer d.Basic.mu.Unlock()

	if !d.Basic.bucket.Take(
		d.Basic.limit,
		d.Basic.rate,
		from,
		from.Add(d.Basic.interval),
	) {
		err = ErrTooManyRequests
	}

	l := len(d.discriminating)
	cerr := make([]error, l+1)
	cerr[l] = err // last cell is basic error
	for i, child := range d.discriminating {
		cerr[i] = child.Take(r, from)
	}

	return errors.Join(cerr...)
}

func (d *MultiDiscriminating) purgeLoop(ctx context.Context, interval time.Duration) {
	var t time.Time
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case t = <-ticker.C:
		}

		d.Basic.mu.Lock()
		for _, child := range d.discriminating {
			child.bucketMap.Purge(t)
		}
		d.Basic.mu.Unlock()
	}
}

func (d *MultiDiscriminating) Middleware() oakhttp.Middleware {
	return NewMiddleware(d, d.Rate())
}

func (d *MultiDiscriminating) ObfuscatedMiddleware(displayRate Rate) oakhttp.Middleware {
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
