package oakrbac

import (
	"context"
	"errors"
	"fmt"
)

// A Role is an [Intent] authorization provider. It returns the [Policy] that granted authorization. The second return value is [AuthorizationError] in place of a generic error.
type Role func(context.Context, *Intent) (policyGrantingAccess Policy, err *AuthorizationError)

func WithNewRole(name string, ps ...Policy) Option {
	return func(set map[string]Role) (err error) {
		if len(ps) == 0 {
			return errors.New("policy set must include at least one policy")
		}
		for i, p := range ps {
			if p == nil {
				return fmt.Errorf("policy set for role %q contains an uninitialized policy at index %d", name, i)
			}
		}
		return WithRole(name, func(ctx context.Context, i *Intent) (Policy, *AuthorizationError) {
			var err error
			var policyGrantingAccess Policy
			for _, policyGrantingAccess = range ps {
				// // check if context is cancelled, cost: mutex operation
				// if err = ctx.Err(); err != nil {
				// 	return nil, err
				// }

				err = policyGrantingAccess(ctx, i)
				if err == nil {
					continue // policy did not match
				} else if errors.Is(err, Allow) {
					return policyGrantingAccess, nil
				} else if errors.Is(err, Deny) {
					return nil, &AuthorizationError{
						Policy: policyGrantingAccess,
						Cause:  nil,
					}
				}
				return nil, &AuthorizationError{
					Policy: policyGrantingAccess,
					Cause:  err,
				}
			}
			return nil, &AuthorizationError{}
		})(set)
	}
}

// WithRole adds a role to [RBAC]. This option is useful if you have implemented the [Role] interface yourself. Otherwise, use [WithNewRole] instead.
func WithRole(name string, r Role) Option {
	return func(set map[string]Role) (err error) {
		if r == nil {
			return errors.New("cannot use an uninitialized role")
		}
		if name == "" {
			return errors.New("cannot use an empty role name")
		}
		if _, ok := set[name]; ok {
			return fmt.Errorf("role %q has already been defined", name)
		}
		set[name] = r
		return nil
	}
}
