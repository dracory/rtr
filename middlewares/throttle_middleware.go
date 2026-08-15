package middlewares

import (
	"time"

	"github.com/dracory/rtr"
)

// ThrottleMiddleware returns a middleware that limits the number of requests per
// time window. It uses the client's RemoteAddr (including port) to track request
// counts.
// This is a stdlib-only reimplementation of go-chi/httprate's Limit function.
func ThrottleMiddleware(requests int, window time.Duration) rtr.MiddlewareInterface {
	rl := newRateLimiter(requests, window, keyByRemoteAddr)
	return rtr.NewMiddleware().
		SetName("Throttle").
		SetHandler(rl.handler)
}
