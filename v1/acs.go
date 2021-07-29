package oakacs

import (
	"context"
	"fmt"
)

type acsContextKeyType string

type backend interface {
	SessionRepository
	PermissionsRepository

	RetrieveIdentity(ctx context.Context, name string) (*Identity, error)
}

// NewAccessControlSystem sets up an access control system.
func NewAccessControlSystem(withOptions ...Option) (*AccessControlSystem, error) {
	var err error
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to setup OakACS: %w", err)
		}
	}()

	acs := &AccessControlSystem{
		hashers: make(map[string]Hasher, 1),
	}
	err = WithOptions(withOptions...)(acs)
	if err != nil {
		return nil, err
	}
	if _, ok := acs.hashers["default"]; !ok {
		// TODO: confirm that those parameters are optimal
		acs.hashers["default"] = NewHasherArgon2id(3, 64*1024, 4)
	}

	// Fill out defaults:
	if acs.subscribers == nil {
		WithSubscribers()(acs)
	}

	return acs, nil
}

// AccessControlSystem manages sessions.
type AccessControlSystem struct {
	sessionContextKey acsContextKeyType
	// Backend?
	//
	TokenValidator *TokenValidator

	subscribers []chan<- (Event)
	hashers     map[string]Hasher
	backend     backend
}
