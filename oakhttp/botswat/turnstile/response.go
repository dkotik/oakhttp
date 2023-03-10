package turnstile

import (
	"errors"
	"fmt"
	"time"

	"github.com/relvacode/iso8601"
)

const responseReadLimit = 1024 * 24

type Response struct {
	// Success is true if the verification passed.
	Success bool `json:"success"`

	// ChallengeTime when the verfication passed in ISO 8601 format.
	ChallengeTime string `json:"challenge_ts"`

	// Action name of the validation, set by the subject. Must match the information inside [Request.Response].
	Action string `json:"action"`

	// Hostname of the site that requested verification.
	Hostname string `json:"hostname"`

	// ErrorCodes of any problems that were encountered.
	//
	// https://developers.cloudflare.com/turnstile/get-started/server-side-validation/#error-codes
	ErrorCodes []string `json:"error-codes"`

	// CData is client subject data.
	CData string `json:"cdata"`
}

func (r *Response) Time() (time.Time, error) {
	return iso8601.ParseString(r.ChallengeTime)
}

func (r *Response) Validate() error {
	if l := len(r.ErrorCodes); l > 0 {
		errs := make([]error, l)
		for i, code := range r.ErrorCodes {
			errs[i] = turnstileError(code)
		}
		return errors.Join(errs...)
	}
	if r.Success == false {
		return errors.New("request rejected")
	}
	return nil
}

func turnstileError(code string) error {
	switch code {
	case "missing-input-secret":
		return errors.New("request is missing the secret key")
	case "invalid-input-secret":
		return errors.New("invalid secret key")
	case "missing-input-response":
		return errors.New("empty client response")
	case "invalid-input-response":
		return errors.New("invalid client response")
	case "bad-request":
		return errors.New("malformed request")
	case "timeout-or-duplicate":
		return errors.New("client response expired")
	case "internal-error":
		return errors.New("request failed")
	default:
		return fmt.Errorf("unknown turnstile error: %s", code)
	}
}
