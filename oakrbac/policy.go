package oakrbac

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

type (
	// Policy returns Allow sentinel error if the session is permitted to interact with the context.
	// Policy returns Deny sentinel error to interrupt the matching loop.
	// Policy returns `nil` if it did not match, but another policy might match.
	Policy func(context.Context, *Intent) error
)

var (
	//revive:disable:error-naming
	// Allow is a sentinel error that explicitly indicates that a [Policy] matched [Intent] and grants access.
	Allow = errors.New("sentinel error: policy granted access")
	// Deny is a sentinel error that explicitly indicates that a [Policy] matched [Intent] and denies access.
	Deny = errors.New("sentinel error: policy denied access")
	//revive:enable:error-naming

	// ErrNoPolicyMatched blocks authorization because every role policy returned a `nil` value.
	// ErrNoPolicyMatched = errors.New("no authorization policy matched")
	// ErrContextRoleNotFound indicates that a context does not include a role that can be retrieved using the package context key. If you see this error, you probably forgot to inject the role using either [ContextWithRole] or [rbac.ContextInjectorWithFallback] early in the execution path. This is typically done using a middleware function like [rbac.ContextMiddleWare].
	// ErrNoPredicates     = errors.New("there are no predicates attached to the Intent")
)

// Name returns the name of the [Policy] function by using reflection.
func (p Policy) Name() string {
	if p == nil {
		return "<no policy matched>"
	}
	return runtime.FuncForPC(reflect.ValueOf(p).Pointer()).Name()
}

// File returns the path to the file containing [Policy] function definition by using reflection.
func (p Policy) File() string {
	if p == nil {
		return "<nil>"
	}
	definition := runtime.FuncForPC(reflect.ValueOf(p).Pointer())
	file, _ := definition.FileLine(definition.Entry())
	return file
}

// Line returns the line number of the [Policy] function in its file by using reflection.
func (p Policy) Line() int {
	if p == nil {
		return 0
	}
	definition := runtime.FuncForPC(reflect.ValueOf(p).Pointer())
	_, line := definition.FileLine(definition.Entry())
	return line
}

func (p Policy) String() string {
	if p == nil {
		return "<no policy matched>"
	}
	definition := runtime.FuncForPC(reflect.ValueOf(p).Pointer())
	name := definition.Name()
	file, line := definition.FileLine(definition.Entry())

	output := strings.Builder{}
	output.WriteString("policy `")
	output.WriteString(name)
	output.WriteString("` defined in ")
	output.WriteString(file)
	output.WriteString(", line ")
	output.WriteString(fmt.Sprintf("%d", line))
	return output.String()
}

// PolicyEither combines a [Policy] list into one that return on first [Allow] or error.
func PolicyEither(ps ...Policy) Policy {
	return func(ctx context.Context, i *Intent) (err error) {
		for _, p := range ps {
			if err = p(ctx, i); err != nil {
				return err
			}
		}
		return nil
	}
}

// PolicyEach combines a [Policy] list into one that succeeds only if each policy returns [Allow]. An empty list never matches.
func PolicyEach(ps ...Policy) Policy {
	return func(ctx context.Context, i *Intent) (err error) {
		if len(ps) == 0 {
			return nil
		}
		for _, p := range ps {
			if err = p(ctx, i); !errors.Is(err, Allow) {
				return err
			}
		}
		return Allow
	}
}

// AllowEverything authorizes any action on any resource. Use cautiously.
func AllowEverything(_ context.Context, _ *Intent) error {
	return Allow
}

// DenyEverything denies authorization for any action on any resource.
func DenyEverything(_ context.Context, _ *Intent) error {
	return Deny
}

// AllowActionsForResourcesMatching authorizes any action from the provided list to any resource matching provided masks. This a helper method for debugging. Prefer to use generated policies.
func AllowActionsForResourcesMatching(actions []Action, resourceMasks [][]string) Policy {
	return func(ctx context.Context, i *Intent) error {
		for _, resourceMask := range resourceMasks {
			if i.MatchResourcePath(resourceMask...) {
				if i.MatchAction(actions...) {
					return Allow
				}
			}
		}
		return nil
	}
}
