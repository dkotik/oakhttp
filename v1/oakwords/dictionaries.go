package oakwords

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"strings"
)

// Dictionary holds 256 words, each corresponding to a byte value.
type Dictionary [256]string

// Validate iterates through every value to check for uniqueness and extra white space characters.
func (d *Dictionary) Validate() error {
	if d == nil {
		return errors.New("provided dictionary is not initialized")
	}

	m := make(map[string]struct{})
	for i, entry := range d {
		w := strings.TrimSpace(entry)
		if w != entry {
			return fmt.Errorf("dictionary value %q has extra white space", entry)
		}
		if w == "" {
			return fmt.Errorf("dictionary value #%d is empty", i)
		}
		if _, ok := m[w]; ok {
			return fmt.Errorf("dictionary value %q is not unique", w)
		}
		m[w] = struct{}{}
	}
	// if missing := 256 - len(m); missing > 0 {
	// 	return fmt.Errorf("please come up with %d more words", missing)
	// }
	return nil
}

// Embedded dictionary names:
const (
	DictionaryEnglishNouns = "english-nouns.txt"
)

var (
	//go:embed dictionaries/english-nouns.txt
	embededDictionaries embed.FS
	defaultDictionary   *Dictionary
)

// Load retrieves an embedded dictionary and parses it.
func Load(name string) *Dictionary {
	b, err := fs.ReadFile(embededDictionaries, path.Join("dictionaries", name))
	if err != nil {
		panic(err)
	}

	var i int
	var word []byte
	var result Dictionary
	for i, word = range bytes.Fields(b) {
		result[i] = string(word)
	}

	return &result
}

// Use sets up the default dictionary.
func Use(dictionary *Dictionary) {
	defaultDictionary = dictionary
	if err := defaultDictionary.Validate(); err != nil {
		panic(err)
	}
}
