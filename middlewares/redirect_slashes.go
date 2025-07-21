package middlewares

import (
	"github.com/dracory/rtr"
	"github.com/go-chi/chi/v5/middleware"
)

func RedirectSlashesMiddleware() rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Redirect Slashes Middleware").
		SetHandler(middleware.RedirectSlashes)
}
