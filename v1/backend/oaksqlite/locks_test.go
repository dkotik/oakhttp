package oaksqlite

import (
	"context"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/xid"
)

func TestLocks(t *testing.T) {
	ctx := context.Background()
	driver := &locks{}
	if err := driver.setup("locks", db); err != nil {
		t.Fatal(err)
	}

	ids := []xid.ID{xid.New(), xid.New(), xid.New(), xid.New(), xid.New()}

	err := driver.Lock(ctx, ids...)
	if err != nil {
		t.Fatal(err)
	}

	err = driver.Lock(ctx, ids[1])
	if err == nil {
		t.Fatal("driver should have errored out on UNIQUE constraint")
	}

	err = driver.Unlock(ctx, ids...)
	if err != nil {
		t.Fatal(err)
	}

	n, err := driver.Clean(ctx, time.Now())
	if err != nil {
		t.Fatal(err, n)
	}
}
