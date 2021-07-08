package oakacs

import (
	"time"

	"github.com/rs/xid"
)

// Identity represents a unique acting entity, a human or a machine.
type Identity struct {
	UUID              xid.ID
	Name              string
	Groups            []Group
	Secrets           []Secret
	HumanityConfirmed time.Time
}
