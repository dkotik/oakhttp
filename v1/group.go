package oakacs

import (
	"context"
	"fmt"

	"github.com/rs/xid"
)

const groupResource = "group"

// GroupRepository persists groups.
type GroupRepository interface {
	CreateGroup(ctx context.Context, name string) (*Group, error)
	RetrieveGroup(ctx context.Context, uuid xid.ID) (*Group, error)
	UpdateGroup(ctx context.Context, uuid xid.ID, update func(*Group) error) error
	DeleteGroup(ctx context.Context, uuid xid.ID) error

	ListGroupMembers(ctx context.Context) ([]*Identity, error)
}

var tempGR GroupRepository

// Group holds roles that identities may assume.
type Group struct {
	UUID            xid.ID
	Name            string
	DefaultRole     Role
	AscendableRoles []Role
}

// DeleteGroup removes the group from the backend.
func (acs *AccessControlSystem) DeleteGroup(ctx context.Context, uuid xid.ID) (err error) {
	if err = acs.Authorize(ctx, ACSService, DomainUniversal, groupResource, "delete"); err != nil {
		return
	}
	members, err := tempGR.ListGroupMembers(ctx)
	if err != nil {
		return
	}
	if l := len(members); l > 0 {
		return fmt.Errorf("cannot delete a group because it has %d members", l)
	}
	return tempGR.DeleteGroup(ctx, uuid)
}
