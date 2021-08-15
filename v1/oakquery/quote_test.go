package oakquery

import (
	"strings"
	"testing"
)

var allowedIdentifiers = []string{"one", "two", "three-four", "five_six", "s.v"}
var deniedIdentifiers = []string{"", "asd;adsd ", "sdf sdf"}

func TestQuoteIdentifier(t *testing.T) {
	for _, a := range allowedIdentifiers {
		if a != strings.Trim(QuoteIdentifier(a), "`") {
			t.Fatal("Identifiers do not match!", strings.Trim(QuoteIdentifier(a), "`"))
		}
	}

	for _, a := range deniedIdentifiers {
		if a == strings.Trim(QuoteIdentifier(a), "`") {
			t.Fatal("Identifiers match when they must not!")
		}
	}
}
