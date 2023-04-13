package oakrouter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dkotik/oakacs/oakhttp"
)

type options struct {
	PathPrefix               string
	TrailingSlashRedirects   bool
	CutPathPrefixFromRequest bool
	Routes                   map[string]oakhttp.Handler
	Middleware               []oakhttp.Middleware
}

type Option func(*options) error

func WithPathPrefix(p string) Option {
	return func(o *options) error {
		if o.PathPrefix != "" {
			return errors.New("path prefix is already set")
		}
		if p == "" {
			return errors.New("cannot use an empty path prefix")
		}
		if !strings.HasPrefix(p, "/") {
			p = "/" + p
		}
		o.PathPrefix = p
		return nil
	}
}

func WithoutTrailingSlashRedirects() Option {
	return func(o *options) error {
		if !o.TrailingSlashRedirects {
			return errors.New("trailing slash redirecting is already turned off")
		}
		o.TrailingSlashRedirects = false
		return nil
	}
}

func WithPathPrefixCutFromRequest() Option {
	return func(o *options) error {
		if o.CutPathPrefixFromRequest {
			return errors.New("path prefix removal is already set")
		}
		o.CutPathPrefixFromRequest = true
		return nil
	}
}

func WithRoute(path string, h oakhttp.Handler) Option {
	return func(o *options) error {
		if path == "" {
			return errors.New("cannot use an empty path")
		}
		if strings.HasSuffix(path, "/") {
			return errors.New("route path should not end with a trailing slash")
		}
		if h == nil {
			return errors.New("cannot use a <nil> handler")
		}
		if _, ok := o.Routes[path]; ok {
			return fmt.Errorf("route %q is already set", path)
		}
		o.Routes[path] = h
		return nil
	}
}

// WithMiddleware wraps handler with [oakhttp.Middleware]. This option can be used multiple times. The final middleware chain is applied in reverse order to preserve logical order. For example, WithMiddleware(first, second, third) provides a handler which is equivalent to manually calling first(second(third(handler))) without this option.
func WithMiddleware(mw ...oakhttp.Middleware) Option {
	return func(o *options) error {
		for _, middleware := range mw {
			if middleware == nil {
				return errors.New("cannot use a <nil> oakhttp.Middleware")
			}
		}
		o.Middleware = append(o.Middleware, mw...)
		return nil
	}
}

// func WithRateLimiterOptions(withOptions ...ratelimiter.Option) Option {
// 	return func(o *options) error {
// 		if o.rateLimiterOptions != nil {
// 			return errors.New("rate limiter options are already set")
// 		}
// 		for _, option := range withOptions {
// 			if option == nil {
// 				return errors.New("cannot use a <nil> rate limiter option")
// 			}
// 		}
// 		o.rateLimiterOptions = withOptions
// 		return nil
// 	}
// }
