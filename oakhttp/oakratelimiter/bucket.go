package oakratelimiter

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"
)

// bucket tracks remaining tokens and limit expiration.
type bucket struct {
	expires time.Time
	tokens  float64
}

// Expires returns true if the bucket is expired at given [time.Time].
func (b *bucket) Expires(at time.Time) bool {
	return b.expires.Before(at)
}

// Take removes one token from the bucket. If the bucket is fresh, some fractional amount of tokens is also replenished according to [Rate] over the transpired time since the previous take from the bucket. Must run inside mutex lock.
func (b *bucket) Take(limit float64, r Rate, from, to time.Time) bool {
	if b.Expires(from) { // reset
		b.tokens = limit - 1
		b.expires = to
		return true
	}

	replenished := b.tokens + r.ReplenishedTokens(b.expires, to)
	b.expires = to
	if replenished < 1 { // nothing to take
		b.tokens = replenished
		return false
	}

	b.tokens = replenished - 1
	return true
}

// Rate is the number of tokens leaked or replenished per nanosecond.
type Rate float64

// NewRate creates a [Rate] based on expected limit and a given time interval.
func NewRate(limit float64, interval time.Duration) Rate {
	if interval == 0 {
		return Rate(math.Inf(1))
	}
	return Rate(limit / float64(interval.Nanoseconds()))
}

// ReplenishedTokens returns fractional token amount based on time passed.
func (r Rate) ReplenishedTokens(from, to time.Time) float64 {
	return float64(to.Sub(from).Nanoseconds()) * float64(r)
}

// bucketMap associates a [Tagger] tag to a [bucket].
type bucketMap map[string]*bucket

// Take locates the proper [bucket] by tag and takes one token from it. If the bucket does not exist, a new one is added to the [bucketMap]. Must run inside mutex lock.
func (m bucketMap) Take(tag string, limit float64, r Rate, from, to time.Time) bool {
	foundBucket, ok := m[tag]
	if !ok {
		foundBucket = &bucket{
			expires: to,
			tokens:  limit - 1,
		}
		m[tag] = foundBucket
		return true
	}
	return foundBucket.Take(limit, r, from, to)
}

// Purge removes all tokens that are expired by given [time.Time].
func (m bucketMap) Purge(to time.Time) {
	for k, bucket := range m {
		if bucket.Expires(to) {
			delete(m, k)
		}
	}
}

// taggedBucketMap enforces a rate limit on a [bucketMap] using a [Tagger].
type taggedBucketMap struct {
	name      string
	interval  time.Duration
	rate      Rate
	limit     float64
	bucketMap bucketMap
	tagger    Tagger
}

// Take locates the proper [bucket] by tag and takes one token from it. If the bucket does not exist, a new one is added to the [bucketMap]. Must run inside mutex lock.
func (d *taggedBucketMap) Take(r *http.Request, from time.Time) (err error) {
	to := from.Add(d.interval)

	tag, err := d.tagger(r)
	if err != nil {
		if errors.Is(err, SkipTagger) {
			return nil
		}
		return fmt.Errorf("tagger %q failed to execute: %w", d.name, err)
	}
	if !d.bucketMap.Take(tag, d.limit, d.rate, from, to) {
		return fmt.Errorf("tagger %q maxed out on tag: %s", d.name, tag)
	}
	return nil
}
