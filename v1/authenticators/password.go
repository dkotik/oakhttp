package authenticators

import "github.com/dkotik/oakacs/v1"

var _ oakacs.Authenticator = (*Password)(nil)

// Register(ctx context.Context, tokenOrPassword string) (*Secret, error)
// Compare(ctx context.Context, tokenOrPassword string, secret *Secret) error
// TODO: MULTIHASH SUPPORT INSIDE THE Authenticator

// if _, ok := acs.hashers["default"]; !ok {
//     // TODO: confirm that those parameters are optimal
//     acs.hashers["default"] = NewHasherArgon2id(3, 64*1024, 4)
// }

// func NewHasherArgon2id(timeCost, memoryCost uint32, threads uint8) Hasher {
// 	return &hasherArgon2id{
// 		TimeCost:   timeCost,
// 		MemoryCost: memoryCost,
// 		Threads:    threads,
// 	}
// }
//
// type hasherArgon2id struct {
// 	TimeCost   uint32
// 	MemoryCost uint32
// 	Threads    uint8
// }
//
// func (h *hasherArgon2id) Hash(secret, salt []byte) []byte {
// 	return argon2.IDKey(secret, salt,
// 		h.TimeCost, h.MemoryCost, h.Threads, uint32(len(salt)))
// }

// Match checks if provided secret is valid.
// func (acs *AccessControlSystem) Match(secret string, against Secret) bool {
// 	if against.Type == SecretDisabled {
// 		// acs.logger.Warn("identity attempted to authenticate with a disabled secret",
// 		// 	zap.String("identity", against.Identity.String()),
// 		// 	zap.String("secret", against.UUID.String()))
// 		return false
// 	}
// 	hasher, ok := acs.hashers[against.HashedWith]
// 	if !ok {
// 		// acs.logger.Error("authentication error",
// 		// 	zap.Error(fmt.Errorf("hasher %q is not registered", against.HashedWith)))
// 		return false
// 	}
// 	// acs.logger.Info("identity authenticated",
// 	// 	zap.String("identity", against.Identity.String()),
// 	// 	zap.String("secret", against.UUID.String()))
// 	return 1 == subtle.ConstantTimeCompare(
// 		[]byte(against.Hash), hasher.Hash([]byte(secret), []byte(against.Salt)))
// }
