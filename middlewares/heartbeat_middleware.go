package middlewares

import (
	"github.com/dracory/rtr"
	"github.com/go-chi/chi/v5/middleware"
)

// HeartbeatMiddleware endpoint middleware
// useful to setting up a path like `/ping` that load balancers or
// uptime testing external services can make a request before hitting any routes.
// It's also convenient to place this above ACL middlewares as well.
func HeartbeatMiddleware(endpoint string) rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Heartbeat Middleware at " + endpoint).
		SetHandler(middleware.Heartbeat(endpoint))
}
