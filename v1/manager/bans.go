package manager

import (
	"context"

	"github.com/dkotik/oakacs/v1"
	"github.com/rs/xid"
)

type BanRepository interface {
	CreateBan(ctx context.Context, *oakacs.Ban) error
	RetrieveBan(ctx context.Context, uuid xid.ID) (*oakacs.Ban, error)
	UpdateBan(ctx context.Context, *oakcs.Ban) error
	DeleteBan(ctx context.Context, uuid xid.ID) error
}
