package oakacs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/xid"
)

// ErrAuthorization happens when access to a resource is denied for any reason.
type ErrAuthorization struct {
	Service  string
	Domain   string
	Resource string
	Action   string
	Cause    error
}

// TODO: add unwrap method?
func (e *ErrAuthorization) Error() string { return "access denied" }

type PermissionsRepository interface {
	PullPermissions(ctx context.Context, roleUUID xid.ID) (deny []Permission, allow []Permission, err error)
	PushPermissions(ctx context.Context, roleUUID xid.ID, deny []Permission, allow []Permission) error
}

// Authorize recovers the role from context, iterates through its permissions. Returns <nil> when one of the permissions matches and satisfies service, domain, resource, and action constraints.
func (acs *AccessControlSystem) Authorize(
	ctx context.Context,
	service, domain, resource, action string,
) (err error) {
	session, err := acs.SessionFrom(ctx)
	if err != nil {
		return err
	}
	event := Event{
		ctx:     ctx,
		Type:    EventTypeAuthorizationAllowed,
		Session: session.UUID,
		Role:    session.Role,
	}
	defer func() {
		if err != nil {
			err = &ErrAuthorization{
				Service:  service,
				Domain:   domain,
				Resource: resource,
				Action:   action,
				Cause:    err,
			}
			event.Error = err
		}
		acs.Broadcast(event)
	}()

	if session.Deadline.After(time.Now()) {
		event.Type = EventTypeSessionExpired
		return errors.New("session expired")
	}

	deny, allow, err := acs.backend.PullPermissions(ctx, session.Role)
	if err != nil {
		return err
	}

	var p Permission
	for _, p = range deny {
		if p.Match(service, domain, resource, action) {
			event.Type = EventTypeAuthorizationDeniedByPermission
			return fmt.Errorf("permission explicitly denied by %s", p)
		}
	}
	for _, p = range allow {
		if p.Match(service, domain, resource, action) {
			event.Type = EventTypeAuthorizationAllowed
			// TODO: need to add context to passing events as well
			return nil
		}
	}
	event.Type = EventTypeAuthorizationDeniedByDefault
	return errors.New("none of the permissions matched")
}

// NewAuthority prepares a function that can authorize actions taken against a resource.
func (acs *AccessControlSystem) NewAuthority(
	service, domain, resource string,
) func(ctx context.Context, action string) error {
	return func(ctx context.Context, action string) error {
		return acs.Authorize(ctx, service, domain, resource, action)
	}
}

// NewAuthorityWithDomainTransience prepares a function that can authorize actions across domains.
func (acs *AccessControlSystem) NewAuthorityWithDomainTransience(
	service, resource string,
) func(ctx context.Context, domain, action string) error {
	return func(ctx context.Context, domain, action string) error {
		return acs.Authorize(ctx, service, domain, resource, action)
	}
}

// // Authority validates that a given session can perform a particular action.
// type Authority interface {
// 	// Authorize(ctx context.Context, sessionUUID, action string) bool
// 	Authorize(ctx context.Context, action string) error
// }
//
// // AuthorityFunc provides a single-method Authority.
// type AuthorityFunc func(ctx context.Context, action string) error
//
// // Authorize satisfies the Authority interface.
// func (a AuthorityFunc) Authorize(ctx context.Context, sessionUUID, action string) error {
// 	return a(ctx, action)
// }
