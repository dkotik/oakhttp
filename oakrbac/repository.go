package oakrbac

import "fmt"

type RoleRepository interface {
	GetRole(name string) (Role, error)
	AddRole(Role) error
	CountRoles() int
}

// List iteration of around five items is faster than `map[string]Role`. Most applications will not have more than five frequently used roles.
type ListRepository struct {
	roles []Role
}

func (l *ListRepository) GetRole(name string) (r Role, err error) {
	for _, r = range l.roles {
		if name == r.Name() {
			return r, nil
		}
	}
	return nil, ErrRoleNotFound
}

func (l *ListRepository) AddRole(r Role) error {
	name := r.Name()
	for _, role := range l.roles {
		if name == role.Name() {
			return fmt.Errorf("duplicate role: %s", name)
		}
	}
	l.roles = append(l.roles, r)
	return nil
}

func (l *ListRepository) CountRoles() int {
	return len(l.roles)
}

type MapRepository struct {
	roles map[string]Role
}

func (m *MapRepository) GetRole(name string) (Role, error) {
	r, ok := m.roles[name]
	if ok {
		return r, nil
	}
	return nil, ErrRoleNotFound
}

func (m *MapRepository) AddRole(r Role) error {
	if m.roles == nil {
		m.roles = make(map[string]Role)
	}
	name := r.Name()
	if _, ok := m.roles[name]; ok {
		return fmt.Errorf("duplicate role: %s", name)
	}
	m.roles[name] = r
	return nil
}

func (m *MapRepository) CountRoles() int {
	return len(m.roles)
}
