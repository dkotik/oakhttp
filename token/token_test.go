package token

import "testing"

func TestToken(t *testing.T) {
	desiredLength := 1024

	factory, err := New(WithTokenLength(desiredLength))
	if err != nil {
		t.Fatal(err)
	}

	token, err := factory()
	if err != nil {
		t.Fatal(err)
	}
	if len(token) != desiredLength {
		t.Fatal("token length mismatch", len(token), desiredLength)
	}
	// panic(token)
}
