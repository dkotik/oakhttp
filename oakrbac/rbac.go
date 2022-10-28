/*

Package oakrbac description...

*/
package oakrbac

import (
	"context"
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
	for _, option := range append(options, func(r RBAC) error {
		if len(r) == 0 {
			return errors.New("provide at least one role")
		}
		return nil
	}) {
		if err = option(rbac); err != nil {
			return nil, fmt.Errorf("cannot create OakRBAC: %w", err)
		}
	}
	return rbac, nil
}

func (r RBAC) Authorize(ctx context.Context, role string, i *Intent) (Policy, error) {
	found, ok := r[role]
	if !ok {
		return nil, ErrAuthorizationDenied
	}
	return found(ctx, i)
}

func (r RBAC) AuthorizeEach(ctx context.Context, role string, i ...*Intent) (p Policy, err error) {
	found, ok := r[role]
	if !ok {
		return nil, ErrAuthorizationDenied
	}
	for _, intent := range i {
		p, err = found(ctx, intent)
		if err != nil {
			return p, err
		}
	}
	return nil, nil
}
