package turnstile

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const DefaultRetention = time.Hour * 24 * 7

type options struct {
	HTTPClient *http.Client
	Endpoint   string
	Hostname   string
	SecretKey  string
}

type (
	Option func(*options) error
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
