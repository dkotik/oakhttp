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

func TestRequestResponseAdaptor(t *testing.T) {
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
