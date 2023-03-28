package oakhttp

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testResponse struct {
	Body string
}

type testRequest struct {
	Body string
}

func (t *testRequest) Validate() error  { return nil }
func (t *testRequest) Normalize() error { return nil }

func TestRequestResponseAdaptor(t *testing.T) {
	panic("http.MaxBytesReader should be added to new adaptors")

	w := httptest.NewRecorder()
	AdaptRequestResponse(
		testAdaptor,
		func(ctx context.Context, r *testRequest) (*testResponse, error) {
			if r.Body != "test body" {
				t.Fatal("request failed", fmt.Errorf("test body does not match: %s", r.Body))
			}
			return &testResponse{
				Body: r.Body + " response",
			}, nil
		},
	).ServeHTTP(
		w,
		httptest.NewRequest(
			http.MethodPost,
			"/",
			bytes.NewReader([]byte(`{"Body":"test body"}`)),
		),
	)

	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if strings.TrimSpace(string(data)) != `{"Body":"test body response"}` {
		t.Errorf("expected encoded response, got %v", string(data))
	}
}

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
