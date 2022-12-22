package ratelimiter

import (
	"net/http"

	"golang.org/x/exp/slog"
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
