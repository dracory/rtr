package middlewares

import (
	"time"

	"github.com/dracory/rtr"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// TimeoutMiddleware returns a middleware that adds a timeout to the request context.
// If the request takes longer than the specified duration, it will be canceled.
// This is a thin wrapper around Chi's Timeout middleware.
func TimeoutMiddleware(timeout time.Duration) rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Timeout").
		SetHandler(chimiddleware.Timeout(timeout))
}
