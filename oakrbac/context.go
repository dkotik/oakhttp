package oakrbac

import (
	"context"
	"net/http"
)

type contextKey struct{}

type contextRole struct {
	RBAC     *RBAC
	RoleName string
}

// ContextWithRole injects the chosen role into [context.Context]. Panics if the role was not registered.
func (r *RBAC) ContextWithRole(
	parent context.Context,
	roleName string,
) context.Context {
	return context.WithValue(parent, contextKey{}, &contextRole{
		RBAC:     r,
		RoleName: roleName,
	})
}

func RoleFromContext(ctx context.Context) (Role, error) {
	contextRole, _ := ctx.Value(contextKey{}).(*contextRole)
	if contextRole == nil {
		return nil, ErrInvalidContext
	}
	return contextRole.RBAC.GetRole(contextRole.RoleName)
}

// Authorize recovers the role associated with a given context and matches the [Intent].
func Authorize(ctx context.Context, i Intent) error {
	contextRole, _ := ctx.Value(contextKey{}).(*contextRole)
	if contextRole == nil {
		return ErrInvalidContext
	}
	return contextRole.RBAC.Authorize(ctx, contextRole.RoleName, i)
}

func AuthorizeEach(ctx context.Context, intents ...Intent) error {
	contextRole, _ := ctx.Value(contextKey{}).(*contextRole)
	if contextRole == nil {
		return ErrInvalidContext
	}
	return contextRole.RBAC.AuthorizeEach(ctx, contextRole.RoleName, intents...)
}

func AuthorizeAny(ctx context.Context, intents ...Intent) error {
	contextRole, _ := ctx.Value(contextKey{}).(*contextRole)
	if contextRole == nil {
		return ErrInvalidContext
	}
	return contextRole.RBAC.AuthorizeAny(ctx, contextRole.RoleName, intents...)
}

// ContextMiddleWare is an example of an HTTP middleware that injects a role into a [context.Context], which can later be recovered using [ContextMount] or [ContextAuthorize] or [ContextAuthorizeEach]. The role here is taken from the HTTP header, but in production it should be taken from a session or token, like JWT, value.
func (r RBAC) ContextMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {

		if true {
			panic("ContextMiddleWare middleware is insecure. It is provided as an example.")
		}

		next.ServeHTTP(w, request.WithContext(
			r.ContextWithRole(
				request.Context(),
				request.Header.Get("role"),
			),
		))
	})
}
