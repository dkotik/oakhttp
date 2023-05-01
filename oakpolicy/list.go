package oakpolicy

import (
	"context"
	"errors"
)

// FirstOf composes a [Policy] list, which returns when the first [Policy] that returns an error, [Allow], or [Deny]. Panics on empty list of a <nil> [Policy] inside the list.
func FirstOf(ps ...Policy) Policy {
	if len(ps) == 0 {
		panic(ErrEmptyPolicyList)
	}
	for _, p := range ps {
		if p == nil {
			panic(ErrNilPolicy)
		}
	}

	return func(ctx context.Context, a Action, r Resource) (err error) {
		for _, p := range ps {
			if err = p(ctx, a, r); err != nil {
				return err
			}
		}
		return nil
	}
}

// EachOf composes a [Policy] list into one that succeeds only if each included [Policy] returns an [Allow]. Panics on empty list of a <nil> [Policy] inside the list.
func EachOf(ps ...Policy) Policy {
	if len(ps) == 0 {
		panic(ErrEmptyPolicyList)
	}
	for _, p := range ps {
		if p == nil {
			panic(ErrNilPolicy)
		}
	}

	return func(ctx context.Context, a Action, r Resource) (err error) {
		for _, p := range ps {
			if err = p(ctx, a, r); !errors.Is(err, Allow) {
				return err
			}
		}
		return Allow
	}
}
