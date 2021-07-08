package oakacs

import (
	"fmt"

	"go.uber.org/zap"
)

// Option sets up the access control system with all its parameters.
type Option func(acs *AccessControlSystem) error

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
