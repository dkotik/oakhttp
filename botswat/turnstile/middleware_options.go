package turnstile

import (
	_ "embed" // for embedding template.html
	"errors"
	"fmt"
)

//go:embed template.html
var templateHTML string

type middlewareOptions struct {
	templateOptions *templateOptions
	challenge       []byte
}

type MiddlewareOption func(*middlewareOptions) error

func WithRenderedTemplate(tmpl []byte) MiddlewareOption {
	return func(o *middlewareOptions) error {
		if len(tmpl) == 0 {
			return errors.New("cannot use an empty template")
		}
		if o.challenge != nil || o.templateOptions != nil {
			return errors.New("template is already set")
		}
		o.challenge = tmpl
		return nil
	}
}

func WithTemplateOptions(options ...TemplateOption) MiddlewareOption {
	return func(o *middlewareOptions) (err error) {
		if len(o.challenge) != 0 {
			return errors.New("template is already set")
		}
		if o.templateOptions == nil {
			o.templateOptions = &templateOptions{}
		}
		for _, option := range append(
			options,
			WithDefaultTemplateTitle(),
			WithDefaultTemplateDescription(),
			// WithDefaultTemplateCookieName(), // applied later
			WithDefaultTemplateCookieDuration(),
			WithDefaultTemplateSiteKey(),
			WithDefaultTemplateSiteAction(),
			WithDefaultTemplateLocale(),
		) {
			if err = option(o.templateOptions); err != nil {
				return fmt.Errorf("cannot apply template option: %w", err)
			}
		}
		return nil
	}
}

func WithDefaultTemplate() MiddlewareOption {
	return func(o *middlewareOptions) error {
		if o.challenge != nil || o.templateOptions != nil {
			return nil
		}
		return WithTemplateOptions()(o)
	}
}
