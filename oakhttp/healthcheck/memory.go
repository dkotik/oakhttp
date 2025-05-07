package healthcheck

import (
	"context"
	"errors"
	"runtime"

	"golang.org/x/exp/slog"
)

var ErrOutOfMemory = errors.New("memory usage exceeded set limits")

type memoryCheck struct {
	warnAt uint64
	failAt uint64
}

func (m memoryCheck) ReportStatus(ctx context.Context, logger *slog.Logger) (report any, err error) {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	if stats.TotalAlloc >= m.failAt {
		return stats.TotalAlloc, ErrOutOfMemory
	} else if stats.TotalAlloc >= m.warnAt {
		// TODO: replace with warning after upgrading slog:
		logger.ErrorCtx(
			ctx,
			"memory usage is worryingly high",
			slog.Uint64("warnAt", m.warnAt),
			slog.Uint64("failAt", m.failAt),
			slog.Uint64("current", stats.TotalAlloc),
		)
	} else {
		logger.DebugCtx(
			ctx,
			"memory usage is normal",
			slog.Uint64("warnAt", m.warnAt),
			slog.Uint64("failAt", m.failAt),
			slog.Uint64("current", stats.TotalAlloc),
		)
	}
	return stats.TotalAlloc, nil
}

func WithMemoryGBLimitOf(warnAtGB, failAtGB float32) Option {
	return func(o *options) error {
		if warnAtGB < 0.2 {
			return errors.New("memory warning limit cannot be less than 0.2GB")
		}
		if failAtGB < 0.2 {
			return errors.New("memory failing limit cannot be less than 0.2GB")
		}
		if warnAtGB > failAtGB {
			return errors.New("memory warning limit cannot be greater than the failure limit")
		}
		return WithCheck("memory", &memoryCheck{
			warnAt: uint64(warnAtGB * 1024 * 1024 * 1024),
			failAt: uint64(failAtGB * 1024 * 1024 * 1024),
		})(o)
	}
}
