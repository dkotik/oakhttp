package oakmanager

import (
	"context"

	"github.com/dkotik/oakacs/v1"
	"github.com/rs/xid"
)

const secretResource = "secrets"

// TODO: secrets are not supposed to be managed?

// CreateSecret creates a new Secret.
func (m *Manager) CreateSecret(ctx context.Context, secret *oakacs.Secret) error {
	if err := m.acs.Authorize(ctx, ACSService, DomainUniversal, secretResource, WR); err != nil {
		return err
	}
	return m.repo.CreateSecret(ctx, name)
}

// RetrieveSecret fetches a Secret.
func (m *Manager) RetrieveSecret(ctx context.Context, uuid xid.ID) (*oakacs.Secret, error) {
	if err := m.acs.Authorize(ctx, ACSService, DomainUniversal, secretResource, RD); err != nil {
		return nil, err
	}
	return m.repo.RetrieveSecret(ctx, uuid)
}

// DeleteSecret removes the Secret from the backend.
func (m *Manager) DeleteSecret(ctx context.Context, uuid xid.ID) error {
	if err := m.acs.Authorize(ctx, ACSService, DomainUniversal, secretResource, WR); err != nil {
		return err
	}
	return m.repo.DeleteSecret(ctx, uuid)
}
