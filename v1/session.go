package oakacs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/xid"
)

type (
	sessionContextKeyType string

	// Session connects an Identity to a combined list of allowed actions accessible to the Identity.
	Session struct {
		UUID           xid.ID
		Differentiator string // to prevent session ID guessing
		Identity       xid.ID
		Role           xid.ID
		Created        time.Time
		LastRetrieved  time.Time
		Values         map[string]interface{}
	}
)

// SessionBind retrieves the Session from ephemeral storage and binds it to context.
func (acs *AccessControlSystem) SessionBind(
	ctx context.Context, id xid.ID, differentiator string,
) (context.Context, error) {
	session, err := acs.ephemeral.Sessions.Retrieve(ctx, id)
	if err != nil {
		return ctx, fmt.Errorf("cannot retrieve session: %w", err)
	}
	if session.Differentiator != differentiator {
		err = errors.New("session breached: differentiator did not match")
		acs.Broadcast(Event{
			ctx:     ctx,
			When:    time.Now(),
			Type:    EventTypeSessionBreached,
			Session: id,
			Role:    session.Role,
			Error:   err,
		})
		if rerr := acs.ephemeral.Sessions.Delete(ctx, id); rerr != nil {
			if rerr = acs.ephemeral.Sessions.Delete(ctx, id); rerr != nil { // retry
				acs.Broadcast(Event{
					ctx:     ctx,
					When:    time.Now(),
					Type:    EventTypeSessionBreached,
					Session: id,
					Role:    session.Role,
					Error:   err,
				})
			}
		}
		return nil, err
	}
	// TODO: check deadline
	// TODO: check activity deadline
	return context.WithValue(ctx, acs.sessionContextKey, session), nil
}

// SessionContinue retrieves the session state from context.
func (acs *AccessControlSystem) SessionContinue(ctx context.Context) (Session, error) {
	switch s := ctx.Value(acs.sessionContextKey).(type) {
	case Session:
		return s, nil
	default:
		// TODO: standardize the error
		return Session{}, errors.New("execution context is not authenticated")
	}
}

// // Bind rolls session into the provided context with deadline.
// func (acs *AccessControlSystem) bind(ctx context.Context, s Session) context.Context {
// 	return context.WithValue(ctx, acs.sessionContextKey, s)
// }

func (acs *AccessControlSystem) PushSession() {}
func (acs *AccessControlSystem) PullSession() {}

// func (acs *AccessControlSystem) PushToken(ctx context.Context, s string, deadline time.Time) error {
// 	return acs.ephemeral.Push(ctx, s, deadline)
// }
//
// func (acs *AccessControlSystem) PullToken(ctx context.Context, s string) (string, error) {
// 	token, err := acs.ephemeral.Pull(ctx, s)
// 	if err != nil {
// 		return "", err
// 	}
// 	return fmt.Sprintf("%s", token), nil
// }
