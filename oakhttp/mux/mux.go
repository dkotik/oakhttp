package mux

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tigerperformanceinstitute/tprograms/service/oakhttp"
)

type mux struct {
	handlers map[*Pattern]oakhttp.Handler
	routes   map[string]*Pattern
	tree     *node
}

func New(withOptions ...Option) oakhttp.Handler {
	o := &options{
		handlers: make(map[*Pattern]oakhttp.Handler, 0),
		routes:   make(map[string]*Pattern),
		tree:     &node{},
	}

	var err error
	for _, option := range withOptions {
		if err = option(o); err != nil {
			panic(fmt.Errorf("cannot initialize a path multiplexer: %w", err))
		}
	}
	return oakhttp.ApplyMiddleware(&mux{
		handlers: o.handlers,
		routes:   o.routes,
		tree:     o.tree,
	}, o.middleware)
}

func (m *mux) ServeHyperText(w http.ResponseWriter, r *http.Request) error {
	pattern, matches := m.tree.matchPath(r.URL.Path, nil)
	// fmt.Println(pattern, matches)
	// pattern, matches := m.tree.matchPath(r.URL.Path)
	// if pat != nil { // TODO: this all needs to be injected into context.
	//   if s.nobind {
	//     return pat, nil
	//   }
	//   return pat, pat.bind(matches)
	// }

	handler, ok := m.handlers[pattern]
	if !ok {
		return ErrNoRouteMatched
	}
	return handler.ServeHyperText(w, r.WithContext(
		context.WithValue(r.Context(), muxContextKey, &RoutingContext{
			mux:     m,
			matched: pattern,
			matches: matches,
		}),
	))
}

type noRouteMatchedError struct{}

func (e *noRouteMatchedError) Error() string {
	return "Not Found"
}

func (e *noRouteMatchedError) HyperTextStatusCode() int {
	return http.StatusNotFound
}

var ErrNoRouteMatched oakhttp.Error = &noRouteMatchedError{}
