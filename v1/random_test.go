package oakacs

import (
	"crypto/rand"
	"regexp"
	"strings"
	"testing"
)

func TestWords(t *testing.T) {
	t.Run("are words unique", func(t *testing.T) {
		m := make(map[string]bool)
		for _, w := range words {
			if w == "" {
				break
			}
			if _, ok := m[w]; ok {
				t.Fatalf("word %q is not unique", w)
			}
			m[w] = true
		}
		if missing := 256 - len(m); missing > 0 {
			t.Fatalf("please come up with %d more words", missing)
		}
	})

	t.Run("humanizing bytes", func(t *testing.T) {
		var b [12]byte
		rand.Read(b[:])
		result := strings.Join(Humanize(b[:]), " ")
		if !regexp.MustCompile(`^(\w\w\w\w\s){11,11}\w\w\w\w$`).MatchString(result) {
			t.Fatal("humanized expected pattern did not match", result)
		}
	})
}
