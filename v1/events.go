package oakacs

import (
	"context"
	"errors"
	"time"
)

type (
	// EventType resprents the kind of events that the ACS may issue.
	EventType           uint8
	eventContextKeyType string
)

const (
	eventContextKeyIP eventContextKeyType = "ip"

	// EventTypeUnknown indicates an unexpected event, should be treated as fatal error.
	EventTypeUnknown = iota
	EventTypeSession
	EventTypeAuthentication
	EventTypeAuthorization
	EventTypeMaintenance
)

type Event struct {
	Context context.Context // important for contextual unpacking
	When    time.Time
	Type    EventType
	Error   error
	// Session Session // it may be or not be already in the context
}

func (e *Event) IP() (string, error) {
	if e.Context != nil { // TODO: is this needed?
		val := e.Context.Value(eventContextKeyIP)
		switch ip := val.(type) {
		case string:
			return ip, nil
		}
	}
	return "", errors.New("ip address is not associated with context")
}

func (e *Event) String() string {
	switch e {
	// TODO: fill out
	// case EventTypeAuthenticationSuccess:
	// 	return "authenticated"
	// case EventTypeAuthenticationFailure:
	// 	return "rejected"
	// case EventTypeAuthorizationAllowed:
	// 	return "authorized"
	// case EventTypeAuthorizationDeniedByPermission:
	// 	return "denied"
	// case EventTypeAuthorizationDeniedByDefault:
	// 	return "denied"
	}
	return "<undocumented-event>"
}

// Broadcast attempts to notify all the subscribers. The dispatch is non-blocking, so if subscriber is busy, the event misses.
func (acs *AccessControlSystem) Broadcast(ctx context.Context, t EventType, err error) {
	for _, c := range acs.subscribers {
		select {
		case c <- Event{
			Context: ctx,
			When:    time.Now(),
			Type:    t,
			Error:   err,
		}:
		default:
			// TODO: issue warning / error about skipped events?
		}
	}
}
