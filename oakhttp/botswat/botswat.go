package botswat

import (
	"context"
	"net/http"
)

type Verifier func(
	ctx context.Context,
	clientResponse string,
	clientIPAddress string,
) (
	userData string,
	err error,
)

type ResponseExtractor func(
	r *http.Request,
) (
	clientResponse string,
	err error,
)
