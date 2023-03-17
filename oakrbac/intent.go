package oakrbac

import (
	"fmt"

	"golang.org/x/exp/slog"
)

// An Intent is the desire of a role to carry out a given [Action] on a resource.
type Intent interface {
	Action() Action
	ResourcePath() ResourcePath
	String() string // for caching
}

type basicIntent struct {
	action       Action
	resourcePath ResourcePath
}

func NewIntent(action Action, resourcePath ResourcePath) Intent {
	return &basicIntent{
		action:       action,
		resourcePath: resourcePath,
	}
}

func (i *basicIntent) Action() Action {
	return i.action
}

func (i *basicIntent) ResourcePath() ResourcePath {
	return i.resourcePath
}

func (i *basicIntent) String() string {
	return fmt.Sprintf("intending action %q on resource %q", i.action, i.resourcePath)
}

func (i *basicIntent) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("action", string(i.action)),
		slog.String("resource", i.resourcePath.String()),
	)
}
