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

	// 3. Authorize an action [Intention]:
	matchedPolicy, err := oakrbac.Authorize(ctx, &oakrbac.Intention{
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

An [Intention] can be created with a set of [Predicate] functions that allow a [Policy] to run code snippets against the resource to examine it during evaluation.

Predicates enable writing incredibly powerful and performant access control policies.

*/
package oakrbac

import (
	"errors"
	"fmt"
)

// RBAC is a simple Role Based Access Control system.
type RBAC struct {
	// policyNames          map[*Policy]string
	// create a separate logger for each policy!
	policies             []Policy
	contextRoleExtractor ContextRoleExtractor
	listeners            []Listener
}

// New builds an [RBAC] using provided [Option] set.
func New(withOptions ...Option) (rbac *RBAC, err error) {
	o := &options{}
	for _, option := range append(
		withOptions,
		WithDefaultOptions(),
		func(o *options) (err error) { // validate
			for _, r := range o.roles {
				if err = o.roleRepository.AddRole(r); err != nil {
					return err
				}
			}
			if o.roleRepository.CountRoles() == 0 {
				return errors.New("at least one role is required")
			}
			return nil
		},
	) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("cannot create OakRBAC: %w", err)
		}
	}
	return &RBAC{
		roleRepository:       o.roleRepository,
		contextRoleExtractor: o.contextRoleExtractor,
		listeners:            o.listeners,
	}, nil
}

// Must panics if an error is associated with [RBAC] constructor. Use together with [New].
func Must(r *RBAC, err error) *RBAC {
	if err != nil {
		panic(err)
	}
	return r
}
