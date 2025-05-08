package store

import "context"

type KeyString interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
	Delete(ctx context.Context, key string) error
}

type KeyKeyString interface {
	Get(ctx context.Context, key1, key2 string) (string, error)
	Set(ctx context.Context, key1, key2 string, value string) error
	Delete(ctx context.Context, key1, key2 string) error
}

type keyString struct {
	kv KeyValue
}

func NewKeyString(kv KeyValue) KeyString {
	if kv == nil {
		panic("cannot use a <nil> key value store")
	}
	return &keyString{kv: kv}
}

func (k *keyString) Get(ctx context.Context, key string) (string, error) {
	value, err := k.kv.Get(ctx, []byte(key))
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func (k *keyString) Set(ctx context.Context, key, value string) error {
	return k.kv.Set(ctx, []byte(key), []byte(value))
}

func (k *keyString) Delete(ctx context.Context, key string) error {
	return k.kv.Delete(ctx, []byte(key))
}

type keyKeyString struct {
	kkv KeyKeyValue
}

func NewKeyKeyString(kkv KeyKeyValue) KeyKeyString {
	if kkv == nil {
		panic("cannot use a <nil> key value store")
	}
	return &keyKeyString{kkv: kkv}
}

func (k *keyKeyString) Get(ctx context.Context, key1, key2 string) (string, error) {
	value, err := k.kkv.Get(ctx, []byte(key1), []byte(key2))
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func (k *keyKeyString) Set(ctx context.Context, key1, key2, value string) error {
	return k.kkv.Set(ctx, []byte(key1), []byte(key2), []byte(value))
}

func (k *keyKeyString) Delete(ctx context.Context, key1, key2 string) error {
	return k.kkv.Delete(ctx, []byte(key1), []byte(key2))
}
