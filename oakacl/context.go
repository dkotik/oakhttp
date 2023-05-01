package oakacl

import "context"

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
		return nil, oakpolicy.ErrInvalidContext
	}
	return permissions, nil
}
