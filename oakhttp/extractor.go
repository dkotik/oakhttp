package oakhttp

import (
	"errors"
	"net/http"
	"path"
)

type StringExtractor func(r *http.Request) (string, error)

func PathTailExtractor(r *http.Request) (string, error) {
	switch tail := path.Base(r.URL.Path); tail {
	case ".", "/":
		return "", NewNotFoundError(errors.New("no value"))
	default:
		return tail, nil
	}
}

func NewCookieExtractor(name string) StringExtractor {
	return func(r *http.Request) (string, error) {
		cookie, err := r.Cookie(name)
		if err != nil {
			return "", NewNotFoundError(err)
		}
		return cookie.Value, nil
	}
}
