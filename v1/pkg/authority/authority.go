package authority

import "errors"

type (
	Action interface {
		Disclose(attribute string) (value string)
		String() string
	}
	Role interface {
		Authorize(Action) error
		String() string
	}
	Authority interface {
		Authorize(Role, Action) error
		Register(Role) error
		Remove(Role) error
		String() string
	}
)

func authorize(roleUUID string, action Action) error {
	actors := make(chan func(map[string]Role))
	errc := make(chan error)
	actors <- func(roles map[string]Role) {
		role, ok := roles[roleUUID]
		if !ok {
			errc <- errors.New("role does not exist")
			return
		}
		errc <- role.Authorize(action)
	}
	return <-errc
}

func newRoleStack() chan<- func(map[string]Role) {
	actors := make(chan func(map[string]Role))
	stack := make(map[string]Role)
	go func() {
		for actor := range actors {
			actor(stack)
		}
	}()
	return actors
}
