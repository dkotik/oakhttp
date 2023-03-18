package oakrbac

import (
	"context"
)

type Listener interface {
	AuthorizationGranted(context.Context, []Intent, []Policy, Role)
	AuthorizationDenied(context.Context, []Intent, []Policy, Role)
	AuthorizationFailed(context.Context, Intent, Policy, Role, error)
}

func (r *RBAC) AuthorizationGranted(ctx context.Context, i []Intent, p []Policy, role Role) {
	for _, listener := range r.listeners {
		listener.AuthorizationGranted(ctx, i, p, role)
	}
}

func (r *RBAC) AuthorizationDenied(ctx context.Context, i []Intent, p []Policy, role Role) {
	for _, listener := range r.listeners {
		listener.AuthorizationDenied(ctx, i, p, role)
	}
}

func (r *RBAC) AuthorizationFailed(ctx context.Context, i Intent, p Policy, role Role, err error) {
	for _, listener := range r.listeners {
		listener.AuthorizationFailed(ctx, i, p, role, err)
	}
}
