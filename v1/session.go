package oakacs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/xid"
)

// SessionRepository persists Sessions.
type SessionRepository interface {
	CreateSession(context.Context, *Session) error
	RetreiveSession(context.Context, xid.ID) (*Session, error)
	UpdateSession(context.Context, xid.ID, func(*Session) error) error
	DeleteSession(context.Context, xid.ID) error
}

// Session connects an Identity to a combined list of allowed actions accessible to the Identity.
type Session struct {
	UUID           xid.ID
	Differentiator [32]byte // to prevent session ID guessing
	Identity       xid.ID
	Role           xid.ID
	Deadline       time.Time
	Values         map[string]interface{}
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

func (acs *AccessControlSystem) PushSession() {}
func (acs *AccessControlSystem) PullSession() {}

func (acs *AccessControlSystem) PushToken(ctx context.Context, s string, deadline time.Time) error {
	return acs.ephemeral.Push(ctx, s, deadline)
}

func (acs *AccessControlSystem) PullToken(ctx context.Context, s string) (string, error) {
	token, err := acs.ephemeral.Pull(ctx, s)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", token), nil
}
