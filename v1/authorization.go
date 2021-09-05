package oakacs

import (
	"context"
	"errors"
)

// Authorize recovers the role from context, iterates through its permissions. Returns <nil> when one of the permissions matches and satisfies service, domain, resource, and action constraints.
func (acs *AccessControlSystem) Authorize(
	ctx context.Context,
	service, domain, resource, action string,
) (err error) {
	// session, err := acs.SessionContinue(ctx)
	// if err != nil {
	// 	return err
	// }

	// if session.Deadline.After(time.Now()) {
	// 	event.Type = EventTypeSessionExpired
	// 	return errors.New("session expired")
	// }

	// deny, allow, err := acs.persistent.PullPermissions(ctx, session.Role)
	// if err != nil {
	// 	return err
	// }
	//
	// var p Permission
	// for _, p = range deny {
	// 	if p.Match(service, domain, resource, action) {
	// 		event.Type = EventTypeAuthorizationDeniedByPermission
	// 		return fmt.Errorf("permission explicitly denied by %s", p)
	// 	}
	// }
	// for _, p = range allow {
	// 	if p.Match(service, domain, resource, action) {
	// 		event.Type = EventTypeAuthorizationAllowed
	// 		// TODO: need to add context to passing events as well
	// 		return nil
	// 	}
	// }
	// event.Type = EventTypeAuthorizationDeniedByDefault
	return errors.New("none of the permissions matched")
}
