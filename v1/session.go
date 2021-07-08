package oakacs

import (
	"time"

	"github.com/rs/xid"
)

// Session connects an Identity to a combined list of allowed actions accessible to the Identity.
type Session struct {
	UUID     xid.ID
	Identity xid.ID
	Role     xid.ID
	Deadline time.Time
}
