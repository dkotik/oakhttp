package botswat

import (
	"errors"
	"net/http"

	"github.com/dkotik/oakacs/oakhttp"
)

func New(withOptions ...Option) (oakhttp.Middleware, error) {
	// verifier, err := New(withOptions)
	// if err != nil {
	// 	return nil, err
	// }
	return func(h oakhttp.Handler) oakhttp.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			return errors.New("unimplemented")
		}
	}, nil
}
