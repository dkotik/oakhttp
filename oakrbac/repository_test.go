package oakrbac

import (
	"testing"
)

func testRepository(t *testing.T, repo RoleRepository) {
	if err := repo.AddRole(&basicRole{name: "test"}); err != nil {
		t.Fatal("cannot add a role", err)
	}
	recovered, err := repo.GetRole("test")
	if err != nil {
		t.Fatal("cannot recover role")
	}
	t.Log("recovered role", recovered.Name())
}

func TestListRepository(t *testing.T) {
	testRepository(t, &ListRepository{})
}

func TestMapRepository(t *testing.T) {
	testRepository(t, &MapRepository{})
}
