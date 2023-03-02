package oakhttp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

// TODO: some excellent suggestions here, regarding UUID extraction: https://haykot.dev/blog/reduce-boilerplate-in-go-http-handlers-with-go-generics/
// TODO: create Adapter with injected RequestStruct builder, so that REST APIs can be build by taking fields out of the URL path. Look at new exposed Match function in Fiber for inspiration of parameter extraction: https://github.com/gofiber/fiber/pull/2142

type Handler[IN any, OUT any] func(context.Context, *IN) (*OUT, error)

type PathTailHandler[OUT any] func(context.Context, string) (*OUT, error)

// Validatable is a generic interface that requires type T to be a pointer and implement the Validate method. It complements the adapter definitions. See https://stackoverflow.com/questions/72090387/what-is-the-generic-type-for-a-pointer-that-implements-an-interface
type Validatable[T any] interface {
	*T
	Validate() error
}

// https://www.youtube.com/watch?v=iWP0ANQ4m7g&list=WL&index=3
// https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#pointer-method-example
// paymentHandler := Adapt[PaymentIntent, *PaymentIntent, PaymentIntentResponse](readLimit, Pay)
// Type parameters are optional, when calling Adapt, because the compiler can infer them from the passed arguments.
func AdaptRequestToJSON[T any, P Validatable[T], OUT any](readLimit int64, h func(context.Context, P) (*OUT, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var (
			in        = new(T) //*IN
			err       error
			errStatus = http.StatusInternalServerError
		)
		defer func() {
			if err != nil {
				w.WriteHeader(errStatus)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"Error": err.Error(),
				})
			}
		}()

		// in = new(IN)
		// const MaxBodyBytes = int64(65536)
		// req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)
		// TODO: close req.Body?
		// err = json.NewDecoder(&io.LimitedReader{R: r.Body, N: readLimit}).Decode(in)
		reader := http.MaxBytesReader(w, r.Body, readLimit)
		defer reader.Close()
		// err = json.NewDecoder(&io.LimitedReader{R: r.Body, N: readLimit}).Decode(in)
		err = json.NewDecoder(reader).Decode(in)
		if err != nil {
			r.Body.Close()
			err = errors.New("JSON decoding failure: " + err.Error())
			return
		}
		r.Body.Close()

		// if in == nil {
		// 	err = errors.New("no post data provided")
		// 	return
		// }

		if err = P(in).Validate(); err != nil {
			err = fmt.Errorf("failed to validate: %w", err)
			return
		}

		out, err := h(r.Context(), in)
		if err != nil {
			return
		}
		if out == nil {
			errStatus = http.StatusNotFound
			err = errors.New("Not Found")
			return
		}

		if err = json.NewEncoder(w).Encode(out); err != nil {
			err = errors.New("JSON encoding failure: " + err.Error())
		}
	}
}

// AdaptRequestPathTailToJSON extracts the tail of a request path to pass to an RPC handler.
func AdaptRequestPathTailToJSON[OUT any](h func(context.Context, string) (*OUT, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func AdaptRequestToHTML[T any, P Validatable[T], OUT any](readLimit int64, t *template.Template, h func(context.Context, P) (*OUT, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// t.Execute(w, in)
	}
}

func AdaptRequestPathTailToHTML[OUT any](h func(context.Context, string) (*OUT, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

type RequestFactory[T any] func(*http.Request) (*T, error)

//
// func AdaptRequestFactoryToJSON[T any, P Validatable[T], OUT any](readLimit int64, f RequestFactory[T], h func(context.Context, P) (*OUT, error)) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 	}
// }
//
// func AdaptRequestFactoryToHTML[T any, P Validatable[T], OUT any](readLimit int64, t *html.Template, f RequestFactory[T], h func(context.Context, P) (*OUT, error)) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// t.Execute(w, in)
// 	}
// }
