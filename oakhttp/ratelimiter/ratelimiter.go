package ratelimiter

import (
	"net/http"

	"golang.org/x/exp/slog"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

func (e Error) HTTPStatusCode() int {
	return http.StatusTooManyRequests
}

const (
	ErrTooManyRequests Error = "too many incoming requests"
)

func Must(h http.Handler, err error) http.Handler {
	if err != nil {
		panic(err)
	}
	return h
}

func writeError(w http.ResponseWriter, r *http.Request) {
	slog.Warn("too many incoming requests", slog.String("ip", r.RemoteAddr))
	http.Error(w, "too many incoming requests", http.StatusTooManyRequests)
}
