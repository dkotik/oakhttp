package store

import (
	"context"
	"encoding/binary"
)

type UpdateUint64 func(uint64) (uint64, error)

func uint64ToBytes(x uint64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, x)
	return buf[:n]
}

func bytesToUint64(buf []byte) uint64 {
	x, _ := binary.Uvarint(buf)
	return x
}

type KeyUint64 interface {
	Get(ctx context.Context, key []byte) (uint64, error)
	Set(ctx context.Context, key []byte, value uint64) error
	Update(ctx context.Context, key []byte, update UpdateUint64) error
	Delete(ctx context.Context, key []byte) error
}

type KeyKeyUint64 interface {
	Get(ctx context.Context, key1, key2 []byte) (uint64, error)
	Set(ctx context.Context, key1, key2 []byte, value uint64) error
	Update(ctx context.Context, key1, key2 []byte, update UpdateUint64) error
	Delete(ctx context.Context, key1, key2 []byte) error
}

type keyUint64 struct {
	kv KeyValue
}

func NewKeyUint64(kv KeyValue) KeyUint64 {
	if kv == nil {
		panic("cannot use a <nil> key value store")
	}
	return &keyUint64{kv: kv}
}

func (k *keyUint64) Get(ctx context.Context, key []byte) (uint64, error) {
	value, err := k.kv.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	return bytesToUint64(value), nil
}

func (k *keyUint64) Set(ctx context.Context, key []byte, value uint64) error {
	return k.kv.Set(ctx, key, uint64ToBytes(value))
}

func (k *keyUint64) Update(ctx context.Context, key []byte, update UpdateUint64) error {
	return k.kv.Update(ctx, key, func(value []byte) ([]byte, error) {
		updated, err := update(bytesToUint64(value))
		if err != nil {
			return nil, err
		}
		return uint64ToBytes(updated), nil
	})
}

func (k *keyUint64) Delete(ctx context.Context, key []byte) error {
	return k.kv.Delete(ctx, key)
}

type keyKeyUint64 struct {
	kkv KeyKeyValue
}

func NewKeyKeyUint64(kkv KeyKeyValue) KeyKeyUint64 {
	if kkv == nil {
		panic("cannot use a <nil> key value store")
	}
	return &keyKeyUint64{kkv: kkv}
}

func (k *keyKeyUint64) Get(ctx context.Context, key1, key2 []byte) (uint64, error) {
	value, err := k.kkv.Get(ctx, key1, key2)
	if err != nil {
		return 0, err
	}
	return bytesToUint64(value), nil
}

func (k *keyKeyUint64) Set(ctx context.Context, key1, key2 []byte, value uint64) error {
	return k.kkv.Set(ctx, key1, key2, uint64ToBytes(value))
}

func (k *keyKeyUint64) Update(ctx context.Context, key1, key2 []byte, update UpdateUint64) error {
	return k.kkv.Update(ctx, key1, key2, func(value []byte) ([]byte, error) {
		updated, err := update(bytesToUint64(value))
		if err != nil {
			return nil, err
		}
		return uint64ToBytes(updated), nil
	})
}

func (k *keyKeyUint64) Delete(ctx context.Context, key1, key2 []byte) error {
	return k.kkv.Delete(ctx, key1, key2)
}
