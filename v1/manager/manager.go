/*

Package oakmanager provides administrative facility for OakACS.

*/
package oakmanager

import "github.com/dkotik/oakacs/v1"

const (
	service = "oakacs"
	domain  = "universal"
)

type backend interface {
	GroupRepository
	SecretRepository
	IntegrityLockRepository
}

// Manager provides the facility to persist all the data necessary for the Oak Access Control System.
type Manager struct {
	acs  *oakacs.AccessControlSystem
	repo backend
}
