package oakacs

import "github.com/rs/xid"

// EventType resprents the kind of events that the ACS may issue.
type EventType uint8

const (
	// EventTypeUnknown indicates an unexpected event, should be treated as fatal error.
	EventTypeUnknown = iota
	// EventTypeAuthentication marks a role being succesffully connected to a session.
	EventTypeAuthentication
	EventTypeAuthorizationAllowed
	EventTypeAuthorizationDenied
)

// TODO: should event be an interface instead? to allow different types of events

type Event struct {
	Type    EventType
	Session xid.ID
	Role    xid.ID
	Error   error
}

// Broadcast attempts to notify all the subscribers. The dispatch is non-blocking, so if subscriber is busy, the event misses.
func (acs *AccessControlSystem) Broadcast(e Event) {
	for _, c := range acs.subscribers {
		select {
		case c <- e:
		default:
		}
	}
}
