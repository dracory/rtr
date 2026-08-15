package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/dracory/rtr"
)

// TimeoutMiddleware returns a middleware that adds a timeout to the request context.
// If the request takes longer than the specified duration, the context is canceled
// and a 504 Gateway Timeout response is written.
//
// It's required that you select the ctx.Done() channel to check for the signal
// if the context has reached its deadline and return, otherwise the timeout
// signal will be just ignored.
//
// Note: the 504 response is only written if the handler has not already written
// headers. If the handler calls WriteHeader before the deadline, the 504 is
// silently suppressed by the standard http.ResponseWriter (which only honors
// the first WriteHeader call). This matches chi's behaviour.
//
// Behaviour mirrors chi's middleware.Timeout.
func TimeoutMiddleware(timeout time.Duration) rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Timeout").
		SetHandler(timeoutHandler(timeout))
}

// timeoutHandler returns the stdlib-only handler that implements the timeout
// behaviour previously provided by chi's middleware.Timeout.
func timeoutHandler(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer func() {
				cancel()
				if ctx.Err() == context.DeadlineExceeded {
					w.WriteHeader(http.StatusGatewayTimeout)
				}
			}()

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
