package ratelimiter

import (
	"io"
	"net/http"
	"net/http/httptest"
)

func captureResponse(h http.Handler, r *http.Request) *http.Response {
	// req := httptest.NewRequest(http.MethodGet, "/upper?word=abc", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w.Result()
}

var testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello world")
})
