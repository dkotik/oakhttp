package oakacs

import (
	"time"

	"github.com/rs/xid"
	"golang.org/x/crypto/argon2"
)

type secretType uint8

// The possible types of secrets:
const (
	SecretDisabled secretType = iota // Disabled secrets are old passwords?
	SecretPrimaryPassword
	SecretExpiredPassword
	SecretRecoveryCode
	SecretOAuthToken
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

// Hasher holds and applies a hashing algorythm with secure paramters. Returned hash should be the same length as the provided salt.
type Hasher interface {
	Hash(secret, salt []byte) []byte
	// Match(hash, salt string) bool
}

func NewHasherArgon2id(timeCost, memoryCost uint32, threads uint8) Hasher {
	return &hasherArgon2id{
		TimeCost:   timeCost,
		MemoryCost: memoryCost,
		Threads:    threads,
	}
}

type hasherArgon2id struct {
	TimeCost   uint32
	MemoryCost uint32
	Threads    uint8
}

func (h *hasherArgon2id) Hash(secret, salt []byte) []byte {
	return argon2.IDKey(secret, salt,
		h.TimeCost, h.MemoryCost, h.Threads, uint32(len(salt)))
}
