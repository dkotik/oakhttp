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
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// Parse parses a string into a Pattern. Pattern consists of slash-separated segments, where each segment is either a literal or a wildcard of the form "{name}", "{name...}", or "{$}".
//
// Wildcard names must be valid Go identifiers.
// The "{$}" and "{name...}" wildcard must occur at the end of PATH.
// PATH may end with a '/'.
// Wildcard names in a path must be distinct.
func Parse(rest string) (*Pattern, error) {
	if len(rest) == 0 {
		return nil, errors.New("empty pattern")
	}

	p := &Pattern{}
	seenNames := map[string]bool{}
	for len(rest) > 0 {
		// Invariant: rest[0] == '/'.
		rest = rest[1:]
		if len(rest) == 0 {
			// Trailing slash.
			p.segments = append(p.segments, segment{wild: true, multi: true})
			break
		}
		i := strings.IndexByte(rest, '/')
		if i == 0 {
			return nil, errors.New("empty path segment")
		}
		if i < 0 {
			i = len(rest)
		}
		var seg string
		seg, rest = rest[:i], rest[i:]
		if i := strings.IndexByte(seg, '{'); i < 0 {
			// Literal
			p.segments = append(p.segments, segment{s: seg})
		} else {
			// Wildcard
			if i != 0 {
				return nil, errors.New("bad wildcard segment (must start with '{')")
			}
			if seg[len(seg)-1] != '}' {
				return nil, errors.New("bad wildcard segment (must end with '}')")
			}
			name := seg[1 : len(seg)-1]
			if name == "$" {
				if len(rest) != 0 {
					return nil, errors.New("{$} not at end")
				}
				p.segments = append(p.segments, segment{s: "/"})
				break
			}
			var multi bool
			if strings.HasSuffix(name, "...") {
				multi = true
				name = name[:len(name)-3]
				if len(rest) != 0 {
					return nil, errors.New("{...} wildcard not at end")
				}
			}
			if name == "" {
				return nil, errors.New("empty wildcard")
			}
			if !isValidWildcardName(name) {
				return nil, fmt.Errorf("bad wildcard name %q", name)
			}
			if seenNames[name] {
				return nil, fmt.Errorf("duplicate wildcard name %q", name)
			}
			seenNames[name] = true
			p.segments = append(p.segments, segment{s: name, wild: true, multi: multi})
		}
	}
	return p, nil
}

func isValidWildcardName(s string) bool {
	if s == "" {
		return false
	}
	// Valid Go identifier.
	for i, c := range s {
		if !unicode.IsLetter(c) && c != '_' && (i == 0 || !unicode.IsDigit(c)) {
			return false
		}
	}
	return true
}

// DescribeRelationship returns a string that describes how pat1 and pat2
// are related, in terms of the paths they match.
func DescribeRelationship(pat1, pat2 string) string {
	p1, err := Parse(pat1)
	if err != nil {
		panic(err)
	}
	p2, err := Parse(pat2)
	if err != nil {
		panic(err)
	}
	return describeRel(p1, p2)
}

// relationship is a relationship between two patterns.
type relationship string

const (
	moreSpecific relationship = "moreSpecific"
	moreGeneral  relationship = "moreGeneral"
	overlaps     relationship = "overlaps"
	equivalent   relationship = "equivalent"
	disjoint     relationship = "disjoint"
)

// comparePaths classifies the paths of the patterns into one of four
// groups:
//
//	moreGeneral: p1 matches all the paths of p2 and more
//	moreSpecific: p2 matches all the paths of p1 and more
//	overlaps: there is a path that both match, but neither is more specific
//	disjoint: there is no path that both match
func comparePaths(p1, p2 *Pattern) relationship {
	// Track whether a single (non-multi) wildcard in p1 matched
	// a literal in p2, and vice versa.
	// We care about these because if a wildcard matches a literal, then the
	// pattern with the wildcard can't be more specific than the one with the
	// literal.
	wild1MatchedLit2 := false
	wild2MatchedLit1 := false
	var segs1, segs2 []segment
	for segs1, segs2 = p1.segments, p2.segments; len(segs1) > 0 && len(segs2) > 0; segs1, segs2 = segs1[1:], segs2[1:] {
		s1 := segs1[0]
		s2 := segs2[0]
		if s1.multi && s2.multi {
			// Two multis match each other.
			continue
		}
		if s1.multi {
			// p1 matches the rest of p2.
			// Does that mean it is more general than p2?
			if !wild2MatchedLit1 {
				// If p2 didn't have any wildcards that matched literals in p1,
				// then yes, p1 is more general.
				return moreGeneral
			}
			// Otherwise neither is more general than the other.
			return overlaps
		}
		if s2.multi {
			// p2 matches the rest of p1. The same logic as above applies.
			if !wild1MatchedLit2 {
				return moreSpecific
			}
			return overlaps
		}
		if s1.s == "/" && s2.s == "/" {
			// Both patterns end in "/{$}"; they match.
			continue
		}
		if s1.s == "/" || s2.s == "/" {
			// One pattern ends in "/{$}", and the other doesn't, nor is the other's
			// corresponding segment a multi. So they are disjoint.
			return disjoint
		}
		if s1.wild && s2.wild {
			// These single-segment wildcards match each other.
		} else if s1.wild {
			// p1's single wildcard matches the corresponding segment of p2.
			wild1MatchedLit2 = true
		} else if s2.wild {
			// p2's single wildcard matches the corresponding segment of p1.
			wild2MatchedLit1 = true
		} else {
			// Two literal segments.
			if s1.s != s2.s {
				return disjoint
			}
		}
	}
	// We've reached the end of the corresponding segments of the patterns.
	if len(segs1) == 0 && len(segs2) == 0 {
		// The patterns matched completely.
		switch {
		case wild1MatchedLit2 && !wild2MatchedLit1:
			return moreGeneral
		case wild2MatchedLit1 && !wild1MatchedLit2:
			return moreSpecific
		case !wild1MatchedLit2 && !wild2MatchedLit1:
			return equivalent
		default:
			return overlaps
		}
	}
	// One pattern has more segments than the other.
	// The only way they can fail to be disjoint is if one ends in a multi, but
	// we handled that case in the loop.
	return disjoint
}

