package oakpolicy

import (
	"errors"
	"net/http"

	"golang.org/x/exp/slog"
)

var (
	//revive:disable:error-naming
	// Allow is a sentinel error that explicitly indicates that a [Policy] matched [Intention] and grants access.
	Allow = errors.New("authorization granted")
	// Deny is a sentinel error that explicitly indicates that a [Policy] matched [Intention] and denies access.
	Deny = &denyError{}
	//revive:enable:error-naming

	ErrNilPolicy       = errors.New("a policy set contains an uninitialized <nil> Policy")
	ErrEmptyPolicyList = errors.New("cannot iterate over an empty policy list")
)

type denyError struct{}

func (e *denyError) Error() string { return "authorization denied" }

func (e *denyError) HTTPStatusCode() int {
	return http.StatusForbidden
}

// AuthorizationError expresses the output of a [Role] as an opaque [Deny] error to prevent attackers from discovering the internals of the access control system by analyzing its error messages. Use [AuthorizationError.Message] for logging and debugging to discover the conditions for authorization failure.
type AuthorizationError struct {
	policy Policy
	cause  error
}

func (e *AuthorizationError) Policy() Policy {
	return e.policy
}

// HTTPStatusCode returns an HTTP status code to satisfy oakhttp.HTTPError interface.
func (e *AuthorizationError) HTTPStatusCode() int {
	return http.StatusForbidden
}

// Unwrap satisfies [errors.Is] and [errors.As] interface requirements.
func (e *AuthorizationError) Unwrap() error {
	return e.cause
}

// Error always returns the value of [Deny] error regardless of the state to prevent attackers from discovering the internals of the access control system by analyzing its error messages.
func (e *AuthorizationError) Error() string {
	return Deny.Error()
}

func (e *AuthorizationError) LogValue() slog.Value {
	return slog.StringValue("authorization denied: " + e.cause.Error())
}
