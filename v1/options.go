package oakacs

import (
	"go.uber.org/zap"
)

// Option sets up the access control system with all its parameters.
type Option func(acs *AccessControlSystem) error

// WithOptions combines a group of options into one. This is a helper for option sets and the constructor.
func WithOptions(options ...Option) Option {
	return func(acs *AccessControlSystem) (err error) {
		for _, option := range options {
			if err = option(acs); err != nil {
				return err
			}
		}
		return nil
	}
}

// WithSubscribers adds subscribers to the ACS broadcasts. Use buffered channels to avoid missing events, because the put is non-blocking.
func WithSubscribers(c ...chan<- (Event)) Option {
	return func(acs *AccessControlSystem) (err error) {
		if acs.subscribers == nil {
			acs.subscribers = make([]chan<- (Event), len(c))
		}
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
		// TODO: subscribe to events using a buffered channel and an adapter
		// acs.logger = logger
		return
	}
}
