package ratelimiter

import (
	"testing"
	"time"
)

func TestBasicRateLimiter(t *testing.T) {
	limit := float64(2)
	interval := time.Second
	rl := NewBasic(limit, interval)

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

	var err error
	for i, c := range cases {
		time.Sleep(c.Sleep)
		err = rl.Take(nil)
		if err != nil && !c.Fails {
			t.Fatal(i+1, "rate limiter failed when not expecting it", rl.tokens, err)
		}
	}

	// rl = NewBasic(0, 0) // test non-sence cases
	// for i := range cases {
	// 	err = rl.Take(nil)
	// 	if err == nil {
	// 		t.Fatal(i+1, "rate limiter passed on a non-sence case")
	// 	}
	// }
}
