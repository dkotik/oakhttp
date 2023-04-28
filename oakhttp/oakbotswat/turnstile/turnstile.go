/*

Package turnstile creates a secure by default humanity verifier backed by Cloudflare's Turnstile service.

On the client include the the following script:

  <script src="https://challenges.cloudflare.com/turnstile/v0/api.js?onload=onloadTurnstileCallback" async defer></script>

Cloudflare Turnstile documentation: <https://developers.cloudflare.com/turnstile/>
*/
package turnstile

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Turnstile struct {
	client         *http.Client
	secretKey      string
	endpoint       string
	hostname       string
	allowedActions []string
}

func (t *Turnstile) IsAllowedAction(name string) bool {
	for _, action := range t.allowedActions {
		if action == name {
			return true
		}
	}
	return false
}

func (t *Turnstile) VerifyHumanityToken(
	ctx context.Context,
	clientResponseToken string,
	clientIPAddress string,
) (
	userData string,
	err error,
) {
	payload, err := json.Marshal(&Request{
		Secret:   t.secretKey,
		Response: clientResponseToken,
		RemoteIP: clientIPAddress,
	})
	if err != nil {
		return "", fmt.Errorf("JSON encoding failed: %w", err)
	}

	request, err := http.NewRequest("POST", t.endpoint, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("invalid HTTP API request: %w", err)
	}
	request = request.WithContext(ctx)
	request.Header.Set("Content-Type", "application/json") // critical
	hr, err := t.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("HTTP API request failed: %w", err)
	}
	defer hr.Body.Close()

	var r *Response
	if err = json.NewDecoder(
		io.LimitReader(hr.Body, responseReadLimit),
	).Decode(&r); err != nil {
		return "", fmt.Errorf("JSON decoding failure: %w", err)
	}
	if err = r.Validate(); err != nil {
		return "", fmt.Errorf("turnstile request failed: %w", err)
	}

	if r.Hostname != t.hostname {
		return "", errors.New("turnstile response and request hostnames do not match")
	}

	if !t.IsAllowedAction(r.Action) {
		return "", errors.New("turnstile response action is not allowed")
	}

	// cData is customer payload that can be used to attach customer data to the challenge throughout its issuance and which is returned upon validation. This can only contain up to 255 alphanumeric characters including _ and -.
	return r.CData, nil
}

func New(withOptions ...Option) (*Turnstile, error) {
	o := &options{}
	var err error
	for _, option := range append(
		withOptions,
		WithDefaultOptions(),
		func(o *options) error { // validate
			if len(o.AllowedActions) == 0 {
				return errors.New("at least one allowed action is required")
			}
			return nil
		},
	) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("cannot initialize Cloudflare Turnstile verifier: %w", err)
		}
	}

	return &Turnstile{
		client:         o.HTTPClient,
		secretKey:      o.SecretKey,
		endpoint:       o.Endpoint,
		hostname:       o.Hostname,
		allowedActions: o.AllowedActions,
	}, nil
}
