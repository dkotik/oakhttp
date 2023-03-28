/*
Package oakhttp holds utility methods, adapters, and builders for hardening the most common elements of standard library `net/http` package. It aims to add to add security by default and resistance to misconfiguration where they are insufficient.
*/
package oakhttp

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

// RequestDecoder parses the contents of an [http.Request] into a struct pointer.
// type RequestDecoder func(toStruct any, r *http.Request) error
//
// type Encoder interface {
// 	Encode(http.ResponseWriter, any) error
// }

type Encoder interface {
	Encode(any) error
}

type Decoder interface {
	Decode(any) error
}

// TODO: can use a generic struct to avoid having to use an adaptor along with a adaptor function
// TODO: .
type DomainAdaptor struct {
	readLimit                  int64
	encoderContentType         string
	encoderFactory             func(io.Writer) Encoder
	decoderFactory             func(io.Reader) Decoder
	errorHandler               func(Handler) http.HandlerFunc
	middlewareFromInnerToOuter []Middleware
}

func New(withOptions ...Option) (*DomainAdaptor, error) {
	o := &options{}

	var err error
	for _, option := range append(
		withOptions,
		WithDefaultOptions(),
		func(o *options) error { // validate
			if o.readLimit == 0 {
				return errors.New("read limit is required")
			}
			if o.encoderFactory == nil {
				return errors.New("encoder factory is required")
			}
			if o.decoderFactory == nil {
				return errors.New("decoder factory is required")
			}
			return nil
		},
	) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("unable to initialize domain adaptor: %w", err)
		}
	}

	return &DomainAdaptor{
		readLimit:                  o.readLimit,
		encoderContentType:         o.encoderContentType,
		encoderFactory:             o.encoderFactory,
		decoderFactory:             o.decoderFactory,
		errorHandler:               o.errorHandler,
		middlewareFromInnerToOuter: o.middlewareFromInnerToOuter,
	}, nil
}

func (d *DomainAdaptor) ReadRequest(
	request any,
	w http.ResponseWriter,
	r *http.Request,
) (err error) {
	defer r.Body.Close()
	reader := http.MaxBytesReader(w, r.Body, d.readLimit)
	defer reader.Close()

	err = d.decoderFactory(reader).Decode(request)
	if err != nil {
		return fmt.Errorf("decoder failed: %w", err)
	}
	return nil
}

func (d *DomainAdaptor) WriteResponse(w http.ResponseWriter, r any) (err error) {
	w.Header().Set("Content-Type", d.encoderContentType)
	err = d.encoderFactory(w).Encode(r)
	if err != nil {
		return fmt.Errorf("encoder failed: %w", err)
	}
	return nil
}

func (d *DomainAdaptor) ApplyMiddleware(h Handler, more ...Middleware) http.HandlerFunc {
	for _, middleware := range append(
		more,
		d.middlewareFromInnerToOuter...,
	) {
		h = middleware(h)
	}
	return d.errorOrPanicHandler(h)
}
