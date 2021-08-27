package oaksqlite

import (
	"database/sql"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxIdleTime(0)
	db.SetConnMaxLifetime(0)
}
