package oakwords

import (
	"fmt"
	"strings"
)

// FromBytes translates bytes into words using the default dictionary.
func FromBytes(b []byte) string {
	result := make([]string, len(b))
	for i, u := range b {
		result[i] = defaultDictionary[u]
	}
	return strings.Join(result, " ")
}

// ToBytes translates words into bytes using the default dictionary.
func ToBytes(words string) ([]byte, error) {
	return NewTranslator(defaultDictionary).Decode(strings.Fields(words))
}

// NewTranslator sets up a dictionary for translating bytes to words and back. If nil dictionary is provided, the EnglishNouns dictionary is used by default.
func NewTranslator(dictionary *Dictionary) *Translator {
	if dictionary == nil {
		dictionary = defaultDictionary
	} else {
		if err := dictionary.Validate(); err != nil {
			panic(err)
		}
	}
	t := &Translator{
		dictionary: dictionary,
		reverse:    make(map[string]byte),
	}

	{ // setup
		for k, v := range dictionary {
			t.reverse[v] = byte(k)
		}
	}

	return t
}

// Translator encodes bytes to words and decodes words to bytes. 256 common four letter nouns can that represent bytes for code recovery.
type Translator struct {
	dictionary *Dictionary
	reverse    map[string]byte
}

// Encode translated a sequence of bytes into words.
func (t *Translator) Encode(b []byte) []string {
	result := make([]string, len(b))
	for i, u := range b {
		result[i] = t.dictionary[u]
	}
	return result
}

// Decode translated a sequence of words into bytes.
func (t *Translator) Decode(w []string) ([]byte, error) {
	result := make([]byte, len(w))
	for i, u := range w {
		if b, ok := t.reverse[u]; ok {
			result[i] = b
		} else {
			return nil, fmt.Errorf("word %q was not found in the dictionary", u)
		}
	}
	return result, nil
}
