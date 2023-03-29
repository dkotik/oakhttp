package oakswitch

import (
	"fmt"
	"net/http"

	"github.com/dkotik/oakacs/oakhttp"
)

func New(withOptions ...Option) (oakhttp.Handler, error) {
	o := &options{}
	for _, option := range withOptions {
		if err := option(o); err != nil {
			return nil, fmt.Errorf("cannot initialize method switch: %w", err)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		switch r.Method {
		case http.MethodGet, http.MethodHead:
			if o.Get != nil {
				return o.Get(w, r)
			}
		case http.MethodPost:
			if o.Post != nil {
				return o.Post(w, r)
			}
		case http.MethodPut:
			if o.Put != nil {
				return o.Put(w, r)
			}
		case http.MethodDelete:
			if o.Delete != nil {
				return o.Delete(w, r)
			}
		case http.MethodPatch:
			if o.Patch != nil {
				return o.Patch(w, r)
			}
		}
		return NewMethodNotAllowedError(r.Method)
	}, nil
}
