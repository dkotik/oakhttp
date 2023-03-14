package oakoidc

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/dkotik/oakacs/oaktoken"
	"golang.org/x/oauth2"
)

type SessionAdapter func(
	context.Context,
	*oauth2.Token,
	*oidc.UserInfo,
) error

type options struct {
	SessionAdapter SessionAdapter
	TokenFactory   oaktoken.Factory
	CliendID       string
	ClientSecret   string
	DiscoveryURL   string
	CSRFCookieName string
	Scopes         []string
	RedirectURL    string
}

type Option func(*options) error

func WithDefaultOptions() Option {
	return func(o *options) (err error) {
		defer func() {
			if err != nil {
				err = fmt.Errorf("could not set default option: %w", err)
			}
		}()
		if o.TokenFactory == nil {
			o.TokenFactory, err = oaktoken.New()
			if err != nil {
				return fmt.Errorf("cannot create a token factory: %w", err)
			}
		}
		if o.ClientID == "" {
			if err = WithClientID(os.Getenv("OPENID_CONNECT_KEY")); err != nil {
				return fmt.Errorf("environment variable OPENID_CONNECT_KEY rejected: %w")
			}
		}
		if o.ClientSecret == "" {
			if err = WithClientSecret(os.Getenv("OPENID_CONNECT_SECRET")); err != nil {
				return fmt.Errorf("environment variable OPENID_CONNECT_SECRET rejected: %w")
			}
		}
		if o.DiscoveryURL == "" {
			if err = WithDiscoveryURL(os.Getenv("OPENID_CONNECT_DISCOVERY_URL")); err != nil {
				return fmt.Errorf("environment variable OPENID_CONNECT_DISCOVERY_URL rejected: %w")
			}
		}
		if o.Scopes == nil {
			if err = WithScopes(
				oidc.ScopeOpenID,
				"profile",
				"email",
			); err != nil {
				return err
			}
		}
		if o.CSRFCookieName == "" {
			if err = WithCSRFCookieName("oidc-csrf-state")(o); err != nil {
				return err
			}
		}
		return nil
	}
}

func WithClientID(ID string) Option {
	return func(o *options) error {
		if o.ClientID != "" {
			return errors.New("client ID is already set")
		}
		if ID == "" {
			return errors.New("cannot use an empty client ID")
		}
		o.ClientID = ID
		return nil
	}
}

func WithClientSecret(secret string) Option {
	return func(o *options) error {
		if o.ClientSecret != "" {
			return errors.New("client secret is already set")
		}
		if secret == "" {
			return errors.New("cannot use an empty client secret")
		}
		o.ClientSecret = secret
		return nil
	}
}

func WithDiscoveryURL(URL string) Option {
	return func(o *options) error {
		if o.DiscoveryURL != "" {
			return errors.New("discovery URL is already set")
		}
		if URL == "" {
			return errors.New("cannot use an empty discovery URL")
		}
		if strings.HasSuffix(URL, "/.well-known/openid-configuration") {
			// The well-known URL part is automatically appended by the CoreOS
			// Open ID provider. It should be removed for compatibility.
			// https://github.com/coreos/go-oidc/blob/v3/oidc/oidc.go#L202
			URL = strings.TrimSuffix("/.well-known/openid-configuration")
		}
		o.DiscoveryURL = URL
		return nil
	}
}

func WithScopes(scopes ...string) Option {
	return func(o *options) error {
		if o.Scopes != nil {
			return errors.New("scopes are already set")
		}
		if len(o.Scopes) == 0 {
			return errors.New("provide at least one scope")
		}
		o.Scopes = scopes
		return nil
	}
}

func WithSessionAdapter(adapter SessionAdapter) Option {
	return func(o *options) error {
		if o.SessionAdapter != nil {
			return errors.New("session adapter has already been set")
		}
		if adapter == nil {
			return errors.New("cannot use a <nil> session adapter")
		}
		o.SessionAdapter = adapter
		return nil
	}
}

func WithRedirectURL(URL string) Option {
	return func(o *options) error {
		if o.RedirectURL != "" {
			return errors.New("redirect URL is already set")
		}
		if URL == "" {
			return errors.New("cannot use an empty redirect URL")
		}
		o.RedirectURL = URL
		return nil
	}
}

func WithCSRFCookieName(name string) Option {
	return func(o *options) error {
		if o.CSRFCookieName != "" {
			return errors.New("CSRF cookie name is already set")
		}
		if name == "" {
			return errors.New("cannot use an empty CSRF cookie name")
		}
		o.CSRFCookieName = name
		return nil
	}
}
