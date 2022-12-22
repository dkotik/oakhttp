package ratelimiter

import (
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/time/rate"
)

type basic struct {
	http.Handler
	limiter *rate.Limiter
}

func (b *basic) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !b.limiter.Allow() {
		writeError(w, r)
		return
	}
	b.Handler.ServeHTTP(w, r)
}

// NewBasic creates a flat atomic rate limiter.
func NewBasic(h http.Handler, withOptions ...Option) (http.Handler, error) {
	o := &options{}
	var err error
	for _, option := range append(
		withOptions,
		WithDefaultLimit(),
		func(o *options) error {
			// validate options as the last step
			if h == nil {
				return errors.New("http.Handler is required")
			}
			if o.detectionPeriod == 0 || o.detectionLimit == 0 {
				return errors.New("WithLimit is a required option")
			}
			if o.cleanUpPeriod != 0 {
				return errors.New("WithCleanupPeriod is not relevant for the basic rate limiter")
			}
			if o.maxAddressCount != 0 {
				return errors.New("WithMaximumAddressCount is not relevant for the basic rate limiter")
			}
			if o.minAddressCount != 0 {
				return errors.New("WithMinimumAddressCount is not relevant for the basic rate limiter")
			}
			return nil
		},
	) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("basic rate limiter initialization failed: %w", err)
		}
	}

	return &basic{
		Handler: h,
		limiter: rate.NewLimiter(rate.Every(o.detectionPeriod), o.detectionLimit),
	}, nil
}
