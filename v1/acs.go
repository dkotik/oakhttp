package oakacs

import (
	"fmt"
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
	TokenValidator    *TokenValidator

	subscribers    []chan<- (Event)
	authenticators map[string]Authenticator
	ephemeral      EphemeralRepository
	persistent     PermissionsRepository
}
