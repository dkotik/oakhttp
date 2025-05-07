package mux

import "golang.org/x/exp/maps"

const maximumListOfChildrenSize = 11

// children abstracts either map or list implementation for child [node]s for performance. When there are more than [maximumListOfChildrenSize], the map implementation is prefered for faster look up.
type children interface {
	get(string) *node
	append(string, *node) children
	keys() []string
}

var _ children = (listOfChildren)(nil) // ensure interface satisfaction
var _ children = (mapOfChildren)(nil)  // ensure interface satisfaction

// list implementation ///////////////////////////////////////////////

type child struct {
	key  string
	node *node
}

type listOfChildren []child

func (l listOfChildren) get(key string) *node {
	for _, c := range l {
		if c.key == key {
			return c.node
		}
	}
	return nil
}

func (l listOfChildren) append(key string, node *node) children {
	if len(l) >= maximumListOfChildrenSize {
		m := make(mapOfChildren)
		for _, c := range l {
			m[c.key] = c.node
		}
		m[key] = node
		return m
	}
	return append(l, child{key, node})
}

func (l listOfChildren) keys() []string {
	keys := make([]string, len(l))
	for i, c := range l {
		keys[i] = c.key
	}
	return keys
}

// map implementation ////////////////////////////////////////////////

type mapOfChildren map[string]*node

func (m mapOfChildren) get(key string) *node {
	return m[key]
}

func (m mapOfChildren) append(key string, node *node) children {
	m[key] = node
	return m
}

func (m mapOfChildren) keys() []string {
	keys := make([]string, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

// legacy hybrid implementation //////////////////////////////////////

type entry struct {
	key   string
	child *node
}

type hybrid struct {
	maxSlice int
	s        []entry
	m        map[string]*node
}

func newHybrid(ms int) *hybrid {
	return &hybrid{
		maxSlice: ms,
	}
}

func (h *hybrid) add(k string, v *node) {
	if h.m == nil && len(h.s) < h.maxSlice {
		h.s = append(h.s, entry{k, v})
	} else {
		if h.m == nil {
			h.m = map[string]*node{}
			for _, e := range h.s {
				h.m[e.key] = e.child
			}
			h.s = nil
		}
		h.m[k] = v
	}
}

func (h *hybrid) get(k string) *node {
	if h == nil {
		return nil
	}
	if h.m != nil {
		return h.m[k]
	}
	for _, e := range h.s {
		if e.key == k {
			return e.child
		}
	}
	return nil
}

func (h *hybrid) keys() []string {
	if h == nil {
		return nil
	}
	if h.m != nil {
		return maps.Keys(h.m)
	}
	var keys []string
	for _, e := range h.s {
		keys = append(keys, e.key)
	}
	return keys
}
