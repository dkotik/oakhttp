package turnstile

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type templateOptions struct {
	Title          string
	Description    string
	CookieName     string
	CookieDuration int
	SiteKey        string
	SiteAction     string
	Locale         string
	DarkTheme      bool
}

type TemplateOption func(*templateOptions) error

func WithTemplateTitle(s string) TemplateOption {
	return func(o *templateOptions) error {
		if o.Title != "" {
			return errors.New("title is already set")
		}
		if s == "" {
			return errors.New("cannot use an empty title")
		}
		o.Title = s
		return nil
	}
}

func WithDefaultTemplateTitle() TemplateOption {
	return func(o *templateOptions) error {
		if o.Title != "" {
			return nil
		}
		return WithTemplateTitle("Humanity Check")(o)
	}
}

func WithTemplateDescription(s string) TemplateOption {
	return func(o *templateOptions) error {
		if o.Description != "" {
			return errors.New("description is already set")
		}
		if s == "" {
			return errors.New("cannot use an empty description")
		}
		o.Description = s
		return nil
	}
}

func WithDefaultTemplateDescription() TemplateOption {
	return func(o *templateOptions) error {
		if o.Description != "" {
			return nil
		}
		return WithTemplateDescription("Complete a series of checks to demonstrate that you are not a bot.")(o)
	}
}

func WithTemplateCookieName(s string) TemplateOption {
	return func(o *templateOptions) error {
		if o.CookieName != "" {
			return errors.New("cookie name is already set")
		}
		if s == "" {
			return errors.New("cannot use an empty cookie name")
		}
		o.CookieName = s
		return nil
	}
}

func WithDefaultTemplateCookieName() TemplateOption {
	return func(o *templateOptions) error {
		if o.CookieName != "" {
			return nil
		}
		return WithTemplateCookieName("turnstile")(o)
	}
}

func WithTemplateCookieDuration(d time.Duration) TemplateOption {
	return func(o *templateOptions) error {
		if o.CookieDuration != 0 {
			return errors.New("cookie duration is already set")
		}
		if d == 0 {
			return errors.New("cannot use an empty cookie duration")
		}
		o.CookieDuration = int(d.Seconds()) * 1000
		return nil
	}
}

func WithDefaultTemplateCookieDuration() TemplateOption {
	return func(o *templateOptions) error {
		if o.CookieDuration != 0 {
			return nil
		}
		return WithTemplateCookieDuration(DefaultRetention)(o)
	}
}

func WithTemplateSiteKey(key string) TemplateOption {
	return func(o *templateOptions) error {
		if o.SiteKey != "" {
			return errors.New("site key is already set")
		}
		if key == "" {
			return errors.New("cannot use an empty site key")
		}
		o.SiteKey = key
		return nil
	}
}

func WithTemplateSiteKeyFromEnvironment(variableName string) TemplateOption {
	return func(o *templateOptions) error {
		key := strings.TrimSpace(os.Getenv(variableName))
		if key == "" {
			return fmt.Errorf("cannot get set key from environment: variable %q is not set", variableName)
		}
		return WithTemplateSiteKey(key)(o)
	}
}

func WithDefaultTemplateSiteKey() TemplateOption {
	return func(o *templateOptions) error {
		if o.SiteKey != "" {
			return nil
		}
		return WithTemplateSiteKeyFromEnvironment("TURNSTILE_SITE_KEY")(o)
	}
}

func WithTemplateSiteAction(action string) TemplateOption {
	return func(o *templateOptions) error {
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

func WithDefaultTemplateSiteAction() TemplateOption {
	return func(o *templateOptions) error {
		if o.SiteAction != "" {
			return nil
		}
		return WithTemplateSiteAction("view")(o)
	}
}

func WithTemplateLocale(locale string) TemplateOption {
	return func(o *templateOptions) error {
		if o.Locale != "" {
			return errors.New("locale is already set")
		}
		if locale == "" {
			return errors.New("cannot use an empty locale")
		}
		o.Locale = locale
		return nil
	}
}

func WithDefaultTemplateLocale() TemplateOption {
	return func(o *templateOptions) error {
		if o.Locale != "" {
			return nil
		}
		return WithTemplateLocale("en")(o)
	}
}

func WithTemplateDarkTheme() TemplateOption {
	return func(o *templateOptions) error {
		if o.DarkTheme {
			return errors.New("dark theme is already set")
		}
		o.DarkTheme = true
		return nil
	}
}
