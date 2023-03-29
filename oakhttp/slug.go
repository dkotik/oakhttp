package oakhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
)

func NewSlugRequestAdaptor(
	h func(context.Context, string) error,
) (Handler, error) {
	if h == nil {
		return nil, errors.New("slug request adaptor cannot use a <nil> handler")
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		slug := path.Base(r.URL.Path)
		if slug == "" {
			return NewInvalidRequestError(errors.New("empty slug"))
		}
		return h(r.Context(), slug)
	}, nil
}

func NewSlugRequestResponseAdaptor[R any](
	h func(context.Context, string) (R, error),
	withOptions ...Option,
) (Handler, error) {
	if h == nil {
		return nil, errors.New("slug request-response adaptor cannot use a <nil> handler")
	}
	o, err := newOptions(append(withOptions, WithDefaultOptions()))
	if err != nil {
		return nil, err
	}

	encoder := o.Encoder
	return func(w http.ResponseWriter, r *http.Request) error {
		slug := path.Base(r.URL.Path)
		if slug == "" {
			return NewInvalidRequestError(errors.New("empty slug"))
		}
		response, err := h(r.Context(), slug)
		if err != nil {
			return err
		}
		if err = encoder(w, response); err != nil {
			return fmt.Errorf("cannot encode response: %w", err)
		}
		return nil
	}, nil
}
