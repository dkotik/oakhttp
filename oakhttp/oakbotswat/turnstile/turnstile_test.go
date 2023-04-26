package turnstile

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/dkotik/oakacs/oakhttp"
	"github.com/dkotik/oakacs/oakhttp/oakbotswat"
	"github.com/dkotik/oakacs/oakhttp/oakserver"
)

func TestIntegration(t *testing.T) {
	verifier, err := New(
		WithHostname("localhost"),
		WithAllowedActions("view"),
	)
	if err != nil {
		t.Fatal("cannot initiate turnstile:", err)
	}
	botswat, err := oakbotswat.New(
		oakbotswat.WithVerifier(verifier),
		oakbotswat.WithEncoder(NewEncoderHTML(os.Getenv("TURNSTILE_SITE_KEY"))),
	)
	if err != nil {
		t.Fatal("cannot initiate OakBotSWAT:", err)
	}

	t.Log("started server on http://localhost:8080")
	err = oakserver.Run(
		context.Background(),
		oakserver.WithOakHandler(
			botswat(
				func(w http.ResponseWriter, r *http.Request) error {
					fmt.Fprintf(w, "though shalt pass")
					return nil
				},
			),
			oakhttp.NewErrorHandlerJSON(nil),
		),
	)
	if err != nil {
		t.Fatal("cannot initiate OakServer:", err)
	}
}
