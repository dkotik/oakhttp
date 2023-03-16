package oakhttp

import "net/http"

type Handler func(w http.ResponseWriter, r *http.Request) error

type Middleware func(Handler) Handler

func AdaptHandler(
	usingDomainAdaptor *DomainAdaptor,
	handler Handler,
	middleware ...Middleware,
) http.Handler {
	return usingDomainAdaptor.ApplyMiddleware(
		handler,
		middleware...,
	)
}
