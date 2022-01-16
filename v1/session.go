package oakacs

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/rs/xid"
)

type SessionID [24]byte

// Bind retrieves the Session and binds it to context. If the session does not exist, it is created.
func (acs *AccessControlSystem) Bind(ctx context.Context, token string) (context.Context, error) {
	b, err := base64.StdEncoding.DecodeString(token)
	if len(b) > 12 {
		session, err := acs.sessions.Retrieve(ctx, b[:12])
		if err == nil {
			if bytes.Compare(session.UUID[12:], b[12:]) == 0 {
				return context.WithValue(ctx, acs.sessionContextKey, &session), nil
			}
			return nil, errors.New("session differentiator did not match")
		}
		return nil, err
	}

	var (
		sid  SessionID
		nxid = xid.New()
	)
	n := copy(sid[:], nxid[:])
	more, err := rand.Read(sid[len(sid)-n:])
	if err != nil {
		return nil, err
	}
	if more+n < len(sid) {
		return nil, errors.New("not enough random bytes")
	}
	// b =

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

// Continue retrieves the session state from context.
func (acs *AccessControlSystem) Continue(ctx context.Context) (*Session, error) {
	switch s := ctx.Value(acs.sessionContextKey).(type) {
	case Session:
		return &s, nil
	default:
		// TODO: standardize the error
		return nil, errors.New("execution context is not authenticated")
	}
}
