package oakhttp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testResponse struct {
	Body string
}

type testRequest struct {
	Body string
}

func (t *testRequest) Validate() error { return nil }

// func (t *testRequest) Normalize() error { return nil }

func newTestRequest(data any, contentType string) *http.Request {
	codec, err := GetCodec(contentType)
	if err != nil {
		panic(err)
	}

	b := &bytes.Buffer{}
	if err = codec.Encode(b, data); err != nil {
		panic(err)
	}
	r := httptest.NewRequest(
		http.MethodPost,
		"/",
		bytes.NewReader(b.Bytes()),
	)
	r.Header.Set("Content-Type", contentType)
	return r
}

func TestRequestResponseAdaptor(t *testing.T) {
	w := httptest.NewRecorder()
	adaptor := NewAdaptor(
		func(ctx context.Context, r *testRequest) (any, error) {
			if r.Body != "test body" {
				t.Fatal("request handler failed", fmt.Errorf("test body does not match: %s", r.Body))
			}
			return nil, nil
		},
	)

	var err error
	if err != nil {
		t.Fatal("could not create handler", err)
	}

	err = adaptor.ServeHyperText(
		w, newTestRequest(
			map[string]string{"Body": "test body"},
			"application/json",
		),
	)
	if err != nil {
		t.Fatal("request handler failed", err)
	}

	response := w.Result()
	defer response.Body.Close()
	var b bytes.Buffer
	if _, err = io.Copy(&b, response.Body); err != nil {
		t.Fatal("failed to read response", err)
	}
	if value := b.String(); value != "null\n" {
		t.Fatal("buffer should be null followed by new line:", value)
	}
}

func TestComplexRequestResponseAdaptor(t *testing.T) {
	w := httptest.NewRecorder()
	adaptor := NewComplexAdaptor(
		func(ctx context.Context, r *testRequest) (any, error) {
			if r.Body != "test body" {
				t.Fatal("request handler failed", fmt.Errorf("test body does not match: %s", r.Body))
			}
			return nil, nil
		},
		func(request *testRequest, r *http.Request) error {
			t.Logf("finalizing %+v request", request)
			return nil
		},
	)

	var err error
	if err != nil {
		t.Fatal("could not create handler", err)
	}

	err = adaptor.ServeHyperText(
		w, newTestRequest(
			map[string]string{"Body": "test body"},
			"application/json",
		),
	)
	if err != nil {
		t.Fatal("request handler failed", err)
	}

	response := w.Result()
	defer response.Body.Close()
	var b bytes.Buffer
	if _, err = io.Copy(&b, response.Body); err != nil {
		t.Fatal("failed to read response", err)
	}
	if value := b.String(); value != "null\n" {
		t.Fatal("buffer should be null followed by new line:", value)
	}
}
