package oakoidc

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

var serve = flag.Bool("serve", false, "listen for HTTP requests")

func TestOIDC(t *testing.T) {
	flag.Parse()
	if !*serve {
		t.Skip("use `clear; go test --v --serve --run=^TestOIDC$ .` command to run the OIDC test")
	}

	login, callback, err := New(
		WithDiscoveryURL("https://accounts.google.com"),
		WithCallbackURL("http://localhost:8081/callback"),
		WithSessionAdapter(func(
			ctx context.Context,
			token *oauth2.Token,
			info *oidc.UserInfo,
		) (string, error) {
			t.Log(fmt.Sprintf("%+v", token))
			t.Log(fmt.Sprintf("%s // %v", info.Email, info.EmailVerified))

			claims, err := NewStandardClaims(info)
			if err != nil {
				return "", err
			}

			t.Log(claims.String())
			return "/finish", nil
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := login(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if err := callback(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	http.HandleFunc("/finish", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "OIDC callback succeded! Check logs to see token and claims.")
	})

	fmt.Println("listening on http://localhost:8081/")
	err = http.ListenAndServe("localhost:8081", nil)
	if err != nil {
		t.Fatal(err)
	}
}
