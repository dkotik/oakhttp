package oakwords

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/big"
	"text/scanner"
)

var defaultReaderDictionary = defaultDictionary.Reverse()

func NewReader(r io.Reader) *Reader {
	reader := &Reader{
		Scanner:    &scanner.Scanner{},
		dictionary: defaultReaderDictionary,
	}
	reader.Scanner.Error = func(s *scanner.Scanner, msg string) {
		reader.scanErr = errors.New(msg)
	}
	reader.Scanner.Init(r)

	return reader
}

type Reader struct {
	*scanner.Scanner
	scanErr    error
	dictionary map[string]byte
}

func (r *Reader) SetDictionary(d *Dictionary) {
	r.dictionary = d.Reverse()
}

func (r *Reader) Read(p []byte) (n int, err error) {
	for i := range p {
		c := r.Scan()
		if c == scanner.EOF {
			return n, io.EOF
		}
		if r.scanErr != nil {
			return n, r.scanErr
		}
		b, ok := r.dictionary[r.TokenText()]
		if !ok {
			return n, fmt.Errorf("word %q is not in chosen dictionary", r.TokenText())
		}
		p[i] = b
		n++
	}
	return
}

func ReadInt(r io.Reader) (n int64, err error) {
	wrapped := NewReader(r)
	b := &bytes.Buffer{}
	_, err = io.CopyN(b, wrapped, 24)
	if err != nil && err != io.EOF {
		return 0, err
	}

	chopped, ok := ChecksumChop(b.Bytes())
	if !ok {
		return 0, errors.New("checksum did not match")
	}

	x := new(big.Int).SetBytes(chopped)
	return x.Int64(), nil
}
