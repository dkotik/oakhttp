package oakhttp

import (
	"context"
	"io"
	"net/http"
)

type Handler func(w http.ResponseWriter, r *http.Request) error

type Middleware func(Handler) Handler

type Encoder func(http.ResponseWriter, any) error

type Decoder func(any, io.Reader) error

type RequestFactory[T any, P ValidatableNormalizable[T]] func(
	w http.ResponseWriter,
	r *http.Request,
) (P, error)

type DomainRequest[T any, P ValidatableNormalizable[T]] func(context.Context, P) error

type DomainRequestResponse[T any, P ValidatableNormalizable[T], O any] func(context.Context, P) (O, error)

type Cache interface {
	Set(ctx context.Context, key, value string) (err error)
	Get(ctx context.Context, key string) (value string, err error)
}
