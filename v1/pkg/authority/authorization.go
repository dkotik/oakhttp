package cueroles

import (
	"context"
	"time"
)

type (
	Role interface {
		HasPermissionFor(Action) error
		GetTimeouts() (time.Time, time.Time)
		String() string
	}

	Repository interface {
		Push(context.Context, string, Role) error
		Pull(context.Context, string) (Role, error)
		Delete(string) error
	}

	Action interface {
		Disclose(attribute string) (value interface{})
		String() string
	}
)
