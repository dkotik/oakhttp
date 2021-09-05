package identity

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// https://zyszys.github.io/awesome-captcha/
// https://github.com/kataras/hcaptcha/blob/master/hcaptcha.go

// TokenValidator checks if a given token is acceptable.
type TokenValidator func(ctx context.Context, token string) error

// HCaptcha1 checks if a certain token identifies a human entity using hCaptcha v1.
func HCaptcha1(secret string) TokenValidator {
	// https://docs.hcaptcha.com/
	return func(ctx context.Context, token string) (err error) {
		defer func() {
			if err != nil {
				err = fmt.Errorf("could not validate humanity: %w", err)
			}
		}()
		client := &http.Client{}
		// @TODO: replace with NewRequestWithContext, then client is not even needed
		resp, err := client.PostForm("https://hcaptcha.com/siteverify",
			url.Values{
				"secret":   {secret},
				"response": {token},
				// remoteip	Optional. The user's IP address.
				// sitekey	Optional. The sitekey you expect to see.
			},
		)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}

		var v struct {
			ChallengeTS string   `json:"challenge_ts"`
			Hostname    string   `json:"hostname"`
			ErrorCodes  []string `json:"error-codes,omitempty"`
			Success     bool     `json:"success"`
			Credit      bool     `json:"credit,omitempty"`
		}
		err = json.Unmarshal(body, &v)
		if err != nil {
			return
		}
		// v.ErrorCodes = append(v.ErrorCodes, err.Error())

		if !v.Success {
			return errors.New("humanity token is busted")
		}
		return nil
	}
}
