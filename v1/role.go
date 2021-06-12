package oakacs

import (
	"github.com/rs/xid"
)

// Role holds a set of allowed actions.
type Role struct {
	UUID    xid.ID
	Name    Name
	Actions []Action
}
