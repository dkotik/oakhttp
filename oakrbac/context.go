package oakrbac

import (
	"context"
	"errors"
	"net/http"
)

type ContextKey int

const (
	ContextKeyRole ContextKey = iota
	ContextKeyRoleName
	contextKeySelf
)

type roleContext struct {
	context.Context
	role     Role
	roleName string
}

func (c *roleContext) Value(key any) (value any) {
	k, ok := key.(ContextKey)
	if ok {
		switch k {
		case ContextKeyRole:
			return c.role
		case ContextKeyRoleName:
			return c.roleName
		case contextKeySelf:
			return c
		}
	}
	return c.Context.Value(key)
}

// ContextWithRole injects the chosen role into [context.Context]. Panics if the role was not registered.
func (r RBAC) ContextWithRole(
	parent context.Context,
	role string,
) context.Context {
	return &roleContext{
		Context:  parent,
		role:     r[role],
		roleName: role,
	}
}

// ContextWithNegotiatedRole injects the chosen role into [context.Context]. If the role was not registered, the defaultRole is used. Panics if the defaultRole was not registered.
func (r RBAC) ContextWithNegotiatedRole(
	parent context.Context,
	role string,
	defaultRole string,
) context.Context {
	found, ok := r[role]
	if ok {
		return &roleContext{
			Context:  parent,
			role:     found,
			roleName: role,
		}
	}
	return &roleContext{
		Context:  parent,
		role:     r[defaultRole],
		roleName: defaultRole,
	}
}

// Authorize recovers the role associated with a given context and checks the intent against the role.
func Authorize(ctx context.Context, i *Intent) (Policy, *AuthorizationError) {
	c, ok := ctx.Value(contextKeySelf).(*roleContext)
	if !ok {
		return nil, &AuthorizationError{Cause: ErrContextRoleNotFound}
	}
	p, err := c.role(ctx, i)
	if errors.Is(err, Allow) {
		return p, nil
	}
	return p, &AuthorizationError{Policy: p, Cause: err}
}

// ContextMiddleWare is an example of an HTTP middleware that injects a role into a [context.Context], which can later be recovered using [ContextMount] or [ContextAuthorize] or [ContextAuthorizeEach]. The role here is taken from the HTTP header, but in production it should be taken from a session or token, like JWT, value.
func (r RBAC) ContextMiddleWare(fallback string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {

		if true {
			panic("ContextMiddleWare middleware is insecure. It is provided as an example.")
		}

		next.ServeHTTP(w, request.WithContext(
			r.ContextWithNegotiatedRole(
				request.Context(),
				request.Header.Get("role"),
				fallback,
			),
		))
	})
}
