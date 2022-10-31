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
	// ErrNoPredicates     = errors.New("there are no predicates attached to the Intent")

	// ErrContextRoleNotFound indicates that a context does not include a role that can be retrieved using the package context key. If you see this error, you probably forgot to inject the role using either [ContextWithRole] or [rbac.ContextInjectorWithFallback] early in the execution path. This is typically done using a middleware function like [rbac.ContextMiddleWare].
	// ErrContextRoleNotFound indicates the absence of [Role] association with [context.Context]. Did you forget to inject the role using [rbac.ContextWithRole] or [rbac.ContextWithNegotiatedRole]?
	ErrContextRoleNotFound = errors.New("OakACS role context not found in context chain")
	ErrRoleNotFound        = errors.New("role is not registered")
	// ErrMissingPredicate is raised when a [Policy] cannot locate the desired predicate in an [Intent].
	ErrPredicateNotFound = errors.New("predicate does not exist")
)

// AuthorizationError expresses the output of a [Role] as an opaque [Deny] error to prevent attackers from discovering the internals of the access control system by analyzing its error messages. Use [AuthorizationError.Message] for logging and debugging to discover the conditions for authorization failure.
type AuthorizationError struct {
	Policy Policy
	Cause  error
}

// Unwrap satisfies [errors.Is] and [errors.As] interface requirements.
func (e *AuthorizationError) Unwrap() error {
	return e.Cause
}

// Error always returns the value of [Deny] error regardless of the state to prevent attackers from discovering the internals of the access control system by analyzing its error messages.
func (e *AuthorizationError) Error() string {
	return Deny.Error()
}

// Message provides the full description of authorization failure. Use the output as the message value of a structured logger entry.
func (e *AuthorizationError) Message() string {
	if e == nil || e.Cause == nil {
		return "access granted by matched policy"
	}
	if errors.Is(e.Cause, Deny) {
		if e.Policy == nil {
			return "access denied: none of the policies matched the intent"
		}
		return "access denied by matched policy"
	}
	return "access denied: policy evaluation failed: " + e.Cause.Error()
}

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
