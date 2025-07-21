package middlewares

import (
	"time"

	"github.com/dracory/rtr"
	"github.com/go-chi/httprate"
)

func RateLimitByIPMiddleware(maxRequests int, seconds int) rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Rate Limit By IP").
		SetHandler(httprate.LimitByIP(maxRequests, time.Duration(seconds)*time.Second))
}
