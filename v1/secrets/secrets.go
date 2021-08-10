package secrets

import (
	"context"
	"errors"
)

var (
	ErrSecretNotFound = errors.New("requested secret value does not exist")
)

type Repository interface {
	SyncMap(name string) func(context.Context) (map[string]interface{}, error)
	SyncString(name string) func(context.Context) (string, error)
	SyncUint(name string) func(context.Context) (uint, error)
	SyncInt(name string) func(context.Context) (int, error)
}
