package store

import "testing"

func TestKeyValueMap(t *testing.T) {
	kv, err := NewMapKeyValue()
	if err != nil {
		t.Fatal(err)
	}
	NewKeyValueTest(kv)(t)
}
