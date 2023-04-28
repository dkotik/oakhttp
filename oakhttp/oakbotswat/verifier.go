package oakbotswat

import (
	"context"
	"fmt"
)

// Verifier returns [Error] if client response was not recognized as valid.
type Verifier interface {
	VerifyHumanityToken(
		ctx context.Context,
		clientResponseToken string,
		clientIPAddress string,
	) (
		userData string,
		err error,
	)
}

type cachedVerifier struct {
	cache   Cache
	backend Verifier
}

func (c *cachedVerifier) VerifyHumanityToken(
	ctx context.Context,
	clientResponseToken string,
	clientIPAddress string,
) (
	userData string,
	err error,
) {
	key := clientResponseToken + "||" + clientIPAddress
	userData, ok, err := c.cache.GetToken(ctx, key)
	if err != nil {
		return "", fmt.Errorf("token cache unreachable: %w", err)
	}
	if ok {
		return userData, nil
	}

	userData, err = c.backend.VerifyHumanityToken(ctx, clientResponseToken, clientIPAddress)
	if err != nil {
		return "", err
	}
	if err = c.cache.SetToken(ctx, key, userData); err != nil {
		return "", fmt.Errorf("cannot write to token cache: %w", err)
	}
	return userData, nil
}
