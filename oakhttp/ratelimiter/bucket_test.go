package ratelimiter

import (
	"testing"
	"time"
)

func TestRate(t *testing.T) {
	r := NewRate(1, time.Second)
	replenishmentOfOne := r.ReplishmentOfOneToken()
	if replenishmentOfOne != time.Second {
		t.Fatal("replenishment does not match", replenishmentOfOne, time.Second)
	}
}

func TestLeakyBucket(t *testing.T) {
	limit := float64(9)
	at := time.Now()
	interval := time.Second
	r := NewRate(limit, interval)

	bucket := &bucket{
		expires: at.Add(interval),
		tokens:  limit,
	}

	cases := []struct {
		Sleep     time.Duration
		Remaining float64
	}{
		{Sleep: 0, Remaining: 8},
		{Sleep: 0, Remaining: 7},
		{Sleep: 0, Remaining: 6},
		{Sleep: time.Second / 9, Remaining: 6},
		{Sleep: time.Second, Remaining: 8},
	}

	comp := floatComparator(0.1)
	for i, c := range cases {
		time.Sleep(c.Sleep)
		at := time.Now()
		bucket.Take(limit, r, at, at.Add(interval))
		if comp(bucket.tokens, c.Remaining) {
			t.Log(i+1, "remaining values match", bucket.tokens, c.Remaining)
		} else {
			t.Fatal(i+1, "remaining values do not match", bucket.tokens, c.Remaining)
		}
	}
}
