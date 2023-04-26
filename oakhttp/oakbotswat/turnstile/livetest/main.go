package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dkotik/oakacs/oakhttp"
	"github.com/dkotik/oakacs/oakhttp/oakbotswat"
	"github.com/dkotik/oakacs/oakhttp/oakbotswat/turnstile"
	"github.com/dkotik/oakacs/oakhttp/oakserver"
	"golang.org/x/exp/slog"
)

func main() {
	err := func() error {
		verifier, err := turnstile.New(
			turnstile.WithHostname("localhost"),
			turnstile.WithAllowedActions("view"),
		)
		if err != nil {
			return fmt.Errorf("cannot initiate turnstile: %w", err)
		}
		botswat, err := oakbotswat.New(
			oakbotswat.WithCache(oakbotswat.NewMapCache(time.Minute, 20)),
			oakbotswat.WithVerifier(verifier),
			oakbotswat.WithEncoder(turnstile.NewEncoderHTML(os.Getenv("TURNSTILE_SITE_KEY"))),
		)
		if err != nil {
			return fmt.Errorf("cannot initiate OakBotSWAT: %w", err)
		}

		slog.Info("started server on http://localhost:8080")
		return oakserver.Run(
			context.Background(),
			oakserver.WithDebugOptions(),
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
	}()

	if err != nil {
		slog.Error("server failed execution", slog.Any("error", err))
	}
}
