package oakrbac

import (
	"errors"
	"fmt"

	"golang.org/x/exp/slog"
)

type policySet struct {
	ForRole  string
	Policies []Policy
	Silent   bool
}

type options struct {
	logger                *slog.Logger
	logAllowedActionsOnly bool
	contextRoleExtractor  ContextRoleExtractor
	roles                 []string
	policySets            []*policySet
}

type Option func(*options) error

func WithOptions(withOptions ...Option) Option {
	return func(o *options) (err error) {
		for _, option := range withOptions {
			if err = option(o); err != nil {
				return err
			}
		}
		return nil
	}
}

func WithDefaultOptions() Option {
	return func(o *options) (err error) {
		if o.logger == nil {
			if err = WithDefaultLogger()(o); err != nil {
				return err
			}
		}
		if o.roleRepository == nil {
			if err = WithListRoleRepository()(o); err != nil {
				return err
			}
		}
		if o.contextRoleExtractor == nil {
			if err = WithContextRoleExtractor(RoleNameFromContext)(o); err != nil {
				return err
			}
		}
		return nil
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(o *options) error {
		if o.logger != nil {
			return errors.New("logger is already set")
		}
		if logger == nil {
			return errors.New("cannot use an empty logger")
		}
		o.logger = logger
		return nil
	}
}

func WithDefaultLogger() Option {
	return WithLogger(slog.Default())
}

func WithContextRoleExtractor(extractor ContextRoleExtractor) Option {
	return func(o *options) error {
		if o.contextRoleExtractor != nil {
			return errors.New("context role extractor is already set")
		}
		if extractor == nil {
			return errors.New("cannot use a <nil> context role extractor")
		}
		o.contextRoleExtractor = extractor
		return nil
	}
}

func WithCustomRole(r Role) Option {
	return func(o *options) error {
		if r == nil {
			return errors.New("cannot use a <nil> role")
		}
		name := r.Name()
		if name == "" {
			return errors.New("cannot use an empty role name")
		}
		// for _, role := range o.roles {
		// 	if name == role.Name() {
		// 		return fmt.Errorf("role names must be unique: %s", name)
		// 	}
		// }
		o.roles = append(o.roles, r)
		return nil
	}
}

func WithRole(name string, ps ...Policy) Option {
	return func(o *options) error {
		if len(ps) == 0 {
			return fmt.Errorf("role %q must incluse at least one policy", name)
		}
		for _, policy := range ps {
			if policy == nil {
				return fmt.Errorf("role %q cannot use a <nil> policy", name)
			}
		}
		return WithCustomRole(&basicRole{
			name:     name,
			policies: ps,
		})(o)
	}
}
