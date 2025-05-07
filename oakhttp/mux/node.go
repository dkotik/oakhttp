package mux

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

type node struct {
	// special children keys:
	//     "/"	trailing slash
	//	   ""   single wildcard
	//	   "*"  multi wildcard
	// children   *hybrid  // map[string]*node // interior node
	children   children
	emptyChild *node    // child with key ""
	pat        *Pattern // leaf
}

func (n *node) addSegments(segs []segment, p *Pattern) {
	if len(segs) == 0 {
		if n.pat != nil {
			panic("n.pat != nil")
		}
		n.pat = p
		return
	}
	seg := segs[0]
	if seg.multi {
		if len(segs) != 1 {
			panic("multi wildcard not last")
		}
		if n.findChild("*") != nil {
			panic("dup multi wildcards")
		}
		c := n.addChild("*")
		c.pat = p
	} else if seg.wild {
		n.addChild("").addSegments(segs[1:], p)
	} else {
		n.addChild(seg.s).addSegments(segs[1:], p)
	}
}

func (n *node) addChild(key string) *node {
	if key == "" {
		if n.emptyChild == nil {
			n.emptyChild = &node{}
		}
		return n.emptyChild
	}
	if c := n.findChild(key); c != nil {
		return c
	}
	c := &node{}
	if n.children == nil {
		// n.children = newHybrid(8)
		n.children = make(listOfChildren, 0)
	}
	// n.children.add(key, c)
	n.children = n.children.append(key, c)
	return c
}

func (n *node) findChild(key string) *node {
	if n.children == nil {
		return nil
	}
	return n.children.get(key)
}

func (n *node) matchPath(path string, matches []string) (*Pattern, []string) {
	// If path is empty, then return the node's pattern, which
	// may be nil.
	if path == "" {
		return n.pat, matches
	}
	seg, rest := nextSegment(path)
	if c := n.findChild(seg); c != nil {
		if p, m := c.matchPath(rest, matches); p != nil {
			return p, m
		}
	}
	// Match single wildcard.
	if c := n.emptyChild; c != nil {
		if p, m := c.matchPath(rest, append(matches, seg)); p != nil {
			return p, m
		}
	}
	// Match multi wildcard to the rest of the pattern.
	if c := n.findChild("*"); c != nil {
		return c.pat, append(matches, path[1:]) // remove initial slash
	}
	return nil, nil
}

// Modifies n; use for testing only.
func (n *node) print(w io.Writer, level int) {
	indent := strings.Repeat("    ", level)
	if n.pat != nil {
		fmt.Fprintf(w, "%s%q\n", indent, n.pat)
	} else {
		fmt.Fprintf(w, "%snil\n", indent)
	}
	if n.emptyChild != nil {
		fmt.Fprintf(w, "%s%q:\n", indent, "")
		n.emptyChild.print(w, level+1)
	}

	if n.children == nil {
		return
	}
	keys := n.children.keys()
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(w, "%s%q:\n", indent, k)
		n.findChild(k).print(w, level+1)
	}
}
