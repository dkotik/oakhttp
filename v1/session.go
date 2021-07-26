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
	switch s := ctx.Value(acs.sessionContextKey).(type) {
	case Session:
		return s, nil
	default:
		// TODO: standardize the error
		return Session{}, errors.New("execution context is not authenticated")
	}
}

// Bind rolls session into the provided context with deadline.
func (acs *AccessControlSystem) bind(ctx context.Context, s Session) context.Context {
	return context.WithValue(ctx, acs.sessionContextKey, s)
}
