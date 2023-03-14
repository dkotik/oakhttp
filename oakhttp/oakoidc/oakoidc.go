/*
Package oakoidc allows logging in using Open ID Connect protocol extension over OAuth2.
*/
package oakoidc

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/dkotik/oakacs/oakhttp"
)

func New(withOptions ...Option) (begin, callback oakhttp.Handler, err error) {
	o := &options{}
	for _, option := range append(
		withOptions,
		WithDefaultOptions(),
		func(o *options) error { //validate
			if o.RedirectURL == "" {
				return errors.New("redirect URL is required")
			}
			if o.SessionAdapter == nil {
				return errors.New("session adapter is required")
			}
			return nil
		},
	) {
		if err = option(o); err != nil {
			return nil, nil, fmt.Errorf("cannot initialize OIDC provider: %w", err)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) error {
			io.WriteString(w, "begin test")
		}, func(w http.ResponseWriter, r *http.Request) error {
			io.WriteString(w, "callback test")
		}, nil
}
