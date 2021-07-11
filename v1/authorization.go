package oakacs

import "context"

// Authority validates that a given session can perform a particular action.
type Authority interface {
	Authorize(ctx context.Context, sessionUUID, action string) bool
}

// AuthorityFunc provides a single-method Authority.
type AuthorityFunc func(ctx context.Context, sessionUUID, action string) bool

// Authorize satisfies the Authority interface.
func (a AuthorityFunc) Authorize(ctx context.Context, sessionUUID, action string) bool {
	return a(ctx, sessionUUID, action)
}

// NewAuthority prepares a function that can authorize actions taken against a resource.
func (acs *AccessControlSystem) NewAuthority(
	service, domain, resource string,
) AuthorityFunc {

	return func(ctx context.Context, sessionUUID, action string) bool {
		// locate role by uuid attached to session located by uuid
		// only need Deny and Allow lists from the role, nothing else
		// GetPermissions(sessionUUID string) ([]Permission, []Permission, error)
		// RBAC Permission provider?
		// Plain Permission provider?
		var p Permission
		for _, p = range r.Deny {
			if p.Match(service, domain, resource, action) {
				// TODO: zap.Logger
				return false
			}
		}
		for _, p = range r.Allow {
			if p.Match(service, domain, resource, action) {
				// TODO: zap.Logger
				return true
			}
		}
		return false
	}
}
