package cors

import (
	"errors"
	"fmt"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"

	"github.com/tigerperformanceinstitute/tprograms/service/oakhttp"
)

const varyKey = "Vary"

var (
	varyHeader          = []string{"Origin"}
	varyPreflightHeader = []string{"Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers"}
)

type corsMiddleware struct {
	next             oakhttp.Handler
	filter           OriginFilter
	allowedMethods   []string
	preflightHeaders textproto.MIMEHeader
	responseHeaders  textproto.MIMEHeader
}

func New(withOptions ...Option) oakhttp.Middleware {
	o := &options{}

	var err error
	for _, option := range append(
		withOptions,
		WithDefaultMethodsGetPostHead(),
		WithDefaultMaxAgeOfOneWeek(),
		func(o *options) error { // finalize filters
			exactOrigins := make([]string, 0)
			var index = 0
			for _, origin := range o.allowedOrigins {
				index = strings.IndexRune(origin, '*')
				if index == -1 {
					exactOrigins = append(exactOrigins, origin)
				} else if origin == "*" {
					o.filters = append(o.filters, func(_ string) bool {
						return true
					})
				} else {
					prefix := origin[:index]
					suffix := origin[index+1:]
					o.filters = append(o.filters, func(s string) bool {
						return len(s) >= len(prefix)+len(suffix) && strings.HasPrefix(s, prefix) && strings.HasSuffix(s, suffix)
					})
				}
			}

			if len(exactOrigins) > 0 {
				o.filters = append(o.filters, func(s string) bool {
					for _, origin := range exactOrigins {
						if origin == s {
							return true
						}
					}
					return false
				})
			}
			if len(o.filters) == 0 {
				return errors.New("cross-origin resource sharing requires at least one allowed origin or origin filter")
			}
			return nil
		},
	) {
		if err = option(o); err != nil {
			panic(fmt.Errorf("cannot initialize cross-origin resource sharing middleware: %w", err))
		}
	}

	return func(next oakhttp.Handler) oakhttp.Handler {
		mw := &corsMiddleware{
			next:           next,
			allowedMethods: o.allowedMethods,
			preflightHeaders: textproto.MIMEHeader(map[string][]string{
				"Access-Control-Allow-Methods": []string{
					strings.Join(o.allowedMethods, ","),
				},
				"Access-Control-Allow-Headers": []string{
					strings.Join(o.allowedHeaders, ","),
				},
				"Access-Control-Max-Age": []string{
					strconv.Itoa(int(o.maxAge.Seconds())),
				},
			}),
			responseHeaders: textproto.MIMEHeader(map[string][]string{
				"Access-Control-Expose-Headers": []string{
					strings.Join(o.exposedHeaders, ","),
				},
			}),
		}

		if o.allowCredentials != nil && *o.allowCredentials {
			truth := []string{"true"}
			mw.preflightHeaders["Access-Control-Allow-Credentials"] = truth
			mw.responseHeaders["Access-Control-Allow-Credentials"] = truth
		}

		if len(o.filters) == 1 {
			mw.filter = o.filters[0]
		} else {
			mw.filter = func(origin string) bool {
				for _, filter := range o.filters {
					if filter(origin) {
						return true
					}
				}
				return false
			}
		}

		return mw
	}
}

func (cors *corsMiddleware) isMethodAllowed(method string) bool {
	for _, m := range cors.allowedMethods {
		if m == method {
			return true
		}
	}
	return false
}

func (cors *corsMiddleware) writeOrigin(h http.Header, o string) {
	// MDN: When responding to a credentialed requests request, the server must specify an origin in the value of the Access-Control-Allow-Origin header, instead of specifying the "*" wildcard.
	h["Access-Control-Allow-Origin"] = []string{o}
}

func (cors *corsMiddleware) ServeHyperText(
	w http.ResponseWriter,
	r *http.Request,
) error {
	headers := w.Header()
	origin := r.Header.Get("Origin")
	if origin == "" {
		headers[varyKey] = varyHeader // always set
	} else if r.Method == http.MethodOptions {
		headers[varyKey] = varyPreflightHeader // always set
		if cors.filter(origin) {
			cors.writeOrigin(headers, origin)
			for key, value := range cors.preflightHeaders {
				headers[key] = value
			}
		}
	} else {
		headers[varyKey] = varyHeader // always set
		if !cors.isMethodAllowed(r.Method) {
			return oakhttp.NewMethodNotAllowedError(r.Method)
		}
		if !cors.filter(origin) {
			return oakhttp.NewAccessDeniedError(fmt.Errorf("origin %q is not available for cross-origin request sharing", origin))
		}
		cors.writeOrigin(headers, origin)
		for key, value := range cors.responseHeaders {
			headers[key] = value
		}
	}

	return cors.next.ServeHyperText(w, r)
}

func (cors *corsMiddleware) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	if err := cors.ServeHyperText(w, r); err != nil {
		var httpError oakhttp.Error
		if errors.As(err, &httpError) {
			http.Error(w, httpError.Error(), httpError.HyperTextStatusCode())
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
