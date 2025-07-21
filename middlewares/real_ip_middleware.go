package middlewares

import (
	"github.com/dracory/rtr"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// RealIPMiddleware returns a middleware that sets the client's real IP address in the request context.
// It uses Chi's RealIP middleware internally.
func RealIPMiddleware() rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Real IP").
		SetHandler(chimiddleware.RealIP)
}
