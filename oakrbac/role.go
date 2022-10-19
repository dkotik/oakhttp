package oakrbac

import (
	"context"
	"errors"
)

// AccessDeniedError represents a [Role] authorization failure.
type AccessDeniedError struct {
	Cause error
}

func (e *AccessDeniedError) Unwrap() error {
	return e.Cause
}

func (e *AccessDeniedError) Error() string {
	return "access denied"
}

// A Role is an [Intent] authorization provider.
type Role interface {
	// Authorize has an identical signature to [Policy] because it also validates an [Intent]. However, Authorize does not return either [Allow] or [Deny] sentinel values but converts them to `nil` and [AccessDeniedError] respectively.
	Authorize(context.Context, *Intent) error
}

type observedRole struct {
	name     string
	logger   Logger
	policies []Policy
}

func (o *observedRole) Authorize(ctx context.Context, i *Intent) (err error) {
	for _, policy := range o.policies {
		if err = policy(ctx, i); err != nil {
			if errors.Is(err, Allow) {
				if err = o.logger(ctx, &Event{
					IsAllowed: true,
					Intent:    i,
					Role:      o.name,
					Policy:    policy,
					Error:     nil,
				}); err != nil {
					return &AccessDeniedError{Cause: err}
				}
				return nil // policy matched
			}
			if errors.Is(err, Deny) {
				if err = o.logger(ctx, &Event{
					IsAllowed: false,
					Intent:    i,
					Role:      o.name,
					Policy:    policy,
					Error:     nil,
				}); err != nil {
					return &AccessDeniedError{Cause: err}
				}
				return &AccessDeniedError{Cause: Deny} // policy blocked
			}
			if lerr := o.logger(ctx, &Event{ // lerr to avoid overwriting err
				IsAllowed: false,
				Intent:    i,
				Role:      o.name,
				Policy:    policy,
				Error:     err,
			}); lerr != nil {
				return &AccessDeniedError{Cause: lerr}
			}
			return err // unexpected error
		}
	}
	if err = o.logger(ctx, &Event{
		IsAllowed: false,
		Intent:    i,
		Role:      o.name,
		Policy:    nil,
		Error:     ErrNoPolicyMatched,
	}); err != nil {
		return &AccessDeniedError{Cause: err}
	}
	return &AccessDeniedError{Cause: ErrNoPolicyMatched}
}

type silentRole struct {
	name     string
	policies []Policy
}

func (s *silentRole) Authorize(ctx context.Context, i *Intent) (err error) {
	for _, policy := range s.policies {
		if err = policy(ctx, i); err != nil {
			if errors.Is(err, Allow) {
				return nil // policy matched
			}
			if errors.Is(err, Deny) {
				return &AccessDeniedError{Cause: Deny} // policy blocked
			}
			return err // unexpected error
		}
	}
	return &AccessDeniedError{Cause: ErrNoPolicyMatched}
}
