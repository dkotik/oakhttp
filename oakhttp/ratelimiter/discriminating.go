package ratelimiter

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// type discriminatingMap struct {
// 	name          string
// 	rate          Rate
// 	limit         float64
// 	interval      time.Duration
// 	buckets       map[string]*bucket
// 	discriminator Discriminator
// }
//
// func (d *discriminatingMap) Take(r *http.Request, at time.Time) error {
// 	token, err := d.discriminator(r)
// 	if err != nil {
// 		return err
// 	}
//
// 	limit = d.limit
// 	reset = at.Add(d.interval)
// 	foundBucket, ok := d.buckets[token]
// 	if !ok {
// 		remaining = limit - 1
// 		foundBucket = &bucket{
// 			expires: reset,
// 			tokens:  remaining,
// 		}
// 		d.buckets[token] = foundBucket
// 		return
// 	}
// 	if foundBucket.Take(limit, d.rate, at, reset) < 0 {
//     return &TooManyRequestsError{
//       cause: fmt.Errorf(""),
//     }
// 	}
// 	return nil
// }

type Discriminating struct {
	Basic

	interval          time.Duration
	rate              Rate
	limit             float64
	bucketMap         bucketMap
	discriminator     Discriminator
	discriminatorName string
}

// Rate returns [Discriminating] [Rate] or global [Rate], whichever is slower.
func (d *Discriminating) Rate() Rate {
	if d.rate < d.Basic.rate {
		return d.rate
	}
	return d.Basic.rate
}

func (d *Discriminating) Take(r *http.Request) (err error) {
	from := time.Now()
	to := from.Add(d.interval)
	d.Basic.mu.Lock()
	defer d.Basic.mu.Unlock()

	if !d.Basic.bucket.Take(d.Basic.limit, d.Basic.rate, from, to) {
		err = ErrTooManyRequests
	}

	token, derr := d.discriminator(r)
	if derr != nil {
		if errors.Is(derr, SkipDiscriminator) {
			return err
		}
		return errors.Join(err, derr)
	}
	if !d.bucketMap.Take(token, d.limit, d.rate, from, to) {
		derr = &TooManyRequestsError{
			cause: fmt.Errorf("discriminator %q ran out of tokens", d.discriminatorName),
		}
	}
	return errors.Join(err, derr)
}

func (d *Discriminating) purgeLoop(ctx context.Context, interval time.Duration) {
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
