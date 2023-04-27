package turnstile

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dkotik/oakacs/oakhttp/oakclient"
)

type options struct {
	HTTPClient     *http.Client
	SecretKey      string
	Endpoint       string
	Hostname       string
	AllowedActions []string
}

func (o *options) IsAllowedAction(a string) bool {
	for _, action := range o.AllowedActions {
		if action == a {
			return true
		}
	}
	return false
}

type Option func(*options) error

func WithDefaultOptions() Option {
	return func(o *options) (err error) {
		if o.HTTPClient == nil {
			client, err := oakclient.New()
			if err != nil {
				return err
			}
			if err = WithHTTPClient(client)(o); err != nil {
				return err
			}
		}
		if o.Hostname == "" {
			if err = WithDefaultHostname()(o); err != nil {
				return err
			}
		}
		if o.SecretKey == "" {
			if err = WithDefaultSecretKey()(o); err != nil {
				return err
			}
		}
		if o.Endpoint == "" {
			if err = WithDefaultEndpoint()(o); err != nil {
				return err
			}
		}
		return nil
	}
}

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
		secret := strings.TrimSpace(os.Getenv(variableName))
		if secret == "" {
			return fmt.Errorf("cannot get secret from environment: variable %q is not set", variableName)
		}
		return WithSecretKey(secret)(o)
	}
}

func WithDefaultSecretKey() Option {
	return WithSecretKeyFromEnvironment("TURNSTILE_SECRET_KEY")
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
	return WithEndpoint("https://challenges.cloudflare.com/turnstile/v0/siteverify")
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
	return WithHostnameFromEnvironment("TURNSTILE_HOST_NAME")
}

func WithAllowedActions(actions ...string) Option {
	return func(o *options) error {
		for _, action := range actions {
			if action == "" {
				return errors.New("cannot use an empty allowed action")
			}
			for _, existing := range o.AllowedActions {
				if action == existing {
					return fmt.Errorf("allowed action %q has already been set", action)
				}
			}
			o.AllowedActions = append(o.AllowedActions, action)
		}
		return nil
	}
}
