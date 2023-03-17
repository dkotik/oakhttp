package oakrbac

import (
	"context"

	"golang.org/x/exp/slog"
)

type EventType string

const (
	EventTypeError                EventType = "authorization error"
	EventTypeAuthorizationGranted EventType = "authorization granted"
	EventTypeAuthorizationDenied  EventType = "authorization denied"
)

type Event struct {
	eventType EventType
	role      Role
	intents   []Intent
	policies  []Policy
	error     error
}

type Listener interface {
	Listen(context.Context, *Event)
}

func NewEvent(t EventType, r Role, intents []Intent, p []Policy, err error) *Event {
	return &Event{
		eventType: t,
		role:      r,
		intents:   intents,
		policies:  p,
		error:     err,
	}
}

func (e *Event) Type() EventType {
	return e.eventType
}

func (e *Event) Role() Role {
	return e.role
}

func (e *Event) Intents() []Intent {
	return e.intents
}

func (e *Event) Policies() []Policy {
	return e.policies
}

func (e *Event) Error() error {
	return e.error
}

func (e *Event) String() string {
	switch e.eventType {
	case EventTypeAuthorizationGranted:
		return "access granted"
	case EventTypeAuthorizationDenied:
		// if e.policies == nil {
		// 	return "access denied: no policy matched"
		// }
		return "access denied"
	case EventTypeError:
		return e.error.Error()
	default:
		return "uknown event"
	}
}

func (e *Event) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("type", string(e.eventType)),
		slog.Any("role", e.role),
		slog.Any("intents", e.intents),
		slog.Any("policies", e.policies),
		slog.Any("error", e.error),
	)
}
