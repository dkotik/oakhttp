package oakquery

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func Normalize(s string) string {
	b := &strings.Builder{}
	io.Copy(b, transform.NewReader(strings.NewReader(s), StringFilter))
	return b.String()
}

func Quote(s string) string {
	b := &strings.Builder{}
	io.Copy(b, transform.NewReader(strings.NewReader(s), StringFilter))
	return fmt.Sprintf("%q", b)
}

func QuoteIdentifier(s string) string {
	b := &strings.Builder{}
	b.WriteRune('`')
	io.Copy(b, transform.NewReader(strings.NewReader(s), IdentifierFilter))
	if b.Len() == 1 {
		return "`unknown`"
	}
	b.WriteRune('`')
	return b.String()
}

var (
	StringFilter = transform.Chain(norm.NFC,
		runes.Remove(runes.Predicate(func(r rune) bool {
			return !unicode.IsOneOf([]*unicode.RangeTable{
				unicode.Letter,
				unicode.Digit,
				unicode.Space,
				unicode.Punct,
				unicode.Symbol,
			}, r)
		})))

	IdentifierFilter = transform.Chain(norm.NFC,
		runes.Remove(runes.Predicate(func(r rune) bool {
			return !unicode.IsOneOf([]*unicode.RangeTable{
				unicode.Letter,
				unicode.Digit,
				&unicode.RangeTable{
					R16: []unicode.Range16{
						{0x002d, 0x002e, 1}, // - dash and . dot
						{0x005f, 0x005f, 1}, // _ underscore
					},
				},
			}, r)
		})))
)
