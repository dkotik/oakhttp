package oakrouter

import (
	"errors"
	"testing"
)

func TestMatcher(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		m := NewMatcher("simple")
		first, err := m.MatchString()
		if err != nil {
			t.Fatal("could not match first string:", err)
		}
		if first != "simple" {
			t.Fatal("first segment does not match `simple`:", first)
		}
		if err = m.MatchEnd(); err != nil {
			t.Fatal("could not match end:", err)
		}
	})

	t.Run("normal", func(t *testing.T) {
		m := NewMatcher("/normal/second")
		first, err := m.MatchString()
		if err != nil {
			t.Fatal("could not match first string:", err)
		}
		if first != "normal" {
			t.Fatal("first segment does not match `normal`:", first)
		}

		second, err := m.MatchString()
		if err != nil {
			t.Fatal("could not match first string:", err)
		}
		if second != "second" {
			t.Fatal("second segment does not match `second`:", second)
		}

		if err = m.MatchEnd(); err != nil {
			t.Fatal("could not match end:", err)
		}
	})

	t.Run("long", func(t *testing.T) {
		m := NewMatcher("/normal/99/second")
		first, err := m.MatchString()
		if err != nil {
			t.Fatal("could not match first string:", err)
		}
		if first != "normal" {
			t.Fatal("first segment does not match `normal`:", first)
		}

		number, err := m.MatchUint()
		if err != nil {
			t.Fatal("could not match the number in the middle:", err)
		}
		if number != 99 {
			t.Fatal("middle number does not match `99`:", number)
		}

		second, err := m.MatchString()
		if err != nil {
			t.Fatal("could not match first string:", err)
		}
		if second != "second" {
			t.Fatal("second segment does not match `second`:", second)
		}

		if err = m.MatchEnd(); err != nil {
			t.Fatal("could not match end:", err)
		}
	})

	t.Run("double slash", func(t *testing.T) {
		m := NewMatcher("/simple//double")
		first, err := m.MatchString()
		if err != nil {
			t.Fatal("could not match first string:", err)
		}
		if first != "simple" {
			t.Fatal("first segment does not match `simple`:", first)
		}
		if _, err = m.MatchBytes(); !errors.Is(err, ErrDoubleSlash) {
			t.Fatal("double slash not found where expected:", err)
		}
	})

	t.Run("trailing slash", func(t *testing.T) {
		m := NewMatcher("simple/")
		first, err := m.MatchString()
		if err != nil {
			t.Fatal("could not match first string:", err)
		}
		if first != "simple" {
			t.Fatal("first segment does not match `simple`:", first)
		}
		if err = m.MatchEnd(); !errors.Is(err, ErrTrailingSlash) {
			t.Fatal("trailing slash not found where expected:", err)
		}
	})
}

func TestChompRouterPath(t *testing.T) {
	cases := []struct{ Path, Head, Tail string }{
		{"", "", ""},
		{"//", "", "/"},
		{"//tail", "", "/tail"},
		{"///tail", "", "//tail"},
		{"///tail/", "", "//tail/"},
		{"/simplest", "simplest", ""},
		{"simplest", "simplest", ""},
		{"/simplest/tail", "simplest", "/tail"},
		{"simplest/tail", "simplest", "/tail"},
		{"/simplest/", "simplest", "/"},
		{"simplest/", "simplest", "/"},
		{"/simplest/x", "simplest", "/x"},
		{"/simplest//", "simplest", "//"},
	}

	for _, c := range cases {
		t.Run(c.Path, func(t *testing.T) {
			head, tail := ChompPath(c.Path)
			if c.Head != head {
				t.Fatal("head does not match", c.Path, "=>", c.Head, "vs", head)
			}
			if c.Tail != tail {
				t.Fatal("tail does not match", c.Path, "=>", c.Tail, "vs", tail)
			}
		})
	}
}
