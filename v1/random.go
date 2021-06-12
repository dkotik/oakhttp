package oakacs

import (
	"crypto/rand"
	"fmt"
	"io"
)

// GenerateRandomBytes is used to generate random bytes of given size.
func GenerateRandomBytes(size int) ([]byte, error) {
	return GenerateRandomBytesWithReader(size, rand.Reader)
}

// GenerateRandomBytesWithReader is used to generate random bytes of given size read from a given reader.
func GenerateRandomBytesWithReader(size int, reader io.Reader) ([]byte, error) {
	if reader == nil {
		return nil, fmt.Errorf("provided reader is nil")
	}
	buf := make([]byte, size)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return nil, fmt.Errorf("failed to read random bytes: %v", err)
	}
	return buf, nil
}

const uuidLen = 16

// GenerateUUID is used to generate a random UUID
// func GenerateUUID() (string, error) {
// 	return GenerateUUIDWithReader(rand.Reader)
// }
