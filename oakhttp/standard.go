package oakhttp

import "net/http"

func NewStandardAdaptor(h http.Handler) Handler {
	return HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		h.ServeHTTP(w, r)
		return nil
	})
}

func writeError(w http.ResponseWriter, err error, code int) {
	if code == http.StatusInternalServerError {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	} else {
		http.Error(w, err.Error(), code)
	}
}
