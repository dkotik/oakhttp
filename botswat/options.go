package botswat

import (
	"errors"
	"net/http"

	"github.com/dkotik/oakhttp"
)

const (
	DefaultCookieName = "HumanityToken"
	DefaultHeaderName = "HumanityToken"
)

type HumanityTokenExtractor func(
	r *http.Request,
) (
	clientResponseToken string,
	err error,
)

type options struct {
	ErrorHandler           oakhttp.ErrorHandler
	Verifier               Verifier
	HumanityTokenExtractor HumanityTokenExtractor
	Cache                  Cache
}

type Option func(*options) error

func WithErrorHandler(eh oakhttp.ErrorHandler) Option {
	return func(o *options) error {
		if o.ErrorHandler != nil {
			return errors.New("error handler is already set")
		}
		if eh == nil {
			return errors.New("nil error handler")
		}
		o.ErrorHandler = eh
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

func WithHumanityTokenExtractor(e HumanityTokenExtractor) Option {
	return func(o *options) error {
		if o.HumanityTokenExtractor != nil {
			return errors.New("response extractor is already set")
		}
		if e == nil {
			return errors.New("cannot use a <nil> response extractor")
		}
		o.HumanityTokenExtractor = e
		return nil
	}
}

func WithCookieHumanityTokenExtractor(name string) Option {
	return func(o *options) error {
		if name == "" {
			return errors.New("cannot use an empty cookie name")
		}
		return WithHumanityTokenExtractor(
			func(r *http.Request) (string, error) {
				c, err := r.Cookie(name)
				if err != nil {
					return "", err
				}
				return c.Value, nil
			},
		)(o)
	}
}

func WithDefaultCookieHumanityTokenExtractor() Option {
	return WithCookieHumanityTokenExtractor(DefaultCookieName)
}

func WithHeaderHumanityTokenExtractor(name string) Option {
	return func(o *options) error {
		if name == "" {
			return errors.New("cannot use an empty HTTP header name")
		}
		return WithHumanityTokenExtractor(
			func(r *http.Request) (string, error) {
				return r.Header.Get(name), nil
			},
		)(o)
	}
}

func WithDefaultHeaderHumanityTokenExtractor() Option {
	return WithHeaderHumanityTokenExtractor(DefaultHeaderName)
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

// func WithCacheAdaptor(adaptor CacheAdaptor) Option {
// 	if adaptor == nil {
// 		return func(o *options) error {
// 			return errors.New("cannot use a <nil> cache adaptor")
// 		}
// 	}
//
// 	return WithCache(
// 		func(v Verifier) Verifier {
// 			return func(
// 				ctx context.Context,
// 				clientResponseToken string,
// 				clientIPAddress string,
// 			) (string, error) {
// 				key := []byte(clientIPAddress + "^" + clientResponseToken)
// 				userData, err := adaptor.Get(ctx, key)
// 				if err == nil {
// 					return string(userData), nil // cache hit
// 				}
// 				freshUserData, err := v(ctx, clientResponseToken, clientIPAddress)
// 				if err != nil {
// 					return "", err
// 				}
// 				if err = adaptor.Set(ctx, key, []byte(freshUserData)); err != nil {
// 					return "", err
// 				}
// 				return freshUserData, err
// 			}
// 		},
// 	)
// }
