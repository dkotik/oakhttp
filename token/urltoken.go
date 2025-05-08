package token

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

func NewURLToken(entropyByteCount int) (Factory, error) {
	if entropyByteCount < 12 {
		return nil, errors.New("URL token factory requires at least twelve entropy bytes")
	}
	return func() (string, error) {
		b := make([]byte, entropyByteCount)
		if _, err := io.ReadFull(rand.Reader, b); err != nil {
			return "", fmt.Errorf("cannot read random source: %w", err)
		}
		return base64.URLEncoding.EncodeToString(b), nil
	}, nil
}
