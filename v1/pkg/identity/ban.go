package identity

import (
	"time"

	"github.com/rs/xid"
)

// use ban event + watcher?

// Ban prevents an account from authenticating.
type Ban struct {
	UUID           xid.ID
	SourceIdentity xid.ID
	Reason         string
	MatchIdentity  string
	MatchRole      string
	MatchGroup     string
	Created        time.Time
	Expires        time.Time
}
