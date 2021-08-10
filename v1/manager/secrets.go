package oakmanager

package oakmanager

import (
	"context"
	"fmt"

	"github.com/dkotik/oakacs/v1"
	"github.com/rs/xid"
)

const secretResource = "Secret"

// TODO: secrets are not supposed to be managed?

// SecretRepository persists secrets.
type SecretRepository interface {
	CreateSecret(context.Context, *oakacs.Secret) error
	RetrieveSecret(context.Context, xid.ID) (*oakacs.Secret, error)
	UpdateSecret(context.Context, *oakacs.Secret) error
	DeleteSecret(context.Context, xid.ID) error
}

// CreateSecret creates a new Secret.
func (m *Manager) CreateSecret(ctx context.Context, name string) (*oakacs.Secret, error) {
	if err = m.acs.Authorize(ctx, ACSService, DomainUniversal, secretResource, "create"); err != nil {
		return
	}
	return m.repo.CreateSecret(ctx, name)
}

// RetrieveSecret fetches a Secret.
func (m *Manager) RetrieveSecret(ctx context.Context, uuid xid.ID) (*oakacs.Secret, error) {
	if err = m.acs.Authorize(ctx, ACSService, DomainUniversal, secretResource, "retrieve"); err != nil {
		return
	}
	return m.repo.RetrieveSecret(ctx, uuid)
}

// DeleteSecret removes the Secret from the backend.
func (m *Manager) DeleteSecret(ctx context.Context, uuid xid.ID) (err error) {
	if err = m.acs.Authorize(ctx, ACSService, DomainUniversal, secretResource, "delete"); err != nil {
		return
	}
	return m.repo.DeleteSecret(ctx, uuid)
}
