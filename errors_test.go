package oakhttp

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sebdah/goldie/v2"
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

func TestErrorRendering(t *testing.T) {
	msg, err := NewError(errors.New("test error"), "testError").Localize(testLocalizer)
	if err != nil {
		t.Fatal(err)
	}

	r := NewErrorRenderer(nil)
	b := &bytes.Buffer{}
	if err = r.RenderError(b, msg); err != nil {
		t.Fatal(err)
	}
	g := goldie.New(t)
	g.Assert(t, "errors/500", b.Bytes())
}
