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

func (t *testRequest) Validate() error  { return nil }
func (t *testRequest) Normalize() error { return nil }

func TestRequestAdaptor(t *testing.T) {
	w := httptest.NewRecorder()
	handler, err := NewRequestAdaptor(
		func(ctx context.Context, r *testRequest) error {
			if r.Body != "test body" {
				t.Fatal("request handler failed", fmt.Errorf("test body does not match: %s", r.Body))
			}
			return nil
		},
	)
	if err != nil {
		t.Fatal("could not create handler", err)
	}

	err = handler(
		w,
		httptest.NewRequest(
			http.MethodPost,
			"/",
			bytes.NewReader([]byte(`{"Body":"test body"}`)),
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
	if b.Len() > 0 {
		t.Fatal("buffer should be empty:", b.String())
	}
}
