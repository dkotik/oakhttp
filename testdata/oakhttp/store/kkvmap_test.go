package store

import "testing"

func TestKeyKeyValueMap(t *testing.T) {
	kkv, err := NewMapKeyKeyValue()
	if err != nil {
		t.Fatal(err)
	}
	NewKeyKeyValueTest(kkv)(t)
}
