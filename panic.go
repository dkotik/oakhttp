package oakhttp

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
)

type PanicError struct {
	Cause any
	Stack string
}

func (e *PanicError) Error() string {
	return fmt.Sprintf("panic: %v", e.Cause)
}

func (e *PanicError) GetHyperTextStatusCode() int {
	return http.StatusInternalServerError
}

func (e *PanicError) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("error", e.Error()),
		slog.Int("status_code", e.GetHyperTextStatusCode()),
		slog.String("stack", e.Stack),
	)
}

type panicHandler struct {
	Next         http.Handler
	ErrorHandler ErrorHandler
}

func NewPanicHandler(eh ErrorHandler) Middleware {
	if eh == nil {
		eh = NewErrorHandler(nil, nil)
	}
	return func(next http.Handler) http.Handler {
		if next == nil {
			panic("next HTTP handler is nil")
		}

		return panicHandler{
			Next:         next,
			ErrorHandler: eh,
		}
	}
}

func (h panicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if recovery := recover(); recovery != nil {
			buf := make([]byte, 10<<10)
			n := runtime.Stack(buf, false)
			h.ErrorHandler.HandleError(w, r, &PanicError{
				// TODO: would debug.Stack() be better?
				Cause: recovery,
				Stack: string(buf[:n]),
			})
		}
	}()
	h.Next.ServeHTTP(w, r)
}
