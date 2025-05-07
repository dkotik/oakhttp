// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// This is taken from Jonathan Amsterdam's alternative
// ServeMux [implementation]. For full discussion
// see the GitHub [proposal].
//
// [implementation]: https://github.com/jba/muxpatterns/blob/main/tree.go
// [proposal]: https://github.com/golang/go/discussions/60227

package mux

import (
	"strconv"
	"strings"
	"testing"

	"golang.org/x/exp/slices"
)

func TestNextSegment(t *testing.T) {
	for _, test := range []struct {
		in   string
		want []string
	}{
		{"/a/b/c", []string{"a", "b", "c"}},
		{"/a/b/", []string{"a", "b", "/"}},
		{"/", []string{"/"}},
	} {
		var got []string
		rest := test.in
		for len(rest) > 0 {
			var seg string
			seg, rest = nextSegment(rest)
			got = append(got, seg)
		}
		if !slices.Equal(got, test.want) {
			t.Errorf("%q: got %v, want %v", test.in, got, test.want)
		}
	}
}

// // TODO: turn this into proper test
var testTree *node

func getTestTree() *node {
	if testTree == nil {
		initTestTree()
	}
	return testTree
}

func initTestTree() {
	testTree = &node{}
	// var ps PatternSet // TODO: return this.
	for _, p := range []string{"/a", "/a/b", "/a/{x}",
		"/g/h/i", "/g/{x}/j",
		"/a/b/{x...}", "/a/b/{y}", "/a/b/{$}"} {
		pat, err := Parse(p)
		if err != nil {
			panic(err)
		}
		// if err := ps.Register(pat); err != nil {
		// 	panic(err)
		// }
		testTree.addSegments(pat.segments, pat)
	}
}

func TestAddPattern(t *testing.T) {
	want := `nil
"a":
    "/a"
    "":
        "/a/{x}"
    "b":
        "/a/b"
        "":
            "/a/b/{y}"
        "*":
            "/a/b/{x...}"
        "/":
            "/a/b/{$}"
"g":
    nil
    "":
        nil
        "j":
            "/g/{x}/j"
    "h":
        nil
        "i":
            "/g/h/i"
`

	var b strings.Builder
	getTestTree().print(&b, 0)
	got := b.String()
	if got != want {
		t.Errorf("got\n%s\nwant\n%s", got, want)
	}
}

func TestNodeMatch(t *testing.T) {
	for _, test := range []struct {
		path        string
		wantPat     string // "" for nil
		wantMatches []string
	}{
		{"/a", "/a", nil},
		{"/b", "", nil},
		{"/a/b", "/a/b", nil},
		{"/a/c", "/a/{x}", []string{"c"}},
		{"/a/b/", "/a/b/{$}", nil},
		{"/a/b/c", "/a/b/{y}", []string{"c"}},
		{"/a/b/c/d", "/a/b/{x...}", []string{"c/d"}},
		{"/g/h/i", "/g/h/i", nil},
		{"/g/h/j", "/g/{x}/j", []string{"h"}},
	} {
		gotPat, gotMatches := getTestTree().matchPath(test.path, nil)
		got := ""
		if gotPat != nil {
			got = gotPat.String()
		}
		if got != test.wantPat {
			t.Errorf("%s: got %q, want %q", test.path, got, test.wantPat)
		}
		if !slices.Equal(gotMatches, test.wantMatches) {
			t.Errorf("%s: got matches %v, want %v", test.path, gotMatches, test.wantMatches)
		}
	}
}

func findChildLinear(key string, entries []entry) *node {
	for _, e := range entries {
		if key == e.key {
			return e.child
		}
	}
	return nil
}

func TestHybrid(t *testing.T) {
	nodes := []*node{&node{}, &node{}, &node{}, &node{}, &node{}}
	h := newHybrid(4)
	for i := 0; i < 4; i++ {
		h.add(strconv.Itoa(i), nodes[i])
	}
	if h.m != nil {
		t.Fatal("h.m != nil")
	}
	for i := 0; i < 4; i++ {
		g := h.get(strconv.Itoa(i))
		if g != nodes[i] {
			t.Fatalf("%d: different", i)
		}
	}
	h.add("4", nodes[4])
	if h.s != nil {
		t.Fatal("h.s != nil")
	}
	if h.m == nil {
		t.Fatal("h.m == nil")
	}
	if g := h.get("4"); g != nodes[4] {
		t.Fatal("4 diff")
	}
}
