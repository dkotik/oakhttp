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
	"fmt"
	"strings"
)

// A segment is a pattern piece that matches one or more path segments, or
// a trailing slash.
// If wild is false, it matches a literal segment, or, if s == "/", a trailing slash.
// If wild is true and multi is false, it matches a single path segment.
// If both wild and multi are true, it matches all remaining path segments.
type segment struct {
	s     string // literal or wildcard name or "/" for "/{$}".
	wild  bool
	multi bool // "..." wildcard
}

func (s segment) String() string {
	switch {
	case s.multi && s.s == "": // Trailing slash.
		return "/"
	case s.multi:
		return fmt.Sprintf("/{%s...}", s.s)
	case s.wild:
		return fmt.Sprintf("/{%s}", s.s)
	case s.s == "/":
		return "/{$}"
	default: // Literal.
		return "/" + s.s
	}
}

type Pattern struct {
	name     string
	segments []segment
}

func (p *Pattern) Path(fields map[string]string) (string, error) {
	var (
		b     strings.Builder
		value string
		ok    bool
	)

	for _, s := range p.segments {
		_ = b.WriteByte('/')
		switch {
		case s.s == "" || s.s == "/": // trailing slash or {$}
		case s.wild || s.multi:
			value, ok = fields[s.s]
			if !ok {
				return "", fmt.Errorf("route pattern %s[%s] does not contain field named %q", p.name, p, s.s)
			}
			_, _ = b.WriteString(value)
		default:
			_, _ = b.WriteString(s.s) // literal value
		}
	}
	return b.String(), nil
}

func (p *Pattern) String() string {
	var b strings.Builder
	for _, s := range p.segments {
		b.WriteString(s.String())
	}
	return b.String()
}

// returns segment, "/" for trailing slash, or "" for done.
// path should start with a "/"
func nextSegment(path string) (seg, rest string) {
	if path == "/" {
		return "/", ""
	}
	path = path[1:] // drop initial slash
	i := strings.IndexByte(path, '/')
	if i < 0 {
		return path, ""
	}
	return path[:i], path[i:]
}
