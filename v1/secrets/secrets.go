package secrets

import (
	"errors"
)

var (
	ErrSecretNotFound = errors.New("requested secret value does not exist")
)

type Repository interface {
	SyncMap(name string) func() (map[string]interface{}, error)
	SyncString(name string) func() (string, error)
	SyncUint(name string) func() (uint, error)
	SyncInt(name string) func() (int, error)
}
