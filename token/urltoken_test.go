package token

import "testing"

func TestURLToken(t *testing.T) {
	factory, err := NewURLToken(36)
	if err != nil {
		t.Fatal(err)
	}
	token, err := factory()
	if err != nil {
		t.Fatal(err)
	}
	if len(token) < 36 {
		t.Fatal("generated token too small")
	}
	// t.Fatal(token)
}
