package oakacs

import (
	"fmt"
	"time"

	"github.com/rs/xid"
)

// Role holds a set of allowed actions.
type Role struct {
	UUID     xid.ID
	Name     string
	Allow    []Permission
	Deny     []Permission
	Duration time.Duration
}

// Permission represents something an Identity can do.
type Permission struct {
	Service  string // Where? -  Namespace
	Domain   string // Where? -  Realm
	Resource string // With?  -  Subject
	Action   string // What?  -  Verb
}

// Match returns true if all parameters are exactly equal to their corresponding fields or by * wildcard.
func (p Permission) Match(service, domain, resource, action string) bool {
	return (p.Service == service || wildcardPatternMatch(p.Service, service)) &&
		(p.Domain == domain || wildcardPatternMatch(p.Domain, domain)) &&
		(p.Resource == resource || wildcardPatternMatch(p.Resource, resource)) &&
		(p.Action == action || wildcardPatternMatch(p.Action, action))
}

func (p Permission) String() string {
	return fmt.Sprintf("%s::%s::s", p.Domain, p.Action, p.Resource)
}
