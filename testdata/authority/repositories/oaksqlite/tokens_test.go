package oaksqlite

import (
	_ "github.com/mattn/go-sqlite3"
)

// func TestTokens(t *testing.T) {
// 	ctx := context.Background()
// 	repo, err := NewTokenRepository("tokens", db)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	driver := repo.(*tokens)
//
// 	key, err := driver.Create(ctx, "payload")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	v, err := driver.RetrieveAndDelete(ctx, key)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	n, err := driver.Clean(ctx, time.Now())
// 	if err != nil {
// 		t.Fatal(err, n)
// 	}
//
// 	if v != "payload" {
// 		t.Fatalf("recovered token does not match expected value 'payload': %s", v)
// 	}
// }
