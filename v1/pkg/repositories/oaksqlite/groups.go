package oaksqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dkotik/oakacs/v1"
	"github.com/rs/xid"
)

func NewGroupRepository(table string, db *sql.DB) (oakacs.GroupRepository, error) {
	// UUID            xid.ID
	// Name            string
	// DefaultRole     Role
	// AscendableRoles []Role

	dr := &groups{}
	var err error
	if _, err = db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (uuid BLOB, name TEXT, default_role BLOB)", table)); err != nil {
		return nil, err
	}
	if _, err = db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s_xidentities` (iuuid BLOB, guuid BLOB)", table)); err != nil {
		return nil, err
	}
	if dr.create, err = db.Prepare(fmt.Sprintf("INSERT INTO `%s` VALUES(?,?,?)", table)); err != nil {
		return nil, err
	}
	if dr.retrieve, err = db.Prepare(fmt.Sprintf("SELECT * FROM `%s` WHERE uuid=?", table)); err != nil {
		return nil, err
	}
	if dr.update, err = db.Prepare(fmt.Sprintf("UPDATE `%s` SET role=? WHERE uuid=?", table)); err != nil {
		return nil, err
	}
	if dr.delete, err = db.Prepare(fmt.Sprintf("DELETE FROM `%s` WHERE uuid=?", table)); err != nil {
		return nil, err
	}
	if dr.query, err = db.Prepare(fmt.Sprintf("SELECT count(role), role FROM `%s` OFFSET ? LIMIT ?", table)); err != nil {
		return nil, err
	}
	return dr, nil
}

type groups struct {
	create   *sql.Stmt
	retrieve *sql.Stmt
	update   *sql.Stmt
	delete   *sql.Stmt
	query    *sql.Stmt
}

func (dr *groups) Create(ctx context.Context, group *oakacs.Group) error {
	if _, err := dr.create.ExecContext(ctx, group.UUID, group.Name, group.DefaultRole.UUID); err != nil {
		return err
	}
	return nil
}

func (dr *groups) Retrieve(ctx context.Context, id xid.ID) (*oakacs.Group, error) {
	row := dr.retrieve.QueryRowContext(ctx, id)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	group := &oakacs.Group{ // TODO: fill out
		UUID: id,
	}
	if err = row.Scan(&group.Name, &group.DefaultRole); err != nil {
		return nil, err
	}
	return group, nil
}

func (dr *groups) Update(ctx context.Context, id xid.ID, update func(*oakacs.Group) error) error {
	group, err := dr.Retrieve(ctx, id)
	if err != nil {
		return err
	}
	if err = update(group); err != nil {
		return err
	}
	result, err := dr.update.ExecContext(ctx, group) // TODO: fill out fields
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("group was not updated, because it was missing")
	}
	return nil
}

func (dr *groups) Delete(ctx context.Context, id xid.ID) (err error) {
	if _, err = dr.delete.ExecContext(ctx, id); err != nil {
		return err
	}
	return err
}
