package oakacs

import (
	"context"
	"crypto/subtle"
	"fmt"
)

// TODO: I should treat all secrets as tokens, hashed passwords should just present a token-like API.
// TODO: AccessControlSystem Authenticate should also indicate the METHOD OF AUTH! and be matched with an authenticator. JUST ONE TOKEN PER AUTHENTICATOR?

type Authenticator interface {
	Authenticate(ctx context.Context, user *Identity, tokenOrPassword string) (err error)
}

// Authenticate matches provided user and password to an indentity.
func (acs *AccessControlSystem) Authenticate(ctx context.Context, user, tokenOrPassword string) (s Session, err error) {
	// throttle attempts

	// retrieve identity
	identity, err := acs.backend.RetrieveIdentity(ctx, user)
	if err != nil {
		// modulate time here to avoid betraying proof of existance
		return Session{}, err
	}
	fmt.Println("got identity", identity)

	// match the secret
	// save against.Used
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
