package oakratelimiter

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type limitOptions struct {
	Name                  string
	Limit                 float64
	Interval              time.Duration
	InitialAllocationSize int
}

func newLimitOptions(withOptions ...LimitOption) (*limitOptions, error) {
	o := &limitOptions{}
	for _, option := range withOptions {
		if err := option(o); err != nil {
			return nil, err
		}
	}
	return o, nil
}

type LimitOption func(*limitOptions) error

func WithName(name string) LimitOption {
	return func(o *limitOptions) error {
		if o.Name != "" {
			return errors.New("name has already been set")
		}
		if name == "" {
			return errors.New("cannot use an empty name")
		}
		o.Name = name
		return nil
	}
}

func WithDefaultName() LimitOption {
	return func(o *limitOptions) error {
		if o.Name == "" {
			return WithName("default")(o)
		}
		return nil
	}
}

func WithRate(limit float64, interval time.Duration) LimitOption {
	return func(o *limitOptions) error {
		if o.Limit != 0 || o.Interval != 0 {
			return errors.New("rate has already been set")
		}
		if limit < 1 {
			return errors.New("limit must be greater than 1")
		}
		if limit > 1<<32 {
			return errors.New("take limit is too large")
		}
		if interval < time.Millisecond*20 {
			return errors.New("interval must be greater than 20ms")
		}
		if interval > time.Hour*24 {
			return errors.New("maximum interval is 24 hours")
		}
		o.Limit = limit
		o.Interval = interval
		return nil
	}
}

func WithInitialAllocationSize(buckets int) LimitOption {
	return func(o *limitOptions) error {
		if o.InitialAllocationSize != 0 {
			return errors.New("initial allocation size is already set")
		}
		if buckets < 64 {
			return errors.New("initial allocation size must not be less than 64")
		}
		if buckets > 1<<32 {
			return errors.New("initial allocation size is too great")
		}
		o.InitialAllocationSize = buckets
		return nil
	}
}

func WithDefaultInitialAllocationSize() LimitOption {
	return func(o *limitOptions) error {
		if o.InitialAllocationSize == 0 {
			return WithInitialAllocationSize(1024)(o)
		}
		return nil
	}
}

type options struct {
	Basic          *basic
	Tagging        []taggedBucketMap
	CleanUpContext context.Context
	CleanUpPeriod  time.Duration
}

type Option func(*options) error

func newOptions(withOptions ...Option) (o *options, err error) {
	o = &options{}
	for _, option := range append(
		withOptions,
		WithDefaultCleanUpPeriod(),
	) {
		if err = option(o); err != nil {
			return nil, err
		}
	}
	if o.Basic == nil {
		return nil, errors.New("global rate limit is required")
	}
	return o, nil
}

func WithGlobalLimit(withOptions ...LimitOption) Option {
	return func(o *options) (err error) {
		if o.Basic != nil {
			return errors.New("global limit is already set")
		}
		if o.Basic, err = newBasic(withOptions...); err != nil {
			return fmt.Errorf("cannot create global limit: %w", err)
		}
		return nil
	}
}

func WithRequestTagger(t Tagger, withOptions ...LimitOption) Option {
	return func(o *options) (err error) {
		if t == nil {
			return errors.New("cannot use a <nil> discriminator")
		}
		limitOptions, err := newLimitOptions(withOptions...)
		if limitOptions.Name == "" {
			limitOptions.Name = fmt.Sprintf("discriminator№%d", len(o.Tagging)+1)
		}
		if err != nil {
			return fmt.Errorf("cannot create %q rate limiter tagger: %w", limitOptions.Name, err)
		}
		for _, existing := range o.Tagging {
			if existing.name == limitOptions.Name {
				return fmt.Errorf("rate limiter %q tagger already exists", limitOptions.Name)
			}
		}
		o.Tagging = append(o.Tagging, taggedBucketMap{
			name:      limitOptions.Name,
			limit:     limitOptions.Limit,
			interval:  limitOptions.Interval,
			rate:      NewRate(limitOptions.Limit, limitOptions.Interval),
			bucketMap: make(bucketMap, limitOptions.InitialAllocationSize),
			tagger:    t,
		})
		return nil
	}
}

func WithIPAddressTagger(withOptions ...LimitOption) Option {
	return WithRequestTagger(
		NewIPAddressTagger(),
		append(
			withOptions,
			func(o *limitOptions) error {
				if o.Name == "" {
					o.Name = "ip-address"
				}
				return nil
			},
		)...,
	)
}

func WithCookieTagger(name, noCookieValue string, withOptions ...LimitOption) Option {
	return WithRequestTagger(
		NewCookieTagger(name, noCookieValue),
		append(
			withOptions,
			func(o *limitOptions) error {
				if o.Name == "" {
					o.Name = "cookie:" + name
				}
				return nil
			},
		)...,
	)
}

func WithCleanUpPeriod(of time.Duration) Option {
	return func(o *options) error {
		if o.CleanUpPeriod != 0 {
			return errors.New("clean up period is already set")
		}
		if of < time.Second {
			return errors.New("clean up period must be greater than 1 second")
		}
		if of > time.Hour {
			return errors.New("clean up period must be less than one hour")
		}
		o.CleanUpPeriod = of
		return nil
	}
}

func WithDefaultCleanUpPeriod() Option {
	return func(o *options) error {
		if o.CleanUpPeriod == 0 {
			return WithCleanUpPeriod(time.Minute * 15)(o)
		}
		return nil
	}
}
