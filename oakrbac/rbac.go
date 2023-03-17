/*

Package oakrbac is a simple flexible Role-Based Access Control (RBAC) implementation.

A role is constructed from a set of policies that are sequentially evaluated until one of the policies returns an [Allow] or a [Deny] sentinel value or an error.

# Usage

OakRBAC leans on [context.Context] as the main mechanism for passing access rights through the execution stack.

	// 1. Initialize the [RBAC]:
	var RBAC = oakrbac.Must(oakrbac.New(
		oakrbac.WithNewRole("administrator", oakrbac.AllowEverything)
	))

	// 2. Inject authorization context:
	ctx := RBAC.ContextWithRole("administrator", context.TODO())

	// 3. Authorize an action [Intent]:
	matchedPolicy, err := oakrbac.Authorize(ctx, &oakrbac.Intent{
		Action: oakrbac.ActionCreate,
		ResourcePath: oakrbac.NewResourcePath(
			"myService",
			"user",
			"userUUID",
		)
	})

	// 4. Act on authorization result:
	if err != nil {
		// access denied, log it using [AuthorizationError.Message] method:
		log.Println(err.Message())
		return err
	}

	// when err == nil, [AuthorizationError.Message] method returns "access granted"
	log.Println(err.Message())

# Policies

OakRBAC comes with only two default policies: [AllowEverything] and [DenyEverything]. You will write or generate policies to match your domain logic.

# Predicates

An [Intent] can be created with a set of [Predicate] functions that allow a [Policy] to run code snippets against the resource to examine it during evaluation.

Predicates enable writing incredibly powerful and performant access control policies.

*/
package oakrbac

import (
	"context"
	"errors"
	"fmt"
	"os"

	"golang.org/x/exp/slog"
)

// RBAC is a simple Role Based Access Control system.
type RBAC struct {
	roles     []Role
	listeners []Listener
}

func (r *RBAC) GetRole(ctx context.Context, name string) (Role, error) {
	for _, r := range r.roles {
		if name == r.Name() {
			return r, nil
		}
	}
	r.Dispatch(
		ctx,
		NewEvent(
			EventTypeError,
			nil,
			nil,
			nil,
			ErrRoleNotFound,
		),
	)
	return nil, ErrRoleNotFound
}

func (r *RBAC) Dispatch(ctx context.Context, e *Event) {
	for _, listener := range r.listeners {
		listener.Listen(ctx, e)
	}
}

// Authorize matches the named [Role] against an [Intent]. It returns the [Policy] that granted authorization. The second return value is [AuthorizationError] in place of a generic error.
func (r *RBAC) Authorize(ctx context.Context, role Role, i Intent) error {
	eventType, policy, err := role.Authorize(ctx, i)
	r.Dispatch(
		ctx,
		NewEvent(
			eventType,
			role,
			[]Intent{i},
			[]Policy{policy},
			err,
		),
	)
	return err
}

func (r *RBAC) AuthorizeEvery(ctx context.Context, role Role, intents ...Intent) error {
	policies := make([]Policy, len(intents))
	for i, intent := range intents {
		eventType, policy, err := role.Authorize(ctx, intent)
		if err != nil {
			r.Dispatch(
				ctx,
				NewEvent(
					eventType,
					role,
					intents,
					[]Policy{policy},
					err,
				),
			)
			return err
		}
		policies[i] = policy
	}
	r.Dispatch(
		ctx,
		NewEvent(
			EventTypeAuthorizationGranted,
			role,
			intents,
			policies,
			nil,
		),
	)
	return nil
}

func (r *RBAC) AuthorizeAny(ctx context.Context, role Role, intents ...Intent) error {
	for _, intent := range intents {
		eventType, policy, err := role.Authorize(ctx, intent)
		switch eventType {
		case EventTypeAuthorizationGranted:
			r.Dispatch(
				ctx,
				NewEvent(EventTypeAuthorizationGranted,
					role,
					intents,
					[]Policy{policy},
					nil,
				),
			)
			return nil
		case EventTypeError:
			r.Dispatch(
				ctx,
				NewEvent(
					eventType,
					role,
					intents,
					[]Policy{policy},
					err,
				),
			)
			return err
		}
	}
	r.Dispatch(
		ctx,
		NewEvent(
			EventTypeAuthorizationDenied,
			role,
			intents,
			nil,
			Deny,
		),
	)
	return nil
}

// New builds an [RBAC] using provided [Option] set.
func New(withOptions ...Option) (rbac *RBAC, err error) {
	o := &options{}
	for _, option := range append(
		withOptions,
		func(o *options) (err error) { // validate
			if len(o.roles) == 0 {
				return errors.New("at least one role is required")
			}
			if len(o.listeners) == 0 {
				if err = WithSlogLogger(
					slog.New(slog.NewTextHandler(os.Stderr)),
					slog.LevelInfo,
				)(o); err != nil {
					return fmt.Errorf("could not setup default logger: %w", err)
				}
			}
			return nil
		},
	) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("cannot create OakRBAC: %w", err)
		}
	}
	return &RBAC{
		listeners: o.listeners,
		roles:     o.roles,
	}, nil
}

// Must panics if an error is associated with [RBAC] constructor. Use together with [New].
func Must(r *RBAC, err error) *RBAC {
	if err != nil {
		panic(err)
	}
	return r
}
