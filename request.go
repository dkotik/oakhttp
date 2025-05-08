package oakhttp

import (
	"errors"
	"net/http"
)

func NewRequestAdaptor[
	T any,
	R ValidatableNormalizable[T],
](
	handler DomainRequest[T, R],
	withOptions ...Option,
) (Handler, error) {
	o, err := newOptions(append(withOptions, WithDefaultOptions()))
	if err != nil {
		return nil, err
	}

	limit := o.ReadLimit
	decoder := o.Decoder
	return NewComplexRequestAdaptor(
		handler,
		func(
			w http.ResponseWriter,
			r *http.Request,
		) (R, error) {
			defer func() {
				err = errors.Join(err, r.Body.Close())
				if err != nil {
					err = NewInvalidRequestError(err)
				}
			}()
			var request R
			if err = decoder(
				&request,
				http.MaxBytesReader(w, r.Body, limit),
			); err != nil {
				return nil, err
			}
			return request, nil
		},
	)
}

func NewComplexRequestAdaptor[
	T any,
	R ValidatableNormalizable[T],
](
	handler DomainRequest[T, R],
	factory RequestFactory[T, R],
) (Handler, error) {
	if handler == nil {
		return nil, errors.New("complex request adaptor cannot use a <nil> request handler")
	}
	if factory == nil {
		return nil, errors.New("complex request adaptor cannot use a <nil> request factory")
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		request, err := valid(factory(w, r))
		if err != nil {
			return err
		}
		return handler(r.Context(), request)
	}, nil
}
