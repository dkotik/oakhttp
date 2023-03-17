package oakrbac

import (
	"fmt"
	"strings"
	"testing"
)

func TestPolicyReflection(t *testing.T) {
	expectedName := "github.com/dkotik/oakacs/oakrbac.AllowEverything"
	p := Policy(AllowEverything)

	n := p.Name()
	if n != expectedName {
		t.Fatalf("policy name %q did not match %q", n, expectedName)
	}

	f := p.File()
	expectedSuffix := "oakacs/oakrbac/policy.go"
	if !strings.HasSuffix(f, expectedSuffix) {
		t.Fatalf("policy file %q does not end with %q", f, expectedSuffix)
	}

	l := p.Line()
	expectedLine := 114
	if l != expectedLine {
		t.Fatalf("policy line `%d` did not match `%d`", l, expectedLine)
	}

	full := p.String()
	if strings.Index(full, n) == -1 {
		t.Fatalf("policy description %q does not contain name %q", full, expectedName)
	}

	if strings.Index(full, expectedSuffix) == -1 {
		t.Fatalf("policy description %q does not contain a file path %q", full, expectedSuffix)
	}

	if strings.Index(full, fmt.Sprintf("%d", expectedLine)) == -1 {
		t.Fatalf("policy description %q does not contain line number %d", full, expectedLine)
	}
}
