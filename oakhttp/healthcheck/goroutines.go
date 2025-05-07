package healthcheck

import (
	"context"
	"errors"
	"runtime"

	"golang.org/x/exp/slog"
)

var ErrTooManyGoRoutines = errors.New("go routine number overflowed")

type goRoutinesCheck struct {
	warnAt int
	failAt int
}

func (g *goRoutinesCheck) ReportStatus(ctx context.Context, logger *slog.Logger) (report any, err error) {
	current := runtime.NumGoroutine()
	if current >= g.failAt {
		return current, ErrTooManyGoRoutines
	} else if current >= g.warnAt {
		// TODO: replace with warning after upgrading slog:
		logger.ErrorCtx(
			ctx,
			"go routine number is worryingly high",
			slog.Int("warnAt", g.warnAt),
			slog.Int("failAt", g.failAt),
			slog.Int("current", current),
		)
	}
	return current, nil
}

func WithGoRoutineLimitOf(warnAt, failAt int) Option {
	return func(o *options) error {
		if warnAt < 5 {
			return errors.New("go routine warning limit cannot be less than 5")
		}
		if failAt < 5 {
			return errors.New("go routine failing limit cannot be less than 5")
		}
		if warnAt > failAt {
			return errors.New("go routine warning limit cannot be greater than failure limit")
		}
		return WithCheck("goroutines", &goRoutinesCheck{
			warnAt: warnAt,
			failAt: failAt,
		})(o)
	}
}
