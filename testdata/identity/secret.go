package identity

import (
	"time"

	"github.com/rs/xid"
)

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
