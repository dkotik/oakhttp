package oakrbac

import (
	"context"
	"testing"
)

func TestContextOperations(t *testing.T) {
	r := Must(New(
		WithNewRole("administrator", AllowEverything),
	))

	ctx := r.ContextWithNegotiatedRole(context.Background(), "manager", "administrator")

	_, err := Authorize(ctx, &Intent{})
	if err != nil {
		t.Error(err)
	}
	// if role != "administrator" {
	// 	t.Fatalf("role %q name did not match", role)
	// }
	// if strings.HasSuffix(policy.Name(), "/AllowEverything") {
	// 	t.Fatal("policy name did not match")
	// }
}
