package oakwords

import (
	"crypto/rand"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestWords(t *testing.T) {
	tr := NewTranslator(nil)

	t.Run("humanizing bytes", func(t *testing.T) {
		var b [12]byte
		rand.Read(b[:])
		result := strings.Join(tr.Encode(b[:]), " ")
		if !regexp.MustCompile(`^(\w\w\w\w\s){11,11}\w\w\w\w$`).MatchString(result) {
			t.Fatal("humanized expected pattern did not match", result)
		}
	})
}

func ExampleTranslator_Encode() {
	tr := NewTranslator(nil)

	fmt.Println(
		tr.Encode([]byte("great")),
	)
	// Output: [iron leak icon hole levy]
}

func ExampleTranslator_Decode() {
	tr := NewTranslator(nil)

	fmt.Println(
		string(
			tr.Decode(
				[]string{"iron", "leak", "icon", "hole", "levy"},
			),
		),
	)
	// Output: great
}
