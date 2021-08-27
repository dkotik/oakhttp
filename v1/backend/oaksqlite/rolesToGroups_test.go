package oaksqlite

import (
	"context"
	"testing"

	"github.com/rs/xid"
)

func TestRolesToGroups(t *testing.T) {
	driver, err := newRolestoGroups(db, "roles_to_groups")
	if err != nil {
		// panic(err)
		t.Fatal(err)
	}
	ctx := context.Background()

	roleid, groupid := xid.New(), xid.New()
	if err = driver.AddRoleToGroup(ctx, roleid, groupid); err != nil {
		t.Fatal(err)
	}
	if err = driver.RemoveRoleFromGroup(ctx, roleid, groupid); err != nil {
		t.Fatal(err)
	}
	if err = driver.CleanupRole(ctx, roleid); err != nil {
		t.Fatal(err)
	}
	if err = driver.CleanupGroup(ctx, groupid); err != nil {
		t.Fatal(err)
	}
}
