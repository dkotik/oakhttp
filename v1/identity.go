package oakacs

import (
	"github.com/rs/xid"
)

// Identity represents a unique acting entity, a human or a machine.
type Identity struct {
	UUID   xid.ID
	Name   Name
	Groups []Group
}
