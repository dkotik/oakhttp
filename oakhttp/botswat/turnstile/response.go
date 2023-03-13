package turnstile

import (
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
	return botswat.NewErrorFromCodes(r.ErrorCodes...)
}
