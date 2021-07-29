package oakacs

import (
	"time"

	"github.com/rs/xid"
)

// Identity represents a unique acting entity, a human or a machine.
type Identity struct {
	UUID              xid.ID
	Name              string
	Groups            []Group // the order matters for default roles
	Secrets           []Secret
	HumanityConfirmed time.Time
}

// Group holds roles that identities may assume.
type Group struct {
	UUID            xid.ID
	Name            string
	DefaultRole     Role
	AscendableRoles []Role
}

// Secret is a password or a recovery code.
type Secret struct {
	UUID          xid.ID
	Identity      xid.ID
	Authenticator string
	Label         string
	Token         string
	Expires       time.Time
	Used          time.Time
}
