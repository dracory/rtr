package middlewares

import (
	"net/http"
	"time"

	"github.com/dracory/rtr"
	"github.com/go-chi/httprate"
)

// ThrottleMiddleware returns a middleware that limits the number of requests per time window.
// It uses the client's IP address to track request counts.
// This is a thin wrapper around go-chi/httprate's Limit function.
func ThrottleMiddleware(requests int, window time.Duration) rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Throttle").
		SetHandler(httprate.Limit(
			requests,
			window,
			httprate.WithKeyFuncs(func(r *http.Request) (string, error) {
				return r.RemoteAddr, nil
			}),
		))
}
