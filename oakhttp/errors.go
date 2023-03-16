package oakhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/exp/slog"
)

func (d *DomainAdaptor) errorOrPanicHandler(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("critical failure: %v", r)
				_ = d.WriteErrors(w, err)
				slog.Error("an HTTP request panicked", err)
			}
		}()

		if err = h(w, r); err != nil {
			// unwrappable, ok := err.Unwrap().(interface{ Unwrap() []error })
			// if ok {
			//   d.WriteErrors(unwrappable.Unwrap()...)
			//   return
			// }
			_ = d.WriteErrors(w, err)
		}
	}
}

func (d *DomainAdaptor) WriteErrors(w http.ResponseWriter, display ...error) (err error) {
	w.Header().Set("Content-Type", d.encoderContentType)
	var httpError Error
	if errors.As(err, &httpError) {
		w.WriteHeader(httpError.HTTPStatusCode())
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	err = d.encoderFactory(w).Encode(map[string][]error{
		"Errors": display,
	})
	if err != nil {
		return fmt.Errorf("encoder failed: %w", err)
	}
	return nil
}

func DefaultErrorHandler(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// if r.Header().Get("Content-Type") == d.encoderContentType {
		//   TODO: encode the error map here after Unwrap() []error
		// }
		if err := h(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// var (
// // not needed, just print
// 	ErrNotFound = errors.New("requested resource was not found")
// )

// Unwrap asserts that some error wraps [Error] and returns it. Use the result together with [http.Error]. Log the original error to prevent leaking potentially sensitive information to the client while recording it on the server.
func Unwrap(wrapped error) (err Error, ok bool) {
	if errors.As(wrapped, &err) {
		return err, true
	}
	return nil, false
}

type Error interface {
	error
	HTTPStatusCode() int
}

type NotFoundError struct {
	resource string
}

func NewNotFoundError(resource string) *NotFoundError {
	return &NotFoundError{resource: resource}
}

func (e *NotFoundError) HTTPStatusCode() int {
	return http.StatusNotFound
}

func (e *NotFoundError) Error() string {
	return "resource \"" + e.resource + "\" was not found"
}

// type EncodingError struct{} // not needed, just print

// type DecodingError struct{} // not needed, just print

// type RequestFactoryError struct{} // not needed, just print

func MarshalErrorToJSON(w http.ResponseWriter, err error) error {
	if httpError, ok := err.(Error); ok {
		w.WriteHeader(httpError.HTTPStatusCode())
	}
	// stringer, ok := str.(fmt.Stringer)
	return json.NewEncoder(w).Encode(map[string]string{
		"Error": err.Error(),
	})
}
