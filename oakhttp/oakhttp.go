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

type Handler func(w http.ResponseWriter, r *http.Request) error

type Middleware func(Handler) Handler

type Encoder interface {
	Encode(any) error
}

type Decoder interface {
	Decode(any) error
}

type ValidatableNormalizable[T comparable] interface {
	*T
	Validate() error
	Normalize() error
}

type DomainAdaptor struct {
	readLimit          int64
	encoderContentType string
	encoderFactory     func(io.Writer) Encoder
	decoderFactory     func(io.Reader) Decoder
	errorHandler       func(Handler) http.Handler
	middleware         []Middleware
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
		readLimit:          o.readLimit,
		encoderContentType: o.encoderContentType,
		encoderFactory:     o.encoderFactory,
		decoderFactory:     o.decoderFactory,
		errorHandler:       o.errorHandler,
		middleware:         o.middleware,
	}, nil
}

func (d *DomainAdaptor) ApplyMiddleware(h Handler) Handler {
	for _, middleware := range d.middleware {
		h = middleware(h)
	}
	return h
}
