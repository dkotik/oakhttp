package cuestore

import "testing"

func TestScope(t *testing.T) {
	if err := scope.Err(); err != nil {
		t.Fatal("could not parse Cuelang scope for Cuestore:", err.Error())
	}
	if scope.Null() == nil {
		t.Fatal("cannot use Cuestore with empty scope")
	}
}
