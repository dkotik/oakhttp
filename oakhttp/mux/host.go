package mux

import (
	"fmt"
	"net/http"

	"github.com/tigerperformanceinstitute/tprograms/service/oakhttp"
)

func NewHostMux(withOptions ...HostMuxOption) oakhttp.HandlerFunc {
	mux := make(map[string]oakhttp.Handler)

	var err error
	for _, option := range withOptions {
		if err = option(mux); err != nil {
			panic(fmt.Errorf("cannot initialize host multiplexer: %w", err))
		}
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		name := r.URL.Hostname()
		h, ok := mux[name]
		if !ok {
			return oakhttp.NewNotFoundError(fmt.Errorf("host multiplexer is not aware of host %q", name))
		}
		return h.ServeHyperText(w, r)
	}
}

type HostMuxOption func(map[string]oakhttp.Handler) error

func WithHost(name string, h oakhttp.Handler) HostMuxOption {
	return func(m map[string]oakhttp.Handler) error {
		if h == nil {
			return fmt.Errorf("cannot set an empty handler for host %q", name)
		}
		if _, ok := m[name]; ok {
			return fmt.Errorf("handler for host %q is already set", name)
		}
		m[name] = h
		return nil
	}
}
