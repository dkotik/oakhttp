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
		return WithListener(&authorizationGrantLogger{
			level:  level,
			logger: l,
		})(o)
	}
}

type authorizationGrantLogger struct {
	level  slog.Level
	logger *slog.Logger
}

func (s *authorizationGrantLogger) AuthorizationGranted(
	ctx context.Context,
	intents []Intent,
	policies []Policy,
	role Role,
) {
	s.logger.Log(
		ctx,
		s.level,
		"authorization granted",
		slog.Any("intents", intents),
		slog.Any("policies", policies),
		slog.Any("role", role),
	)
}

func (s *authorizationGrantLogger) AuthorizationDenied(
	ctx context.Context,
	intents []Intent,
	policies []Policy,
	role Role,
) {
	// do nothing, because errors should be logged upsteam when they are handled
}

func (s *authorizationGrantLogger) AuthorizationFailed(
	ctx context.Context,
	intents Intent,
	policies Policy,
	role Role,
	err error,
) {
	// do nothing, because errors should be logged upsteam when they are handled
}
