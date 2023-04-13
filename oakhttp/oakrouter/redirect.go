package oakrouter

import (
	"errors"
	"net/http"

	"github.com/dkotik/oakacs/oakhttp"
)

var trailingSlashRedirect = func(w http.ResponseWriter, r *http.Request) error {
	URL := r.URL.String()
	if URL == "" {
		return errors.New("trailing slash redirect received request with empty URL")
	}
	http.Redirect(w, r, URL[:len(URL)-1], http.StatusTemporaryRedirect)
	return nil
}

func NewRedirect(URL string, statusCode int) oakhttp.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		http.Redirect(w, r, URL, statusCode)
		return nil
	}
}
