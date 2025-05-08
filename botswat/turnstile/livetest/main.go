package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dkotik/oakacs/oakhttp"
	"github.com/dkotik/oakacs/oakhttp/botswat"
	"github.com/dkotik/oakacs/oakhttp/botswat/turnstile"
	"github.com/dkotik/oakacs/oakhttp/oakserver"
	"golang.org/x/exp/slog"
)

func main() {
	logger := oakserver.NewDebugLogger()
	slog.SetDefault(logger)
	err := func() error {
		verifier, err := turnstile.New(
			turnstile.WithHostname("localhost"),
			turnstile.WithAllowedActions("view"),
		)
		if err != nil {
			return fmt.Errorf("cannot initialize turnstile: %w", err)
		}
		gate, err := turnstile.NewGate(
			turnstile.WithSiteAction("view"),
		)
		if err != nil {
			return fmt.Errorf("cannot initialize encoder: %w", err)
		}
		botswat, err := botswat.New(
			botswat.WithCache(botswat.NewMapCache(time.Minute, 20)),
			botswat.WithVerifier(verifier),
		)
		if err != nil {
			return fmt.Errorf("cannot initialize botswat: %w", err)
		}

		slog.Info("started server on http://localhost:8080")
		return oakserver.Run(
			context.Background(),
			// oakserver.WithDebugOptions(),
			oakserver.WithOakHandler(
				botswat.Gate(gate)(
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
