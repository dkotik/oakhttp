package client

import "testing"

func TestClientCreation(t *testing.T) {
	_, err := New()
	if err != nil {
		t.Fatal("cannot create HTTP client with default settings:", err)
	}
}
