package turnstile

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/tigerperformanceinstitute/tprograms/service/oakhttp"
	"github.com/tigerperformanceinstitute/tprograms/service/oakhttp/store"
)

type contextKey struct{}

type authenticator struct {
	client     *Turnstile
	siteAction string
	extractor  oakhttp.StringExtractor
	store      store.KeyKeyValue
	storeKey   []byte
}

func NewAuthenticator(withOptions ...AuthenticatorOption) (oakhttp.Authenticator, error) {
	o, err := newAuthenticatorOptions(withOptions)
	if err != nil {
		return nil, err
	}
	return &authenticator{
		client:     o.Client,
		siteAction: o.SiteAction,
		extractor:  o.Extractor,
		store:      o.Store,
		storeKey:   append([]byte("turnstile:"), []byte(o.SiteAction)...),
	}, nil
}

func (a *authenticator) IsAllowed(r *http.Request) (*http.Request, error) {
	token, err := a.extractor(r)
	if err != nil {
		return nil, err
	}
	ctx := r.Context()
	tokenKey := []byte(token)
	data, err := a.store.Get(ctx, a.storeKey, tokenKey)
	if err != nil {
		if !errors.Is(err, store.ErrValueNotFound) {
			return nil, err
		}
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		userData, err := a.client.Challenge(ctx, token, ip, a.siteAction)
		if err != nil {
			return nil, err
		}
		if err = a.store.Set(ctx, a.storeKey, tokenKey, []byte(userData)); err != nil {
			return nil, err
		}
	}
	return r.WithContext(context.WithValue(ctx, contextKey{}, data)), nil
}

// UserDataFromContext recovers Turnstile CData from context.
func UserDataFromContext(ctx context.Context) []byte {
	data, _ := ctx.Value(contextKey{}).([]byte)
	return data
}
