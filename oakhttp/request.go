package oakhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

func AdaptRequest[T comparable, R ValidatableNormalizable[T]](
	usingDomainAdaptor *DomainAdaptor,
	domainRequestHandler func(context.Context, R) error,
	middleware ...Middleware,
) http.Handler {
	return usingDomainAdaptor.ApplyMiddleware(
		func(w http.ResponseWriter, r *http.Request) error {
			var request R

			defer r.Body.Close()
			reader := http.MaxBytesReader(
				w,
				r.Body,
				usingDomainAdaptor.readLimit,
			)
			defer reader.Close()

			err := usingDomainAdaptor.decoderFactory(reader).Decode(&request)
			if err != nil {
				return fmt.Errorf("decoder failed: %w", err)
			}

			var zero R
			if request == zero {
				return errors.New("empty request")
			}
			if err = request.Validate(); err != nil {
				return fmt.Errorf("invalid request: %w", err)
			}
			if err = request.Normalize(); err != nil {
				return fmt.Errorf("failed to normalize request: %w", err)
			}

			return domainRequestHandler(r.Context(), request)
		},
		middleware...)
}

func AdaptCustomRequest[T comparable](
	usingDomainAdaptor *DomainAdaptor,
	requestDecoderValidatorNormalizer func(*http.Request) (T, error),
	domainRequestHandler func(context.Context, T) error,
	middleware ...Middleware,
) http.Handler {
	return usingDomainAdaptor.ApplyMiddleware(
		func(w http.ResponseWriter, r *http.Request) error {
			request, err := requestDecoderValidatorNormalizer(r)
			if err != nil {
				return err
			}
			var zero T
			if request == zero {
				return errors.New("empty request")
			}
			return domainRequestHandler(r.Context(), request)
		},
		middleware...,
	)
}
