package oaksqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dkotik/oakacs/v1"
	"github.com/rs/xid"
)

func NewIdentityRepository(table string, db *sql.DB) (oakacs.IdentityRepository, error) {
	// UUID              xid.ID
	// Name              string
	// Groups            []Group // the order matters for default roles
	// Secrets           []Secret
	// HumanityConfirmed time.Time

	dr := &identities{}
	var err error
	if _, err = db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (uuid BLOB, role BLOB)", table)); err != nil {
		return nil, err
	}
	if dr.create, err = db.Prepare(fmt.Sprintf("INSERT INTO `%s`(uuid, name) VALUES(?,?)", table)); err != nil {
		return nil, err
	}
	if dr.retrieve, err = db.Prepare(fmt.Sprintf("SELECT role FROM `%s` WHERE uuid=?", table)); err != nil {
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

type identities struct {
	create   *sql.Stmt
	retrieve *sql.Stmt
	update   *sql.Stmt
	delete   *sql.Stmt
	query    *sql.Stmt
}

func (dr *identities) Create(ctx context.Context, identity *oakacs.Identity) error {
	_, err := dr.create.ExecContext(ctx, identity.UUID, identity.Name)
	return err
}

func (dr *identities) Retrieve(ctx context.Context, id xid.ID) (*oakacs.Identity, error) {
	row := dr.retrieve.QueryRowContext(ctx, id)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	identity := &oakacs.Identity{}
	if err = row.Scan(&identity); err != nil { // TODO: fill out fields
		return nil, err
	}
	return identity, nil
}

func (dr *identities) Update(ctx context.Context, id xid.ID, update func(*oakacs.Identity) error) error {
	identity, err := dr.Retrieve(ctx, id)
	if err != nil {
		return err
	}
	if err = update(identity); err != nil {
		return err
	}
	result, err := dr.update.ExecContext(ctx, identity) // TODO: fill out fields
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("identity was not updated, because it was missing")
	}
	return nil
}

func (dr *identities) Delete(ctx context.Context, id xid.ID) (err error) {
	if _, err = dr.delete.ExecContext(ctx, id); err != nil {
		return err
	}
	return err
}
