package oakwords

// NewTranslator sets up a dictionary for translating bytes to words and back. If nil dictionary is provided, the English dictionary is used by default.
func NewTranslator(dictionary *[256]string) *Translator {
	if dictionary == nil {
		dictionary = &English
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
	dictionary *[256]string
	reverse    map[string]byte
}

// Decode translated a sequence of words into bytes.
func (t *Translator) Decode(w []string) []byte {
	result := make([]byte, len(w))
	for i, u := range w {
		result[i] = t.reverse[u]
	}
	return result
}

// Encode translated a sequence of bytes into words.
func (t *Translator) Encode(b []byte) []string {
	result := make([]string, len(b))
	for i, u := range b {
		result[i] = t.dictionary[u]
	}
	return result
}
