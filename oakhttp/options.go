package oakhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type options struct {
	ReadLimit int64
	Decoder   Decoder
	Encoder   Encoder
}

func newOptions(withOptions []Option) (o *options, err error) {
	o = &options{}
	for _, option := range withOptions {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("cannot initialize the request domain adaptor: %w", err)
		}
	}
	return o, nil
}

type Option func(*options) error

func WithDefaultOptions() Option {
	return func(o *options) (err error) {
		if o.ReadLimit == 0 {
			if err = WithReadLimit(1024 * 1024)(o); err != nil {
				return err
			}
		}
		if o.Decoder == nil || o.Encoder == nil {
			if err = WithEncodingJSON()(o); err != nil {
				return nil
			}
		}
		return nil
	}
}

var DecoderJSON = func(v any, r io.Reader) error {
	return json.NewDecoder(r).Decode(&v)
}

var EncoderJSON = func(w http.ResponseWriter, v any) error {
	w.Header().Set("Context-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func WithEncodingJSON() Option {
	return func(o *options) (err error) {
		if err = WithDecoder(DecoderJSON)(o); err != nil {
			return err
		}
		return WithEncoder(EncoderJSON)(o)
	}
}

func WithReadLimit(l int64) Option {
	return func(o *options) error {
		if o.ReadLimit != 0 {
			return errors.New("read limit is already set")
		}
		if l < 1 {
			return errors.New("cannot read less than 1 byte")
		}
		if l > 1<<32 {
			return errors.New("read limit is too large")
		}
		o.ReadLimit = l
		return nil
	}
}

func WithDecoder(d Decoder) Option {
	return func(o *options) error {
		if o.Decoder != nil {
			return errors.New("decoder is already set")
		}
		if d == nil {
			return errors.New("cannot use a <nil> decoder")
		}
		o.Decoder = d
		return nil
	}
}

func WithEncoder(e Encoder) Option {
	return func(o *options) error {
		if o.Encoder != nil {
			return errors.New("encoder is already set")
		}
		if e == nil {
			return errors.New("cannot use a <nil> encoder")
		}
		o.Encoder = e
		return nil
	}
}
