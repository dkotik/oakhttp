package oakacs

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"time"
)

// TODO: I should treat all secrets as tokens, hashed passwords should just present a token-like API.
// TODO: AccessControlSystem Authenticate should also indicate the METHOD OF AUTH! and be matched with an authenticator. JUST ONE TOKEN PER AUTHENTICATOR?

type Authenticator interface {
	Generate(ctx context.Context, tokenOrPassword string) (*Secret, error)
	Compare(ctx context.Context, tokenOrPassword string, secret *Secret) error
}

// Authenticate matches provided user and password to an indentity.
func (acs *AccessControlSystem) Authenticate(ctx context.Context, user, tokenOrPassword, authenticator string) (s *Session, err error) {
	defer func() {
		event := Event{
			ctx:  ctx,
			Type: EventTypeAuthenticationSuccess,
		}
		if err != nil {
			err = fmt.Errorf("authentication error: %w", err)
			event.Error = err
			event.Type = EventTypeAuthenticationFailure
		}
		acs.Broadcast(event)
	}()
	// throttle attempts by user

	auth, ok := acs.authenticators[authenticator]
	if !ok {
		return nil, errors.New("chosen authenticator is not active")
	}

	// retrieve identity
	identity, err := acs.backend.RetrieveIdentity(ctx, user)
	if err != nil {
		// modulate time here to avoid betraying proof of existance
		return nil, err
	}

	secret, err := acs.backend.RetrieveSecret(ctx, identity.UUID, authenticator)
	if err != nil {
		return nil, err
	}
	if err = auth.Compare(ctx, tokenOrPassword, secret); err != nil {
		return nil, err
	}
	secret.Used = time.Now()
	if err = acs.backend.UpdateSecret(ctx, secret); err != nil {
		return
	}

	// start session
	// attach role to session
	// issue event
	return
}

// Match checks if provided secret is valid.
func (acs *AccessControlSystem) Match(secret string, against Secret) bool {
	if against.Type == SecretDisabled {
		// acs.logger.Warn("identity attempted to authenticate with a disabled secret",
		// 	zap.String("identity", against.Identity.String()),
		// 	zap.String("secret", against.UUID.String()))
		return false
	}
	hasher, ok := acs.hashers[against.HashedWith]
	if !ok {
		// acs.logger.Error("authentication error",
		// 	zap.Error(fmt.Errorf("hasher %q is not registered", against.HashedWith)))
		return false
	}
	// acs.logger.Info("identity authenticated",
	// 	zap.String("identity", against.Identity.String()),
	// 	zap.String("secret", against.UUID.String()))
	return 1 == subtle.ConstantTimeCompare(
		[]byte(against.Hash), hasher.Hash([]byte(secret), []byte(against.Salt)))
}
