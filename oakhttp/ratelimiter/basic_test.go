package ratelimiter

import (
	"testing"
	"time"
)

func TestBasicRateLimiter(t *testing.T) {
	limit := float64(2)
	interval := time.Second
	rl, err := NewBasic(WithRate(limit, interval))
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		Sleep time.Duration
		Fails bool
	}{
		{Sleep: 0, Fails: false},
		{Sleep: 0, Fails: false},
		{Sleep: 0, Fails: true},
		{Sleep: 0, Fails: true},
		{Sleep: time.Millisecond * 500, Fails: false},
		{Sleep: time.Millisecond * 500, Fails: false},
		{Sleep: 0, Fails: true},
	}

	for i, c := range cases {
		time.Sleep(c.Sleep)
		err = rl.Take(nil)
		if err != nil && !c.Fails {
			t.Fatal(i+1, "rate limiter failed when not expecting it", rl.tokens, err)
		}
	}
}
