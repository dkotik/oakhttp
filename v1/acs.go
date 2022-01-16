package oakacs

import (
	"context"
	"fmt"
	"time"
)

type (
	Identity interface {
		TokensFor(authenticator string) [][]byte
		AvailableRoles() [][]byte
	}
	IdentityRepository interface {
		Create(context.Context, AuthenticationRequest) (Identity, []byte, error)
		Retrieve(context.Context, []byte) (Identity, error) // identity can be banned
	}

	Action interface {
		Disclose(attribute string) (value interface{})
		String() string
	}
	Role           func(Action) error
	RoleRepository interface {
		Retrieve(context.Context, []byte) (Role, error) // role can be banned
		GetDuration(context.Context, []byte) (time.Duration, error)
	}

	sessionContextKeyType string

	// Session connects an Identity to a combined list of allowed actions accessible to the Identity.
	Session struct {
		// Differentiator [24]byte // to prevent session ID guessing
		UUID     []byte
		Role     []byte
		Identity []byte
		Deadline time.Time
		// Identity       xid.ID
		// Role           xid.ID
		// Created        time.Time
		// LastRetrieved  time.Time
		// Values         map[string]interface{}
	}
	SessionRepository interface {
		Create(context.Context, *Session) error
		// Update(context.Context, []byte, func(*Session) error) error
		Retrieve(context.Context, []byte) (*Session, error)
		Delete(context.Context, []byte) error
	}
)

// NewAccessControlSystem sets up an access control system.
func NewAccessControlSystem(withOptions ...Option) (*AccessControlSystem, error) {
	var err error
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to setup OakACS: %w", err)
		}
	}()

	acs := &AccessControlSystem{}
	err = WithOptions(withOptions...)(acs)
	if err != nil {
		return nil, err
	}

	// Fill out defaults:
	if acs.subscribers == nil {
		WithSubscribers()(acs)
	}

	return acs, nil
}

// AccessControlSystem manages sessions.
type AccessControlSystem struct {
	sessionContextKey sessionContextKeyType
	sessions          SessionRepository
	identities        IdentityRepository
	roles             RoleRepository
	// TokenValidator    *TokenValidator

	subscribers []chan<- (Event)
	// ephemeral      EphemeralRepository
	// persistent     PermissionsRepository
}
