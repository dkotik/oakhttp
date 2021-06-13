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

// Action represents something an Identity can do.
type Action struct {
	Name   Name
	Verb   Name
	Target Name
}
