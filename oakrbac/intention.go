package oakrbac

import (
	"fmt"

	"golang.org/x/exp/slog"
)

// An Intention is the desire of a [Role] to carry out a given [Action] on a [ResourcePath].
type Intention interface {
	Action() Action
	ResourcePath() ResourcePath
	String() string // for caching
}

type basicIntention struct {
	action       Action
	resourcePath ResourcePath
}

func NewIntention(action Action, resourcePath ResourcePath) Intention {
	return &basicIntention{
		action:       action,
		resourcePath: resourcePath,
	}
}

func (i *basicIntention) Action() Action {
	return i.action
}

func (i *basicIntention) ResourcePath() ResourcePath {
	return i.resourcePath
}

func (i *basicIntention) String() string {
	return fmt.Sprintf("intending action %q on resource %q", i.action, i.resourcePath)
}

func (i *basicIntention) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("action", string(i.action)),
		slog.String("resource", i.resourcePath.String()),
	)
}
