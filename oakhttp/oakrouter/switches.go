package oakrouter

import (
	"net/http"
	"strings"

	"github.com/dkotik/oakacs/oakhttp"
)

type mapSwitch struct {
	Prefix   string
	Handlers map[string]oakhttp.Handler
}

type mapSwitchOption func(*mapSwitch) error

func NewMapSwitch() (oakhttp.Handler, error) {
	mapSwitch := &mapSwitch{}

	return func(w http.ResponseWriter, r *http.Request) error {
		tail, ok := strings.CutPrefix(r.URL.Path, mapSwitch.Prefix)
		if !ok {
			return ErrPrefixMismatch
		}

		h, ok := mapSwitch.Handlers[tail]
		if !ok {
			return &MatchError{oakhttp.NewNotFoundError(r.URL.Path)}
		}
		// r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix) // cut
		return h(w, r)
	}, nil
}
