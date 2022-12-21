package oakhttp

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func captureResponse(h http.Handler, r *http.Request) *http.Response {
	// req := httptest.NewRequest(http.MethodGet, "/upper?word=abc", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w.Result()
}

func TestLimit(t *testing.T) {
	t.Parallel()

	cases := []struct {
		StatusCode int
		Sleep      time.Duration
	}{
		{StatusCode: http.StatusOK, Sleep: 0},
		{StatusCode: http.StatusOK, Sleep: 0},
		{StatusCode: http.StatusTooManyRequests, Sleep: time.Millisecond * 1050},
		{StatusCode: http.StatusOK, Sleep: 0},
		{StatusCode: http.StatusTooManyRequests, Sleep: time.Millisecond * 1050},
		{StatusCode: http.StatusOK, Sleep: 0},
		{StatusCode: http.StatusTooManyRequests, Sleep: 0},
	}

	limiter := Limit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello world")
	}), rate.NewLimiter(rate.Every(time.Second), 2))

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	for i, c := range cases {
		r := captureResponse(limiter, request)
		if r.StatusCode != c.StatusCode {
			t.Fatalf("rate limiter step %d failed: %d does not match %d",
				i+1, r.StatusCode, c.StatusCode)
		}
		time.Sleep(c.Sleep)
	}
}

func TestLimitByClientAddress(t *testing.T) {
	t.Parallel()

	cases := []struct {
		StatusCode int
		Sleep      time.Duration
		RemoteAddr string
	}{
		{
			StatusCode: http.StatusOK,
			Sleep:      0,
			RemoteAddr: "localhost",
		},
		{
			StatusCode: http.StatusOK,
			Sleep:      0,
			RemoteAddr: "localhost",
		},
		{
			StatusCode: http.StatusTooManyRequests,
			Sleep:      time.Millisecond * 1050,
			RemoteAddr: "localhost",
		},
		{
			StatusCode: http.StatusOK,
			Sleep:      0,
			RemoteAddr: "localhost",
		},
		{
			StatusCode: http.StatusTooManyRequests,
			Sleep:      time.Millisecond * 1050,
			RemoteAddr: "localhost",
		},
		{
			StatusCode: http.StatusOK,
			Sleep:      0,
			RemoteAddr: "localhost",
		},
		{
			StatusCode: http.StatusTooManyRequests,
			Sleep:      0,
			RemoteAddr: "localhost",
		},
	}

	limiter := LimitByClientAddress(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello world")
	}), time.Second, 2)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	for i, c := range cases {
		request.RemoteAddr = c.RemoteAddr
		r := captureResponse(limiter, request)
		if r.StatusCode != c.StatusCode {
			t.Fatalf("rate limiter step %d failed: %d does not match %d",
				i+1, r.StatusCode, c.StatusCode)
		}
		time.Sleep(c.Sleep)
	}
}
