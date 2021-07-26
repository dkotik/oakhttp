package oakacs

import (
	"context"
	"errors"
	"time"

	"github.com/rs/xid"
)

// Session connects an Identity to a combined list of allowed actions accessible to the Identity.
type Session struct {
	UUID     xid.ID
	Identity xid.ID
	Role     xid.ID
	Deadline time.Time
}

// SessionFrom retrieves the session state from context.
func (acs *AccessControlSystem) SessionFrom(ctx context.Context) (Session, error) {
	s := ctx.Value(acs.sessionContextKey)
	if s == nil {
		// TODO: standardize the error
		return nil, errors.New("execution context is not authenticated")
	}
	return s, nil
}

// Bind rolls session into the provided context with deadline.
func (acs *AccessControlSystem) bind(ctx context.Context, s Session) context.Context {
	cd := context.WithDeadline(ctx, s.Deadline)
	return context.WithValue(cd, acs.sessionContextKey, s)
}
