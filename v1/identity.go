package oakacs

import (
	"time"

	"github.com/rs/xid"
)

// Identity represents a unique acting entity, a human or a machine.
type Identity struct {
	UUID              xid.ID
	Name              Name
	Salt, Password    [hashSize]byte
	Groups            []Group
	HumanityConfirmed time.Time
}
