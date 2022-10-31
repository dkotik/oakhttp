package oakrbac

import (
	"context"
	"fmt"
	"strings"
)

type (
	// An Action specifies the verb of an [Intent]. OakRBAC comes with a set of most frequently occuring actions. Specify custom actions as constants.
	Action string

	// A Predicate characterizes a resource and gives a Policy the ability to make an access control assertion regarding a certain resource property.
	//
	// A Policy can embed the predicate code. However, predicates keep the data layer calls opaque to the Access Control System. This simplifies writing and testing the policies and allows sharing resource evaluation code between multiple policies, lazy evaluation, selective overriding, and consistent caching. If you cache the output of policies, you can expect inconsistent state between two different policies.
	//
	// Include a predicate list in your Intent constructors.
	Predicate func(ctx context.Context, desiredValues ...string) (bool, error)
)

const (
	// ResourcePathSeparator divides the sections of resource paths and masks. The value is used for parsing and printing their plain string reprsentations.
	ResourcePathSeparator = "/"

	// ResourcePathWildCardSegment matches any single value of a resource path. The segment must be present to match.
	ResourcePathWildCardSegment = "*"

	// ResourcePathWildCardTail matches any resource path ending segments, present or not. Any resource path mask value after ResourcePathWildCardSegment is ignored.
	ResourcePathWildCardTail = ">"

	ActionCreate   Action = "create"
	ActionRetrieve Action = "retrieve"
	ActionUpdate   Action = "update"
	ActionDelete   Action = "delete"
	ActionQuery    Action = "query"

	ActionAssign    Action = "assign"
	ActionUnassign  Action = "unassign"
	ActionBlock     Action = "block"
	ActionUnblock   Action = "unblock"
	ActionReset     Action = "reset"
	ActionRecover   Action = "recover"
	ActionPromote   Action = "promote"
	ActionDemote    Action = "demote"
	ActionUpgrade   Action = "upgrade"
	ActionDowngrade Action = "downgrade"
	ActionCommit    Action = "commit"
	ActionClear     Action = "clear"
	ActionInstall   Action = "install"
)

// An Intent is the desire of a role to carry out a given [Action] on a resource.
type Intent struct {
	Action       Action
	ResourcePath []string
	Predicates   map[string]Predicate
}

func (i *Intent) String() string {
	return fmt.Sprintf("perform action %q on resource %q",
		i.Action,
		strings.Join(i.ResourcePath, ResourcePathSeparator))
}

// MatchAction returns true if the [Intent] [Action] matches exactly any one of the provided actions.
func (i *Intent) MatchAction(actions ...Action) bool {
	for _, action := range actions {
		if action == i.Action {
			return true
		}
	}
	return false
}

// MatchResourcePath returns true if the [Intent] resource path matches mask segments. The [ResourcePathWildCardSegment] matches any present value. The [ResourcePathWildCardTail] matches any values to the end of the path.
func (i *Intent) MatchResourcePath(mask ...string) bool {
	// available := strings.Split(i.ResourcePath, ResourcePathSeparator)
	lenA, lenB := len(i.ResourcePath), len(mask)

	if lenB > lenA {
		if mask[lenA] != ResourcePathWildCardTail {
			return false
		}
		mask = mask[:lenA] // chop tail // TODO: test this.
	}
	for position, p := range mask {
		switch p {
		case ResourcePathWildCardTail:
			return true
		case ResourcePathWildCardSegment, i.ResourcePath[position]:
			continue
		default:
			return false
		}
	}
	if lenB < lenA {
		return false // must match every segment
	}
	return true
}

// MatchPredicate returns true if [Intent] [Predicate] was satisfied with provided desired values.
func (i *Intent) MatchPredicate(ctx context.Context, name string, desiredValues ...string) (result bool, err error) {
	// initiliazed map can be checked using `_, ok :=`
	// if i.Predicates == nil {
	// 	return false, fmt.Errorf("cannot evaluate predicate for property %q: %w", property, ErrNoPredicates)
	// }

	predicate, ok := i.Predicates[name]
	if !ok {
		return false, &PredicateError{
			Name:  name,
			Cause: ErrPredicateNotFound,
		}
	}

	if result, err = predicate(ctx, desiredValues...); err != nil {
		return false, &PredicateError{
			Name:  name,
			Cause: err,
		}
	}

	return
}

// PredicateEither combines a [Predicate] list into one that succeeds on first positive match.
func PredicateEither(ps ...Predicate) Predicate {
	return func(ctx context.Context, desiredValues ...string) (ok bool, err error) {
		for _, p := range ps {
			ok, err = p(ctx, desiredValues...)
			if err != nil {
				return false, err
			}
			if ok {
				return
			}
		}
		return false, nil
	}
}

// PredicateEach combines a [Predicate] list into one that succeeds only if each predicate matches. An empty predicate list always return false.
func PredicateEach(ps ...Predicate) Predicate {
	return func(ctx context.Context, desiredValues ...string) (ok bool, err error) {
		for _, p := range ps {
			ok, err = p(ctx, desiredValues...)
			if err != nil {
				return false, err
			}
			if !ok {
				return
			}
		}
		return len(ps) > 0, nil
	}
}
