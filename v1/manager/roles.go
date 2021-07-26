package oakmanager

import (
	"context"

	"github.com/dkotik/oakacs/v1"
	"github.com/rs/xid"
)

// RoleRepository persists the roles.
type RoleRepository interface {
	CreateRole(ctx context.Context, name string) (*oakacs.Role, error)
	RetrieveRole(ctx context.Context, uuid xid.ID) (*oakacs.Role, error)
	UpdateRole(ctx context.Context, uuid xid.ID, update func(*oakacs.Role) error) error
	DeleteRole(ctx context.Context, uuid xid.ID) error
}
