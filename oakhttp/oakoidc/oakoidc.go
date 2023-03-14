/*
Package oakoidc allows logging in using Open ID Connect protocol extension over OAuth2.
*/
package oakoidc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/dkotik/oakacs/oakhttp"
	"golang.org/x/oauth2"
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

	provider, err := oidc.NewProvider(context.Background(), o.DiscoveryURL)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot initialize OIDC provider: %w", err)
	}

	oauth := oauth2.Config{
		ClientID:     o.ClientID,
		ClientSecret: o.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  o.RedirectURL,
		Scopes:       o.Scopes,
	}

	return func(w http.ResponseWriter, r *http.Request) error {
			CSRFToken, err := o.TokenFactory()
			if err != nil {
				return fmt.Errorf("cannot generate a CSRF token: %w", err)
			}
			http.SetCookie(w, &http.Cookie{
				Name:     o.CookieName,
				Value:    CSRFToken,
				MaxAge:   int((time.Minute * 5).Seconds()),
				Secure:   r.TLS != nil,
				HttpOnly: true,
			})

			http.Redirect(w, r, oauth.AuthCodeURL(CSRFToken), http.StatusFound)
			return nil
		}, func(w http.ResponseWriter, r *http.Request) error {
			CSRFToken, err := r.Cookie(o.CookieName)
			if err != nil {
				return fmt.Errorf("CSRF state token cannot be recovered: %w", err)
			}
			q := r.URL.Query()
			if q.Get("state") != CSRFToken.Value {
				return errors.New("CSRF state token does not match the last set token")
			}

			ctx := r.Context()
			oauth2Token, err := oauth.Exchange(ctx, q.Get("code"))
			if err != nil {
				return fmt.Errorf("token exchange failed: %w", err)
			}
			userInfo, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(oauth2Token))
			if err != nil {
				return fmt.Errorf("could not retrieve user profile: %w", err)
			}

			return o.SessionAdapter(ctx, oauth2Token, userInfo)
		}, nil
}
