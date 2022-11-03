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
)

type (
	// RBAC is a simple Role Based Access Control system.
	RBAC map[string]Role

	// Option customizes the [RBAC] constructor [New].
	Option func(RBAC) error
)

// Must panics if an error is associated with [RBAC] constructor. Use together with [New].
func Must(r RBAC, err error) RBAC {
	if err != nil {
		panic(err)
	}
	return r
}

// New builds an [RBAC] using provided [Option] set.
func New(options ...Option) (rbac RBAC, err error) {
	rbac = make(RBAC)
	for _, option := range append(options, func(r RBAC) error {
		if len(r) == 0 {
			return errors.New("provide at least one role")
		}
		return nil
	}) {
		if err = option(rbac); err != nil {
			return nil, fmt.Errorf("cannot create OakRBAC: %w", err)
		}
	}
	return rbac, nil
}

func (r RBAC) Authorize(ctx context.Context, role string, i *Intent) (Policy, *AuthorizationError) {
	found, ok := r[role]
	if !ok {
		return nil, &AuthorizationError{Cause: ErrRoleNotFound}
	}
	p, err := found(ctx, i)
	if errors.Is(err, Allow) {
		return p, nil
	}
	return p, &AuthorizationError{Policy: p, Cause: err}
}
