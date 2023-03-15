package oakhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type options struct {
	readLimit                  int64
	encoderContentType         string
	encoderFactory             func(io.Writer) Encoder
	decoderFactory             func(io.Reader) Decoder
	errorHandler               func(Handler) http.HandlerFunc
	middlewareFromInnerToOuter []Middleware
}

type Option func(*options) error

func WithDefaultOptions() Option {
	return func(o *options) (err error) {
		defer func() {
			if err != nil {
				err = fmt.Errorf("failed to apply default option settion: %w", err)
			}
		}()

		if o.readLimit == 0 {
			if err = WithReadLimit(1024 * 64)(o); err != nil {
				return err
			}
		}
		if o.encoderFactory == nil || o.decoderFactory == nil {
			if err = WithEncoderDecoder(
				"application/json",
				func(w io.Writer) Encoder {
					return json.NewEncoder(w)
				},
				func(r io.Reader) Decoder {
					return json.NewDecoder(r)
				},
			)(o); err != nil {
				return err
			}
		}
		if o.errorHandler == nil {
			if err = WithErrorHandler(DefaultErrorHandler)(o); err != nil {
				return err
			}
		}
		return nil
	}
}

func WithReadLimit(n int64) Option {
	return func(o *options) error {
		if o.readLimit != 0 {
			return errors.New("read limit is already set")
		}
		if n < 1 {
			return errors.New("minimum read limit is 1")
		}
		if n > 1024*1024*96 {
			return errors.New("maximum read limit is 1024*1024*96")
		}
		o.readLimit = n
		return nil
	}
}

func WithEncoderDecoder(
	encoderContentType string,
	encoderFactory func(io.Writer) Encoder,
	decoderFactory func(io.Reader) Decoder,
) Option {
	return func(o *options) error {
		if o.encoderFactory != nil || o.decoderFactory != nil {
			return errors.New("encoder and decoder are already set")
		}
		if encoderContentType == "" {
			return errors.New("cannot use an empty encoder content type")
		}
		if encoderFactory == nil {
			return errors.New("cannot use a <nil> encoder factory")
		}
		if o.decoderFactory != nil {
			return errors.New("decoder factory is already set")
		}
		if decoderFactory == nil {
			return errors.New("cannot use a <nil> decoder factory")
		}
		o.encoderContentType = encoderContentType
		o.encoderFactory = encoderFactory
		o.decoderFactory = decoderFactory
		return nil
	}
}

func WithErrorHandler(h func(Handler) http.HandlerFunc) Option {
	return func(o *options) error {
		if o.errorHandler != nil {
			return errors.New("error handler is already set")
		}
		if h == nil {
			return errors.New("cannot use a <nil> error handler")
		}
		o.errorHandler = h
		return nil
	}
}
