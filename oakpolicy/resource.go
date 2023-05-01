package oakpolicy

import "strings"

type DomainPath []string

type Resource interface {
	DomainPath() DomainPath
}

const (
	// DomainPathSeparator divides the sections of resource paths and masks. The value is used for parsing and printing their plain string reprsentations.
	DomainPathSeparator = "/"

	// DomainPathWildCardSegment matches any single value of a resource path. The segment must be present to match.
	DomainPathWildCardSegment = "*"

	// DomainPathWildCardTail matches any resource path ending segments, present or not. Any resource path mask value after DomainPathWildCardSegment is ignored.
	DomainPathWildCardTail = ">"
)

func NewDomainPath(p ...string) DomainPath {
	return DomainPath(p)
}

func (d DomainPath) String() string {
	return strings.Join(d, DomainPathSeparator)
}

// Match returns true if the resource path matches mask segments. The [DomainPathWildCardSegment] matches any present value. The [DomainPathWildCardTail] matches any values to the end of the path.
func (d DomainPath) Match(mask ...string) bool {
	// available := strings.Split(i.DomainPath, DomainPathSeparator)
	lenA, lenB := len(d), len(mask)

	if lenB > lenA {
		if mask[lenA] != DomainPathWildCardTail {
			return false
		}
		mask = mask[:lenA] // chop tail // TODO: test this.
	}
	for position, p := range mask {
		switch p {
		case DomainPathWildCardTail:
			return true
		case DomainPathWildCardSegment, d[position]:
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
