package oakrouter

import (
	"context"
	"reflect"
	"runtime"

	"github.com/dkotik/oakacs/oakhttp"
)

func handlerName(h any) string {
	return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
}

func WithRequestResponseHandler[
	T any,
	R oakhttp.ValidatableNormalizable[T],
	P any,
](
	h func(context.Context, R) (P, error),
	withAdaptorOptions ...oakhttp.Option,
) Option {
	return func(o *options) error {
		name := handlerName(h)
		adaptedRequest, err := oakhttp.NewRequestResponseAdaptor[T, R, P](
			h, withAdaptorOptions...,
		)
		if err != nil {
			return &adaptorError{name: name, cause: err}
		}
		return WithRoute(name, adaptedRequest)(o)
	}
}

func WithRequestHandler[
	T any,
	R oakhttp.ValidatableNormalizable[T],
](
	h func(context.Context, R) error,
	withAdaptorOptions ...oakhttp.Option,
) Option {
	return func(o *options) error {
		name := handlerName(h)
		adaptedRequest, err := oakhttp.NewRequestAdaptor[T, R](
			h, withAdaptorOptions...,
		)
		if err != nil {
			return &adaptorError{name: name, cause: err}
		}
		return WithRoute(name, adaptedRequest)(o)
	}
}

func WithSlugRequestHandler(
	h func(context.Context, string) error,
) Option {
	return func(o *options) error {
		name := handlerName(h)
		adaptedRequest, err := oakhttp.NewSlugRequestAdaptor(h)
		if err != nil {
			return &adaptorError{name: name, cause: err}
		}
		return WithRoute(name, adaptedRequest)(o)
	}
}
