package oakswitch

import (
	"errors"
	"fmt"

	"github.com/dkotik/oakacs/oakhttp"
)

type options struct {
	Get    oakhttp.Handler
	Post   oakhttp.Handler
	Put    oakhttp.Handler
	Delete oakhttp.Handler
	Patch  oakhttp.Handler
}

type Option func(*options) error

func WithGet(h oakhttp.Handler, err error) Option {
	return func(o *options) error {
		if err != nil {
			return fmt.Errorf("invalid get request handler: %w", err)
		}
		if h == nil {
			return errors.New("cannot use a <nil> get request handler")
		}
		if o.Get != nil {
			return errors.New("get request handler is already set")
		}
		o.Get = h
		return nil
	}
}

func WithPost(h oakhttp.Handler, err error) Option {
	return func(o *options) error {
		if err != nil {
			return fmt.Errorf("invalid post request handler: %w", err)
		}
		if h == nil {
			return errors.New("cannot use a <nil> post request handler")
		}
		if o.Post != nil {
			return errors.New("post request handler is already set")
		}
		o.Post = h
		return nil
	}
}

func WithPut(h oakhttp.Handler, err error) Option {
	return func(o *options) error {
		if err != nil {
			return fmt.Errorf("invalid put request handler: %w", err)
		}
		if h == nil {
			return errors.New("cannot use a <nil> put request handler")
		}
		if o.Put != nil {
			return errors.New("put request handler is already set")
		}
		o.Put = h
		return nil
	}
}

func WithDelete(h oakhttp.Handler, err error) Option {
	return func(o *options) error {
		if err != nil {
			return fmt.Errorf("invalid delete request handler: %w", err)
		}
		if h == nil {
			return errors.New("cannot use a <nil> delete request handler")
		}
		if o.Delete != nil {
			return errors.New("delete request handler is already set")
		}
		o.Delete = h
		return nil
	}
}

func WithPatch(h oakhttp.Handler, err error) Option {
	return func(o *options) error {
		if err != nil {
			return fmt.Errorf("invalid patch request handler: %w", err)
		}
		if h == nil {
			return errors.New("cannot use a <nil> patch request handler")
		}
		if o.Patch != nil {
			return errors.New("patch request handler is already set")
		}
		o.Patch = h
		return nil
	}
}
