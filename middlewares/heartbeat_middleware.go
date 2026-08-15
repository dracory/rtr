package middlewares

import (
	"net/http"
	"strings"

	"github.com/dracory/rtr"
)

// HeartbeatMiddleware endpoint middleware
// useful to setting up a path like `/ping` that load balancers or
// uptime testing external services can make a request before hitting any routes.
// It's also convenient to place this above ACL middlewares as well.
//
// Behaviour mirrors chi's middleware.Heartbeat:
//   - Only GET and HEAD methods trigger the heartbeat response.
//   - Path matching is case-insensitive.
//   - The response is 200 OK with Content-Type: text/plain and body ".".
func HeartbeatMiddleware(endpoint string) rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Heartbeat Middleware at " + endpoint).
		SetHandler(heartbeatHandler(endpoint))
}

// heartbeatHandler returns the stdlib-only handler that implements the heartbeat
// behaviour previously provided by chi's middleware.Heartbeat.
func heartbeatHandler(endpoint string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if (r.Method == http.MethodGet || r.Method == http.MethodHead) && strings.EqualFold(r.URL.Path, endpoint) {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("."))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
