package turnstile

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/tigerperformanceinstitute/tprograms/service/oakhttp"
	"github.com/tigerperformanceinstitute/tprograms/service/oakhttp/store"
)

const DefaultRetention = time.Hour * 24 * 7

type options struct {
	HTTPClient *http.Client
	Endpoint   string
	Hostname   string
	SecretKey  string
}

type authenticatorOptions struct {
	Client     *Turnstile
	CookieName string
	SiteAction string
	Extractor  oakhttp.StringExtractor
	Store      store.KeyKeyValue
}

func newAuthenticatorOptions(withOptions []AuthenticatorOption) (*authenticatorOptions, error) {
	o := &authenticatorOptions{}
	var err error
	for _, option := range append(
		withOptions,
		WithDefaultClient(),
		WithDefaultSiteAction(),
		WithDefaultExtractor(),
		WithDefaultStore(),
		func(o *authenticatorOptions) error {
			if o.Client == nil {
				return errors.New("Turnstile client is required")
			}
			if o.SiteAction == "" {
				return errors.New("site action is required")
			}
			if o.Extractor == nil {
				return errors.New("extractor or cookie name are required")
			}
			if o.Store == nil {
				return errors.New("store is required")
			}
			return nil
		},
	) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("cannot create Turnstile authenticator: %w", err)
		}
	}
	return o, nil
}

type (
	Option              func(*options) error
	AuthenticatorOption func(*authenticatorOptions) error
)

func WithHTTPClient(client *http.Client) Option {
	return func(o *options) error {
		if o.HTTPClient != nil {
			return errors.New("HTTP client is already set")
		}
		if client == nil {
			return errors.New("cannot use a <nil> HTTP client")
		}
		o.HTTPClient = client
		return nil
	}
}

func WithDefaultHTTPClient() Option {
	return func(o *options) error {
		if o.HTTPClient != nil {
			return nil
		}
		return WithHTTPClient(&http.Client{
			Timeout: time.Second * 2,
		})(o)
	}
}

func WithSecretKey(key string) Option {
	return func(o *options) error {
		if o.SecretKey != "" {
			return errors.New("secret key is already set")
		}
		if key == "" {
			return errors.New("cannot use an empty secret key")
		}
		o.SecretKey = key
		return nil
	}
}

func WithSecretKeyFromEnvironment(variableName string) Option {
	return func(o *options) error {
		key := strings.TrimSpace(os.Getenv(variableName))
		if key == "" {
			return fmt.Errorf("cannot get secret key from environment: variable %q is not set", variableName)
		}
		return WithSecretKey(key)(o)
	}
}

func WithDefaultSecretKey() Option {
	return func(o *options) error {
		if o.SecretKey != "" {
			return nil
		}
		return WithSecretKeyFromEnvironment("TURNSTILE_SECRET_KEY")(o)
	}
}

func WithEndpoint(URL string) Option {
	return func(o *options) error {
		if o.Endpoint != "" {
			return errors.New("endpoint is already set")
		}
		if URL == "" {
			return errors.New("cannot use an empty endpoint URL")
		}
		o.Endpoint = URL
		return nil
	}
}

func WithDefaultEndpoint() Option {
	return func(o *options) error {
		if o.Endpoint != "" {
			return nil
		}
		return WithEndpoint("https://challenges.cloudflare.com/turnstile/v0/siteverify")(o)
	}
}

func WithHostname(name string) Option {
	return func(o *options) error {
		if o.Hostname != "" {
			return errors.New("host name is already set")
		}
		if name == "" {
			return errors.New("cannot use an empty host name")
		}
		o.Hostname = name
		return nil
	}
}

func WithHostnameFromEnvironment(variableName string) Option {
	return func(o *options) error {
		hostname := strings.TrimSpace(os.Getenv(variableName))
		if hostname == "" {
			return fmt.Errorf("cannot get hostname from the environment: variable %q is not set", variableName)
		}
		return WithHostname(hostname)(o)
	}
}

func WithDefaultHostname() Option {
	return func(o *options) error {
		if o.Hostname != "" {
			return nil
		}
		return WithHostnameFromEnvironment("TURNSTILE_HOSTNAME")(o)
	}
}

func WithClient(client *Turnstile) AuthenticatorOption {
	return func(o *authenticatorOptions) error {
		if o.Client != nil {
			return errors.New("client is already set")
		}
		if client == nil {
			return errors.New("cannot use a <nil> client")
		}
		o.Client = client
		return nil
	}
}

func WithClientOptions(options ...Option) AuthenticatorOption {
	return func(o *authenticatorOptions) (err error) {
		if o.Client != nil {
			return errors.New("client is already set")
		}
		o.Client, err = New(options...)
		if err != nil {
			return err
		}
		return nil
	}
}

func WithDefaultClient() AuthenticatorOption {
	return func(o *authenticatorOptions) (err error) {
		if o.Client != nil {
			return nil
		}
		o.Client, err = New()
		if err != nil {
			return err
		}
		return nil
	}
}

func WithSiteAction(action string) AuthenticatorOption {
	return func(o *authenticatorOptions) error {
		if o.SiteAction != "" {
			return errors.New("site action is already set")
		}
		if action == "" {
			return errors.New("cannot use an empty site action")
		}
		o.SiteAction = action
		return nil
	}
}

func WithDefaultSiteAction() AuthenticatorOption {
	return func(o *authenticatorOptions) error {
		if o.SiteAction != "" {
			return nil
		}
		return WithSiteAction("view")(o)
	}
}

func WithExtractor(e oakhttp.StringExtractor) AuthenticatorOption {
	return func(o *authenticatorOptions) error {
		if o.Extractor != nil {
			return errors.New("extractor is already set")
		}
		if e == nil {
			return errors.New("cannot use a <nil> extractor")
		}
		o.Extractor = e
		return nil
	}
}

func WithCookieName(name string) AuthenticatorOption {
	return func(o *authenticatorOptions) error {
		if o.Extractor != nil {
			return errors.New("extractor is already set")
		}
		if name == "" {
			return errors.New("cannot use an empty cookie name")
		}
		o.CookieName = name
		return WithExtractor(func(r *http.Request) (string, error) {
			cookie, err := r.Cookie(name)
			if err != nil {
				return "", ErrNoCookie
			}
			return cookie.Value, nil
		})(o)
	}
}

func WithDefaultExtractor() AuthenticatorOption {
	return func(o *authenticatorOptions) error {
		if o.Extractor != nil {
			return nil
		}
		return WithCookieName("turnstile_" + o.SiteAction)(o)
	}
}

func WithStore(s store.KeyKeyValue) AuthenticatorOption {
	return func(o *authenticatorOptions) error {
		if o.Store != nil {
			return errors.New("data store is already set")
		}
		if s == nil {
			return errors.New("cannot use a <nil> data store")
		}
		o.Store = s
		return nil
	}
}

func WithDefaultStore() AuthenticatorOption {
	return func(o *authenticatorOptions) error {
		if o.Store != nil {
			return nil
		}
		store, err := store.NewMapKeyKeyValue(
			store.WithValueRetentionFor(DefaultRetention),
			store.WithMaximumValueCount(12000),
		)
		if err != nil {
			return err
		}
		return WithStore(store)(o)
	}
}
