package healthcheck

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/exp/slog"
)

type Authenticator func(*http.Request) error

type options struct {
	authenticator Authenticator
	frequency     time.Duration
	limit         time.Duration
	names         []string
	checks        []HealthCheck
	logger        *slog.Logger
}

type Option func(*options) error

func WithDefaultLogger() Option {
	return func(o *options) error {
		if o.logger != nil {
			return nil
		}
		return WithLogger(slog.Default())(o)
	}
}

func WithDefaultFrequencyEveryFiveMinutes() Option {
	return func(o *options) error {
		if o.frequency > 0 {
			return nil
		}
		return WithFrequency(time.Minute * 5)(o)
	}
}

func WithDefaultLimitHalfOfFrequency() Option {
	return func(o *options) error {
		if o.limit > 0 {
			return nil
		}
		if o.frequency == 0 {
			return errors.New("default frequency must be set before the default limit")
		}
		return WithLimit(o.frequency / 2)(o)
	}
}

func WithFrequency(d time.Duration) Option {
	return func(o *options) error {
		if d < time.Second {
			return errors.New("health check frequency higher than once per second is unsafe")
		}
		if d > time.Hour*24 {
			return errors.New("health check frequency lower that once per day is unsafe")
		}
		if o.frequency != 0 {
			return errors.New("health check frequency is already set")
		}
		o.frequency = d
		return nil
	}
}

func WithLimit(d time.Duration) Option {
	return func(o *options) error {
		if d < time.Second {
			return errors.New("health check limit cannot be less than one second")
		}
		if d > time.Hour {
			return errors.New("health check limit greater than one hour is unsafe")
		}
		if o.limit != 0 {
			return errors.New("health check frequency is already set")
		}
		o.limit = d
		return nil
	}
}

func WithCheck(name string, h HealthCheck) Option {
	return func(o *options) error {
		if name == "" {
			return errors.New("cannot use a health with an empty name")
		}
		if h == nil {
			return errors.New("cannot use a <nil> health check")
		}
		for _, existing := range o.names {
			if existing == name {
				return fmt.Errorf("health check %q is already set", name)
			}
		}
		o.names = append(o.names, name)
		o.checks = append(o.checks, h)
		return nil
	}
}

func WithAuthenticator(f Authenticator) Option {
	return func(o *options) error {
		if f == nil {
			return errors.New("cannot use a <nil> authenticator")
		}
		if o.authenticator != nil {
			return errors.New("authenticator is already set")
		}
		o.authenticator = f
		return nil
	}
}

func WithAuthorizationToken(t string) Option {
	return WithAuthenticator(
		func(r *http.Request) error {
			_, token, _ := strings.Cut(r.Header.Get("Authorization"), " ")
			if token != t {
				return ErrTokenRejected
			}
			return nil
		},
	)
}

func WithLogger(logger *slog.Logger) Option {
	return func(o *options) error {
		if logger == nil {
			return errors.New("cannot use a <nil> logger")
		}
		if o.logger != nil {
			return errors.New("logger is already set")
		}
		o.logger = logger
		return nil
	}
}
