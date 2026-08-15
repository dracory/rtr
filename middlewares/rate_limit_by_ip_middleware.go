package middlewares

import (
	"time"

	"github.com/dracory/rtr"
)

// RateLimitByIPMiddleware returns a middleware that limits the number of
// requests per time window from the same client IP address. The IP is extracted
// from r.RemoteAddr (the TCP peer address). IPv6 addresses are bucketed by
// their /64 prefix.
// This is a stdlib-only reimplementation of go-chi/httprate's LimitByIP.
func RateLimitByIPMiddleware(maxRequests int, seconds int) rtr.MiddlewareInterface {
	rl := newRateLimiter(maxRequests, time.Duration(seconds)*time.Second, keyByIP)
	return rtr.NewMiddleware().
		SetName("Rate Limit By IP").
		SetHandler(rl.handler)
}
