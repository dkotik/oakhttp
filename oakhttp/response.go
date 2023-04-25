package oakhttp

import (
	"errors"
	"net/http"
)

func NewRequestResponseAdaptor[
	T any,
	R ValidatableNormalizable[T],
	P any,
](
	handler DomainRequestResponse[T, R, P],
	withOptions ...Option,
) (Handler, error) {
	o, err := newOptions(append(withOptions, WithDefaultOptions()))
	if err != nil {
		return nil, err
	}

	limit := o.ReadLimit
	decoder := o.Decoder
	encoder := o.Encoder
	return NewComplexRequestResponseAdaptor(
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

			// TODO: check Content-Size header.
			// do this for all other places before MaxBytesReader

			var request R
			if err = decoder(
				&request,
				http.MaxBytesReader(w, r.Body, limit),
			); err != nil {
				return nil, err
			}
			return request, nil
		},
		encoder,
	)
}

func NewComplexRequestResponseAdaptor[
	T any,
	R ValidatableNormalizable[T],
	P any,
](
	handler DomainRequestResponse[T, R, P],
	factory RequestFactory[T, R],
	encoder Encoder,
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
		response, err := handler(r.Context(), request)
		if err != nil {
			return err
		}
		return encoder(w, response)
	}, nil
}
