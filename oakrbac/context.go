package oakrbac

import (
	"context"
	"errors"
	"net/http"
)

type contextKey struct{}

// ContextWithRole injects the chosen role into [context.Context]. Panics if the role was not registered.
func (r RBAC) ContextWithRole(parent context.Context, role string) context.Context {
	return context.WithValue(parent, contextKey{}, r(role))
}

// Authorize recovers the role associated with a given context and matches the [Intent]. It returns the [Policy] that granted authorization. The second return value is [AuthorizationError] in place of a generic error.
func Authorize(ctx context.Context, i *Intent) (policyGrantingAccess Policy, err *AuthorizationError) {
	role, _ := ctx.Value(contextKey{}).(Role)
	if role == nil {
		return nil, &AuthorizationError{
			Cause: errors.New("role context not found"),
		}
	}
	return role(ctx, i)
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
