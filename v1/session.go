package oakacs

import (
	"context"
	"errors"

	"github.com/rs/xid"
)

// SessionBind retrieves the Session from ephemeral storage and binds it to context.
func (acs *AccessControlSystem) SessionBind(
	ctx context.Context, id xid.ID, differentiator string,
) (context.Context, error) {
	// TODO: could I request by prefix, then check if the rest of the id matches?
	// session, err := acs.ephemeral.Sessions.Retrieve(ctx, id)
	// if err != nil {
	// 	// TODO: can I find a session with certain prefix? with badger I could, close them all - differentiator did not match?
	// 	return ctx, fmt.Errorf("cannot retrieve session: %w", err)
	// }
	// if bytes.Compare(session.Differentiator[:], []byte(differentiator)) != 0 {
	// 	err = errors.New("session breached: differentiator did not match")
	// 	if rerr := acs.ephemeral.Sessions.Delete(ctx, id); rerr != nil {
	// 		if rerr = acs.ephemeral.Sessions.Delete(ctx, id); rerr != nil { // retry
	// 		}
	// 	}
	// 	return nil, err
	// }
	// TODO: check deadline
	// TODO: check activity deadline
	return context.WithValue(ctx, acs.sessionContextKey, nil), nil
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

func (acs *AccessControlSystem) PushSession() {}
func (acs *AccessControlSystem) PullSession() {}
