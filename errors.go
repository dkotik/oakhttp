package oakhttp

import (
	"cmp"
	"io"
	"log/slog"
	"net/http"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Error struct {
	StatusCode  int
	Title       *i18n.LocalizeConfig
	Description *i18n.LocalizeConfig
	Message     *i18n.LocalizeConfig
	Cause       error
}

func (e Error) GetHyperTextStatusCode() int {
	return e.StatusCode
}

func (e Error) Unwrap() error {
	return e.Cause
}

func (e Error) Error() string {
	if e.Cause == nil {
		return http.StatusText(e.StatusCode)
	}
	return http.StatusText(e.StatusCode) + ": " + e.Cause.Error()
}

func (e Error) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("status_code", e.StatusCode),
		slog.String("message", e.Error()),
		slog.Any("cause", e.Cause),
	)
}

type ErrorWithStatusCode interface {
	GetHyperTextStatusCode() int
}

type ErrorHandler interface {
	HandleError(http.ResponseWriter, *http.Request, error)
}

type ErrorHandlerFunc func(http.ResponseWriter, *http.Request, error)

func (f ErrorHandlerFunc) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	f(w, r, err)
}

type ErrorHandlerMiddleware interface {
	WrapErrorHandler(ErrorHandler) ErrorHandler
}

type ErrorRenderer interface {
	RenderError(io.Writer, Error) error
}

type errorHandler struct {
	Renderer ErrorRenderer
	Logger   *slog.Logger
}

func NewErrorHandler(r ErrorRenderer, logger *slog.Logger) ErrorHandler {
	return errorHandler{
		Renderer: r,
		Logger:   cmp.Or(logger, slog.Default()),
	}
}

func (h errorHandler) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	h.Logger.Log(
		r.Context(),
		slog.LevelError,
		err.Error(),
		// slog.Any("details", err),
		slog.Group("request",
			slog.String("host", r.URL.Hostname()),
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
		),
	)
}

type bufferedErrorHandler struct {
	Next   ErrorHandler
	Cache  map[int]map[language.Tag][]byte
	Logger *slog.Logger
}

func NewBufferedErrorHandler(
	b *i18n.Bundle,
	r ErrorRenderer,
	l *slog.Logger,
	override ...Error,
) interface {
	ErrorHandler
	ErrorHandlerMiddleware
} {
	// if len(override) == 0 {
	// 	return NewErrorHandler(r, l)
	// }

	return bufferedErrorHandler{
		Logger: cmp.Or(l, slog.Default()),
	}
}

func (h bufferedErrorHandler) HandleError(w http.ResponseWriter, r *http.Request, err error) {

}

func (h bufferedErrorHandler) WrapErrorHandler(e ErrorHandler) ErrorHandler {
	return bufferedErrorHandler{
		Next:   e, // TODO: cmp.Or(e, http.Error(w http.ResponseWriter, error string, code int))
		Cache:  h.Cache,
		Logger: h.Logger,
	}
}
