package oakpolicy

import "testing"

func TestResourcePathSuccessfulMatches(t *testing.T) {
	cases := []struct {
		path DomainPath
		mask []string
	}{
		{path: nil, mask: nil},
		{path: nil, mask: []string{">"}},
		{path: DomainPath{"service", "resource", "uuid"}, mask: []string{"*", "*", "*"}},
		{path: DomainPath{"service", "resource", "uuid"}, mask: []string{"service", "*", "*"}},
		{path: DomainPath{"service", "resource", "uuid"}, mask: []string{"*", "resource", "*"}},
		{path: DomainPath{"service", "resource", "uuid"}, mask: []string{"service", "*", "uuid"}},
		{path: DomainPath{"service", "resource", "uuid"}, mask: []string{"service", ">", "resource"}},
		{path: DomainPath{"service", "resource", "uuid"}, mask: []string{">"}},
	}

	for i, c := range cases {
		if !c.path.Match(c.mask...) {
			t.Errorf("case %d: path %q did not match mask %q", i, c.path, c.mask)
		}
	}
}

func TestResourcePathFailingMatches(t *testing.T) {
	cases := []struct {
		path DomainPath
		mask []string
	}{
		{path: nil, mask: []string{"1"}},
		{path: nil, mask: []string{"*"}},
		{path: DomainPath{"1"}, mask: nil},
		{path: DomainPath{"service", "resource", "uuid"}, mask: []string{"service", "1", "resource"}},
		{path: DomainPath{"service", "resource", "uuid"}, mask: []string{"service", "resource", "*", "*"}},
		{path: DomainPath{"service", "resource", "uuid"}, mask: nil},
		// {path: nil, mask: nil},
	}

	for i, c := range cases {
		if c.path.Match(c.mask...) {
			t.Errorf("case %d: path %q did matched mask %q", i, c.path, c.mask)
		}
	}
}