func describeRel(p1, p2 *Pattern) string {
	rel := comparePaths(p1, p2)
	switch rel {
	case disjoint:
		return fmt.Sprintf("%s has no paths in common with %s.", p1, p2)
	case equivalent:
		return fmt.Sprintf("%s matches the same paths as %s.", p1, p2)
	case moreSpecific:
		over := matchingPath(p1)
		diff := differencePath(p2, p1)
		return fmt.Sprintf(`%s is more specific than %s.
Both match %q.
Only %[2]s matches %[4]q.`,
			p1, p2, over, diff)
	case moreGeneral:
		over := matchingPath(p2)
		diff := differencePath(p1, p2)
		return fmt.Sprintf(`%s is more general than %s.
Both match %q.
Only %[1]s matches %[4]q.`,
			p1, p2, over, diff)
	default: // overlap
		return fmt.Sprintf(`%[1]s and %[2]s both match some paths, like %[3]q.
But neither is more specific than the other.
%[1]s matches %[4]q, but %[2]s doesn't.
%[2]s matches %[5]q, but %[1]s doesn't.`,
			p1, p2, overlapPath(p1, p2), differencePath(p1, p2), differencePath(p2, p1))
	}
}

func matchingPath(p *Pattern) string {
	var b strings.Builder
	writeMatchingPath(&b, p.segments)
	return b.String()
}

// writeMatchingPath writes to b a path that matches the segments.
func writeMatchingPath(b *strings.Builder, segs []segment) {
	for _, s := range segs {
		writeSegment(b, s)
	}
}

func writeSegment(b *strings.Builder, s segment) {
	b.WriteByte('/')
	if !s.multi && s.s != "/" {
		b.WriteString(s.s)
	}
}

// overlapPath returns a path that both p1 and p2 match.
// It assumes there is such a path.
func overlapPath(p1, p2 *Pattern) string {
	var b strings.Builder
	var segs1, segs2 []segment
	for segs1, segs2 = p1.segments, p2.segments; len(segs1) > 0 && len(segs2) > 0; segs1, segs2 = segs1[1:], segs2[1:] {
		s1 := segs1[0]
		s2 := segs2[0]
		if s1.wild {
			writeSegment(&b, s2)
		} else {
			writeSegment(&b, s1)
		}
	}
	if len(segs1) > 0 {
		writeMatchingPath(&b, segs1)
	} else if len(segs2) > 0 {
		writeMatchingPath(&b, segs2)
	}
	return b.String()
}

// differencePath returns a path that p1 matches and p2 doesn't.
// It assumes there is such a path.
func differencePath(p1, p2 *Pattern) string {
	b := new(strings.Builder)

	var segs1, segs2 []segment
	for segs1, segs2 = p1.segments, p2.segments; len(segs1) > 0 && len(segs2) > 0; segs1, segs2 = segs1[1:], segs2[1:] {
		s1 := segs1[0]
		s2 := segs2[0]
		if s1.multi && s2.multi {
			// From here the patterns match the same paths, so we must have found a difference earlier.
			b.WriteByte('/')
			return b.String()

		}
		if s1.multi && !s2.multi {
			// s1 ends in a "..." wildcard but s2 does not.
			// A trailing slash will distinguish them, unless s2 ends in "{$}",
			// in which case any segment will do; prefer the wildcard name if
			// it has one.
			b.WriteByte('/')
			if s2.s == "/" {
				if s1.s != "" {
					b.WriteString(s1.s)
				} else {
					b.WriteString("x")
				}
			}
			return b.String()
		}
		if !s1.multi && s2.multi {
			writeSegment(b, s1)
		} else if s1.wild && s2.wild {
			// Both patterns will match whatever we put here; use
			// the first wildcard name.
			writeSegment(b, s1)
		} else if s1.wild && !s2.wild {
			// s1 is a wildcard, s2 is a literal.
			// Any segment other than s2.s will work.
			// Prefer the wildcard name, but if it's the same as the literal,
			// tweak the literal.
			if s1.s != s2.s {
				writeSegment(b, s1)
			} else {
				b.WriteByte('/')
				b.WriteString(s2.s + "x")
			}
		} else if !s1.wild && s2.wild {
			writeSegment(b, s1)
		} else {
			// Both are literals. A precondition of this function is that the
			// patterns overlap, so they must be the same literal. Use it.
			if s1.s != s2.s {
				fmt.Printf("%q, %q\n", s1.s, s2.s)
				panic("literals differ")
			}
			writeSegment(b, s1)
		}
	}
	if len(segs1) > 0 {
		// p1 is longer than p2, and p2 does not end in a multi.
		// Anything that matches the rest of p1 will do.
		writeMatchingPath(b, segs1)
	} else if len(segs2) > 0 {
		writeMatchingPath(b, segs2)
	}
	return b.String()
}
