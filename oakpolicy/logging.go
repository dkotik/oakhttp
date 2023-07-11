package oakpolicy

import (
	"context"
	"errors"

	"golang.org/x/exp/slog"
)

func Log(l *slog.Logger, p Policy) Policy {
	if l == nil {
		l = slog.Default()
	}
	if p == nil {
		panic(ErrNilPolicy)
	}
	l = p.Logger(l)
	return func(ctx context.Context, a Action, r Resource) (err error) {
		err = p(ctx, a, r)
		if errors.Is(err, Allow) {
			l.InfoCtx(ctx,
				"allowed by policy",
				slog.String("action", string(a)),
				slog.Any("policy", p),
				slog.Any("resource", r.DomainPath()),
			)
		} else if errors.Is(err, Deny) {
			l.ErrorCtx(ctx,
				"denied by policy",
				slog.String("action", string(a)),
				slog.Any("policy", p),
				slog.Any("resource", r.DomainPath()),
			)
		}
		return err
	}
}

func LogAllowedActions(l *slog.Logger, level slog.Level, p Policy) Policy {
	if l == nil {
		l = slog.Default()
	}
	if p == nil {
		panic(ErrNilPolicy)
	}
	l = p.Logger(l)
	return func(ctx context.Context, a Action, r Resource) (err error) {
		err = p(ctx, a, r)
		if errors.Is(err, Allow) {
			l.Log(ctx, level,
				"allowed by policy",
				slog.String("action", string(a)),
				slog.Any("policy", p),
				slog.Any("resource", r.DomainPath()),
			)
		}
		return err
	}
}

func LogDeniedActions(l *slog.Logger, level slog.Level, p Policy) Policy {
	if l == nil {
		l = slog.Default()
	}
	if p == nil {
		panic(ErrNilPolicy)
	}
	l = p.Logger(l)
	return func(ctx context.Context, a Action, r Resource) (err error) {
		err = p(ctx, a, r)
		if errors.Is(err, Deny) {
			l.Log(ctx, level,
				"denied by policy",
				slog.String("action", string(a)),
				slog.Any("policy", p),
				slog.Any("resource", r.DomainPath()),
			)
		}
		return err
	}
}
