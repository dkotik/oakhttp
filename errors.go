package oakhttp

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"maps"
	"net/http"
	"slices"

	"github.com/dkotik/oakhttp/internal/msg"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Error struct {
	StatusCode    int
	KnowledgeCode string
	Title         *i18n.LocalizeConfig
	Description   *i18n.LocalizeConfig
	Message       *i18n.LocalizeConfig
	Cause         error
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

type LocalizedError struct {
	LanguageTag   string
	StatusCode    int
	KnowledgeCode string
	Title         string
	Description   string
	Message       string
	// Cause       func() string
}

func (e Error) Localize(lc *i18n.Localizer) (localized LocalizedError, err error) {
	title, tag, err := lc.LocalizeWithTag(e.Title)
	if err != nil {
		return localized, fmt.Errorf("unable to localize error title: %w", err)
	}
	localized.Title = title
	localized.LanguageTag = tag.String()

	localized.Description, err = lc.Localize(e.Description)
	if err != nil {
		return localized, fmt.Errorf("unable to localize error description: %w", err)
	}
	localized.Message, err = lc.Localize(e.Message)
	if err != nil {
		return localized, fmt.Errorf("unable to localize error message: %w", err)
	}
	localized.StatusCode = e.StatusCode
	localized.KnowledgeCode = e.KnowledgeCode
	// localized.Cause = e.Cause.Error
	return localized, nil
}

type ErrorWithStatusCode interface {
	error
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
	RenderError(io.Writer, LocalizedError) error
}

type ErrorRendererFunc func(io.Writer, LocalizedError) error

func (f ErrorRendererFunc) RenderError(w io.Writer, err LocalizedError) error {
	return f(w, err)
}

func NewErrorRenderer(t *template.Template) ErrorRenderer {
	if t == nil {
		et, err := Templates.ReadFile("internal/templates/page/error.html")
		if err != nil {
			panic(fmt.Sprintf("unable to load error page template: %v", err))
		}
		t, err = template.New("error").Parse(string(et))
		if err != nil {
			panic(fmt.Sprintf("unable to parse error page template: %v", err))
		}
	}
	return ErrorRendererFunc(func(w io.Writer, err LocalizedError) error {
		return t.Execute(w, err)
	})
}

type staticErrorHandler struct {
	StatusCode int
	Content    []byte
	Logger     *slog.Logger
}

func (eh staticErrorHandler) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(eh.StatusCode)
	_, _ = io.Copy(w, bytes.NewReader(eh.Content))
	eh.Logger.Log(
		r.Context(),
		slog.LevelError,
		"HTTP request failed",
		slog.Any("error", err),
		slog.Group("request",
			slog.String("host", r.URL.Hostname()),
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
		),
	)
}

type errorHandler struct {
	LocalizerBundle *i18n.Bundle
	Renderer        ErrorRenderer
	Logger          *slog.Logger
}

func NewErrorHandler(
	localizationBunble *i18n.Bundle,
	r ErrorRenderer,
	logger *slog.Logger,
	static ...Error,
) ErrorHandler {
	if localizationBunble == nil {
		// TODO: replace with package level bundle
		localizationBunble = i18n.NewBundle(language.English)
	}
	languages := localizationBunble.LanguageTags()
	if len(languages) == 0 {
		panic("localization bundle contains zero translation languages")
	}
	if r == nil {
		r = NewErrorRenderer(nil)
	}
	if logger == nil {
		logger = slog.Default()
	}
	if len(static) == 0 {
		return errorHandler{
			LocalizerBundle: localizationBunble,
			Renderer:        r,
			Logger:          logger,
		}
	}

	byStatusCode := make(map[int]ErrorHandler)
	var (
		statusCode int
		ok         bool
	)

	for _, errWithStatusCode := range static {
		statusCode = errWithStatusCode.StatusCode
		if _, ok = byStatusCode[statusCode]; ok {
			panic(fmt.Sprintf("error with status code %d occurs twice", statusCode))
		}

		byLanguageTag := make(map[language.Tag]ErrorHandler)
		for _, lang := range languages {
			localized, err := errWithStatusCode.Localize(i18n.NewLocalizer(localizationBunble, lang.String()))
			if err != nil {
				panic(fmt.Errorf("unable to localize error: %w", err))
			}
			b := &bytes.Buffer{}
			if err = r.RenderError(b, localized); err != nil {
				panic(fmt.Errorf("unable to render a localized error: %w", err))
			}
			byLanguageTag[lang] = staticErrorHandler{
				StatusCode: statusCode,
				Content:    b.Bytes(),
			}
		}
		byStatusCode[statusCode] = NewErrorHandlerSwitchByLanguage(byLanguageTag, language.PreferSameScript(true))
	}
	if _, ok = byStatusCode[http.StatusInternalServerError]; !ok {
		byStatusCode[http.StatusInternalServerError] = errorHandler{
			LocalizerBundle: localizationBunble,
			Renderer:        r,
			Logger:          logger,
		}
	}
	return NewErrorHandlerSwitchByStatusCode(byStatusCode)
}

