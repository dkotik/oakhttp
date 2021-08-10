package oakmanager

import (
	"context"

	"github.com/dkotik/oakacs/v1"
	"github.com/rs/xid"
)

// RoleRepository persists the roles.
type RoleRepository interface {
	CreateRole(context.Context, *oakacs.Role) error
	RetrieveRole(context.Context, xid.ID) (*oakacs.Role, error)
	UpdateRole(ctx context.Context, uuid xid.ID, update func(*oakacs.Role) error) error
	DeleteRole(context.Context, xid.ID) error
}
