package middlewares

import (
	"github.com/dracory/rtr"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// LoggerMiddleware returns a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the status code was, and how long it took.
// This is a thin wrapper around Chi's Logger middleware.
func LoggerMiddleware() rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Logger").
		SetHandler(chimiddleware.Logger)
}
