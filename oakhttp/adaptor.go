package oakhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type adaptor[
	T any,
	V ValidatableNormalizable[T],
	O any,
] struct {
	readLimit int64
	call      DomainRequestResponse[T, V, O]
}

func NewAdaptor[
	T any,
	V ValidatableNormalizable[T],
	O any,
](
	call DomainRequestResponse[T, V, O],
	withOptions ...AdaptorOption,
) Handler {
	o := newAdaptorOptions(withOptions)
	return ApplyMiddleware(&adaptor[T, V, O]{
		readLimit: o.readLimit,
		call:      call,
	}, o.middlewares)
}

func (a *adaptor[T, V, O]) ServeHyperText(
	w http.ResponseWriter,
	r *http.Request,
) error {
	header := r.Header
	contentType := header.Get("Content-Type")
	codec, err := GetCodec(contentType)
	if err != nil {
		return err
	}

	var request V
	if err = codec.Decode(
		&request,
		http.MaxBytesReader(w, r.Body, a.readLimit),
	); err != nil {
		return NewInvalidRequestError(fmt.Errorf("unable to decode: %w", err))
	}
	if err = request.Validate(); err != nil {
		return NewInvalidRequestError(err)
	}

	response, err := a.call(r.Context(), request)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", contentType)
	if err = codec.Encode(w, response); err != nil {
		return fmt.Errorf("unable to encode: %w", err)
	}
	return nil
}

type terminalAdaptor[
	T any,
	V ValidatableNormalizable[T],
] struct {
	readLimit int64
	call      DomainRequest[T, V]
}

func NewTerminalAdaptor[
	T any,
	V ValidatableNormalizable[T],
](
	call DomainRequest[T, V],
	withOptions ...AdaptorOption,
) Handler {
	o := newAdaptorOptions(withOptions)
	return ApplyMiddleware(&terminalAdaptor[T, V]{
		readLimit: o.readLimit,
		call:      call,
	}, o.middlewares)
}

func (a *terminalAdaptor[T, V]) ServeHyperText(
	w http.ResponseWriter,
	r *http.Request,
) error {
	header := r.Header
	contentType := header.Get("Content-Type")
	codec, err := GetCodec(contentType)
	if err != nil {
		return err
	}

	var request V
	if err = codec.Decode(
		&request,
		http.MaxBytesReader(w, r.Body, a.readLimit),
	); err != nil {
		return NewInvalidRequestError(fmt.Errorf("unable to decode: %w", err))
	}
	if err = request.Validate(); err != nil {
		return NewInvalidRequestError(err)
	}
	if err = a.call(r.Context(), request); err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func NewStaticAdaptor[O any](
	call func(context.Context) (O, error),
	withOptions ...AdaptorOption,
) Handler {
	o := newAdaptorOptions(withOptions)
	return ApplyMiddleware(HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) (err error) {
			header := r.Header
			contentType := header.Get("Content-Type")
			codec, err := GetCodec(contentType)
			if err != nil {
				return err
			}
			response, err := call(r.Context())
			if err != nil {
				return err
			}
			w.Header().Set("Content-Type", contentType)
			if err = codec.Encode(w, response); err != nil {
				return fmt.Errorf("unable to encode: %w", err)
			}
			return nil
		},
	), o.middlewares)
}

type complexAdaptor[
	T any,
	V ValidatableNormalizable[T],
	O any,
] struct {
	readLimit int64
	call      DomainRequestResponse[T, V, O]
	finalizer func(V, *http.Request) error
}

func NewComplexAdaptor[
	T any,
	V ValidatableNormalizable[T],
	O any,
](
	call DomainRequestResponse[T, V, O],
	finalizer func(V, *http.Request) error,
	withOptions ...AdaptorOption,
) Handler {
	o := newAdaptorOptions(withOptions)
	return ApplyMiddleware(&complexAdaptor[T, V, O]{
		readLimit: o.readLimit,
		call:      call,
		finalizer: finalizer,
	}, o.middlewares)
}

