package oakmanager

import (
	"context"
	"fmt"

	"github.com/dkotik/oakacs/v1"
	"github.com/rs/xid"
)

const groupResource = "group"

// CreateGroup creates a new group.
func (m *Manager) CreateGroup(ctx context.Context, name string) (*oakacs.Group, error) {
	if err = m.acs.Authorize(ctx, ACSService, DomainUniversal, groupResource, "create"); err != nil {
		return
	}
	return m.repo.CreateGroup(ctx, name)
}

// RetrieveGroup fetches a group.
func (m *Manager) RetrieveGroup(ctx context.Context, uuid xid.ID) (*oakacs.Group, error) {
	if err = m.acs.Authorize(ctx, ACSService, DomainUniversal, groupResource, "retrieve"); err != nil {
		return
	}
	return m.repo.RetrieveGroup(ctx, uuid)
}

// DeleteGroup removes the group from the backend.
func (m *Manager) DeleteGroup(ctx context.Context, uuid xid.ID) (err error) {
	if err = m.acs.Authorize(ctx, ACSService, DomainUniversal, groupResource, "delete"); err != nil {
		return
	}
	members, err := m.repo.ListGroupMembers(ctx)
	if err != nil {
		return
	}
	if l := len(members); l > 0 {
		return fmt.Errorf("cannot delete a group because it has %d members", l)
	}
	return m.repo.DeleteGroup(ctx, uuid)
}
