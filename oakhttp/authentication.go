package oakhttp

import (
	"errors"
	"net/http"

	"github.com/tigerperformanceinstitute/tprograms/service/oakhttp/store"
)

type Authenticator interface {
	// IsAllowed must return an [Error] with [http.StatusForbidden] status code, if operation succeeded but access was denied.
	IsAllowed(*http.Request) (*http.Request, error)
}

type knownRequestAuthenticator struct {
	extractor StringExtractor
	store     store.KeyValue
}

func NewKnownRequestAuthenticator(
	extractor StringExtractor,
	store store.KeyValue,
) Authenticator {
	if extractor == nil {
		panic("cannot use a <nil> string extractor")
	}
	if store == nil {
		panic("cannot use a <nil> key value store")
	}
	return &knownRequestAuthenticator{
		extractor: extractor,
		store:     store,
	}
}

func (k *knownRequestAuthenticator) IsAllowed(r *http.Request) (*http.Request, error) {
	requestTag, err := k.extractor(r)
	if err != nil {
		return nil, NewAccessDeniedError(err)
	}
	if _, err = k.store.Get(r.Context(), []byte(requestTag)); err != nil {
		if errors.Is(err, store.ErrValueNotFound) {
			return nil, NewAccessDeniedError(err)
		}
		return nil, err
	}
	return r, nil
}

type sequentialAuthenticator []Authenticator

func NewSequentialAuthenticator(using ...Authenticator) Authenticator {
	for _, f := range using {
		if f == nil {
			panic("cannot use a <nil> filter")
		}
	}
	if len(using) == 0 {
		panic("at least one authenticator is required")
	}
	return sequentialAuthenticator(using)
}

func (s sequentialAuthenticator) IsAllowed(r *http.Request) (updated *http.Request, err error) {
	for _, f := range s {
		updated, err = f.IsAllowed(r)
		if err == nil {
			return updated, nil
		} else if !IsAccessDeniedError(err) {
			return nil, err
		}
	}
	return nil, NewAccessDeniedError(errors.New("every authenticator failed"))
}

type gate struct {
	authenticator Authenticator
	allowed       Handler
	forbidden     Handler
}

func NewGate(a Authenticator, forbidden Handler) Middleware {
	if a == nil {
		panic("cannot use a <nil> authenticator")
	}
	if forbidden == nil {
		panic("cannot use a <nil> forbidden handler")
	}
	return func(next Handler) Handler {
		if next == nil {
			panic("cannot use a <nil> handler")
		}
		return &gate{
			authenticator: a,
			allowed:       next,
			forbidden:     forbidden,
		}
	}
}

func (g *gate) ServeHyperText(
	w http.ResponseWriter,
	r *http.Request,
) (err error) {
	r, err = g.authenticator.IsAllowed(r)
	if err != nil {
		if IsAccessDeniedError(err) {
			return g.forbidden.ServeHyperText(w, r)
		}
		return err
	}
	return g.allowed.ServeHyperText(w, r)
}
