package botswat

import (
	"errors"

	"github.com/dkotik/oakhttp"
)

const errorKnowledgeCodePrefix = "botswat:"

var ErrTokenEmpty = oakhttp.NewAccessDeniedError(errors.New("humanity token is empty"), errorKnowledgeCodePrefix+"emptyToken")

type Error string

func NewErrorFromCodes(codes ...string) error {
	errs := make([]error, len(codes))
	for i, code := range codes {
		errs[i] = Error(code)
	}
	if len(errs) == 0 {
		return nil
	}
	return oakhttp.NewAccessDeniedError(errors.Join(errs...), errorKnowledgeCodePrefix+codes[0])
}

const (
	ErrMissingInputSecret   Error = "missing-input-secret"
	ErrInvalidInputSecret   Error = "invalid-input-secret"
	ErrMissingInputResponse Error = "missing-input-response"
	ErrInvalidInputResponse Error = "invalid-input-response"
	ErrBadRequest           Error = "bad-request"
	ErrTimeoutOrDuplicate   Error = "timeout-or-duplicate"
	ErrInternalError        Error = "internal-error"
)

func (err Error) Error() string {
	switch err {
	case ErrMissingInputSecret:
		return "request is missing the secret key"
	case ErrInvalidInputSecret:
		return "invalid secret key"
	case ErrMissingInputResponse:
		return "empty client response"
	case ErrInvalidInputResponse:
		return "invalid client response"
	case ErrBadRequest:
		return "malformed request"
	case ErrTimeoutOrDuplicate:
		return "client response expired"
	case ErrInternalError:
		return "request failed"
	case "":
		return "unexpected empty error message"
	default:
		return "unexpected error: " + string(err)
	}
}
