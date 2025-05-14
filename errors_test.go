package oakhttp

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorHandler(t *testing.T) {
	eh := NewErrorHandler(nil, nil, nil)
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	eh.HandleError(w, r, errors.New("test"))
}
