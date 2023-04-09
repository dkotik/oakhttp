package ratelimiter

import (
	"net/http"
	"testing"

	"github.com/dkotik/oakacs/oakhttp"
)

func LoadTest(t *testing.T, fn func() error, rate Rate, times int) {
	// interval := rate.ReplishmentOfOneToken()
	// for i := 0; i < times; i++ {
	// 	requests[i] = httptest.NewRequest(http.MethodGet, "/", nil)
	// }
	// return nil
}

func MiddlewareTest(m oakhttp.Middleware, good, bad Rate) error {
	// cases := []struct {
	// 	StatusCode int
	// 	Sleep      time.Duration
	// 	RemoteAddr string
	// }{
	// 	{
	// 		StatusCode: http.StatusOK,
	// 		Sleep:      0,
	// 		RemoteAddr: "localhost",
	// 	},
	// 	{
	// 		StatusCode: http.StatusOK,
	// 		Sleep:      0,
	// 		RemoteAddr: "localhost",
	// 	},
	// 	{
	// 		StatusCode: http.StatusTooManyRequests,
	// 		Sleep:      time.Millisecond * 1050,
	// 		RemoteAddr: "localhost",
	// 	},
	// 	{
	// 		StatusCode: http.StatusOK,
	// 		Sleep:      0,
	// 		RemoteAddr: "localhost",
	// 	},
	// 	{
	// 		StatusCode: http.StatusTooManyRequests,
	// 		Sleep:      time.Millisecond * 1050,
	// 		RemoteAddr: "localhost",
	// 	},
	// 	{
	// 		StatusCode: http.StatusOK,
	// 		Sleep:      0,
	// 		RemoteAddr: "localhost",
	// 	},
	// 	{
	// 		StatusCode: http.StatusTooManyRequests,
	// 		Sleep:      0,
	// 		RemoteAddr: "localhost",
	// 	},
	// }
	//
	// handler := Must(New(
	// 	testHandler,
	// 	WithLimit(2, time.Second),
	// ))
	//
	// request := httptest.NewRequest(http.MethodGet, "/", nil)
	// for i, c := range cases {
	// 	request.RemoteAddr = c.RemoteAddr
	// 	r := captureResponse(handler, request)
	// 	if r.StatusCode != c.StatusCode {
	// 		t.Fatalf("rate limiter step %d failed: %d does not match %d",
	// 			i+1, r.StatusCode, c.StatusCode)
	// 	}
	// 	time.Sleep(c.Sleep)
	// }

	return nil
}

func RateLimiterTest(r RateLimiter, good, bad Rate, requests ...*http.Request) error {
	// interval := good.ReplishmentOfOneToken()
	// if len(requests) == 0 {
	// 	requests = make([]*http.Request, 10)
	// 	for i := 0; i < 10; i++ {
	// 		requests[i] = httptest.NewRequest(http.MethodGet, "/", nil)
	// 	}
	// }
	//
	// var err error
	// for i, request := range requests {
	// 	if err = r.Take(request); err != nil {
	// 		return fmt.Errorf("case %d: rate limited did not pass on good rate: %w", i+1, err)
	// 	}
	// 	time.Sleep(interval)
	// }

	// interval = bad.ReplishmentOfOneToken()
	// for i, request := range requests {
	//   if err = r.Take(request); err != nil {
	//     return fmt.Errorf("case %d: rate limited did not pass on good rate: %w", i+1, err)
	//   }
	//   time.Sleep(interval)
	// }

	return nil
}
