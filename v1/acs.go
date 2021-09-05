package oakacs

import (
	"context"
	"fmt"
)

type (
	Secret   func(token []byte) error
	Identity interface {
		ListSecrets(authenticator string) ([]Secret, error)
		ListRoles() ([]byte, error)
		ListGroups() ([]byte, error)
	}
	IdentityRepository interface {
		Retrieve(context.Context, []byte) (Identity, error)
	}

	Action interface {
		Disclose(attribute string) (value interface{})
		String() string
	}
	Role           func(Action) error
	RoleRepository interface {
		Retrieve(context.Context, []byte) (Role, error)
	}

	sessionContextKeyType string

	// Session connects an Identity to a combined list of allowed actions accessible to the Identity.
	Session struct {
		// Differentiator [24]byte // to prevent session ID guessing
		UUID     []byte
		Role     []byte
		Identity []byte
		// Deadline is already set
		// Identity       xid.ID
		// Role           xid.ID
		// Created        time.Time
		// LastRetrieved  time.Time
		// Values         map[string]interface{}
	}
	SessionRepository interface {
		Create(context.Context, *Session) error
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

	acs := &AccessControlSystem{
		authenticators: make(map[string]Authenticator),
	}
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
	// TokenValidator    *TokenValidator

	subscribers    []chan<- (Event)
	authenticators map[string]Authenticator
	// ephemeral      EphemeralRepository
	// persistent     PermissionsRepository
}
