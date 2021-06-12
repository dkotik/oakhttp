package oakacs

import (
	"github.com/rs/xid"
)

const recoveryCodeLength = 64

// Recovery holds a code that can restore access to an Identity.
type Recovery struct {
	Identity xid.ID
	Code     [recoveryCodeLength]byte
}
