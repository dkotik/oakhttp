package oakbotswat

import (
	"context"
	"net/http"
)

type Botswat interface {
	// VerifyResponseToken should return `nil` for valid tokens, [ErrNotHuman] for rejected tokens, or an [Error] for any other condition.
	//
	// A response may only be validated once. If the same response is presented twice, the second and each subsequent request will generate an error stating that the response has already been consumed.
	VerifyResponseToken(
		ctx context.Context,
		clientToken string,
		clientIPAddress string,
	) error
}

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
