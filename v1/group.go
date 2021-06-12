package oakacs

import (
	"github.com/rs/xid"
)

// Group holds roles that identities may assume.
type Group struct {
	UUID    xid.ID
	Name    Name
	Primary Role
	Roles   []Role
}
