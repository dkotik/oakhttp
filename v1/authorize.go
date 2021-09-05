package oakacs

import (
	"context"
)

// Authorize recovers the role from context and applies it to action.
func (acs *AccessControlSystem) Authorize(ctx context.Context, a Action) (err error) {
	defer acs.Broadcast(ctx, EventTypeAuthorization, err) // wrap error
	session, err := acs.Continue(ctx)
	if err != nil {
		return err
	}
	role, err := acs.roles.Retrieve(ctx, session.Role)
	if err != nil {
		return err
	}
	return role(a)
}
