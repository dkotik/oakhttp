package oakrbac

import (
	"errors"
	"fmt"
)

type options struct {
	roles     []Role
	listeners []Listener
}

type Option func(*options) error

func WithCustomRole(r Role) Option {
	return func(o *options) error {
		if r == nil {
			return errors.New("cannot use a <nil> role")
		}
		name := r.Name()
		if name == "" {
			return errors.New("cannot use an empty role name")
		}
		for _, role := range o.roles {
			if name == role.Name() {
				return fmt.Errorf("role names must be unique: %s", name)
			}
		}
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

func WithOmnipotentRole(name string) Option {
	return WithCustomRole(&omnipotentRole{name: name})
}

func WithImpotentRole(name string) Option {
	return WithCustomRole(&impotentRole{name: name})
}

func WithListener(l Listener) Option {
	return func(o *options) error {
		if l == nil {
			return errors.New("cannot use a <nil> event listener")
		}
		o.listeners = append(o.listeners, l)
		return nil
	}
}