func (a *complexAdaptor[T, V, O]) ServeHyperText(
	w http.ResponseWriter,
	r *http.Request,
) error {
	header := r.Header
	contentType := header.Get("Content-Type")
	codec, err := GetCodec(contentType)
	if err != nil {
		return err
	}

	var request V
	if err = codec.Decode(
		&request,
		http.MaxBytesReader(w, r.Body, a.readLimit),
	); err != nil {
		return NewInvalidRequestError(fmt.Errorf("unable to decode: %w", err))
	}
	if err := a.finalizer(request, r); err != nil {
		return NewInvalidRequestError(err)
	}
	if err = request.Validate(); err != nil {
		return NewInvalidRequestError(err)
	}

	response, err := a.call(r.Context(), request)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", contentType)
	if err = codec.Encode(w, response); err != nil {
		return fmt.Errorf("unable to encode: %w", err)
	}
	return nil
}

type complexTerminalAdaptor[
	T any,
	V ValidatableNormalizable[T],
] struct {
	readLimit int64
	call      DomainRequest[T, V]
	finalizer func(V, *http.Request) error
}

func NewComplexTerminalAdaptor[
	T any,
	V ValidatableNormalizable[T],
](
	call DomainRequest[T, V],
	finalizer func(V, *http.Request) error,
	withOptions ...AdaptorOption,
) Handler {
	o := newAdaptorOptions(withOptions)
	return ApplyMiddleware(&complexTerminalAdaptor[T, V]{
		readLimit: o.readLimit,
		call:      call,
		finalizer: finalizer,
	}, o.middlewares)
}

func (a *complexTerminalAdaptor[T, V]) ServeHyperText(
	w http.ResponseWriter,
	r *http.Request,
) error {
	header := r.Header
	contentType := header.Get("Content-Type")
	codec, err := GetCodec(contentType)
	if err != nil {
		return err
	}

	var request V
	if err = codec.Decode(
		&request,
		http.MaxBytesReader(w, r.Body, a.readLimit),
	); err != nil {
		return NewInvalidRequestError(fmt.Errorf("unable to decode: %w", err))
	}
	if err := a.finalizer(request, r); err != nil {
		return NewInvalidRequestError(err)
	}
	if err = request.Validate(); err != nil {
		return NewInvalidRequestError(err)
	}
	if err = a.call(r.Context(), request); err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

type stringAdaptor[O any] struct {
	readLimit int
	call      func(context.Context, string) (O, error)
	extractor StringExtractor
}

func NewStringAdaptor[O any](
	call func(context.Context, string) (O, error),
	extractor StringExtractor,
	withOptions ...AdaptorOption,
) Handler {
	o := newAdaptorOptions(withOptions)
	return ApplyMiddleware(&stringAdaptor[O]{
		readLimit: int(o.readLimit),
		call:      call,
		extractor: extractor,
	}, o.middlewares)
}

func (a *stringAdaptor[O]) ServeHyperText(
	w http.ResponseWriter,
	r *http.Request,
) error {
	header := r.Header
	contentType := header.Get("Content-Type")
	codec, err := GetCodec(contentType)
	if err != nil {
		return err
	}

	request, err := a.extractor(r)
	if err != nil {
		return NewInvalidRequestError(fmt.Errorf("unable to decode: %w", err))
	}
	if len(request) > a.readLimit {
		return NewInvalidRequestError(errors.New("request value is too long"))
	}

	response, err := a.call(r.Context(), request)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", contentType)
	if err = codec.Encode(w, response); err != nil {
		return fmt.Errorf("unable to encode: %w", err)
	}
	return nil
}

type terminalStringAdaptor struct {
	readLimit int
	call      func(context.Context, string) error
	extractor StringExtractor
}

func NewTerminalStringAdaptor(
	call func(context.Context, string) error,
	extractor StringExtractor,
	withOptions ...AdaptorOption,
) Handler {
	o := newAdaptorOptions(withOptions)
	return ApplyMiddleware(&terminalStringAdaptor{
		readLimit: int(o.readLimit),
		call:      call,
		extractor: extractor,
	}, o.middlewares)
}

func (a *terminalStringAdaptor) ServeHyperText(
	w http.ResponseWriter,
	r *http.Request,
) error {
	request, err := a.extractor(r)
	if err != nil {
		return NewInvalidRequestError(fmt.Errorf("unable to decode: %w", err))
	}
	if len(request) > a.readLimit {
		return NewInvalidRequestError(errors.New("request value is too long"))
	}

	if err = a.call(r.Context(), request); err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}
