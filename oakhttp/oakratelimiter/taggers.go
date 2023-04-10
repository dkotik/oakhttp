package oakratelimiter

import (
	"context"
	"errors"
	"net"
	"net/http"
)

// SkipTagger is a sentinel error used to indicate
// that a certain [http.Request] must be disregarded.
//
// revive:disable-next-line:error-naming
var SkipTagger = errors.New("discriminator did not match")

// Tagger associates tags to [http.Request]s in order to
// group related requests for a discriminating rate limiter.
// Requests with the same association tag will be tracked
// together by the [RateLimiter]. Return [SkipTagger]
// sentinel value to disregard the [http.Request].
type Tagger func(*http.Request) (string, error)

type ContextTagger func(context.Context) (string, error)

func NewRequestTaggerFromContextTagger(t ContextTagger) Tagger {
	if t == nil {
		panic("cannot use a <nil> context tagger")
	}
	return func(r *http.Request) (string, error) {
		return t(r.Context())
	}
}

func NewIPAddressTagger() Tagger {
	return func(r *http.Request) (string, error) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return "", err
		}
		return ip, nil
	}
}

func NewIPAddressWithPortTagger() Tagger {
	return func(r *http.Request) (string, error) {
		return r.RemoteAddr, nil
	}
}

func NewCookieTagger(name, noCookieValue string) Tagger {
	if noCookieValue == "" {
		return func(r *http.Request) (string, error) {
			cookie, err := r.Cookie(name)
			if err == http.ErrNoCookie {
				return "", SkipTagger
			} else if err != nil {
				return "", err
			}
			return cookie.Value, nil
		}
	}

	return func(r *http.Request) (string, error) {
		cookie, err := r.Cookie(name)
		if err == http.ErrNoCookie {
			return noCookieValue, nil
		} else if err != nil {
			return "", err
		}
		return cookie.Value, nil
	}
}

func NewHeaderTagger(name, noHeaderValue string) Tagger {
	if noHeaderValue == "" {
		return func(r *http.Request) (string, error) {
			value := r.Header.Get(name)
			if value == "" {
				return "", SkipTagger
			}
			return value, nil
		}
	}

	return func(r *http.Request) (string, error) {
		value := r.Header.Get(name)
		if value == "" {
			return noHeaderValue, nil
		}
		return value, nil
	}
}
