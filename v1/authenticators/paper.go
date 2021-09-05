package authenticators

import (
	"context"
	"errors"
	"fmt"

	"github.com/dkotik/oakacs/v1/oakwords"
)

// var _ oakacs.Authenticator = (*Paper)(nil)

// Paper stores a long printable key that can be used for account recovery.
type Paper struct {
	Length     int
	Visible    int
	Translator *oakwords.Translator
	// Labeler func([]string) string ?
}

func (p *Paper) Prepare() (string, error) {
	if p.Length < 12 || p.Translator == nil {
		return "", errors.New("paper authenticator is not correctly setup")
	}
	b, err := random(p.Length)
	if err != nil {
		return "", err
	}
	return oakwords.FromBytes(b), nil
}

func (p *Paper) Generate(ctx context.Context, tokenOrPassword string) ([]byte, error) {
	// fields := strings.Fields(tokenOrPassword)
	b, err := oakwords.ToBytes(tokenOrPassword)
	if len(b) != p.Length {
		return nil, errors.New("provided code is not the right size")
	}
	if err != nil {
		return nil, fmt.Errorf("could not parse code: %w", err)
	}
	return b, nil
}

// func (p *Paper) Compare(ctx context.Context, tokenOrPassword string, secret oakacs.Secret) error {
// 	// fields := strings.Fields(tokenOrPassword)
// 	b, err := oakwords.ToBytes(tokenOrPassword)
// 	if len(b) != p.Length {
// 		return errors.New("provided code is not the right size")
// 	}
// 	if err != nil {
// 		return fmt.Errorf("could not parse code: %w", err)
// 	}
// 	return secret(b)
// }
