package turnstile

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"

	_ "embed"

	"github.com/dkotik/oakacs/oakhttp"
	"github.com/dkotik/oakacs/oakhttp/botswat"
)

//go:embed gate.html
var templateAdjustments string

type gateOptions struct {
	*botswat.TemplateOptions
	SiteKey    string
	SiteAction string
}

type GateOption func(*gateOptions) error

func WithTemplateOptions(options ...botswat.TemplateOption) GateOption {
	return func(o *gateOptions) (err error) {
		if o.TemplateOptions != nil {
			return errors.New("template options are already set")
		}
		o.TemplateOptions, err = botswat.NewTemplateOptions(options...)
		if err != nil {
			return fmt.Errorf("could not setup HTML template options: %w", err)
		}
		return nil
	}
}

func WithSiteKey(key string) GateOption {
	return func(o *gateOptions) (err error) {
		if o.SiteKey != "" {
			return errors.New("site key is already set")
		}
		key = strings.TrimSpace(key)
		if key == "" {
			return errors.New("cannot use an empty site key")
		}
		o.SiteKey = key
		return nil
	}
}

func WithSiteKeyFromEnvironment(variableName string) GateOption {
	return func(o *gateOptions) (err error) {
		key := strings.TrimSpace(os.Getenv(variableName))
		if key == "" {
			return fmt.Errorf("cannot set up site key: environment variable %q is empty", variableName)
		}
		return WithSiteKey(key)(o)
	}
}

func WithDefaultSiteKey() GateOption {
	return WithSiteKeyFromEnvironment("TURNSTILE_SITE_KEY")
}

func WithSiteAction(action string) GateOption {
	return func(o *gateOptions) (err error) {
		if o.SiteAction != "" {
			return errors.New("site action is already set")
		}
		action = strings.TrimSpace(action)
		if action == "" {
			return errors.New("cannot use an empty site action")
		}
		o.SiteAction = action
		return nil
	}
}

func NewGate(withOptions ...GateOption) (oakhttp.Encoder, error) {
	o := &gateOptions{}

	var err error
	for _, option := range append(
		withOptions,
		func(o *gateOptions) (err error) {
			if o.TemplateOptions == nil {
				if err = WithTemplateOptions()(o); err != nil {
					return err
				}
			}
			if o.SiteKey == "" {
				if err = WithDefaultSiteKey()(o); err != nil {
					return err
				}
			}
			if o.SiteAction == "" {
				return errors.New("WithSiteAction option is required")
			}
			return nil
		},
	) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("cannot initialize encoder: %w", err)
		}
	}

	t, err := template.New("turnstile").Parse(botswat.Template)
	if err != nil {
		return nil, fmt.Errorf("could not parse default HTML template: %w", err)
	}

	t, err = t.Parse(templateAdjustments)
	if err != nil {
		return nil, fmt.Errorf("could not make Turnstile HTML template adjustments: %w", err)
	}
	b := &bytes.Buffer{}
	if err = t.Execute(b, &gateOptions{
		TemplateOptions: o.TemplateOptions,
		SiteKey:         o.SiteKey,
		SiteAction:      o.SiteAction,
	}); err != nil {
		return nil, fmt.Errorf("could not render Turnstile template: %w", err)
	}
	rendered := b.Bytes()

	return func(w http.ResponseWriter, any any) (err error) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusForbidden)
		_, err = io.Copy(w, bytes.NewReader(rendered))
		return err
	}, nil
}
