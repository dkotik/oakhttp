package msg

import (
	"net/http"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	ErrorInternalTitle = &i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ErrorInternalTitle",
			Other: http.StatusText(http.StatusInternalServerError),
		},
	}
	ErrorInternalDescription = &i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ErrorInternalDescription",
			Other: "Service encountered internal error. It is unable to complete the desired operation.",
		},
	}

	ErrorNotFoundTitle = &i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ErrorNotFoundTitle",
			Other: http.StatusText(http.StatusNotFound),
		},
	}
	ErrorNotFoundDescription = &i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ErrorNotFoundDescription",
			Other: "Requested content was not found.",
		},
	}

	ErrorAccessDeniedTitle = &i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ErrorAccessDeniedTitle",
			Other: http.StatusText(http.StatusForbidden),
		},
	}
	ErrorAccessDeniedDescription = &i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ErrorAccessDeniedDescription",
			Other: "Service is unable to complete the desired operation due to insufficient access level.",
		},
	}
)
