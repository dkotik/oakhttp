/*

Package oakmanager provides administrative facility for OakACS.

*/
package oakmanager

import "github.com/dkotik/oakacs/v1"

const (
	service = "oakacs"
	RD      = "read" // TODO: move RD and WR to oakacs?
	WR      = "write"
	domain  = "universal"
)

// Manager provides the facility to persist all the data necessary for the Oak Access Control System.
type Manager struct {
	acs *oakacs.AccessControlSystem
	// persistent *oakacs.PersistentRepository
}
