package oakacs

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/xid"
	"go.uber.org/zap"
)

const (
	ACSService = "oakacs"

	DomainUniversal = "universal"
)

type acsContextKeyType string

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
	if acs.logger == nil {
		if err = WithLogger(nil)(acs); err != nil {
			return nil, err
		}
	}

	return acs, nil
}

// AccessControlSystem manages sessions.
type AccessControlSystem struct {
	// Backend?
	//
	TokenValidator *TokenValidator

	subscribers       []chan (Event)
	sessionContextKey acsContextKeyType
	hashers           map[string]Hasher
	logger            *zap.Logger
}

// GetSession recovers session from a request.
func (a *AccessControlSystem) GetSession(r *http.Request) (*Session, error) {
	val, ok := r.Context().Value(sessionContextKey).(xid.ID)
	if !ok {
		return nil, errors.New("session UUID is not in context")
	}
	// spew.Dump(val)
	// if err != nil {
	// 	return nil, err
	// }
	return nil, nil
}

// Protect wraps a handler and injects session into its context after checking throttling and access.
func (a *AccessControlSystem) Protect(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create session
		// inject new context
		h.ServeHTTP(w, r.WithContext(
			context.WithValue(r.Context(), sessionContextKey, "xid.ID")))
	})
}

// Close cleans up loose ends.
func (a *AccessControlSystem) Close() error {
	return a.logger.Sync()
}
