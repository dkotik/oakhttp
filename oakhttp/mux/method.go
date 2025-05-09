package mux

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/tigerperformanceinstitute/tprograms/service/oakhttp"
)

func NewMethodMux(withOptions ...MethodMuxOption) oakhttp.HandlerFunc {
	o := &muxByMethodOptions{
		allowed: http.MethodOptions,
	}
	for _, option := range withOptions {
		if err := option(o); err != nil {
			panic(fmt.Errorf("cannot initialize method switch: %w", err))
		}
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		switch r.Method {
		case http.MethodGet, http.MethodHead:
			if o.Get != nil {
				return o.Get.ServeHyperText(w, r)
			}
		case http.MethodPost:
			if o.Post != nil {
				return o.Post.ServeHyperText(w, r)
			}
		case http.MethodPut:
			if o.Put != nil {
				return o.Put.ServeHyperText(w, r)
			}
		case http.MethodPatch:
			if o.Patch != nil {
				return o.Patch.ServeHyperText(w, r)
			}
		case http.MethodDelete:
			if o.Delete != nil {
				return o.Delete.ServeHyperText(w, r)
			}
		case http.MethodOptions:
			w.Header().Set("Allow", o.allowed)
			w.WriteHeader(http.StatusOK)
			return nil
		}
		return oakhttp.NewMethodNotAllowedError(r.Method)
	}
}

type muxByMethodOptions struct {
	Get     oakhttp.Handler
	Post    oakhttp.Handler
	Put     oakhttp.Handler
	Delete  oakhttp.Handler
	Patch   oakhttp.Handler
	allowed string
}

type MethodMuxOption func(*muxByMethodOptions) error

func WithGet(h oakhttp.Handler) MethodMuxOption {
	return func(o *muxByMethodOptions) error {
		if h == nil {
			return errors.New("cannot use a <nil> get request handler")
		}
		if o.Get != nil {
			return errors.New("get request handler is already set")
		}
		o.Get = h
		o.allowed += "," + http.MethodGet
		return nil
	}
}

func WithPost(h oakhttp.Handler) MethodMuxOption {
	return func(o *muxByMethodOptions) error {
		if h == nil {
			return errors.New("cannot use a <nil> post request handler")
		}
		if o.Post != nil {
			return errors.New("post request handler is already set")
		}
		o.Post = h
		o.allowed += "," + http.MethodPost
		return nil
	}
}

func WithPut(h oakhttp.Handler) MethodMuxOption {
	return func(o *muxByMethodOptions) error {
		if h == nil {
			return errors.New("cannot use a <nil> put request handler")
		}
		if o.Put != nil {
			return errors.New("put request handler is already set")
		}
		o.Put = h
		o.allowed += "," + http.MethodPut
		return nil
	}
}

func WithPatch(h oakhttp.Handler) MethodMuxOption {
	return func(o *muxByMethodOptions) error {
		if h == nil {
			return errors.New("cannot use a <nil> patch request handler")
		}
		if o.Patch != nil {
			return errors.New("patch request handler is already set")
		}
		o.Patch = h
		o.allowed += "," + http.MethodPatch
		return nil
	}
}

func WithDelete(h oakhttp.Handler) MethodMuxOption {
	return func(o *muxByMethodOptions) error {
		if h == nil {
			return errors.New("cannot use a <nil> delete request handler")
		}
		if o.Delete != nil {
			return errors.New("delete request handler is already set")
		}
		o.Delete = h
		o.allowed += "," + http.MethodDelete
		return nil
	}
}
