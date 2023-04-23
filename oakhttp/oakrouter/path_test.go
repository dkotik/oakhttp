package oakrouter

import (
	"testing"
)

func TestChompRouterPath(t *testing.T) {
	cases := []struct{ Path, Head, Tail string }{
		{"", "", ""},
		{"//", "", "/"},
		{"//tail", "", "/tail"},
		{"///tail", "", "//tail"},
		{"///tail/", "", "//tail/"},
		{"/simplest", "simplest", ""},
		{"simplest", "simplest", ""},
		{"/simplest/tail", "simplest", "/tail"},
		{"simplest/tail", "simplest", "/tail"},
		{"/simplest/", "simplest", "/"},
		{"simplest/", "simplest", "/"},
		{"/simplest/x", "simplest", "/x"},
		{"/simplest//", "simplest", "//"},
	}

	for _, c := range cases {
		t.Run(c.Path, func(t *testing.T) {
			head, tail := ChompPath(c.Path)
			if c.Head != head {
				t.Fatal("head does not match", c.Path, "=>", c.Head, "vs", head)
			}
			if c.Tail != tail {
				t.Fatal("tail does not match", c.Path, "=>", c.Tail, "vs", tail)
			}
		})
	}
}
