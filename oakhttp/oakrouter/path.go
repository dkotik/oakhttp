package oakrouter

import (
	"errors"
	"net/http"

	"github.com/dkotik/oakacs/oakhttp"
)

// ChompPath splits off the first segment of the request
// URL path as head from the remaining tail.
// The head segment never contains a slash.
// Path is not cleaned because it creates ambiguities.
// The same source may be fetched by using different paths.
//
// This function allows building branching routing handlers without
// having to rely on a multiplexer. Chomp multiple times in the
// same routing handler as needed. Before passing request to
// the next handler, overwrite [http.request.URL.Path].
//
// Based on ShiftPath: https://blog.merovius.de/posts/2017-06-18-how-not-to-use-an-http-router/
// See also: https://github.com/benhoyt/go-routing/blob/master/shiftpath/route.go
// And: https://benhoyt.com/writings/go-routing/
func ChompPath(p string) (head, tail string) {
	length := len(p)
	if length == 0 {
		return "", ""
	}

	var (
		to        int
		character rune
	)

	if p[0] != '/' {
		for to, character = range p {
			if character == '/' {
				return p[0:to], p[to:length]
			}
		}
		return p[:to+1], ""
	}

	for to, character = range p[1:length] {
		if character == '/' {
			to++
			return p[1:to], p[to:length]
		}
	}
	to += 2
	return p[1:to], p[to:length]
}

func TailTrailingSlashRedirectOrNotFound(w http.ResponseWriter, r *http.Request, tail string) error {
	for _, character := range tail {
		if character != '/' {
			return oakhttp.NewNotFoundError(r.URL.Path)
		}
	}

	URL := r.URL.String()
	if URL == "" {
		return errors.New("trailing slash redirect received request with empty URL")
	}
	http.Redirect(w, r, URL[:len(URL)-1], http.StatusTemporaryRedirect)
	return nil
}
