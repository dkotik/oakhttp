package oakbotswat

import "errors"

type Error string

func NewErrorFromCodes(codes ...string) error {
	errs := make([]error, len(codes))
	for i, code := range codes {
		errs[i] = Error(code)
	}
	return errors.Join(errs...)
}

var ErrNotHuman = errors.New("robot detected")

const (
	ErrMissingInputSecret   Error = "missing-input-secret"
	ErrInvalidInputSecret   Error = "invalid-input-secret"
	ErrMissingInputResponse Error = "missing-input-response"
	ErrInvalidInputResponse Error = "invalid-input-response"
	ErrBadRequest           Error = "bad-request"
	ErrTimeoutOrDuplicate   Error = "timeout-or-duplicate"
	ErrInternalError        Error = "internal-error"

	errorPrefix = "humanity check failed: "
)

func (err Error) Error() string {
	switch err {
	case ErrMissingInputSecret:
		return errorPrefix + "request is missing the secret key"
	case ErrInvalidInputSecret:
		return errorPrefix + "invalid secret key"
	case ErrMissingInputResponse:
		return errorPrefix + "empty client response"
	case ErrInvalidInputResponse:
		return errorPrefix + "invalid client response"
	case ErrBadRequest:
		return errorPrefix + "malformed request"
	case ErrTimeoutOrDuplicate:
		return errorPrefix + "client response expired"
	case ErrInternalError:
		return errorPrefix + "request failed"
	case "":
		return errorPrefix + "unexpected empty error message"
	default:
		return errorPrefix + "unexpected error: " + string(err)
	}
}
