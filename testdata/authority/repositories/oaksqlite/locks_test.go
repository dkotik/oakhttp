package oaksqlite

import (
	_ "github.com/mattn/go-sqlite3"
)

// func TestLocks(t *testing.T) {
// 	ctx := context.Background()
// 	repo, err := NewIntegrityLockRepository("locks", db)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	driver := repo.(*locks)
//
// 	ids := []xid.ID{xid.New(), xid.New(), xid.New(), xid.New(), xid.New()}
//
// 	err = driver.Lock(ctx, time.Minute, ids...)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	err = driver.Lock(ctx, time.Minute, ids[1])
// 	if err == nil {
// 		t.Fatal("driver should have errored out on UNIQUE constraint")
// 	}
//
// 	err = driver.Unlock(ctx, ids...)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	n, err := driver.Clean(ctx, time.Now())
// 	if err != nil {
// 		t.Fatal(err, n)
// 	}
// }
