package oakacs

import (
	"github.com/rs/xid"
)

// Entity is a fragment designed for composing into models.
type Entity struct {
	UUID xid.ID
	Name Name
}
