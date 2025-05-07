package cors

import (
	"errors"
	"fmt"
	"net/textproto"
	"strings"
	"time"
)

type OriginFilter func(s string) bool

type options struct {
	allowedOrigins   []string
	allowedMethods   []string
	allowedHeaders   []string
	exposedHeaders   []string
	maxAge           time.Duration
	allowCredentials *bool
	filters          []OriginFilter
}

type Option func(*options) error

func WithOptions(given ...Option) Option {
	return func(o *options) (err error) {
		for _, option := range given {
			if err = option(o); err != nil {
				return err
			}
		}
		return nil
	}
}

func WithDefaultMaxAgeOfOneWeek() Option {
	return func(o *options) error {
		if o.maxAge > 0 {
			return nil
		}
		return WithMaxAge(time.Hour * 24 * 7)(o)
	}
}

func WithDefaultMethodsGetPostHead() Option {
	return func(o *options) error {
		if len(o.allowedMethods) > 0 {
			return nil
		}
		return WithOptions(
			WithMethod("GET"),
			WithMethod("POST"),
			WithMethod("HEAD"),
		)(o)
	}
}

func WithOrigin(s string) Option {
	return func(o *options) error {
		s = strings.TrimSpace(s)
		if s == "" {
			return errors.New("cannot use an empty origin")
		}
		for _, allowed := range o.allowedOrigins {
			if allowed == s {
				return fmt.Errorf("origin %q is allowed twice", s)
			}
		}
		o.allowedOrigins = append(o.allowedOrigins, s)
		return nil
	}
}

func WithMethod(s string) Option {
	return func(o *options) error {
		s = strings.ToUpper(strings.TrimSpace(s))
		if s == "" {
			return errors.New("cannot use an empty method name")
		}
		for _, allowed := range o.allowedMethods {
			if allowed == s {
				return fmt.Errorf("method %q is allowed twice", s)
			}
		}
		o.allowedMethods = append(o.allowedMethods, s)
		return nil
	}
}

func WithHeader(s string) Option {
	return func(o *options) error {
		s = textproto.CanonicalMIMEHeaderKey(strings.TrimSpace(s))
		if s == "" {
			return errors.New("cannot use an empty header name")
		}
		for _, allowed := range o.allowedHeaders {
			if allowed == s {
				return fmt.Errorf("header %q is allowed twice", s)
			}
		}
		o.allowedHeaders = append(o.allowedHeaders, s)
		return nil
	}
}

func WithExposedHeader(s string) Option {
	return func(o *options) error {
		s = textproto.CanonicalMIMEHeaderKey(strings.TrimSpace(s))
		if s == "" {
			return errors.New("cannot use an empty header name")
		}
		for _, allowed := range o.exposedHeaders {
			if allowed == s {
				return fmt.Errorf("header %q is exposed twice", s)
			}
		}
		o.exposedHeaders = append(o.exposedHeaders, s)
		return nil
	}
}

func WithMaxAge(d time.Duration) Option {
	return func(o *options) error {
		if d < time.Second {
			return errors.New("maximum age cannot be less than 1 second")
		}
		if d > time.Hour*24*30*3 {
			return errors.New("maximum age longer than three months is insecure")
		}
		if o.maxAge != 0 {
			return errors.New("maximum age is already set")
		}
		o.maxAge = d
		return nil
	}
}

func WithCredentials() Option {
	return func(o *options) error {
		if o.allowCredentials != nil {
			return errors.New("credentials are already allowed")
		}
		allow := true
		o.allowCredentials = &allow
		return nil
	}
}

func WithOriginFilter(f OriginFilter) Option {
	return func(o *options) error {
		if f == nil {
			return errors.New("cannot use a <nil> origin filter")
		}
		o.filters = append(o.filters, f)
		return nil
	}
}
