package oakhttp

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testRequest struct {
	Body string
}

func (t *testRequest) Validate() error  { return nil }
func (t *testRequest) Normalize() error { return nil }

var testAdaptor = func() *DomainAdaptor {
	adaptor, err := New()
	if err != nil {
		panic(err)
	}
	return adaptor
}()

func TestRequestAdaptor(t *testing.T) {
	w := httptest.NewRecorder()
	AdaptRequest(
		testAdaptor,
		func(ctx context.Context, r *testRequest) error {
			if r.Body != "test body" {
				t.Fatal("request failed", fmt.Errorf("test body does not match: %s", r.Body))
			}
			return nil
		},
	).ServeHTTP(
		w,
		httptest.NewRequest(
			http.MethodPost,
			"/",
			bytes.NewReader([]byte(`{"Body":"test body"}`)),
		),
	)
}

func TestCustomRequestAdaptor(t *testing.T) {
	w := httptest.NewRecorder()
	AdaptCustomRequest(
		testAdaptor,
		func(r *http.Request) (string, error) {
			return r.URL.Path, nil
		},
		func(ctx context.Context, p string) error {
			if p != "/path" {
				t.Fatal(fmt.Errorf("test path does not match: %s", p))
			}
			return nil
		},
	).ServeHTTP(
		w,
		httptest.NewRequest(
			http.MethodPost,
			"/path",
			bytes.NewReader([]byte(`{"Body":"test body"}`)),
		),
	)
}
