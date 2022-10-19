package oakrbac

import (
	"errors"
	"fmt"
)

type (
	// Option customizes the [RBAC] constructor [New].
	Option func(RBAC) error

	// RBAC is a simple Role Based Access Control system.
	RBAC map[string]Role
)

// Must panics if an error is associated with [RBAC] constructor. Use together with [New].
func Must(r RBAC, err error) RBAC {
	if err != nil {
		panic(err)
	}
	return r
}

// New builds an [RBAC] using provided [Option] set.
func New(options ...Option) (rbac RBAC, err error) {
	rbac = make(RBAC)
	for _, option := range options {
		if err = option(rbac); err != nil {
			return nil, fmt.Errorf("cannot create OakRBAC: %w", err)
		}
	}
	return rbac, nil
}

// GetRole returns an authorization role by name. If the role is not registered, default role is returned instead.
func (r RBAC) GetRole(name string) (Role, error) {
	role, ok := r[name]
	if ok {
		return role, nil
	}
	return nil, &AccessDeniedError{
		Cause: &ErrRoleNotFound{Name: name},
	}
}

// WithRole adds a role to [RBAC]. This option is useful if you have implemented the [Role] interface yourself. Otherwise, use [WithRoles] and [WithSilentRoles] instead.
func WithRole(name string, r Role) Option {
	return func(rb RBAC) (err error) {
		if name == "" {
			return errors.New("cannot use an empty role name")
		}
		if _, ok := rb[name]; ok {
			return fmt.Errorf("role %q has already been defined", name)
		}
		rb[name] = r
		return nil
	}
}

// Uninitialized [Logger] will be replaced with a default system logger relying on slow [Event] reflection. Do not use uninitialized loggers in production.
func WithRoles(definition map[string][]Policy, logger Logger) Option {
	return func(r RBAC) (err error) {
		if logger == nil {
			return errors.New("cannot use uninitialized logger")
		}
		for name, policies := range definition {
			if err = ValidatePolicySet(policies); err != nil {
				return fmt.Errorf("invalid policy set for role %q: %w", name, err)
			}
			if err = WithRole(name, &observedRole{
				name:     name,
				policies: policies,
				logger:   logger,
			})(r); err != nil {
				return fmt.Errorf("failed to add role %q: %w", name, err)
			}
		}
		return nil
	}
}

// WithSilentRoles creates a role set equivalent to [WithRoles] except for trading observability for faster execution. Use with caution for roles exposed to heavy traffic.
func WithSilentRoles(definition map[string][]Policy) Option {
	return func(r RBAC) (err error) {
		for name, policies := range definition {
			if err = ValidatePolicySet(policies); err != nil {
				return fmt.Errorf("invalid policy set for role %q: %w", name, err)
			}
			if err = WithRole(name, &silentRole{
				name:     name,
				policies: policies,
			})(r); err != nil {
				return fmt.Errorf("failed to add role %q: %w", name, err)
			}
		}
		return nil
	}
}
