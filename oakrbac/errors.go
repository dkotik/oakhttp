package oakrbac

import (
	"errors"
	"fmt"
)

var (
	//revive:disable:error-naming
	// Allow is a sentinel error that explicitly indicates that a [Policy] matched [Intent] and grants access.
	Allow = errors.New("authorization granted")
	// Deny is a sentinel error that explicitly indicates that a [Policy] matched [Intent] and denies access.
	Deny = errors.New("authorization denied")
	//revive:enable:error-naming

	// ErrNoPolicyMatched blocks authorization because every role policy returned a `nil` value.
	// ErrNoPolicyMatched = errors.New("no authorization policy matched")
	// ErrContextRoleNotFound indicates that a context does not include a role that can be retrieved using the package context key. If you see this error, you probably forgot to inject the role using either [ContextWithRole] or [rbac.ContextInjectorWithFallback] early in the execution path. This is typically done using a middleware function like [rbac.ContextMiddleWare].
	// ErrNoPredicates     = errors.New("there are no predicates attached to the Intent")

	// ErrContextRoleNotFound indicates the absence of [Role] association with [context.Context]. Did you forget to inject the role using [rbac.ContextWithRole] or [rbac.ContextWithNegotiatedRole]?
	// ErrContextRoleNotFound = errors.New("context does not include an OakACS role value")
)

var (
	// ErrMissingPredicate is raised when a [Policy] cannot locate the desired predicate in an [Intent].
	ErrPredicateNotFound = errors.New("predicate does not exist")
)

// PredicateError is raised when a [Policy] cannot execute an [Intent] [Predicate].
type PredicateError struct {
	Name  string
	Cause error
}

func (e *PredicateError) Error() string {
	return fmt.Sprintf("predicate %q failure: %s", e.Name, e.Cause)
}

func (e *PredicateError) Unwrap() error {
	return e.Cause
}

type AccessDeniedError struct {
	Role   string
	Policy Policy
	Cause  error
}

func (e *AccessDeniedError) Declassify() string {
	if e.Role == "" {
		return "role not found"
	}
	if e.Policy == nil {
		return "none of the role policies matched the intent"
	}
	if e.Cause == nil {
		return "denied by policy " + e.Policy.Name()
	}
	return fmt.Sprintf("policy %q failure: %s", e.Policy.Name(), e.Cause)
}

func (e *AccessDeniedError) Error() string {
	return "access denied"
}

func (e *AccessDeniedError) Unwrap() error {
	return e.Cause
}
