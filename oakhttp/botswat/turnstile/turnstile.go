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

	"github.com/dkotik/oakacs/oakhttp/botswat"
)

func New(withOptions ...Option) (botswat.Verifier, error) {
	o := &options{}
	var err error
	for _, option := range append(
		withOptions,
		WithDefaultOptions(),
		func(o *options) error { // validate
			if o.Hostname == "" {
				return errors.New("host name is required")
			}
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

	return func(ctx context.Context, clientResponseToken, IP string) (data string, err error) {
		payload, err := json.Marshal(&Request{
			Secret:   o.Secret,
			Response: clientResponseToken,
			RemoteIP: IP,
		})
		if err != nil {
			return "", fmt.Errorf("JSON encoding failed: %w", err)
		}

		request, err := http.NewRequest("POST", o.Endpoint, bytes.NewReader(payload))
		if err != nil {
			return "", fmt.Errorf("invalid HTTP API request: %w", err)
		}
		request = request.WithContext(ctx)
		hr, err := o.HTTPClient.Do(request)
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

		if r.Hostname != o.Hostname {
			return "", errors.New("turnstile response and request hostnames do not match")
		}

		if !o.IsAllowedAction(r.Action) {
			return "", errors.New("turnstile response action is not allowed")
		}

		// cData is customer payload that can be used to attach customer data to the challenge throughout its issuance and which is returned upon validation. This can only contain up to 255 alphanumeric characters including _ and -.
		return r.CData, nil
	}, nil
}
