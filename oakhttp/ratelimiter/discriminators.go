package ratelimiter

import (
	"errors"
	"net"
	"net/http"
)

// SkipDiscriminator is a sentinel error used to indicate
// that a certain [http.Request] must be disregarded.
//
// revive:disable-next-line:error-naming
var SkipDiscriminator = errors.New("discriminator did not match")

// Discriminator associates tags to [http.Request]s in order to
// group related requests for a discriminating rate limiter.
// Requests with the same association tag will be tracked
// together by the [RateLimiter]. Return [SkipDiscriminator]
// sentinel value to disregard the [http.Request].
type Discriminator func(*http.Request) (string, error)

func NewIPAddressDiscriminator() Discriminator {
	return func(r *http.Request) (string, error) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return "", err
		}
		return ip, nil
	}
}

func NewRemoteAddrDiscriminator() Discriminator {
	return func(r *http.Request) (string, error) {
		return r.RemoteAddr, nil
	}
}

func NewCookieDiscriminator(name, noCookieValue string) Discriminator {
	if noCookieValue == "" {
		return func(r *http.Request) (string, error) {
			cookie, err := r.Cookie(name)
			if err == http.ErrNoCookie {
				return "", SkipDiscriminator
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

func NewHeaderDiscriminator(name, noHeaderValue string) Discriminator {
	if noHeaderValue == "" {
		return func(r *http.Request) (string, error) {
			value := r.Header.Get(name)
			if value == "" {
				return "", SkipDiscriminator
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
