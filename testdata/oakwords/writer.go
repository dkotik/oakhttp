package oakwords

import (
	"errors"
	"io"
	"math/big"
)

type SeparatorFunc func(wordCount int) []byte

func SeparatorHTML(perRow int) SeparatorFunc {
	return func(wordCount int) []byte {
		if wordCount%perRow == 0 {
			if wordCount == 0 {
				return nil
			}
			return []byte("<br />")
		}
		return []byte(" ")
	}
}

type Writer struct {
	io.Writer
	separator  SeparatorFunc
	wordCount  int
	dictionary *Dictionary
}

func (w *Writer) Write(p []byte) (n int, err error) {
	var (
		j, l int
		sep  []byte
	)
	for _, c := range p {
		sep = w.separator(w.wordCount)
		l = len(sep)
		if l > 0 {
			j, err = w.Writer.Write(sep)
			if err != nil {
				return
			}
			if j != l {
				return n, io.ErrShortWrite
			}
		}

		j, err = w.Writer.Write([]byte(w.dictionary[c]))
		if err != nil {
			return
		}
		if j != len(w.dictionary[c]) {
			return n, io.ErrShortWrite
		}
		w.wordCount++
		n++
	}
	return n, nil
}

func (w *Writer) SetDictionary(d *Dictionary) {
	w.dictionary = d
}

func (w *Writer) SetSeparator(s SeparatorFunc) {
	w.separator = s
}

func NewWriter(w io.Writer) io.WriteCloser {
	return ChecksumWriter(&Writer{
		Writer: w,
		separator: func(wordCount int) []byte {
			if wordCount%4 == 0 {
				if wordCount == 0 {
					return nil
				}
				return []byte("\n")
			}
			return []byte(" ")
		},
		dictionary: defaultDictionary,
	})
}

func WriteInt(w io.Writer, n int64) (err error) {
	if n < 0 {
		return errors.New("negative integers cannot be reliabily encoded to bytes")
	}

	x := new(big.Int).SetInt64(n)
	wrapped := NewWriter(w)
	_, err = wrapped.Write(x.Bytes())
	if err != nil {
		return
	}
	return wrapped.Close()
}
