package oakacs

import (
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

const (
	hashSize = 32
)

// Hash applies Argon2id over secure salt with decent default settings.
func Hash(s []byte) ([]byte, error) {
	var salt [hashSize]byte
	if _, err := rand.Read(salt[:]); err != nil {
		return nil, err
	}
	// TODO: those need to be tweaked
	return argon2.IDKey(s, salt[:], 3, 64*1024, 4, hashSize), nil
}
