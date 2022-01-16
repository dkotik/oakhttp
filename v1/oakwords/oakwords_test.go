package oakwords

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

var reLegitWordString = regexp.MustCompile(`^(\w\w\w\w[\s\-\,\.]+){3,}\w\w\w\w$`)

func legitWords(t *testing.T, s string) {
	t.Helper()
	if !reLegitWordString.MatchString(s) {
		t.Fatalf("value %q does not constitute a legitimate word list", s)
	}
}

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

func TestIntTransformations(t *testing.T) {
	cases := []int64{0, 9, 16, 32, 999999, 38729387428974, 2374761653249823, 88999999}
	for _, i := range cases {
		t.Run(fmt.Sprintf("transforming integer: %d", i), func(t *testing.T) {
			b := &bytes.Buffer{}
			if err := WriteInt(b, i); err != nil {
				t.Fatal(err)
			}
			legitWords(t, b.String())
			t.Log("words:", b.String())

			j, err := ReadInt(bytes.NewReader(b.Bytes()))
			if err != nil {
				t.Fatal(err)
			}
			if j != i {
				t.Fatalf("%d does not match %d", j, i)
			}
		})
	}
}

func ExampleFromBytes() {
	fmt.Println(
		FromBytes([]byte("marvel")),
	)
	// Output: type chat flop milk tone aunt
}

func ExampleToBytes() {
	b, err := ToBytes("  type chat     flop  \n\n milk     tone   aunt    ")
	fmt.Println(string(b), err)
	// Output: marvel <nil>
}

func ExampleTranslator_Encode() {
	tr := NewTranslator(nil)

	fmt.Println(
		tr.Encode([]byte("great")),
	)
	// Output: [cape flop tone chat year]
}

func ExampleTranslator_Decode() {
	tr := NewTranslator(nil)

	b, err := tr.Decode(
		[]string{"cape", "flop", "tone", "chat", "year"},
	)
	fmt.Println(string(b), err)
	// Output: great <nil>
}
