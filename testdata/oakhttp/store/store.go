package store

import (
	"context"
	"errors"
	"time"
)

var (
	ErrFull          = errors.New("there are too many store records")
	ErrValueNotFound = errors.New("value not found in store")
)

type Update func([]byte) ([]byte, error)

type KeyValue interface {
	Get(ctx context.Context, key []byte) ([]byte, error)
	Set(ctx context.Context, key, value []byte) error
	Update(ctx context.Context, key []byte, update Update) error
	Delete(ctx context.Context, key []byte) error
}

type KeyKeyValue interface {
	Get(ctx context.Context, key1, key2 []byte) ([]byte, error)
	Set(ctx context.Context, key1, key2, value []byte) error
	Update(ctx context.Context, key1, key2 []byte, update Update) error
	Delete(ctx context.Context, key1, key2 []byte) error
}

type expiringValue struct {
	Data    []byte
	Expires time.Time
}
