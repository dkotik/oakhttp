/*
Package token provides a secure random string factory that can be used for session or other types of tokens.
*/
package token

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math/big"
)

type Factory func() (string, error)

func New(withOptions ...Option) (Factory, error) {
	o := &options{}
	for _, option := range append(
		withOptions,
		WithDefaultOptions(),
		func(o *options) error { // validate
			if o.EdgeLength*2+4 > o.TokenLength {
				return errors.New("token edge size is too great")
			}

			buf := make([]byte, o.TokenLength) // read at least one-token worth
			if _, err := io.ReadFull(rand.Reader, buf); err != nil {
				return fmt.Errorf("crypto/rand is unavailable: %w", err)
			}
			return nil
		},
	) {
		if err := option(o); err != nil {
			return nil, fmt.Errorf("cannot initialize token factory: %w", err)
		}
	}

	o.bodyStop = o.TokenLength - o.EdgeLength
	o.bodyCharacterSetLength = big.NewInt(int64(len(o.BodyCharacterSet)))
	o.edgeCharacterSetLength = big.NewInt(int64(len(o.EdgeCharacterSet)))
	return func() (string, error) {
		var (
			b     = make([]byte, o.TokenLength)
			i     = 0
			order *big.Int
			err   error
		)

		for i = 0; i < o.EdgeLength; i++ {
			order, err = rand.Int(rand.Reader, o.edgeCharacterSetLength)
			if err != nil {
				return "", err
			}
			b[i] = o.EdgeCharacterSet[order.Int64()]
		}

		for i = o.EdgeLength; i < o.bodyStop; i++ {
			order, err = rand.Int(rand.Reader, o.bodyCharacterSetLength)
			if err != nil {
				return "", err
			}
			b[i] = o.BodyCharacterSet[order.Int64()]
		}

		for i = o.bodyStop; i < o.TokenLength; i++ {
			order, err = rand.Int(rand.Reader, o.edgeCharacterSetLength)
			if err != nil {
				return "", err
			}
			b[i] = o.EdgeCharacterSet[order.Int64()]
		}

		return string(b), nil
	}, nil
}
