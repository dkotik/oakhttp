package oakrbac

import (
	"context"
)

// Listener tracks [RBAC] authorization events. Multiple [Intention]s or [Policy]s may be involved in making the decision. Passed [Intention], [Policy], or [Role] may be nil.
//
// When implementing a logger, you should not record every [Listener.AuthorizationDenied] or [Listener.AuthorizationFailed] event as they are passed back through the execution stack and should be logged upstream. [AuthorizationGranted] events should always be logged somewhere, because a <nil> error value replaces the sentinel [Allow].
type Listener interface {
	AuthorizationGranted(context.Context, []Intention, []Policy, Role)
	AuthorizationDenied(context.Context, []Intention, []Policy, Role)
	AuthorizationFailed(context.Context, Intention, Policy, Role, error)
}

// AuthorizationGranted broadcasts an [Allow] event to all [RBAC] [Listener]s.
func (r *RBAC) AuthorizationGranted(ctx context.Context, i []Intention, p []Policy, role Role) {
	for _, listener := range r.listeners {
		listener.AuthorizationGranted(ctx, i, p, role)
	}
}

// AuthorizationDenied broadcasts a [Deny] event to all [RBAC] [Listener]s.
func (r *RBAC) AuthorizationDenied(ctx context.Context, i []Intention, p []Policy, role Role) {
	for _, listener := range r.listeners {
		listener.AuthorizationDenied(ctx, i, p, role)
	}
}

// AuthorizationFailed broadcasts an [error] event to all [RBAC] [Listener]s.
func (r *RBAC) AuthorizationFailed(ctx context.Context, i Intention, p Policy, role Role, err error) {
	for _, listener := range r.listeners {
		listener.AuthorizationFailed(ctx, i, p, role, err)
	}
}
