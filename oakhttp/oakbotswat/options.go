package oakbotswat

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dkotik/oakacs/oakhttp"
)

type options struct {
	Verifier          Verifier
	ResponseExtractor ResponseExtractor
	Cache             Cache
	Encoder           oakhttp.Encoder
}

type Option func(*options) error

func WithDefaultOptions() Option {
	return func(o *options) (err error) {
		defer func() {
			if err != nil {
				err = fmt.Errorf("cannot apply a default setting: %w", err)
			}
		}()

		if o.ResponseExtractor == nil {
			if err = WithCookieResponseExtractor("botswat_token")(o); err != nil {
				return err
			}
		}

		if o.Encoder == nil {
			o.Encoder = oakhttp.EncoderJSON
		}

		return nil
	}
}

func WithVerifier(v Verifier) Option {
	return func(o *options) error {
		if o.Verifier != nil {
			return errors.New("verifier is already set")
		}
		if v == nil {
			return errors.New("cannot use a <nil> verifier")
		}
		o.Verifier = v
		return nil
	}
}

func WithResponseExtractor(e ResponseExtractor) Option {
	return func(o *options) error {
		if o.ResponseExtractor != nil {
			return errors.New("response extractor is already set")
		}
		if e == nil {
			return errors.New("cannot use a <nil> response extractor")
		}
		o.ResponseExtractor = e
		return nil
	}
}

func WithCookieResponseExtractor(name string) Option {
	if name == "" {
		return func(o *options) error {
			return errors.New("cannot use an empty cookie name")
		}
	}

	return WithResponseExtractor(
		func(r *http.Request) (string, error) {
			c, err := r.Cookie(name)
			if err != nil {
				return "", err
			}
			return c.Value, nil
		},
	)
}

func WithEncoder(e oakhttp.Encoder) Option {
	return func(o *options) error {
		if o.Encoder != nil {
			return errors.New("request encoder is already set")
		}
		if e == nil {
			return errors.New("cannot use a <nil> request encoder")
		}
		o.Encoder = e
		return nil
	}
}

func WithCache(c Cache) Option {
	return func(o *options) error {
		if o.Cache != nil {
			return errors.New("Cache is already set")
		}
		if c == nil {
			return errors.New("cannot use a <nil> Cache")
		}
		o.Cache = c
		return nil
	}
}

func WithCacheAdaptor(adaptor CacheAdaptor) Option {
	if adaptor == nil {
		return func(o *options) error {
			return errors.New("cannot use a <nil> cache adaptor")
		}
	}

	return WithCache(
		func(v Verifier) Verifier {
			return func(
				ctx context.Context,
				clientResponseToken string,
				clientIPAddress string,
			) (string, error) {
				key := []byte(clientIPAddress + "^" + clientResponseToken)
				userData, err := adaptor.Get(ctx, key)
				if err == nil {
					return string(userData), nil // cache hit
				}
				freshUserData, err := v(ctx, clientResponseToken, clientIPAddress)
				if err != nil {
					return "", err
				}
				if err = adaptor.Set(ctx, key, []byte(freshUserData)); err != nil {
					return "", err
				}
				return freshUserData, err
			}
		},
	)
}
