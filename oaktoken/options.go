package oaktoken

import (
	"errors"
	"fmt"
	"math/big"
)

type options struct {
	EdgeLength       int
	TokenLength      int
	EdgeCharacterSet []byte
	BodyCharacterSet []byte

	// derrived properties
	bodyStop               int // token length - edge length
	edgeCharacterSetLength *big.Int
	bodyCharacterSetLength *big.Int
}

type Option func(*options) error

func WithDefaultOptions() Option {
	return func(o *options) (err error) {
		defer func() {
			if err != nil {
				err = fmt.Errorf("failed using default setting: %w", err)
			}
		}()

		if o.TokenLength == 0 {
			if err = WithTokenLength(24)(o); err != nil {
				return err
			}
		}
		if o.EdgeLength == 0 {
			target := int(float32(o.TokenLength)*0.1) + 1
			if target > 24 {
				target = 24
			}
			if err = WithTokenEdgeLength(target)(o); err != nil {
				return err
			}
		}
		if o.BodyCharacterSet == nil {
			if err = WithBodyCharacterSet([]byte(`abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`))(o); err != nil {
				return err
			}
		}
		if o.EdgeCharacterSet == nil {
			if err = WithEdgeCharacterSet(append(o.BodyCharacterSet, '-'))(o); err != nil {
				return err
			}
		}
		return nil
	}
}

func WithTokenLength(l int) Option {
	return func(o *options) error {
		if o.TokenLength != 0 {
			return errors.New("token length is already set")
		}
		if l < 24 {
			return errors.New("token length of less than 24 characters is not secure")
		}
		if l > 4096 {
			return errors.New("token length of more than 4096 users too much memory")
		}
		o.TokenLength = l
		return nil
	}
}

func WithTokenEdgeLength(l int) Option {
	return func(o *options) error {
		if o.EdgeLength != 0 {
			return errors.New("token edge length is already set")
		}
		if l < 1 {
			return errors.New("token edge length cannot be less than 1")
		}
		if l > 24 {
			return errors.New("token edge length cannot exceed 24")
		}
		o.EdgeLength = l
		return nil
	}
}

func WithBodyCharacterSet(set []byte) Option {
	return func(o *options) error {
		if o.BodyCharacterSet != nil {
			return errors.New("body character set is already set")
		}
		length := len(set)
		if length < 10 {
			return errors.New("body character set cannot contain less than 10 items")
		}
		if length > 256 {
			return errors.New("body character set cannot exceed 255 entries")
		}
		o.BodyCharacterSet = set
		return nil
	}
}

func WithEdgeCharacterSet(set []byte) Option {
	return func(o *options) error {
		if o.EdgeCharacterSet != nil {
			return errors.New("edge character set is already set")
		}
		length := len(set)
		if length < 10 {
			return errors.New("edge character set cannot contain less than 10 items")
		}
		if length > 256 {
			return errors.New("edge character set cannot exceed 255 entries")
		}
		o.EdgeCharacterSet = set
		return nil
	}
}
