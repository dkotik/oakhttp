package oakhttp

import (
	"encoding/json"
	"net/http"
)

// var (
// // not needed, just print
// 	ErrNotFound = errors.New("requested resource was not found")
// )

type HTTPError interface {
	HTTPStatusCode() int
	error
	// WriteResponse(w http.ResponseWriter) error
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

type ValidationError struct{}

func (e *ValidationError) HTTPStatusCode() int {
	return http.StatusUnprocessableEntity
}

// type EncodingError struct{} // not needed, just print

// type DecodingError struct{} // not needed, just print

// type RequestFactoryError struct{} // not needed, just print

func MarshalErrorToJSON(w http.ResponseWriter, err error) error {
	if httpError, ok := err.(HTTPError); ok {
		w.WriteHeader(httpError.HTTPStatusCode)
	}
	// stringer, ok := str.(fmt.Stringer)
	return json.NewEncoder(w).Encode(map[string]string{
		"Error": err.Error(),
	})
}
