package oakserver

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/dkotik/oakacs/oakhttp"
)

func TestServer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	go func() {
		err := Run(
			ctx,
			WithOakHandler(
				func(w http.ResponseWriter, r *http.Request) error {
					return oakhttp.NewNotFoundError("test page")
				},
				oakhttp.NewErrorHandlerJSON(nil),
			),
		)
		if err != nil {
			t.Fatal(err)
		}
	}()

	time.Sleep(time.Millisecond * 100)

	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatal("unexpected status code:", resp.StatusCode)
	}
}
