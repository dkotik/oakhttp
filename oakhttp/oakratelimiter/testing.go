package oakratelimiter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dkotik/oakacs/oakhttp"
)

type RequestFactory func(context.Context) *http.Request

func GetRequestFactory(ctx context.Context) *http.Request {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	return r.WithContext(ctx)
}

func MiddlewareLoadTest(
	ctx context.Context,
	m oakhttp.Middleware,
	r Rate,
	rf RequestFactory,
	expectedRejectionRate float64,
) func(t *testing.T) {
	return func(t *testing.T) {
		handler := m(func(w http.ResponseWriter, r *http.Request) error {
			return nil // do nothing
		})

		requests := make(chan *http.Request, 1)
		oneTokenWindow := time.Nanosecond * time.Duration(1/r)
		ticker := time.NewTicker(oneTokenWindow)
		defer ticker.Stop()

		passed := 0
		rejected := 0

		var err error
		for {
			select {
			case <-ctx.Done():
				if passed == 0 {
					t.Fatal("no requests succeeded:", passed, "out of", rejected)
					return
				}
				if expectedRejectionRate == 0 && rejected > 0 {
					t.Fatalf("%d requests were rejected when 0%% rejection rate was expected", rejected)
				}
				actualRejectionRate := float64(rejected) / float64(passed+rejected)
				if !floatComparator(0.1)(expectedRejectionRate, actualRejectionRate) {
					t.Fatal(
						"expected rejection rate is not close enough to the actual",
						expectedRejectionRate,
						"vs",
						actualRejectionRate,
					)
				}
				return
			case <-ticker.C:
				requests <- rf(ctx)
			case request := <-requests:
				if request == nil {
					t.Error("received a <nil> request")
					continue
				}
				w := httptest.NewRecorder()
				err = handler(w, request)
				if err == nil {
					passed++
					continue
				}

				httpError, ok := err.(oakhttp.Error)
				if !ok {
					t.Fatal("unexpected error:", err)
					return
				}
				if code := httpError.HTTPStatusCode(); code != http.StatusTooManyRequests {
					t.Fatal("status code mismatch:", code, "vs", http.StatusTooManyRequests)
					return
				}
				rejected++
			}
		}
	}
	// cases := []struct {
	// 	StatusCode int
	// 	Sleep      time.Duration
	// 	RemoteAddr string
	// }{
	// 	{
	// 		StatusCode: http.StatusOK,
	// 		Sleep:      0,
	// 		RemoteAddr: "localhost",
	// 	},
	// 	{
	// 		StatusCode: http.StatusOK,
	// 		Sleep:      0,
	// 		RemoteAddr: "localhost",
	// 	},
	// 	{
	// 		StatusCode: http.StatusTooManyRequests,
	// 		Sleep:      time.Millisecond * 1050,
	// 		RemoteAddr: "localhost",
	// 	},
	// 	{
	// 		StatusCode: http.StatusOK,
	// 		Sleep:      0,
	// 		RemoteAddr: "localhost",
	// 	},
	// 	{
	// 		StatusCode: http.StatusTooManyRequests,
	// 		Sleep:      time.Millisecond * 1050,
	// 		RemoteAddr: "localhost",
	// 	},
	// 	{
	// 		StatusCode: http.StatusOK,
	// 		Sleep:      0,
	// 		RemoteAddr: "localhost",
	// 	},
	// 	{
	// 		StatusCode: http.StatusTooManyRequests,
	// 		Sleep:      0,
	// 		RemoteAddr: "localhost",
	// 	},
	// }
	//
	// handler := Must(New(
	// 	testHandler,
	// 	WithLimit(2, time.Second),
	// ))
	//
	// request := httptest.NewRequest(http.MethodGet, "/", nil)
	// for i, c := range cases {
	// 	request.RemoteAddr = c.RemoteAddr
	// 	r := captureResponse(handler, request)
	// 	if r.StatusCode != c.StatusCode {
	// 		t.Fatalf("rate limiter step %d failed: %d does not match %d",
	// 			i+1, r.StatusCode, c.StatusCode)
	// 	}
	// 	time.Sleep(c.Sleep)
	// }

	return nil
}
