package oakwords

import (
	"fmt"
	"testing"
)

func TestDictionaries(t *testing.T) {
	uniqueness := func(words *[256]string) error {
		m := make(map[string]bool)
		for _, w := range words {
			if w == "" {
				break
			}
			if _, ok := m[w]; ok {
				return fmt.Errorf("word %q is not unique", w)
			}
			m[w] = true
		}
		if missing := 256 - len(m); missing > 0 {
			return fmt.Errorf("please come up with %d more words", missing)
		}
		return nil
	}

	for _, d := range map[string]*[256]string{
		"English": &EnglishNouns,
	} {
		t.Run("are words unique?", func(t *testing.T) {
			if err := uniqueness(d); err != nil {
				t.Fatal(err)
			}
		})
	}
}
