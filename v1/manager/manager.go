package oakmanager

import "github.com/dkotik/oakacs/v1"

const (
	service = "oakacs"
	domain  = "universal"
)

// Manager provides the facility to persist all the data necessary for the Oak Access Control System.
type Manager struct {
	acs *oakacs.AccessControlSystem
}
