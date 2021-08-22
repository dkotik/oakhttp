package oaksqlite

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"time"

	// "github.com/dkotik/oakacs/v1"

	"github.com/rs/xid"
)

// var _ oakacs.TokenRepository = (*oakacs.TokenRepository)(nil)

type tokens struct {
	create   *sql.Stmt
	retrieve *sql.Stmt
	delete   *sql.Stmt
	clean    *sql.Stmt
}

func (t *tokens) setup(table string, db *sql.DB) (err error) {
	if _, err = db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (key TEXT, value TEXT, deadline INTEGER)", table)); err != nil {
		return
	}
	if t.create, err = db.Prepare(fmt.Sprintf("INSERT INTO `%s` VALUES(?,?,?)", table)); err != nil {
		return
	}
	if t.retrieve, err = db.Prepare(fmt.Sprintf("SELECT value FROM `%s` WHERE key=?", table)); err != nil {
		return
	}
	if t.delete, err = db.Prepare(fmt.Sprintf("DELETE FROM `%s` WHERE key=?", table)); err != nil {
		return
	}
	if t.clean, err = db.Prepare(fmt.Sprintf("DELETE FROM `%s` WHERE deadline<?", table)); err != nil {
		return
	}
	return nil
}

func (t *tokens) CreateToken(ctx context.Context, v string) (string, error) {
	x := make([]byte, 16)
	if n, err := rand.Reader.Read(x); err != nil {
		return "", err
	} else if n < 16 {
		return "", errors.New("not enough random bytes")
	}
	id := fmt.Sprintf("%s-%x-%x", xid.New(), x[:8], x[8:])
	if _, err := t.create.ExecContext(ctx, id, v, time.Now().Unix()); err != nil {
		return "", err
	}
	return id, nil
}

func (t *tokens) RetrieveAndDeleteToken(ctx context.Context, key string) (string, error) {
	row := t.retrieve.QueryRowContext(ctx, key)
	err := row.Err()
	if err != nil {
		return "", err
	}
	var value string
	if err = row.Scan(&value); err != nil {
		return "", err
	}
	if _, err = t.delete.ExecContext(ctx, key); err != nil {
		return "", err
	}
	return value, nil
}

func (t *tokens) Clean(ctx context.Context, deadline time.Time) (int64, error) {
	result, err := t.delete.ExecContext(ctx, deadline)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
