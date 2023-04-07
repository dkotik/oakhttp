package ratelimiter

import (
	"errors"
	"time"
)

type options struct {
	detectionLimit  int
	detectionPeriod time.Duration
	cleanUpPeriod   time.Duration
	maxAddressCount int
	minAddressCount int
	GlobalRate      Rate
	ObfuscatedRate  Rate
}

type Option func(*options) error

func WithLimit(count int, per time.Duration) Option {
	return func(o *options) error {
		if o.detectionLimit != 0 || o.detectionPeriod != 0 {
			return errors.New("max frequency has already been set")
		}
		if count < 1 {
			return errors.New("max frequency count must be greater than 1")
		}
		if per < time.Millisecond*20 {
			return errors.New("max frequency detection period must be greater than 20ms")
		}
		o.detectionLimit = count
		o.detectionPeriod = per
		return nil
	}
}

func WithDefaultLimit() Option {
	return func(o *options) error {
		if o.detectionLimit == 0 {
			o.detectionLimit = 50
		}
		if o.detectionPeriod == 0 {
			o.detectionPeriod = time.Minute
		}
		return nil
	}
}

func WithCleanUpPeriod(of time.Duration) Option {
	return func(o *options) error {
		if o.cleanUpPeriod != 0 {
			return errors.New("clean up period is already set")
		}
		if of < time.Second {
			return errors.New("clean up period must be greater than 1 second")
		}
		if of > time.Minute*15 {
			return errors.New("clean up period must be less than 15 minutes")
		}
		o.cleanUpPeriod = of
		return nil
	}
}

func WithDefaultCleanUpPeriod() Option {
	return func(o *options) error {
		if o.cleanUpPeriod == 0 {
			o.cleanUpPeriod = time.Minute * 15
		}
		return nil
	}
}

func WithMaximumAddressCount(of int) Option {
	return func(o *options) error {
		if o.maxAddressCount != 0 {
			return errors.New("maximum address count is already set")
		}
		if of < 1 {
			return errors.New("maximum address count must be greater than 1")
		}
		o.maxAddressCount = of
		return nil
	}
}

func WithDefaultMaximumAddressCount() Option {
	return func(o *options) error {
		if o.maxAddressCount == 0 {
			o.maxAddressCount = 500_000
		}
		return nil
	}
}

func WithMinimumAddressCount(of int) Option {
	return func(o *options) error {
		if o.minAddressCount != 0 {
			return errors.New("minimum address count is already set")
		}
		if of < 1 {
			return errors.New("minimum address count must be greater than 1")
		}
		o.maxAddressCount = of
		return nil
	}
}

func WithDefaultMinimumAddressCount() Option {
	return func(o *options) error {
		if o.minAddressCount != 0 {
			return nil
		}
		if o.maxAddressCount != 0 {
			o.minAddressCount = o.maxAddressCount / 5
			return nil
		}
		o.minAddressCount = 5_000
		return nil
	}
}
