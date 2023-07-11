package oakacl

import (
	"context"
	"errors"
)

// ErrInvalidContext indicates that a context does not include a set of permissions that can be retrieved using the package context key. If you see this error, you probably forgot to inject the permissions using either [ContextWithPermissions] early in the execution path. This is typically done using a middleware function like [ContextMiddleWare].
var ErrInvalidContext = errors.New("context chain contains no ACL permissions")

type contextKey struct{}

type ContextPermissionsExtractor func(context.Context) ([]Permission, error)

func ContextWithPermissions(
	parent context.Context,
	permissions []Permission,
) context.Context {
	return context.WithValue(parent, contextKey{}, permissions)
}

func PermissionsFromContext(ctx context.Context) ([]Permission, error) {
	permissions, _ := ctx.Value(contextKey{}).([]Permission)
	if len(permissions) == 0 {
		return nil, ErrInvalidContext
	}
	return permissions, nil
}

type authorizedContext []Permission

// func (a *authorizedContext) To(action string, r any) error {
// 	dpath := r.DomainPath()
// 	for _, p := range ps {
// 		if p.Action.Match(action) && dpath.Match(p.ResourceMask) {
// 			return nil
// 		}
// 	}
// 	return oakpolicy.Deny
// }
