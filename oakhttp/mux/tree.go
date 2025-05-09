package mux

import (
	"errors"
	"fmt"
)

// Branches abstracts either map or list implementation for child [Node]s for performance. When there are more than [maximumListOfChildrenSize], the map implementation is prefered for faster look up.
type Branches interface {
	Get(string) *Node
	Append(string, *Node) Branches
	Keys() []string
}

// Node is the core routing tree component.
type Node struct {
	Leaf              *Route
	TrailingSlashLeaf *Route
	TerminalLeaf      *Route
	Branches          Branches
	DynamicBranch     *Node
}

func (n *Node) MatchPath(path string) (route *Route, matches []string) {
	switch path {
	case "":
		return n.Leaf, nil
	case "/":
		return n.TrailingSlashLeaf, nil
	}

	segment, remainder, err := munchPath(path) // speed this up
	if err != nil {                            // drop this check
		return nil, nil // double slash error possible only
	}

	if n.Branches != nil {
		if branch := n.Branches.Get(segment); branch != nil {
			route, matches = branch.MatchPath(remainder)
			if route != nil {
				return route, matches
			}
		}
	}

	if n.DynamicBranch != nil {
		// TODO: pass in array with pre-initialized len instead?
		route, matches = n.DynamicBranch.MatchPath(remainder)
		if route != nil {
			return route, append([]string{segment}, matches...)
		}
	}

	if n.TerminalLeaf != nil {
		// if n.TerminalLeaf.segments[len(n.TerminalLeaf.segments)-1].Name() == "" {
		// 	return n.TerminalLeaf, nil // deal with the {...}
		// }
		return n.TerminalLeaf, []string{path[1:]}
	}
	return nil, nil
}

func (n *Node) Grow(route *Route, remaining []Segment) (err error) {
	if len(remaining) == 0 { // leaf
		if n.Leaf != nil {
			return fmt.Errorf("routes %q and %q overlap: %s resolves to the same static tree node as %s", n.Leaf.Name(), route.Name(), n.Leaf.String(), route.String())
		}
		n.Leaf = route
		return nil
	}
	current := remaining[0]
	switch current.Type() {
	case SegmentTypeTrailingSlash: // leaf
		if n.TrailingSlashLeaf != nil {
			return fmt.Errorf("routes %q and %q overlap: %s resolves to the same trailing slash tree node as %s", n.TrailingSlashLeaf.Name(), route.Name(), n.TrailingSlashLeaf.String(), route.String())
		}
		n.TrailingSlashLeaf = route
	case SegmentTypeTerminal: // leaf
		if n.TerminalLeaf != nil {
			return fmt.Errorf("routes %q and %q overlap: %s resolves to the same terminal tree node as %s", n.TerminalLeaf.Name(), route.Name(), n.TerminalLeaf.String(), route.String())
		}
		n.TerminalLeaf = route
	case SegmentTypeStatic: // branch
		if n.Branches == nil {
			return errors.New("must create branch list first")
		}
		name := current.Name()
		node := n.Branches.Get(name)
		if node == nil {
			node = &Node{}
			n.Branches = n.Branches.Append(name, node)
		}
		return node.Grow(route, remaining[1:])
	case SegmentTypeDynamic: // branch
		if n.DynamicBranch == nil {
			n.DynamicBranch = &Node{}
		}
		return n.DynamicBranch.Grow(route, remaining[1:])
	default:
		return fmt.Errorf("cannot grow tree using a segment %q of unknown type %q", current.Name(), current.Type())
	}
	return nil
}
