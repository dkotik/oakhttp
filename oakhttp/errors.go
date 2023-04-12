package oakhttp

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"

	"golang.org/x/exp/slog"
)

type Error interface {
	error
	HTTPStatusCode() int
}

type ErrorHandler func(http.ResponseWriter, *http.Request, error)

func NewErrorHandlerMiddleware(h ErrorHandler) Middleware {
	if h == nil {
		h = NewErrorHandlerJSON(slog.Default())
	}
	return func(wrapped Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			if err := wrapped(w, r); err != nil {
				h(w, r, err)
			}
			return nil
		}
	}
}

func NewErrorHandlerJSON(l *slog.Logger) ErrorHandler {
	if l == nil {
		l = slog.Default()
	}
	return func(w http.ResponseWriter, r *http.Request, err error) {
		if err == nil {
			return
		}
		l.ErrorCtx(r.Context(), "OakHTTP request failed", slog.Any("error", err))
		w.Header().Set("Content-Type", "application/json")

		var httpError Error
		if errors.As(err, &httpError) {
			w.WriteHeader(httpError.HTTPStatusCode())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		multi, ok := err.(interface {
			Unwrap() []error // multi error
		})

		if ok {
			err = json.NewEncoder(w).Encode(multi.Unwrap())
		} else {
			err = json.NewEncoder(w).Encode(err)
		}

		if err != nil { // encoding failed
			l.ErrorCtx(r.Context(), "OakHTTP error handler encoder failed", slog.Any("error", err))
		}
	}
}

func NewErrorHandlerHTML(l *slog.Logger, t *template.Template) ErrorHandler {
	if l == nil {
		l = slog.Default()
	}
	if t == nil {
		t = template.Must(template.New("error").Parse(`
<html>
<head>
<title>Error Encountered</title>
<style type=text/css>
</style>
</head>

<body>
  <h1>Server Error</h1>
  <p>{{.}}</p>
  <p>Server failed to complete your request. If you believe this is a mistake, please contact support.</p>
</body>
</html>
    `))
	}

	return func(w http.ResponseWriter, r *http.Request, err error) {
		if err == nil {
			return
		}
		l.Error("OakHTTP request failed", slog.Any("error", err))
		w.Header().Set("Content-Type", "text/html")

		var httpError Error
		if errors.As(err, &httpError) {
			w.WriteHeader(httpError.HTTPStatusCode())
			err = t.Execute(w, httpError)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			err = t.Execute(w, err)
		}

		if err != nil { // encoding failed
			l.Error("OakHTTP error handler encoder failed", slog.Any("error", err))
		}
	}
}

type NotFoundError struct {
	resource string
}

func NewNotFoundError(resource string) *NotFoundError {
	return &NotFoundError{resource: resource}
}

func (e *NotFoundError) HTTPStatusCode() int {
	return http.StatusNotFound
}

func (e *NotFoundError) Error() string {
	return "resource \"" + e.resource + "\" was not found"
}

// func (e *NotFoundError) LogValue() slog.Value {
// 	return slog.GroupValue(
// 		slog.Any("resource", e.resource),
// 		slog.String("status_code", e.HTTPStatusCode()),
// 	)
// }
