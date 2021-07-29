package oakacs

import (
	"crypto/rand"
	"errors"
)

func random(length int) (b []byte, err error) {
	b = make([]byte, length)
	n, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	if n < length {
		return nil, errors.New("source of crypto is compromised")
	}
	return b, nil
}
