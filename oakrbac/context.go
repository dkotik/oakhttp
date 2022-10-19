package oakrbac

import (
	"context"
	"errors"
	"net/http"
)

type contextKey int

var (
	ErrContextRoleNotFound = errors.New("context does not include an OakACS role value")
)

const (
	contextKeyForRole contextKey = iota
)

// ContextWithRole injects the chosen role into [context.Context].
func ContextWithRole(parent context.Context, role Role) context.Context {
	return context.WithValue(parent, contextKeyForRole, role)
}

// ContextMount recovers an RBAC role from [context.Context].
func ContextMount(ctx context.Context) (Role, error) {
	role, ok := ctx.Value(contextKeyForRole).(Role)
	if ok {
		return role, nil
	}
	return nil, &AccessDeniedError{Cause: ErrContextRoleNotFound}
}

// ContextAuthorize recovers the role associated with a given context and checks the intent against the role. It is a helper method for [ContextMount].
func ContextAuthorize(ctx context.Context, i *Intent) error {
	role, err := ContextMount(ctx)
	if err != nil {
		return err
	}
	return role.Authorize(ctx, i)
}

// ContextAuthorizeEach recovers the role associated with a given context and checks each provided intent against the role. It is a helper method for [ContextMount].
func ContextAuthorizeEach(ctx context.Context, i ...*Intent) (err error) {
	role, err := ContextMount(ctx)
	if err != nil {
		return err
	}
	for _, i := range i {
		if err = role.Authorize(ctx, i); err != nil {
			return err
		}
	}
	return nil
}

// ContextWithRole is a helper method for [ContextWithRole] which first locates the RBAC role by name.
func (r *RBAC) ContextWithRole(parent context.Context, role string) (context.Context, error) {
	found, err := r.GetRole(role)
	if err != nil {
		return nil, err
	}
	return ContextWithRole(parent, found), nil
}

// ContextMiddleWare is an example of an HTTP middleware that injects a role into a [context.Context], which can later be recovered using [ContextMount] or [ContextAuthorize] or [ContextAuthorizeEach]. The role here is taken from the HTTP header, but in production it should be taken from a session or token, like JWT, value.
func (r *RBAC) ContextMiddleWare(fallback string, next http.Handler) http.Handler {
	injector := r.ContextInjectorWithFallback(fallback)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if true {
			panic("ContextMiddleWare middleware is insecure. It is provided as an example.")
		}

		next.ServeHTTP(w, r.WithContext(
			injector(r.Context(), r.Header.Get("role")),
		))
	})
}

// ContextWithRole is a builder for [ContextWithRole] which first locates the RBAC role by name. If the desired role cannot be located, the fallback role is used instead.
func (r *RBAC) ContextInjectorWithFallback(fallbackRole string) func(context.Context, string) context.Context {
	fallback, err := r.GetRole(fallbackRole)
	if err != nil {
		panic(err)
	}
	return func(parent context.Context, role string) context.Context {
		found, err := r.GetRole(role)
		if err != nil {
			return ContextWithRole(parent, fallback)
		}
		return ContextWithRole(parent, found)
	}
}
