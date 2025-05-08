package oaksqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rs/xid"
)

type rolesToGroups struct {
	// selectRoles         *sql.Stmt // do not mix concerns
	// selectGroups        *sql.Stmt // do not mix concerns
	table               string
	addRoleToGroup      *sql.Stmt
	removeRoleFromGroup *sql.Stmt
	cleanupRole         *sql.Stmt
	cleanupGroup        *sql.Stmt
}

func newRolestoGroups(db *sql.DB, table string) (*rolesToGroups, error) {
	var (
		err error
		dr  = &rolesToGroups{
			table: table,
		}
	)
	if _, err = db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (ruuid BLOB, guuid BLOB)", table)); err != nil {
		return nil, err
	}
	// if _, err = db.Exec(fmt.Sprintf("CREATE TRIGGER preventRolesToGroupsDuplicates BEFORE INSERT ON `%s` BEGIN DELETE FROM `%s` WHERE ruuid=NEW.ruuid AND guuid=NEW.guuid END", table, table)); err != nil {
	// 	return nil, err
	// }
	if dr.addRoleToGroup, err = db.Prepare(fmt.Sprintf("INSERT INTO `%s` VALUES(?, ?)", table)); err != nil {
		return nil, err
	}
	if dr.removeRoleFromGroup, err = db.Prepare(fmt.Sprintf("DELETE FROM `%s` WHERE ruuid=? AND guuid=?", table)); err != nil {
		return nil, err
	}
	if dr.cleanupRole, err = db.Prepare(fmt.Sprintf("DELETE FROM `%s` WHERE ruuid=?", table)); err != nil {
		return nil, err
	}
	if dr.cleanupGroup, err = db.Prepare(fmt.Sprintf("DELETE FROM `%s` WHERE guuid=?", table)); err != nil {
		return nil, err
	}
	return dr, nil
}

func (dr *rolesToGroups) AddRoleToGroup(ctx context.Context, role, group xid.ID) error {
	_, err := dr.addRoleToGroup.ExecContext(ctx, role, group)
	return err
}

func (dr *rolesToGroups) RemoveRoleFromGroup(ctx context.Context, role, group xid.ID) error {
	_, err := dr.removeRoleFromGroup.ExecContext(ctx, role, group)
	return err
}

func (dr *rolesToGroups) CleanupRole(ctx context.Context, role xid.ID) error {
	_, err := dr.cleanupRole.ExecContext(ctx, role)
	return err
}

func (dr *rolesToGroups) CleanupGroup(ctx context.Context, group xid.ID) error {
	_, err := dr.cleanupGroup.ExecContext(ctx, group)
	return err
}
