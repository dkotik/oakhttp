package oakrbac

import (
	"context"
	"errors"

	"golang.org/x/exp/slog"
)

func WithSlogLogger(l *slog.Logger, level slog.Level) Option {
	return func(o *options) error {
		if l == nil {
			return errors.New("cannot use a <nil> slog.Logger")
		}
		return WithListener(&slogAdapter{
			level:  level,
			logger: l,
		})(o)
	}
}

type slogAdapter struct {
	level  slog.Level
	logger *slog.Logger
}

func (s *slogAdapter) Listen(
	ctx context.Context,
	e *Event,
) {
	s.logger.Log(
		ctx,
		s.level,
		e.String(),
		slog.Any("authorization", e),
	)
}
