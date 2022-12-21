package oakhttp

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func Limit(h http.Handler, l *rate.Limiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// time.Sleep(time.Millisecond * 2020)
		// fmt.Println("tokens left", l.Tokens(), l.Burst())
		if !l.Allow() {
			http.Error(w, "there are too many requests", http.StatusTooManyRequests)
			return
		}
		h.ServeHTTP(w, r)
	}
}

func LimitByClientAddress(h http.Handler, d time.Duration, n int) http.HandlerFunc {
	tracker := make(map[string]*rate.Limiter)
	mu := &sync.Mutex{}
	go func() { // clean up in parallel
		// t := time.NewTicker(d * 10)
		t := time.NewTicker(d / 2)
		tokenLimit := float64(n)
		for {
			select {
			case <-t.C:
				var eliminationQueue []string
				mu.Lock()
				if len(eliminationQueue) >= n || true {
					for key, limiter := range tracker {
						if limiter.Tokens() == tokenLimit {
							eliminationQueue = append(eliminationQueue, key)
						}
					}
				}

				for _, key := range eliminationQueue {
					delete(tracker, key)
				}
				mu.Unlock()
			}
		}
	}()

	tryLimit := func(IP string) bool {
		mu.Lock()
		defer mu.Unlock()

		limit, ok := tracker[IP]
		if !ok {
			limit = rate.NewLimiter(rate.Every(d), n)
			tracker[IP] = limit
		}
		return limit.Allow()
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// time.Sleep(time.Millisecond * 2020)
		// fmt.Println("tokens left", l.Tokens(), l.Burst())
		if !tryLimit(r.RemoteAddr) {
			http.Error(w, "there are too many requests", http.StatusTooManyRequests)
			return
		}
		h.ServeHTTP(w, r)
	}
}
