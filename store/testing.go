package store

import (
	"bytes"
	"context"
	"testing"
	"time"
)

func NewKeyKeyValueTest(kkv KeyKeyValue) func(t *testing.T) {
	return NewKeyValueTest(NewKeyKeyValueToKeyValueAdaptor(
		kkv,
		[]byte("key1"),
	))
}

func NewKeyValueTest(kv KeyValue) func(t *testing.T) {
	return func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		t.Cleanup(cancel)

		key := []byte("testKey")
		value := []byte("testValue")
		t.Helper()

		t.Run("set", func(t *testing.T) {
			if err := kv.Set(ctx, key, value); err != nil {
				t.Fatal(err)
			}
		})

		t.Run("get", func(t *testing.T) {
			retrieved, err := kv.Get(ctx, key)
			if err != nil {
				t.Fatal(key, err)
			}
			if bytes.Compare(retrieved, value) != 0 {
				t.Fatalf("values %q and %q do not match", string(retrieved), string(value))
			}
		})

		t.Run("missing", func(t *testing.T) {
			_, err := kv.Get(ctx, []byte("missingKey"))
			if err != ErrValueNotFound {
				t.Fatal("missingKey did not return ErrValueNotFound")
			}
		})

		t.Run("update", func(t *testing.T) {
			err := kv.Update(ctx, key, func(v []byte) ([]byte, error) {
				return append(v, byte('1')), nil
			})
			if err != nil {
				t.Fatal(key, err)
			}
			retrieved, err := kv.Get(ctx, key)
			if err != nil {
				t.Fatal(key)
			}
			if string(retrieved) != "testValue1" {
				t.Fatal("updated value does not match:", string(retrieved), "testValue1")
			}
		})

		t.Run("delete", func(t *testing.T) {
			err := kv.Delete(ctx, key)
			if err != nil {
				t.Fatal(key, err)
			}
			_, err = kv.Get(ctx, key)
			if err == nil {
				t.Fatal("value recovered for deleted key:", key)
			}
		})
	}
}
