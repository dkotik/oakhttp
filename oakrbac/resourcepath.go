package oakrbac

import "strings"

type ResourcePath []string

const (
	// ResourcePathSeparator divides the sections of resource paths and masks. The value is used for parsing and printing their plain string reprsentations.
	ResourcePathSeparator = "/"

	// ResourcePathWildCardSegment matches any single value of a resource path. The segment must be present to match.
	ResourcePathWildCardSegment = "*"

	// ResourcePathWildCardTail matches any resource path ending segments, present or not. Any resource path mask value after ResourcePathWildCardSegment is ignored.
	ResourcePathWildCardTail = ">"
)

func (r ResourcePath) String() string {
	if r == nil {
		return "<nil>"
	}
	return strings.Join(r, ResourcePathSeparator)
}

// Match returns true if the resource path matches mask segments. The [ResourcePathWildCardSegment] matches any present value. The [ResourcePathWildCardTail] matches any values to the end of the path.
func (r ResourcePath) Match(mask ...string) bool {
	// available := strings.Split(i.ResourcePath, ResourcePathSeparator)
	lenA, lenB := len(r), len(mask)

	if lenB > lenA {
		if mask[lenA] != ResourcePathWildCardTail {
			return false
		}
		mask = mask[:lenA] // chop tail // TODO: test this.
	}
	for position, p := range mask {
		switch p {
		case ResourcePathWildCardTail:
			return true
		case ResourcePathWildCardSegment, r[position]:
			continue
		default:
			return false
		}
	}
	if lenB < lenA {
		return false // must match every segment
	}
	return true
}
