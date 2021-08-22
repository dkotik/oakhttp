package oakacs

import (
	"context"

	"github.com/dkotik/oakacs/v1/oakquery"
	"github.com/rs/xid"
)

// EphemeralRepository keeps short-lived expiring objects like Sessions.
type EphemeralRepository interface {
	// TODO: should this be an abstraction over key/value pair storage for serialized push/pull?
	// Push(context.Context, string, time.Time, interface{}) error
	// Pull(context.Context, string) (interface{}, error)
	SessionRepository
	IntegrityLockRepository
	TokenRepository
}

// PersistentRepository keep permanent objects like Identities.
type PersistentRepository interface {
	IdentityRepository
	GroupRepository
	RoleRepository
	SecretRepository
	BanRepository
}

// IntegrityLock preserves data integrity by making sure relevant resources do not disappear. For example, an Identity cannot be added to a Group, if that Group has been removed right away. The lock helps prevent such conditions.
type IntegrityLockRepository interface {
	Lock(context.Context, ...xid.ID) error // requires unique constraint on the table
	Unlock(context.Context, ...xid.ID) error
	// PurgeLocks(context.Context) error // clean up implementation is up to driver
}

// IdentityRepository persists identities.
type IdentityRepository interface {
	CreateIdentity(context.Context, *Identity) error
	RetrieveIdentity(context.Context, xid.ID) (*Identity, error)
	UpdateIdentity(context.Context, xid.ID, func(*Identity) error) error
	DeleteIdentity(context.Context, xid.ID) error

	ListIdentities(context.Context, *oakquery.Query) ([]Identity, error)
}

// SessionRepository persists Sessions.
type SessionRepository interface {
	CreateSession(context.Context, *Session) error
	RetrieveSession(context.Context, xid.ID) (*Session, error)
	// Only role, last retrieved, and values can actually change.
	UpdateSession(context.Context, xid.ID, func(*Session) error) error
	// UpdateSessionRole(context.Context, xid.ID, xid.ID) error
	// UpdateSessionValues(context.Context, xid.ID, map[string]interface{}) error
	DeleteSession(context.Context, xid.ID) error
}

// GroupRepository persists groups.
type GroupRepository interface {
	CreateGroup(context.Context, *Group) error
	RetrieveGroup(context.Context, xid.ID) (*Group, error)
	UpdateGroup(context.Context, xid.ID, func(*Group) error) error
	DeleteGroup(context.Context, xid.ID) error

	ListGroups(context.Context, *oakquery.Query) ([]Group, error)
	ListGroupMembers(context.Context, *oakquery.Query) ([]Identity, error)
}

// RoleRepository persists the roles.
type RoleRepository interface {
	CreateRole(context.Context, *Role) (xid.ID, error)
	RetrieveRole(context.Context, xid.ID) (*Role, error)
	UpdateRole(context.Context, xid.ID, func(*Role) error) error
	DeleteRole(context.Context, xid.ID) error

	// QueryRoles(context.Context, *oakquery.Query) ([]Role, error)
}

// SecretRepository persists secrets.
type SecretRepository interface {
	CreateSecret(context.Context, *Secret) error
	RetrieveSecret(context.Context, xid.ID) (*Secret, error)
	UpdateSecret(context.Context, xid.ID, func(*Secret) error) error
	DeleteSecret(context.Context, xid.ID) error

	ListSecrets(context.Context, *oakquery.Query) ([]Secret, error)
}

type BanRepository interface {
	CreateBan(context.Context, *Ban) error
	RetrieveBan(context.Context, xid.ID) (*Ban, error)
	// UpdateBan(context.Context, xid.ID, func(*Ban) error) error
	DeleteBan(context.Context, xid.ID) error

	QueryBans(context.Context, *oakquery.Query) ([]Ban, error)
}

type TokenRepository interface {
	CreateToken(ctx context.Context, value string) (key string, err error)
	RetrieveAndDeleteToken(ctx context.Context, key string) (value string, err error)
}
