package server

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/dkotik/oakhttp"
)

func TestServer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	logger := oakhttp.NewDebugLogger()

	errChannel := make(chan error)
	go func() {
		errChannel <- Run(
			ctx,
			WithHandler(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				},
			)),
			WithLogger(logger),
			WithDebugOptions(),
		)
	}()

	time.Sleep(time.Millisecond * 100)

	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatal("unexpected status code:", resp.StatusCode)
	}

	select {
	case <-time.After(time.Second):
		t.Fatal("server did not shut down within one second limit")
	case err := <-errChannel:
		if err != nil {
			t.Fatal("server shut down with an error:", err)
		}
	}
}
