package store

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const maximumRemovalFrequency = time.Hour * 24

type options struct {
	retainValuesFor          time.Duration
	removeExpiredValuesEvery time.Duration
	removalContext           context.Context
	valueLimit               int
}

type Option func(*options) error

func newOptions(from []Option) (*options, error) {
	o := &options{}
	var err error
	for i, option := range append(
		from,
		WithDefaultValueRetention(),
		WithDefaultRemovalFrequency(),
		WithDefaultMaximumValueCount(),
		WithDefaultRemovalContext(),
	) {
		if option == nil {
			return nil, fmt.Errorf("cannot use store option #%d: <nil>", i)
		}
		if err = option(o); err != nil {
			return nil, fmt.Errorf("cannot use store option #%d: %w", i, err)
		}
	}
	return o, nil
}

func WithValueRetentionFor(d time.Duration) Option {
	return func(o *options) error {
		if d < time.Millisecond*100 {
			return errors.New("value retention cannot be less than 100ms")
		}
		if d > time.Hour*24*7*4 {
			return errors.New("value retention cannot exceed a month")
		}
		if o.retainValuesFor != 0 {
			return errors.New("value retention is already set")
		}
		o.retainValuesFor = d
		return nil
	}
}

func WithDefaultValueRetention() Option {
	return func(o *options) error {
		if o.retainValuesFor != 0 {
			return nil
		}
		return WithValueRetentionFor(time.Minute * 5)(o)
	}
}

func WithRemovalFrequencyOf(d time.Duration) Option {
	return func(o *options) error {
		if d < time.Second {
			return errors.New("value removal cannot be less than a second")
		}
		if d > maximumRemovalFrequency {
			return errors.New("value removal cannot exceed a day")
		}
		if o.removeExpiredValuesEvery != 0 {
			return errors.New("value removal is already set")
		}
		o.removeExpiredValuesEvery = d
		return nil
	}
}

func WithDefaultRemovalFrequency() Option {
	return func(o *options) error {
		if o.removeExpiredValuesEvery != 0 {
			return nil
		}
		suggested := o.retainValuesFor * 2
		if suggested > maximumRemovalFrequency {
			suggested = maximumRemovalFrequency
		}
		return WithRemovalFrequencyOf(suggested)(o)
	}
}

func WithRemovalContext(ctx context.Context) Option {
	return func(o *options) error {
		if ctx == nil {
			return errors.New("cannot use a <nil> context")
		}
		if o.removalContext != nil {
			return errors.New("removal context is already set")
		}
		o.removalContext = ctx
		return nil
	}
}

func WithDefaultRemovalContext() Option {
	return func(o *options) error {
		if o.removalContext != nil {
			return nil
		}
		o.removalContext = context.Background()
		return nil
	}
}

func WithMaximumValueCount(limit int) Option {
	return func(o *options) error {
		if limit < 1 {
			return errors.New("cannot have a value limit of less than 1")
		}
		if o.valueLimit != 0 {
			return errors.New("value limit is already set")
		}
		o.valueLimit = limit
		return nil
	}
}

func WithDefaultMaximumValueCount() Option {
	return func(o *options) error {
		if o.valueLimit != 0 {
			return nil
		}
		return WithMaximumValueCount(1 << 16)(o)
	}
}
