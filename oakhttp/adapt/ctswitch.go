package adapt

import (
	"fmt"
	"net/http"

	"github.com/dkotik/oakmux"
)

type contentTypeSwitch struct {
	handlers map[string]oakmux.Handler
}

func (ct *contentTypeSwitch) ServeHyperText(
	w http.ResponseWriter,
	r http.Request,
) error {
	contentType := r.Header().Get("Content-Type")
	handler, ok := ct.handlers[contentType]
	if !ok {
		// TODO: there is probably an HTTP code for this, use custom error:
		return fmt.Errorf("content type %q is not supported", contentType)
	}
	return handler.ServeHyperText(w, r)
}
