package hcaptcha

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// TODO: model all errors: https://docs.hcaptcha.com/#siteverify-error-codes-table

// Should be TokenStore! // TODO: fill out
// type TokenStore interface {
// 	IsValid(context.Context, string) error
// 	Commit(context.Context, string, time.Duration) error
// }
// type KeyValueCache interface {
// 	Get(ctx context.Context, key string) (value string, err error)
// 	Set(ctx context.Context, key, value string) (err error)
// }
//
// type MemoryCache struct {
// 	duration time.Duration
// 	values   map[string]time.Time
// 	mu       *sync.Mutex
// }
//
// func (m *MemoryCache) Set(ctx context.Context, key, value string) error {
//
// }
//
// func (m *MemoryCache) Get(ctx context.Context, key string) (string, error) {
// 	mu.Lock()
// 	defer mu.Unlock()
// 	when, ok := m.values[key]
// 	if !ok {
// 		return "", nil
// 	}
//
// 	t := time.Now()
// 	if when.Before(t) {
// 		delete(m.cache, key)
// 		return "", nil
// 	}
// 	return
// }
//
// func NewMemoryCache(d time.Duration) *MemoryCache {
// 	return &MemoryCache{
// 		duration: d,
// 		values:   make(map[string]time.Time),
// 		mu:       &sync.Mutex{},
// 	}
// }

type HCaptchaValidator func(ctx context.Context, token, personIP string) error

type HCaptchaError struct {
	SiteKey string
	Cause   error
}

func (h *HCaptchaError) Error() string {
	return fmt.Sprintf("HCaptcha humanity validation failed for site key %q: %v", h.SiteKey, h.Cause)
}

func NewMemoryCachedValidator(d time.Duration, v HCaptchaValidator) HCaptchaValidator {
	cache := make(map[string]time.Time)
	mu := &sync.Mutex{}
	commit := func(token string) {
		mu.Lock()
		defer mu.Unlock()

		t := time.Now()
		filter := make([]string, 0)
		for key, exp := range cache {
			if exp.Before(t) {
				filter = append(filter, key)
			}
		}
		for _, expiredToken := range filter {
			delete(cache, expiredToken)
		}

		if len(cache) < 10000 { // maximum number of cached tokens
			cache[token] = time.Now().Add(d)
		}
	}

	return func(ctx context.Context, token, personIP string) (err error) {
		mu.Lock()

		when, ok := cache[token]
		if !ok {
			mu.Unlock()
			if err = v(ctx, token, personIP); err != nil {
				return err
			}
			commit(token)
			return nil
		}

		t := time.Now()
		if when.Before(t) {
			delete(cache, token)
			mu.Unlock()
			if err = v(ctx, token, personIP); err != nil {
				return err
			}
			commit(token)
			return nil
		}

		mu.Unlock()
		return nil
	}
}

func NewHCaptchaValidator(siteKey, secretKey string) HCaptchaValidator {
	client := http.DefaultClient // TODO: harden

	return func(ctx context.Context, token, personIP string) (err error) {
		defer func() {
			if err != nil {
				// leakKey := "none"
				// if len(secretKey) > 6 {
				// 	leakKey = secretKey[0:5] + "..."
				// }
				err = &HCaptchaError{
					SiteKey: siteKey,
					Cause:   err,
				}
			}
		}()

		// if token == "passthrough" {
		// 	return nil // TODO: REMOVE passthrough
		// }

		raw, err := client.PostForm(
			"hCaptchaAPIEndpoint", url.Values{ // TODO: fix variable.
				"secret":   {secretKey},
				"sitekey":  {siteKey},
				"remoteip": {personIP},
				"response": {token},
			})
		if err != nil {
			return fmt.Errorf("failed to reach HCaptcha server: %w", err)
		}

		body, err := io.ReadAll(raw.Body)
		raw.Body.Close()
		if err != nil {
			return fmt.Errorf("failed to read HCaptcha response body: %w", err)
		}

		var response struct {
			ChallengeTS string   `json:"challenge_ts"`
			Hostname    string   `json:"hostname"`
			ErrorCodes  []string `json:"error-codes,omitempty"`
			Success     bool     `json:"success"`
			Credit      bool     `json:"credit,omitempty"`
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			return fmt.Errorf("failed to parse HCaptcha response body: %w", err)
		}

		if len(response.ErrorCodes) > 0 {
			return fmt.Errorf("error codes: %+v", response.ErrorCodes)
		}

		if !response.Success {
			return errors.New("unknown cause")
		}

		return nil
	}
}

/*

fly secrets set HCAPTCHA_SECRET_KEY=... HCAPTCHA_SITE_KEY=...

Content-Security-Policy Settings#

Content Security Policy (CSP) headers are an added layer of security that help to mitigate certain types of attacks, including Cross Site Scripting (XSS), clickjacking, and data injection attacks.

If you use CSP headers, please add the following to your configuration:

    script-src should include https://hcaptcha.com, https://*.hcaptcha.com
    frame-src should include https://hcaptcha.com, https://*.hcaptcha.com
    style-src should include https://hcaptcha.com, https://*.hcaptcha.com
    connect-src should include https://hcaptcha.com, https://*.hcaptcha.com

Please do not hard-code specific subdomains, like newassets.hcaptcha.com, into your CSP: asset subdomains used may vary over time or by region.

If you are an enterprise customer and would like to enable additional verification to be performed, you can optionally choose the following CSP strategy:

    unsafe-eval and unsafe-inline should include https://hcaptcha.com, https://*.hcaptcha.com

*/
