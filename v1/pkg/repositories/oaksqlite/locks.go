package oaksqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/dkotik/oakacs/v1"
	"github.com/rs/xid"
)

func NewIntegrityLockRepository(table string, db *sql.DB) (oakacs.IntegrityLockRepository, error) {
	l := &locks{db: db}
	var err error
	if _, err = db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (id BLOB UNIQUE, deadline INTEGER)", table)); err != nil {
		return nil, err
	}
	if l.create, err = db.Prepare(fmt.Sprintf("INSERT INTO `%s` VALUES(?,?)", table)); err != nil {
		return nil, err
	}
	if l.delete, err = db.Prepare(fmt.Sprintf("DELETE FROM `%s` WHERE id=?", table)); err != nil {
		return nil, err
	}
	if l.clean, err = db.Prepare(fmt.Sprintf("DELETE FROM `%s` WHERE deadline<?", table)); err != nil {
		return nil, err
	}
	return l, nil
}

type locks struct {
	db     *sql.DB
	create *sql.Stmt
	delete *sql.Stmt
	clean  *sql.Stmt
}

func (l *locks) Lock(ctx context.Context, timeout time.Duration, ids ...xid.ID) error {
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, id := range ids {
		create := tx.StmtContext(ctx, l.create)
		if _, err := create.ExecContext(ctx, id, timeout); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (l *locks) Unlock(ctx context.Context, ids ...xid.ID) error {
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, id := range ids {
		delete := tx.StmtContext(ctx, l.delete)
		if _, err := delete.ExecContext(ctx, id); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (l *locks) Clean(ctx context.Context, deadline time.Time) (int64, error) {
	result, err := l.clean.ExecContext(ctx, deadline)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
