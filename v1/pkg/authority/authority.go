package cueroles

// TODO: most significant:

/*
RoleCandidate interface {
	IsAllowedToPerform(Action) error
	String() string
}
Action interface {
	Disclose(attribute string) (value interface{})
	String() string
}

type (
	Role struct {
		Allow []Match
		Deny  []Match
	}

	Match struct {
		Attribute string
		Method    string
		Pattern   string
	}

	Repository interface {
		Load(context.Context, string) (Role, error)
	}
)

*/

type (
	// Method is not needed
	Method           func(value interface{}, pattern string) error
	roleSetOperation func(map[string]*Role)

	// AuthorityCache func(ctx context.Context, role string, action Action) error

	Authority struct {
		operations chan (roleSetOperation)
		cache      map[string]map[string]error
		roleSet    map[string]*Role
		methodSet  map[string]Method
		repository Repository
	}
)

// func (a *Authority) matchAction(m *Match, action Action) error {
// 	matcher, ok := a.methodSet[m.Method]
// 	if !ok {
// 		return fmt.Errorf("method %q is not registered", m.Method)
// 	}
// 	return matcher(action.Disclose(m.Attribute), m.Pattern)
// }
//
// func (a *Authority) matchRoleToAction(r *Role, action Action) (err error) {
// 	for _, match := range r.Deny {
// 		if err = a.matchAction(&match, action); err != nil {
// 			return fmt.Errorf("denied explicitly: %w", err)
// 		}
// 	}
// 	for _, match := range r.Allow {
// 		if err = a.matchAction(&match, action); err == nil {
// 			return nil
// 		}
// 	}
// 	return errors.New("denined by default")
// }
//
// func (a *Authority) IsAllowed(ctx context.Context, role string, action Action) error {
// 	errc := make(chan error, 1) // so does not block on cancel
// 	a.operations <- func(m map[string]*Role) {
// 		// first check cache here?
// 		// using maps can just overfill the cache?
// 		// need something more solid?
// 		r, ok := m[role]
// 		if !ok {
// 			r, err := a.repository.Load(ctx, role)
// 			if err != nil {
// 				errc <- err
// 				return
// 			}
// 			m[role] = &r
// 		}
// 		go func() { // match in parallel
// 			errc <- a.matchRoleToAction(r, action)
// 		}()
// 	}
// 	select {
// 	case <-ctx.Done():
// 		return ctx.Err()
// 	case err := <-errc:
// 		return err
// 	}
// }

func (a *Authority) Expire(role string) {
	a.operations <- func(m map[string]*Role) {
		delete(m, role)
	}
}