func (h errorHandler) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	h.Logger.Log(
		r.Context(),
		slog.LevelError,
		"HTTP request failed",
		slog.Any("error", err),
		slog.Group("request",
			slog.String("host", r.URL.Hostname()),
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
		),
	)

	var localizableError Error
	if !errors.As(err, &localizableError) {
		localizableError = NewError(err, "")
	}
	w.WriteHeader(localizableError.StatusCode)
	lc := i18n.NewLocalizer(h.LocalizerBundle, r.Header.Get("Accept-Language"))
	localized, err := localizableError.Localize(lc)
	if err != nil {
		// panic(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		h.Logger.Log(
			r.Context(),
			slog.LevelError,
			"error localization failed",
			slog.Any("error", err),
			slog.Group("request",
				slog.String("host", r.URL.Hostname()),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
			),
		)
		return
	}

	if err = h.Renderer.RenderError(w, localized); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		h.Logger.Log(
			r.Context(),
			slog.LevelError,
			"error rendering failed",
			slog.Any("error", err),
			slog.Group("request",
				slog.String("host", r.URL.Hostname()),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
			),
		)
	}
}

type ehByStatusCode map[int]ErrorHandler

func NewErrorHandlerSwitchByStatusCode(sw map[int]ErrorHandler) ErrorHandler {
	if len(sw) == 0 {
		panic("empty error handler switch")
	}
	if _, ok := sw[http.StatusInternalServerError]; !ok {
		panic(fmt.Sprintf("default status handler is missing: %d", http.StatusInternalServerError))
	}

	return ehByStatusCode(maps.Clone(sw))
}

func (eh ehByStatusCode) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	var errorWithStatusCode ErrorWithStatusCode
	if errors.As(err, &errorWithStatusCode) {
		h, ok := eh[errorWithStatusCode.GetHyperTextStatusCode()]
		if ok {
			h.HandleError(w, r, err)
			return
		}
	}
	eh[http.StatusInternalServerError].HandleError(w, r, err)
}

type ehByLanguage struct {
	Matcher  language.Matcher
	Handlers map[language.Tag]ErrorHandler
}

func NewErrorHandlerSwitchByLanguage(sw map[language.Tag]ErrorHandler, options ...language.MatchOption) ErrorHandler {
	if len(sw) == 0 {
		panic("empty error handler switch")
	}

	return ehByLanguage{
		Matcher:  language.NewMatcher(slices.Collect(maps.Keys(sw)), options...),
		Handlers: maps.Clone(sw),
	}
}

func (eh ehByLanguage) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	tag, _, _ := eh.Matcher.Match(language.Make(r.Header.Get("Accept-Language")))
	eh.Handlers[tag].HandleError(w, r, err)
}

func NewError(from error, knowledgeCode string) Error {
	return Error{
		StatusCode:  http.StatusInternalServerError,
		Title:       msg.ErrorInternalTitle,
		Description: msg.ErrorInternalDescription,
		Message:     msg.ErrorInternalDescription,
		Cause:       from,
	}
}

func NewNotFoundError(from error, knowledgeCode string) Error {
	return Error{
		StatusCode:    http.StatusNotFound,
		KnowledgeCode: knowledgeCode,
		Title:         msg.ErrorNotFoundTitle,
		Description:   msg.ErrorNotFoundDescription,
		Message:       msg.ErrorNotFoundDescription,
		Cause:         from,
	}
}

func NewAccessDeniedError(from error, knowledgeCode string) Error {
	return Error{
		StatusCode:    http.StatusForbidden,
		KnowledgeCode: knowledgeCode,
		Title:         msg.ErrorAccessDeniedTitle,
		Description:   msg.ErrorAccessDeniedDescription,
		Message:       msg.ErrorAccessDeniedDescription,
		Cause:         from,
	}
}
