package oakbotswat

import (
	"errors"
	"fmt"
	"strings"
	"time"

	_ "embed"
)

//go:embed template.html
var Template string

type TemplateOptions struct {
	Locale         string
	Title          string
	Description    string
	ErrorURL       string
	CookieName     string
	CookieDuration time.Duration
}

type TemplateOption func(*TemplateOptions) error

func NewTemplateOptions(all ...TemplateOption) (*TemplateOptions, error) {
	o := &TemplateOptions{}
	var err error
	for _, option := range append(all, func(o *TemplateOptions) (err error) {
		if o.Locale == "" {
			if err = WithLocale("en")(o); err != nil {
				return err
			}
		}
		if o.Title == "" {
			if err = WithTitle("Humanity Check")(o); err != nil {
				return err
			}
		}
		if o.Description == "" {
			if err = WithDescription("Please confirm your humanity to access this resource.")(o); err != nil {
				return err
			}
		}
		if o.ErrorURL == "" {
			if err = WithErrorURL("/")(o); err != nil {
				return err
			}
		}
		if o.CookieName == "" {
			if err = WithCookieName(DefaultCookieName)(o); err != nil {
				return err
			}
		}
		if o.CookieDuration == 0 {
			if err = WithCookieDuration(time.Minute * 5)(o); err != nil {
				return err
			}
		}
		return nil
	}) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("cannot initialize OakBotSWAT template options: %w", err)
		}
	}
	return o, nil
}

func WithLocale(key string) TemplateOption {
	return func(o *TemplateOptions) error {
		if o.Locale != "" {
			return errors.New("site key is already set")
		}
		key = strings.TrimSpace(key)
		if key == "" {
			return errors.New("cannot use an empty site key")
		}
		o.Locale = key
		return nil
	}
}

func WithTitle(title string) TemplateOption {
	return func(o *TemplateOptions) error {
		if o.Title != "" {
			return errors.New("title is already set")
		}
		title = strings.TrimSpace(title)
		if title == "" {
			return errors.New("cannot use an empty title")
		}
		o.Title = title
		return nil
	}
}

func WithDescription(description string) TemplateOption {
	return func(o *TemplateOptions) error {
		if o.Description != "" {
			return errors.New("description is already set")
		}
		description = strings.TrimSpace(description)
		if description == "" {
			return errors.New("cannot use an empty description")
		}
		o.Description = description
		return nil
	}
}

func WithErrorURL(URL string) TemplateOption {
	return func(o *TemplateOptions) error {
		if o.ErrorURL != "" {
			return errors.New("error URL is already set")
		}
		URL = strings.TrimSpace(URL)
		if URL == "" {
			return errors.New("cannot use an empty error URL")
		}
		o.ErrorURL = URL
		return nil
	}
}

func WithCookieName(name string) TemplateOption {
	return func(o *TemplateOptions) error {
		if o.CookieName != "" {
			return errors.New("cookie name is already set")
		}
		name = strings.TrimSpace(name)
		if name == "" {
			return errors.New("cannot use an empty cookie name")
		}
		o.CookieName = name
		return nil
	}
}

func WithCookieDuration(d time.Duration) TemplateOption {
	return func(o *TemplateOptions) error {
		if o.CookieDuration != 0 {
			return errors.New("cookie duration is already set")
		}
		if d < time.Minute {
			return errors.New("cookie duration must be greater than 1 minute")
		}
		o.CookieDuration = d
		return nil
	}
}

// func NewTemplate(withOptions ...TemplateOption) (*template.Template, error) {
// 	o := &templateOptions{}
//
// 	var err error
// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("cannot create OakBotSWAT template: %w", err)
// 		}
// 	}()
//
// 	for _, option := range append(withOptions, func(o *templateOptions) error {
// 		// if o.SiteKey == "" {
// 		// 	return errors.New("site key is required")
// 		// }
// 		return nil
// 	}) {
// 		if err = option(o); err != nil {
// 			return nil, err
// 		}
// 	}
//
// 	t, err := template.New("botswat").Parse(templateHTML)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot parse template: %w", err)
// 	}
// 	return t, nil
// }
