package oakserver

import (
	"context"
	"math"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"golang.org/x/exp/slog"
)

func NewDebugLogger() *slog.Logger {
	return slog.New(tint.Options{
		// Level:      slog.LevelDebug,
		Level:      -math.MaxInt, // log everything
		TimeFormat: time.Kitchen,
	}.NewHandler(os.Stderr))
}

type slogAdaptor struct {
	logger *slog.Logger
	level  slog.Level
}

func (s *slogAdaptor) Write(b []byte) (n int, err error) {
	s.logger.Log(context.Background(), s.level, string(b))
	return len(b), nil
}
