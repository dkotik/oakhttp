package oakrbac

import (
	"context"
	"errors"
)

// Authorize matches the named [Role] against an intended [Action] aimed at a [Resource].
func (r *RBAC) Authorize(ctx context.Context, a oakpolicy.Action, r oakpolicy.Resource) error {
	roleName, err := r.contextRoleExtractor(ctx)
	if err != nil {
		return &AuthorizationError{cause: err}
	}
	role, err := r.GetRole(roleName)
	if err != nil {
		return &AuthorizationError{cause: err}
	}
	policy, err := role.Authorize(ctx, i)
	if errors.Is(err, Allow) {
		return nil
	} else if errors.Is(err, Deny) {
		return err
	}
	return &AuthorizationError{
		policy: policy,
		cause:  err,
	}
}

func (r *RBAC) AuthorizeEach(ctx context.Context, roleName string, intents ...Intention) (err error) {

}

func (r *RBAC) AuthorizeAny(ctx context.Context, roleName string, intents ...Intention) error {

}
