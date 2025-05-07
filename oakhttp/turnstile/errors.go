package turnstile

import (
	"errors"
	"net/http"

	"github.com/tigerperformanceinstitute/tprograms/service/oakhttp"
)

var ErrNoCookie = oakhttp.NewAccessDeniedError(http.ErrNoCookie)

// var ErrTokenEmpty = errors.New("cannot recover request token proving humanity: token is empty")

type Error string

func NewErrorFromCodes(codes ...string) error {
	errs := make([]error, len(codes))
	for i, code := range codes {
		errs[i] = Error(code)
	}
	return errors.Join(errs...)
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

func (err Error) HyperTextStatusCode() int {
	return http.StatusForbidden
}

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
