package oakhttp

import (
	"embed"
	"log/slog"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

//go:embed templates
var templates embed.FS

type Middleware func(http.Handler) http.Handler

// ApplyMiddleware applies [Middleware] in reverse to preserve logical order.
func ApplyMiddleware(h http.Handler, mws []Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func NewDebugLogger() *slog.Logger {
	return slog.New(NewTracingHandler(
		tint.NewHandler(os.Stderr, &tint.Options{
			// Level:      slog.LevelDebug,
			Level:      slog.Level(-math.MaxInt), // log everything
			TimeFormat: time.Kitchen,
		}))).With(
		slog.String("commit", vcsCommit()),
	)
}
