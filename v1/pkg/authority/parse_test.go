package cueroles

import (
	"fmt"
	"testing"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"github.com/davecgh/go-spew/spew"
)

type (
	Match struct {
		Attribute string
		Operation string
		Pattern   string
	}

	RolePlan struct {
		Name  string
		Allow []Match
		Deny  []Match
	}
)

func TestParseFromCue(t *testing.T) {
	ctx := cuecontext.New()

	scope := ctx.CompileString(`
        #permissions: [...#permission]
        #permission: [...#match]
        #match: close({
        	attribute: string
            operation?: string
        	pattern: string
        })
        #timeouts: {
            session: int
            idle: int
        }
    `)
	if err := scope.Err(); err != nil {
		t.Fatal(err)
	}

	// spew.Dump(scope)

	role := ctx.CompileString(`
        allow: #permission & [
            {attribute: "attr1", operation: "test", pattern: "pattern1" }
        ]
        deny: #permission & [
            {attribute: "attr1", operation: "test", pattern: "pattern1" }
        ]
        timeouts: #timeouts & {
            session: 1
            idle: 1
        }
    `, cue.Scope(scope))
	if err := role.Err(); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", role)

	var rplan RolePlan
	if err := role.Decode(&rplan); err != nil {
		t.Fatal(err)
	}
	spew.Dump(rplan)
}
