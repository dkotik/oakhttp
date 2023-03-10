package botswat

import (
	"context"
	"errors"
	"net/http"
)

var ErrNotHuman = errors.New("robot detected")

// Verifier returns [ErrNotHuman] if client response was not recognized as valid.
//
// A response may only be validated once. If the same response is presented twice, the second and each subsequent request will generate an error stating that the response has already been consumed.
type Verifier func(
	ctx context.Context,
	clientResponseToken string,
	clientIPAddress string,
) (
	userData string,
	err error,
)

type Cache func(Verifier) Verifier

type CacheAdaptor interface {
	Get(context.Context, []byte) ([]byte, error)
	Set(context.Context, []byte, []byte) error
}

type ResponseExtractor func(
	r *http.Request,
) (
	clientResponseToken string,
	err error,
)
