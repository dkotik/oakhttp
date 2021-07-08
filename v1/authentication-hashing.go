package oakacs

import (
	"golang.org/x/crypto/argon2"
)

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
