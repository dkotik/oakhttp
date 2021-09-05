package oakacs

import (
	"bytes"
	"context"
	"errors"
	"time"
)

type (
	AuthenticationRequest struct {
		Authenticator string
		Identity      []byte
		Secret        []byte
		Role          []byte
		// MFA token?
		// SAML token?
	}

	Authenticator interface {
		NewSecretFrom([]byte) ([]byte, error)
		TrySecret(secret []byte, tokens [][]byte) error
		// TrySecret([]byte, Identity) error
	}
)

func (r *AuthenticationRequest) MatchingRole(i Identity) []byte {
	roles := i.AvailableRoles()
	if len(roles) == 0 {
		return nil
	}
	if len(r.Role) == 0 {
		return roles[0]
	}
	for _, role := range roles {
		if bytes.Compare(role, r.Role) == 0 {
			return r.Role
		}
	}
	return nil
}

func (acs *AccessControlSystem) Authenticate(ctx context.Context, r *AuthenticationRequest) (session *Session, err error) {
	defer func() {
		// TODO: modulate time here to avoid betraying proof of existance
		// compare to a random string to modulate?
		acs.Broadcast(ctx, EventTypeAuthentication, err) // wrap error
	}()
	// TODO: throttle attempts by user
	auth, ok := acs.authenticators[r.Authenticator]
	if !ok {
		return nil, errors.New("chosen authenticator is not active") // TODO: standardize
	}

	// 1. Identify: locate the identity
	identity, err := acs.identities.Retrieve(ctx, r.Identity)
	if err != nil {
		return nil, err
	}

	// 2. Authenticate: confirm identity using a secret
	if err = auth.TrySecret(r.Secret, identity.TokensFor(r.Authenticator)); err != nil {
		return nil, err
	}

	// 3. Assume: select the correct role.
	role := r.MatchingRole(identity)
	if len(role) == 0 {
		return nil, errors.New("this identity appears to be banned") // standardize
	}
	duration, err := acs.roles.GetDuration(ctx, role)
	if err != nil {
		return nil, err
	}

	// 4. Issue: establish a session.
	session = &Session{
		Identity: r.Identity,
		Role:     role,
		Deadline: time.Now().Add(duration),
	}
	if err = acs.sessions.Create(ctx, session); err != nil {
		return nil, err
	}
	if len(session.UUID) == 0 {
		return nil, errors.New("new session UUID is nil") // standardize? rare
	}
	return session, nil
}
