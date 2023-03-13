package oakhttp

import "context"

type Cache interface {
	Set(ctx context.Context, key, value string) (err error)
	Get(ctx context.Context, key string) (value string, err error)
}
