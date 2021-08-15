package oakquery

import (
	"fmt"
	"strings"
	"unicode"
)

var (
	underscoreDashDot = &unicode.RangeTable{
		R16: []unicode.Range16{
			{0x002d, 0x002e, 1}, // - dash and . dot
			{0x005f, 0x005f, 1}, // _ underscore
		},
	}

	unicodeIdentifierRange = []*unicode.RangeTable{
		unicode.Letter,
		unicode.Digit,
		underscoreDashDot,
	}

	unicodeStringRange = []*unicode.RangeTable{
		unicode.Letter,
		unicode.Digit,
		unicode.Space,
		unicode.Punct,
		unicode.Symbol,
	}
)

func Quote(s string) string {
	b := &strings.Builder{}
	for _, char := range s {
		if unicode.IsOneOf(unicodeIdentifierRange, char) {
			b.WriteRune(char)
		}
	}
	return fmt.Sprintf("%q", b)
}

func QuoteIdentifier(s string) string {
	b := &strings.Builder{}
	b.WriteRune('`')
	for _, char := range s {
		if unicode.IsOneOf(unicodeIdentifierRange, char) {
			b.WriteRune(char)
		}
	}
	if b.Len() == 1 {
		return "`unknown`"
	}
	b.WriteRune('`')
	return b.String()
}
