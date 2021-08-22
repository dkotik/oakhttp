package oaksqlite

import (
	"context"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestTokens(t *testing.T) {
	ctx := context.Background()
	driver := &tokens{}
	if err := driver.setup("tokens", db); err != nil {
		t.Fatal(err)
	}

	key, err := driver.CreateToken(ctx, "payload")
	if err != nil {
		t.Fatal(err)
	}

	v, err := driver.RetrieveAndDeleteToken(ctx, key)
	if err != nil {
		t.Fatal(err)
	}

	n, err := driver.Clean(ctx, time.Now())
	if err != nil {
		t.Fatal(err, n)
	}

	if v != "payload" {
		t.Fatalf("recovered token does not match expected value 'payload': %s", v)
	}
}
