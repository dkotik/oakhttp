package ratelimiter

import (
	"math"
	"time"
)

type bucket struct {
	expires time.Time
	tokens  float64
}

func (b *bucket) Expires(at time.Time) bool {
	return b.expires.Before(at)
}

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

type Rate float64

func NewRate(limit float64, interval time.Duration) Rate {
	if interval == 0 {
		return Rate(math.Inf(1))
	}
	return Rate(limit / float64(interval.Nanoseconds()))
}

func (r Rate) ReplenishedTokens(from, to time.Time) float64 {
	return float64(to.Sub(from).Nanoseconds()) * float64(r)
}

func (r Rate) ReplishmentOfOneToken() time.Duration {
	return time.Nanosecond*time.Duration(1/r) + 1
}

type bucketMap map[string]*bucket

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

func (m bucketMap) Purge(to time.Time) {
	for k, bucket := range m {
		if bucket.Expires(to) {
			delete(m, k)
		}
	}
}
