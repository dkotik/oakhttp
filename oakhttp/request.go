package oakhttp

import (
	"context"
	"fmt"
	"net/http"
	"path"
)

func AdaptRequestResponse[
	T comparable,
	R ValidatableNormalizable[T],
	P comparable,
](
	usingDomainAdaptor *DomainAdaptor,
	domainRequestHandler func(context.Context, R) (P, error),
	middleware ...Middleware,
) http.Handler {
	return usingDomainAdaptor.ApplyMiddleware(
		func(w http.ResponseWriter, r *http.Request) error {
			var request R
			err := usingDomainAdaptor.ReadRequest(&request, w, r)
			if err != nil {
				return err
			}
			if err = request.Validate(); err != nil {
				return fmt.Errorf("invalid request: %w", err)
			}
			if err = request.Normalize(); err != nil {
				return fmt.Errorf("failed to normalize request: %w", err)
			}

			response, err := domainRequestHandler(r.Context(), request)
			if err != nil {
				return err
			}
			return usingDomainAdaptor.WriteResponse(w, response)
		},
		middleware...,
	)
}

func AdaptCustomRequestResponse[T, P comparable](
	usingDomainAdaptor *DomainAdaptor,
	requestDecoderValidatorNormalizer func(*http.Request) (T, error),
	domainRequestHandler func(context.Context, T) (P, error),
	middleware ...Middleware,
) http.Handler {
	return usingDomainAdaptor.ApplyMiddleware(
		func(w http.ResponseWriter, r *http.Request) error {
			request, err := requestDecoderValidatorNormalizer(r)
			if err != nil {
				return err
			}
			response, err := domainRequestHandler(r.Context(), request)
			if err != nil {
				return err
			}
			return usingDomainAdaptor.WriteResponse(w, response)
		},
		middleware...,
	)
}

func AdaptRequest[T comparable, R ValidatableNormalizable[T]](
	usingDomainAdaptor *DomainAdaptor,
	domainRequestHandler func(context.Context, R) error,
	middleware ...Middleware,
) http.Handler {
	return usingDomainAdaptor.ApplyMiddleware(
		func(w http.ResponseWriter, r *http.Request) error {
			var request R
			err := usingDomainAdaptor.ReadRequest(&request, w, r)
			if err != nil {
				return err
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
			return domainRequestHandler(r.Context(), request)
		},
		middleware...,
	)
}

func AdaptURLPathTail(
	usingDomainAdaptor *DomainAdaptor,
	domainRequestHandler func(context.Context, string) error,
	middleware ...Middleware,
) http.Handler {
	return AdaptCustomRequest(
		usingDomainAdaptor,
		func(r *http.Request) (string, error) {
			return path.Base(r.URL.Path), nil
		},
		domainRequestHandler,
		middleware...,
	)
}

func AdaptURLPathTailResponse[P comparable](
	usingDomainAdaptor *DomainAdaptor,
	domainRequestHandler func(context.Context, string) (P, error),
	middleware ...Middleware,
) http.Handler {
	return AdaptCustomRequestResponse(
		usingDomainAdaptor,
		func(r *http.Request) (string, error) {
			return path.Base(r.URL.Path), nil
		},
		domainRequestHandler,
		middleware...,
	)
}
