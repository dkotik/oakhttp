package oakacs

import (
	"context"
	"crypto/subtle"
	"fmt"
	"time"

	"github.com/rs/xid"
	"go.uber.org/zap"
)

type secretType uint8

// The possible types of secrets:
const (
	SecretDisabled secretType = iota
	SecretPrimaryPassword
	SecretRecoveryCode
)

// const recoveryCodeLength = 64
//
// // Recovery holds a code that can restore access to an Identity.
// type Recovery struct {
// 	Code [recoveryCodeLength]byte
// }

// Secret is a password or a recovery code.
type Secret struct {
	UUID       xid.ID
	Identity   xid.ID
	Salt       string
	Hash       string
	HashedWith string
	Type       secretType
	Used       time.Time
}

// Authenticate matches provided user and password to an indentity.
func (acs *AccessControlSystem) Authenticate(ctx context.Context, user, password string) (i Identity, s Session, err error) {
	// throttle attempts
	// modulate time
	// retrieve identity
	// match the secret
	// save against.Used
	// start session
	return
}

// Match checks if provided secret is valid.
func (acs *AccessControlSystem) Match(secret string, against Secret) bool {
	if against.Type == SecretDisabled {
		acs.logger.Warn("identity attempted to authenticate with a disabled secret",
			zap.String("identity", against.Identity.String()),
			zap.String("secret", against.UUID.String()))
		return false
	}
	hasher, ok := acs.hashers[against.HashedWith]
	if !ok {
		acs.logger.Error("authentication error",
			zap.Error(fmt.Errorf("hasher %q is not registered", against.HashedWith)))
		return false
	}
	acs.logger.Info("identity authenticated",
		zap.String("identity", against.Identity.String()),
		zap.String("secret", against.UUID.String()))
	return 1 == subtle.ConstantTimeCompare(
		[]byte(against.Hash), hasher.Hash([]byte(secret), []byte(against.Salt)))
}
