package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/dkotik/oakhttp/botswat/turnstile"
	"github.com/dkotik/oakhttp/server"
)

func main() {
	fmt.Println("export TURNSTILE_SITE_KEY and TURNSTILE_SECRET_KEY before running")

	err := func() error {
		gate, err := turnstile.NewMiddleware(
		// turnstile.WithAuthenticatorOptions(
		// 	turnstile.WithClientOptions(
		// 		turnstile.WithHostname("localhost"),
		// 	),
		// ),
		)
		if err != nil {
			return err
		}

		return server.Run(
			context.Background(),
			// oakhttp.WithDebugOptions(),
			server.WithHandler(
				gate(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprintf(w, "though shalt pass")
					},
				)),
			),
		)
	}()

	if err != nil {
		slog.Error("server failed execution", slog.Any("error", err))
	}
}
