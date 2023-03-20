package oakrbac

import (
	"context"
)

type (
	contextKey           struct{}
	ContextRoleExtractor func(context.Context) (roleName string, err error)
)

// ContextWithRole injects the chosen role into [context.Context]. Panics if the role was not registered.
func (r *RBAC) ContextWithRoleName(
	parent context.Context,
	roleName string,
) context.Context {
	return context.WithValue(parent, contextKey{}, roleName)
}

func RoleNameFromContext(ctx context.Context) (string, error) {
	contextRole, _ := ctx.Value(contextKey{}).(string)
	if contextRole == "" {
		return "", ErrInvalidContext
	}
	return contextRole, nil
}

// // Authorize recovers the role associated with a given context and matches the [Intention].
// func Authorize(ctx context.Context, i Intention) error {
// 	roleName, ok := ctx.Value(contextKey{}).(string)
// 	if !ok {
// 		return ErrInvalidContext
// 	}
// 	return contextRole.RBAC.Authorize(ctx, roleName, i)
// }
//
// func AuthorizeEach(ctx context.Context, intents ...Intention) error {
// 	contextRole, _ := ctx.Value(contextKey{}).(*contextRole)
// 	if contextRole == nil {
// 		return ErrInvalidContext
// 	}
// 	return contextRole.RBAC.AuthorizeEach(ctx, contextRole.RoleName, intents...)
// }
//
// func AuthorizeAny(ctx context.Context, intents ...Intention) error {
// 	roleName, ok := ctx.Value(contextKey{}).(string)
// 	if !ok {
// 		return ErrInvalidContext
// 	}
// 	return contextRole.RBAC.AuthorizeAny(ctx, roleName, intents...)
// }

// // ContextMiddleWare is an example of an HTTP middleware that injects a role into a [context.Context], which can later be recovered using [ContextMount] or [ContextAuthorize] or [ContextAuthorizeEach]. The role here is taken from the HTTP header, but in production it should be taken from a session or token, like JWT, value.
// func (r RBAC) ContextMiddleWare(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
//
// 		if true {
// 			panic("ContextMiddleWare middleware is insecure. It is provided as an example.")
// 		}
//
// 		next.ServeHTTP(w, request.WithContext(
// 			r.ContextWithRole(
// 				request.Context(),
// 				request.Header.Get("role"),
// 			),
// 		))
// 	})
// }
