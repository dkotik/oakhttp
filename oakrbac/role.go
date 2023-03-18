package oakrbac

import (
	"context"
)

// A Role is an [Intention] authorization provider. It returns the [Policy] that granted authorization. The second return value is [AuthorizationError] in place of a generic error.
type Role interface {
	Name() string
	Authorize(context.Context, Intention) (Policy, error)
}

type basicRole struct {
	name     string
	policies []Policy
}

func (r *basicRole) Name() string {
	return r.name
}

func (r *basicRole) Authorize(ctx context.Context, i Intention) (policy Policy, err error) {
	for _, policy = range r.policies {
		if err = policy(ctx, i); err != nil {
			if err == nil {
				// policy did not match
				// expecting an [Allow], [Deny], or error
				continue
			}
			return policy, err
		}
	}
	return nil, Deny
}

type omnipotentRole struct {
	name string
}

func (o *omnipotentRole) Name() string {
	return o.name
}

func (o *omnipotentRole) Authorize(ctx context.Context, i Intention) (Policy, error) {
	return AllowEverything, Allow
}

type impotentRole struct {
	name string
}

func (o *impotentRole) Name() string {
	return o.name
}

func (o *impotentRole) Authorize(ctx context.Context, i Intention) (Policy, error) {
	return DenyEverything, Deny
}
