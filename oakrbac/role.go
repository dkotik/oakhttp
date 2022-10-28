package oakrbac

import (
	"context"
	"errors"
	"fmt"
)

var ErrAuthorizationDenied = errors.New("authorization denied")

// A Role is an [Intent] authorization provider. It returns true when authorization is granted. The returned policy points to [Policy] that granted or denied authorization or was interrupted by an error. The pointer can be used for observability using either [Policy.Name] and [Policy.NameFileLine] methods, which can handle `nil` values.
//
// TODO: retype the comment above
type Role func(context.Context, *Intent) (Policy, error)

func newRole(ps ...Policy) Role {
	return func(ctx context.Context, i *Intent) (
		policy Policy,
		err error,
	) {
		for _, policy = range ps {
			if err = policy(ctx, i); err != nil {
				if errors.Is(err, Allow) {
					return policy, nil // policy matched
				}
				if errors.Is(err, Deny) {
					return policy, ErrAuthorizationDenied // policy blocked
				}
				return policy, err // unexpected error
			}
		}
		return nil, ErrAuthorizationDenied
	}
}

func WithNewRole(name string, ps ...Policy) Option {
	return func(r RBAC) (err error) {
		if len(ps) == 0 {
			return errors.New("policy set must include at least one policy")
		}
		for i, p := range ps {
			if p == nil {
				return fmt.Errorf("policy set for role %q contains an uninitialized policy at index %d", name, i)
			}
		}
		return WithRole(name, newRole(ps...))(r)
	}
}

// WithRole adds a role to [RBAC]. This option is useful if you have implemented the [Role] interface yourself. Otherwise, use [WithNewRole] instead.
func WithRole(name string, r Role) Option {
	return func(rb RBAC) (err error) {
		if r == nil {
			return errors.New("cannot use an uninitialized role")
		}
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
