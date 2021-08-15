package oakmanager

import (
	"context"

	"github.com/dkotik/oakacs/v1"
	"github.com/dkotik/oakacs/v1/oakquery"
	"github.com/rs/xid"
)

// IdentityRepository persists identities.
type IdentityRepository interface {
	CreateIdentity(context.Context, *oakacs.Identity) error
	RetrieveIdentity(context.Context, xid.ID) (*oakacs.Identity, error)
	UpdateIdentity(context.Context, xid.ID, func(*oakacs.Identity) error) error
	DeleteIdentity(context.Context, xid.ID) error

	ListIdentities(context.Context, *oakquery.Query) ([]oakacs.Identity, error)
}
