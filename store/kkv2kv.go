package store

import "context"

type keyKeyValueToKeyValueAdaptor struct {
	key1 []byte
	kkv  KeyKeyValue
}

func NewKeyKeyValueToKeyValueAdaptor(kkv KeyKeyValue, key1 []byte) KeyValue {
	if kkv == nil {
		panic("canot use a <nil> key-key-value store")
	}
	return &keyKeyValueToKeyValueAdaptor{
		key1: key1,
		kkv:  kkv,
	}
}

func (k *keyKeyValueToKeyValueAdaptor) Get(
	ctx context.Context,
	key2 []byte,
) ([]byte, error) {
	return k.kkv.Get(ctx, k.key1, key2)
}

func (k *keyKeyValueToKeyValueAdaptor) Set(
	ctx context.Context,
	key2, value []byte,
) error {
	return k.kkv.Set(ctx, k.key1, key2, value)
}

func (k *keyKeyValueToKeyValueAdaptor) Update(
	ctx context.Context,
	key2 []byte, update Update,
) error {
	return k.kkv.Update(ctx, k.key1, key2, update)
}

func (k *keyKeyValueToKeyValueAdaptor) Delete(
	ctx context.Context,
	key2 []byte,
) error {
	return k.kkv.Delete(ctx, k.key1, key2)
}
