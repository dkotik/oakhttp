package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tigerperformanceinstitute/tprograms/service/oakhttp"
	"github.com/tigerperformanceinstitute/tprograms/service/oakhttp/turnstile"
	"golang.org/x/exp/slog"
)

func main() {
	fmt.Println("export TURNSTILE_SITE_KEY and TURNSTILE_SECRET_KEY before running")

	err := func() error {
		gate, err := turnstile.NewMiddleware(
			turnstile.WithAuthenticatorOptions(
				turnstile.WithClientOptions(
					turnstile.WithHostname("localhost"),
				),
			),
		)
		if err != nil {
			return err
		}

		return oakhttp.Run(
			context.Background(),
			// oakhttp.WithDebugOptions(),
			oakhttp.WithHandler(
				gate(oakhttp.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) error {
						fmt.Fprintf(w, "though shalt pass")
						return nil
					},
				)),
			),
		)
	}()

	if err != nil {
		slog.Error("server failed execution", slog.Any("error", err))
	}
}
