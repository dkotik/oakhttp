package oakwords

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"text/scanner"
)

//go:generate go run dictionaries_generate.go --source dictionaries/english-nouns.txt --destination dict_english_nouns.gen.go --variable EnglishFourLetterNouns

// Dictionary holds 256 words, each corresponding to a byte value.
type Dictionary [256]string

func (d *Dictionary) Reverse() map[string]byte {
	m := make(map[string]byte)
	for i, w := range d {
		m[w] = byte(i)
	}
	return m
}

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
	defaultDictionary = &EnglishFourLetterNouns
)

// Load parses the first 256 words of a dictionary. Commented lines are ignored.
func Load(r io.Reader) (d *Dictionary, err error) {
	s := &scanner.Scanner{}
	s.Init(r)
	s.Error = func(s *scanner.Scanner, msg string) {
		err = errors.New(msg)
	}
	unique := make(map[string]struct{})

	cursor := 0
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		if err != nil || cursor > 255 {
			break
		}
		word := s.TokenText()
		if _, ok := unique[word]; ok {
			return d, fmt.Errorf("word %q is not unique in chosen dictionary", word)
		}
		unique[word] = struct{}{}

		d[cursor] = word
		cursor++
	}
	return
}

// Use sets up the default dictionary.
func Use(dictionary *Dictionary) {
	defaultDictionary = dictionary
	if err := defaultDictionary.Validate(); err != nil {
		panic(err)
	}
}
