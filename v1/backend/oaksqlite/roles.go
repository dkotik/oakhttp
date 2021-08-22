package oaksqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dkotik/oakacs/v1"
	"github.com/rs/xid"
)

var _ oakacs.RoleRepository = (*roles)(nil)

type roles struct {
	create   *sql.Stmt
	retrieve *sql.Stmt
	update   *sql.Stmt
	delete   *sql.Stmt
}

func (r *roles) setup(table string, db *sql.DB) (err error) {
	if _, err = db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (uuid BLOB, role BLOB)", table)); err != nil {
		return
	}
	if r.create, err = db.Prepare(fmt.Sprintf("INSERT INTO `%s` VALUES(?,?)", table)); err != nil {
		return
	}
	if r.retrieve, err = db.Prepare(fmt.Sprintf("SELECT * FROM `%s` WHERE uuid=?", table)); err != nil {
		return
	}
	if r.update, err = db.Prepare(fmt.Sprintf("UPDATE `%s` SET role=? WHERE uuid=?", table)); err != nil {
		return
	}
	if r.delete, err = db.Prepare(fmt.Sprintf("DELETE FROM `%s` WHERE uuid=?", table)); err != nil {
		return
	}
	return nil
}

func (r *roles) CreateRole(ctx context.Context, role *oakacs.Role) (xid.ID, error) {
	serialized, err := json.Marshal(role)
	if err != nil {
		return xid.NilID(), err
	}
	id := xid.New()
	if _, err = r.create.ExecContext(ctx, id, serialized); err != nil {
		return xid.NilID(), err
	}
	return id, nil
}

func (r *roles) RetrieveRole(ctx context.Context, id xid.ID) (*oakacs.Role, error) {
	row := r.retrieve.QueryRowContext(ctx, id)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	var data []byte
	if err = row.Scan(&data); err != nil {
		return nil, err
	}
	role := &oakacs.Role{}
	if err = json.Unmarshal(data, role); err != nil {
		return nil, err
	}
	return role, nil
}

func (r *roles) UpdateRole(ctx context.Context, id xid.ID, update func(*oakacs.Role) error) error {
	role, err := r.RetrieveRole(ctx, id)
	if err != nil {
		return err
	}
	if err = update(role); err != nil {
		return err
	}
	serialized, err := json.Marshal(role)
	if err != nil {
		return err
	}
	result, err := r.update.ExecContext(ctx, serialized, id)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("roles was not updated, because it was missing")
	}
	return nil
}

func (r *roles) DeleteRole(ctx context.Context, id xid.ID) (err error) {
	if _, err = r.delete.ExecContext(ctx, id); err != nil {
		return err
	}
	return err
}
