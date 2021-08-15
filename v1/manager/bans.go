package manager

import (
	"context"

	"github.com/dkotik/oakacs/v1"
	"github.com/dkotik/oakacs/v1/oakquery"
	"github.com/rs/xid"
)

type BanRepository interface {
	CreateBan(context.Context, *oakacs.Ban) error
	RetrieveBan(context.Context, xid.ID) (*oakacs.Ban, error)
	UpdateBan(context.Context, xid.ID, func(*oakcs.Ban) error) error
	DeleteBan(context.Context, xid.ID) error

	ListBans(context.Context, *oakquery.Query) ([]oakacs.Ban, error)
}
