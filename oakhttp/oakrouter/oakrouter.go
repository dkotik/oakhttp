/*
Package oakrouter provides an ugly [oakhttp.Handler] multiplexer as a hashmap router paired with a [http.Method] switch.

The pairing is entirely sufficient for building small services, assuming that most endpoints should extract request values from request.Body only without touching URL.Path or the HTTP header.
*/
package oakrouter

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dkotik/oakacs/oakhttp"
)

func New(withOptions ...Option) (oakhttp.Handler, error) {
	o := &options{
		Routes: make(map[string]oakhttp.Handler),
	}

	for _, option := range append(
		withOptions,
		func(o *options) error { // validate
			if len(o.Routes) == 0 {
				return errors.New("cannot use 0 routing paths")
			}
			return nil
		},
	) {
		if err := option(o); err != nil {
			return nil, fmt.Errorf("cannot create OakRouter: %w", err)
		}
	}

	routes := make(map[string]oakhttp.Handler)
	for p, handler := range o.Routes {
		routes[o.PathPrefix+p] = handler
	}

	if o.CutPathPrefixFromRequest {
		prefix := o.PathPrefix
		return func(w http.ResponseWriter, r *http.Request) error {
			h, ok := routes[r.URL.Path]
			if !ok {
				return oakhttp.NewNotFoundError(r.URL.Path)
			}
			r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix) // cut
			return h(w, r)
		}, nil
	}

	handler := func(w http.ResponseWriter, r *http.Request) error {
		h, ok := routes[r.URL.Path]
		if !ok {
			return oakhttp.NewNotFoundError(r.URL.Path)
		}
		return h(w, r)
	}
	for i := len(o.Middleware) - 1; i > 0; i-- {
		handler = o.Middleware[i](handler)
	}

	return handler, nil
}
