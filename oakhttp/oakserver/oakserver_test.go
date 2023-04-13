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

	logger := NewDebugLogger()

	go func() {
		err := Run(
			ctx,
			WithOakHandler(
				func(w http.ResponseWriter, r *http.Request) error {
					logger.InfoCtx(r.Context(), "trace IDs must match between two entries")
					panic("boo")
					return oakhttp.NewNotFoundError("test page")
				},
				oakhttp.NewErrorHandlerJSON(NewDebugLogger()),
			),
			WithLogger(logger),
			WithDebugOptions(),
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

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatal("unexpected status code:", resp.StatusCode)
	}
}
