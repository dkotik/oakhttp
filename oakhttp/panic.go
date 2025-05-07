package oakhttp

import (
	"net/http"
	"runtime"

	"golang.org/x/exp/slog"
)

type PanicHandler struct {
	next   Handler
	logger *slog.Logger
}

func NewPanicHandler(next Handler, l *slog.Logger) *PanicHandler {
	if next == nil {
		panic("panic middleware given a <nil> handler")
	}
	if l == nil {
		l = slog.Default()
	}
	return &PanicHandler{next: next, logger: l}
}

func (p *PanicHandler) ServeHyperText(
	w http.ResponseWriter, r *http.Request,
) (err error) {
	defer func() {
		if recovery := recover(); recovery != nil {
			// TODO: would debug.Stack() be better?
			buf := make([]byte, 10<<10)
			n := runtime.Stack(buf, false)
			p.logger.Error(
				"recovered from panic",
				slog.Int("code", http.StatusInternalServerError),
				slog.Any("error", recovery),
				// slog.Any("handler_error", err),
				slog.String("trace", string(buf[:n])),
				slog.Any("address", r.RemoteAddr),
				slog.Group("request",
					slog.String("host", r.URL.Hostname()),
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
				),
			)
			http.Error( // obfuscate error message to the user
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
		} else if err != nil {
			code, _ := UnwrapError(err)
			p.logger.Error(
				"HTTP request failed",
				slog.Int("code", code),
				slog.Any("error", err),
				slog.Any("address", r.RemoteAddr),
				slog.Group("request",
					slog.String("host", r.URL.Hostname()),
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
				),
			)
			http.Error( // obfuscate error message to the user
				w,
				http.StatusText(code),
				code,
			)
		}
	}()
	return p.next.ServeHyperText(w, r)
}

func (p *PanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := p.ServeHyperText(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
