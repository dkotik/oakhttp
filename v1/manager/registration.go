package oakmanager

import (
	"context"

	"github.com/dkotik/oakacs/v1"
)

// Registration creates an Entity within ACS backend.
type Registration func(context.Context, *oakacs.Identity) error

// var salt [hashSize]byte
// if _, err := rand.Read(salt[:]); err != nil {
// 	return "", "", err
// }

// // Register creates a new identity and establishes a session.
// func (acs *AccessControlSystem) Register(ctx context.Context, user, password string) (i Identity, s Session, err error) {
// 	return
// }
//
// // AddRecoveryCode generates a new random paper code.
// func (acs *AccessControlSystem) AddRecoveryCode(ctx context.Context, i *Identity) (s *Secret, err error) {
// 	// s = &Secret{}
//
// 	return
// }
