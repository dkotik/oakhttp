/*

Package turnstile creates a secure by default humanity verifier backed by Cloudflare's Turnstile service.

On the client include the the following script:

  <script src="https://challenges.cloudflare.com/turnstile/v0/api.js?onload=onloadTurnstileCallback" async defer></script>
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

type Verifier func(ctx context.Context, response, IP string) (string, error)

type ResponseExtractor func(r *http.Request) (string, error)

func New(withOptions ...Option) (Verifier, error) {
	options := &options{}
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
		if err = option(options); err != nil {
			return nil, fmt.Errorf("cannot initialize Cloudflare Turnstile verifier: %w", err)
		}
	}

	return func(ctx context.Context, response, IP string) (data string, err error) {
		payload, err := json.Marshal(&Request{
			Secret:   options.Secret,
			Response: response,
			RemoteIP: IP,
		})
		if err != nil {
			return "", fmt.Errorf("JSON encoding failed: %w", err)
		}

		request, err := http.NewRequest("POST", options.Endpoint, payload)
		if err != nil {
			return "", fmt.Errorf("invalid HTTP API request: %w", err)
		}
		request = request.WithContext(ctx)
		response, err := options.HTTPClient.Do(request)
		if err != nil {
			return "", fmt.Errorf("HTTP API request failed: %w", err)
		}
		defer response.Body.Close()

		var r *Response
		err := json.NewDecoder(
			io.LimitReader(bytes.NewReader(response.Body), responseReadLimit),
		).Decode(&r)
		if err != nil {
			return "", fmt.Errorf("JSON decoding failure: %w", err)
		}
		if err = r.Validate(); err != nil {
			return "", fmt.Errorf("turnstile request failed: %w", err)
		}

		if r.Hostname != options.Hostname {
			return "", errors.New("turnstile response and request hostnames do not match")
		}

		if !options.IsAllowedAction(r.Action) {
			return "", errors.New("turnstile response action is not allowed")
		}

		return r.CData, nil
	}
}
