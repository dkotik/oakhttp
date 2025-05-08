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
	for _, option := range append(
		withOptions,
		func(o *limitOptions) error { // validate
			if o.Limit == 0 || o.Interval == 0 {
				return errors.New("WithRate option is required")
			}
			return nil
		},
	) {
		if err := option(o); err != nil {
			return nil, err
		}
	}
	return o, nil
}

func newSupervisingLimitOptions(withOptions ...LimitOption) (*limitOptions, error) {
	return newLimitOptions(append(
		withOptions,
		WithDefaultName(),
		func(o *limitOptions) error {
			if o.InitialAllocationSize != 0 {
				return errors.New("initial allocation size option cannot be applied to the supervising rate limiter")
			}
			return nil
		},
	)...)
}

// LimitOption configures a rate limitter. [Basic] relies on one set of [LimitOption]s. [SingleTagging] and [MultiTagging] use a set for superivising rate limit and additional sets for each [Tagger].
type LimitOption func(*limitOptions) error

// WithName associates a name with a rate limiter. It is displayed only in the logs.
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

// WithDefaultName sets rate limiter name to "default."
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

// WithInitialAllocationSize sets the number of pre-allocated items for a tagged bucket map. Higher number can improve initial performance at the cost of using more memory.
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

// WithDefaultInitialAllocationSize sets initial map allocation to 1024.
func WithDefaultInitialAllocationSize() LimitOption {
	return func(o *limitOptions) error {
		if o.InitialAllocationSize == 0 {
			return WithInitialAllocationSize(1024)(o)
		}
		return nil
	}
}

type options struct {
	Supervising    *limitOptions
	Tagging        []taggedBucketMap
	CleanUpContext context.Context
	CleanUpPeriod  time.Duration
}

// Option initializes a [RateLimiter].
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
	if o.Supervising == nil {
		return nil, errors.New("WithSupervisingLimit option is required")
	}
	return o, nil
}

// WithSupervisingLimit sets the top rate limit for either [SingleTagging] or [MultiTagging] rate limiters to prevent [Tagger]s from consuming too much memory. Should be higher than the limit of any request tagger.
func WithSupervisingLimit(withOptions ...LimitOption) Option {
	return func(o *options) (err error) {
		if o.Supervising != nil {
			return errors.New("supervising limit is already set")
		}
		if o.Supervising, err = newLimitOptions(withOptions...); err != nil {
			return fmt.Errorf("cannot create supervising limit: %w", err)
		}
		return nil
	}
}

// WithRequestTagger associates a request [Tagger] with a [RateLimiter]. [Tagger]s allow you to differentiate and group requests based on their properties.
func WithRequestTagger(t Tagger, withOptions ...LimitOption) Option {
	return func(o *options) (err error) {
		if t == nil {
			return errors.New("cannot use a <nil> discriminator")
		}
		limitOptions, err := newLimitOptions(append(
			withOptions,
			WithDefaultInitialAllocationSize(),
		)...)
		if err != nil {
			return fmt.Errorf("cannot create rate limiter with tagger %+v: %w", t, err)
		}
		if limitOptions.Name == "" {
			limitOptions.Name = fmt.Sprintf("discriminator№%d", len(o.Tagging)+1)
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

// WithIPAddressTagger configures rate limiter to track requests based on client IP addresses.
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

// WithCookieTagger configures rate limiter to track requests based on a certain cookie. If [noCookieValue] is an empty string, this [Tagger] issues a [SkipTagger] sentinel value.
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

// WithCleanUpPeriod sets the frequency of map clean up. Lower value frees up more memory at the cost of CPU cycles.
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

// WithDefaultCleanUpPeriod sets clean up period to 15 minutes.
func WithDefaultCleanUpPeriod() Option {
	return func(o *options) error {
		if o.CleanUpPeriod == 0 {
			return WithCleanUpPeriod(time.Minute * 15)(o)
		}
		return nil
	}
}
