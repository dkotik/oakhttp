package turnstile

import (
	_ "embed" // for the template
	"fmt"
	"net/http"
	"text/template"

	"github.com/dkotik/oakacs/oakhttp"
)

//go:embed htmlencoder.html
var templateHTML string

func NewEncoderHTML(siteKey string) oakhttp.Encoder {
	t, err := template.New("root").Parse(templateHTML)
	if err != nil {
		panic(fmt.Errorf("could not render Turnstile HTML template: %w", err))
	}
	return func(w http.ResponseWriter, any any) (err error) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusForbidden)
		return t.Execute(w, siteKey)
	}
}
