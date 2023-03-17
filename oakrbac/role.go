package oakrbac

import (
	"context"
	"errors"
)

// A Role is an [Intent] authorization provider. It returns the [Policy] that granted authorization. The second return value is [AuthorizationError] in place of a generic error.
type Role interface {
	Name() string
	Authorize(context.Context, Intent) (EventType, Policy, error)
}

type basicRole struct {
	name     string
	policies []Policy
}

func (r *basicRole) Name() string {
	return r.name
}

func (r *basicRole) Authorize(ctx context.Context, i Intent) (EventType, Policy, error) {
	var (
		policy Policy
		err    error
	)
	for _, policy = range r.policies {
		err = policy(ctx, i)
		if err == nil {
			continue // policy did not match
		} else if errors.Is(err, Allow) {
			return EventTypeAuthorizationGranted, policy, nil
		} else if errors.Is(err, Deny) {
			return EventTypeAuthorizationDenied, policy, Deny
		}
		return EventTypeError, policy, &AuthorizationError{
			policy: policy,
			cause:  err,
		}
	}
	return EventTypeAuthorizationDenied, nil, Deny
}
