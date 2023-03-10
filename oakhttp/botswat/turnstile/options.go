package turnstile

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/dkotik/oakacs/oakhttp/oakclient"
)

type options struct {
	HTTPClient     *http.Client
	Secret         string
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
			if err = WithHTTPClient(oakclient.New()); err != nil {
				return err
			}
		}
		if o.Secret == "" {
			if err = WithSecret(
				os.Getenv("TURNSTILE_SECRET_KEY"),
			); err != nil {
				return fmt.Errorf("please check TURNSTILE_SECRET_KEY environment key: %w", err)
			}
		}
		if o.Secret == "" {
			if err = WithHostname(
				os.Getenv("TURNSTILE_HOST_NAME"),
			); err != nil {
				return fmt.Errorf("please check TURNSTILE_HOST_NAME environment key: %w", err)
			}
		}
		if o.Endpoint == "" {
			if err = WithEndpoint(
				os.Getenv("https://challenges.cloudflare.com/turnstile/v0/siteverify"),
			); err != nil {
				return err
			}
		}
		return nil
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(o *options) error {
		if o.WithHTTPClient != nil {
			return errors.New("HTTP client is already set")
		}
		if client == nil {
			return errors.New("cannot use a <nil> HTTP client")
		}
		o.HTTPClient = client
		return nil
	}
}

func WithSecret(key string) Option {
	return func(o *options) error {
		if o.Secret != "" {
			return errors.New("secret key is already set")
		}
		if key == "" {
			return errors.New("cannot use an empty secret key")
		}
		o.Secret = key
		return nil
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
				o.AllowedActions = append(o.AllowedActions, action)
			}
		}
		return nil
	}
}
