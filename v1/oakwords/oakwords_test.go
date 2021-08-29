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

func ExampleFromBytes() {
	fmt.Println(
		FromBytes([]byte("marvel")),
	)
	// Output: math hole song sage flow cash
}

func ExampleToBytes() {
	b, err := ToBytes("  math hole     song  \n\n sage     flow   cash    ")
	fmt.Println(string(b), err)
	// Output: marvel <nil>
}

func ExampleTranslator_Encode() {
	tr := NewTranslator(nil)

	fmt.Println(
		tr.Encode([]byte("great")),
	)
	// Output: [wish song flow hole team]
}

func ExampleTranslator_Decode() {
	tr := NewTranslator(nil)

	b, err := tr.Decode(
		[]string{"wish", "song", "flow", "hole", "team"},
	)
	fmt.Println(string(b), err)
	// Output: great <nil>
}
