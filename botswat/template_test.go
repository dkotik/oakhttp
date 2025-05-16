package botswat

import (
	"bytes"
	"testing"

	"github.com/sebdah/goldie/v2"
)

func TestErrorRendering(t *testing.T) {
	tmpl, err := NewTemplate()
	if err != nil {
		t.Fatal(err)
	}
	b := &bytes.Buffer{}
	if err = tmpl.Execute(b, templateOptions{}); err != nil {
		t.Fatal("template rendering failed:", err)
	}
	g := goldie.New(t)
	g.Assert(t, "template", b.Bytes())
}
