package mux

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tigerperformanceinstitute/tprograms/service/oakhttp"
)

func TestMux(t *testing.T) {
	mux := New(
		WithRoute(
			"firstRoute",
			"/test/{pattern}/yep/{$}",
			oakhttp.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
				t.Log("test request suceeded")
				return nil
			}),
		),
		WithRoute(
			"secondRoute",
			"/test/{wild}/{pattern1}/yep/{$}",
			oakhttp.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
				t.Log("test request suceeded")
				return nil
			}),
		),
	)

	r := httptest.NewRequest(
		http.MethodPost,
		"/test/something/yep/",
		nil,
	)
	w := httptest.NewRecorder()

	if err := mux.ServeHyperText(w, r); err != nil {
		t.Fatal("mux route failed:", err)
	}
}
