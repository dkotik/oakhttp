package oakrouter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dkotik/oakacs/oakhttp"
)

type options struct {
	PathPrefix               string
	CutPathPrefixFromRequest bool
	Routes                   map[string]oakhttp.Handler
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
