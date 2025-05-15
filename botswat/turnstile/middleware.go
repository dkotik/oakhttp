package turnstile

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"text/template"

	"github.com/tigerperformanceinstitute/tprograms/service/oakhttp"
)

type Middleware struct {
	authenticator *authenticator
	next          oakhttp.Handler
	challenge     []byte
}

func NewMiddleware(
	withOptions ...MiddlewareOption,
) (
	oakhttp.Middleware, error,
) {
	o := &middlewareOptions{}
	var err error
	for _, option := range append(
		withOptions,
		WithDefaultAuthenticatorOptions(),
		WithDefaultTemplate(),
		func(o *middlewareOptions) (err error) { // render
			if o.templateOptions == nil {
				return nil
			}
			if o.templateOptions.SiteAction == "" {
				if err = WithTemplateCookieName(o.authenticatorOptions.SiteAction)(o.templateOptions); err != nil {
					return fmt.Errorf("cannot synchronize template and authenticator site action: %w", err)
				}
			} else if o.templateOptions.SiteAction != o.authenticatorOptions.SiteAction {
				return fmt.Errorf("authenticator cookie name does not match template site action: %q vs %q", o.authenticatorOptions.SiteAction, o.templateOptions.SiteAction)
			}

			if o.templateOptions.CookieName == "" {
				if err = WithTemplateCookieName(o.authenticatorOptions.CookieName)(o.templateOptions); err != nil {
					return fmt.Errorf("cannot synchronize template and authenticator cookie name: %w", err)
				}
			} else if o.templateOptions.CookieName != o.authenticatorOptions.CookieName {
				return fmt.Errorf("authenticator cookie name does not match template cookie name: %q vs %q", o.authenticatorOptions.CookieName, o.templateOptions.CookieName)
			}

			b := &bytes.Buffer{}
			t, err := template.New("challenge").Parse(templateHTML)
			if err != nil {
				return fmt.Errorf("failed to parse default challenge template: %w", err)
			}
			if err = t.Execute(b, o.templateOptions); err != nil {
				return fmt.Errorf("failed to render challenge template: %w", err)
			}
			o.challenge = b.Bytes()
			return nil
		},
	) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("cannot initialize Turnstile middleware: %w", err)
		}
	}

	challenge := o.challenge
	authenticator := &authenticator{
		client:     o.authenticatorOptions.Client,
		siteAction: o.authenticatorOptions.SiteAction,
		extractor:  o.authenticatorOptions.Extractor,
		store:      o.authenticatorOptions.Store,
		storeKey:   append([]byte("turnstile:"), []byte(o.authenticatorOptions.SiteAction)...),
	}
	return func(next oakhttp.Handler) oakhttp.Handler {
		if next == nil {
			panic("cannot use a <nil> handler")
		}
		return &Middleware{
			authenticator: authenticator,
			next:          next,
			challenge:     challenge,
		}
	}, nil
}

func (m *Middleware) ServeHyperText(
	w http.ResponseWriter,
	r *http.Request,
) (err error) {
	r, err = m.authenticator.IsAllowed(r)
	if err != nil {
		slog.Debug(
			"turnstile authentication failed",
			slog.Any("error", err),
		)
		w.WriteHeader(http.StatusForbidden)
		_, err = io.Copy(w, bytes.NewReader(m.challenge))
		return err
	}
	return m.next.ServeHyperText(w, r)
}

func (m *Middleware) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	if err := m.ServeHyperText(w, r); err != nil {
		slog.Debug(
			"turnstile middleware failed",
			slog.Any("error", err),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
