package oakacs

import (
	"context"
	"time"

	"github.com/dkotik/oakacs/v1/oakquery"
	"github.com/rs/xid"
)

// EphemeralRepository keeps short-lived expiring objects like Sessions.
type EphemeralRepository struct {
	Sessions SessionRepository
	Locks    IntegrityLockRepository
	Tokens   TokenRepository
	// TODO: should this be an abstraction over key/value pair storage for serialized push/pull?
	// Push(context.Context, string, time.Time, interface{}) error
	// Pull(context.Context, string) (interface{}, error)
}

// PersistentRepository keep permanent objects like Identities.
type PersistentRepository struct {
	Identities IdentityRepository
	Groups     GroupRepository
	Roles      RoleRepository
	Secrets    SecretRepository
	Bans       BanRepository
}

// IntegrityLock preserves data integrity by making sure relevant resources do not disappear. For example, an Identity cannot be added to a Group, if that Group has been removed right away. The lock helps prevent such conditions.
type IntegrityLockRepository interface {
	Lock(context.Context, time.Duration, ...xid.ID) error // requires unique constraint on the table
	Unlock(context.Context, ...xid.ID) error
	// Purge(context.Context) error // clean up implementation is up to driver
}

// IdentityRepository persists identities.
type IdentityRepository interface {
	Create(context.Context, *Identity) error
	Retrieve(context.Context, xid.ID) (*Identity, error)
	Update(context.Context, xid.ID, func(*Identity) error) error
	Delete(context.Context, xid.ID) error

	Query(context.Context, *oakquery.Query) ([]Identity, error)
}

// SessionRepository persists Sessions.
type SessionRepository interface {
	Create(context.Context, *Session) error
	Retrieve(context.Context, xid.ID) (*Session, error)
	// Only role, last retrieved, and values can actually change.
	Update(context.Context, xid.ID, func(*Session) error) error
	// UpdateRole(context.Context, xid.ID, xid.ID) error
	// UpdateValues(context.Context, xid.ID, map[string]interface{}) error
	Delete(context.Context, xid.ID) error
}

// GroupRepository persists groups.
type GroupRepository interface {
	Create(context.Context, *Group) error
	Retrieve(context.Context, xid.ID) (*Group, error)
	Update(context.Context, xid.ID, func(*Group) error) error
	Delete(context.Context, xid.ID) error

	Query(context.Context, *oakquery.Query) ([]Group, error)
	ListMembers(context.Context, *oakquery.Query) ([]Identity, error)
}

// RoleRepository persists the roles.
type RoleRepository interface {
	Create(context.Context, *Role) (xid.ID, error)
	Retrieve(context.Context, xid.ID) (*Role, error)
	Update(context.Context, xid.ID, func(*Role) error) error
	Delete(context.Context, xid.ID) error
	Query(context.Context, *oakquery.Query) ([]Role, error)
}

// SecretRepository persists secrets.
type SecretRepository interface {
	Create(context.Context, *Secret) error
	Retrieve(context.Context, xid.ID) (*Secret, error)
	Update(context.Context, xid.ID, func(*Secret) error) error
	Delete(context.Context, xid.ID) error

	Query(context.Context, *oakquery.Query) ([]Secret, error)
}

type BanRepository interface {
	Create(context.Context, *Ban) error
	Retrieve(context.Context, xid.ID) (*Ban, error)
	// UpdateBan(context.Context, xid.ID, func(*Ban) error) error
	Delete(context.Context, xid.ID) error

	Query(context.Context, *oakquery.Query) ([]Ban, error)
}

type TokenRepository interface {
	Create(ctx context.Context, value string) (key string, err error)
	RetrieveAndDelete(ctx context.Context, key string) (value string, err error)
}
