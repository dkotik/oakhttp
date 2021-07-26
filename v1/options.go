package oakacs

import (
	"fmt"

	"go.uber.org/zap"
)

// Option sets up the access control system with all its parameters.
type Option func(acs *AccessControlSystem) error

// WithOptions combines a group of options into one. This is a helper for option sets and the constructor.
func WithOptions(options ...Option) Option {
	return func(acs *AccessControlSystem) (err error) {
		for _, option := range withOptions {
			if err = option(acs); err != nil {
				return err
			}
		}
		return nil
	}
}

// WithHasher attaches possible secret hashers to the ACS.
func WithHasher(name string, hasher Hasher) Option {
	return func(acs *AccessControlSystem) (err error) {
		if _, ok := acs.hashers[name]; ok {
			return fmt.Errorf("hasher %q is already set", name)
		}
		acs.hashers[name] = hasher
		return nil
	}
}

// WithSubscribers adds a list of event listeners that the ACS broadcasts.
func WithSubscribers(c ...chan<- (Event)) Option {
	return func(acs *AccessControlSystem) (err error) {
		acs.subscribers = make([]chan (Event), len(c))
		for i, channel := range c {
			acs.subscribers[i] = channel
		}
		return nil
	}
}

// WithLogger attaches a zap logger to the ACS.
func WithLogger(logger *zap.Logger) Option {
	return func(acs *AccessControlSystem) (err error) {
		if logger == nil {
			logger, err = zap.NewDevelopment()
			if err != nil {
				return
			}
		}
		acs.logger = logger
		return
	}
}
