package middlewares

import (
	"net/http"
	"strings"

	"github.com/dracory/rtr"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// ProfilerMiddleware returns a middleware that serves the Go pprof profiler.
// It's a thin wrapper around Chi's Profiler middleware.
// The profiler will be available at the specified path (e.g., "/debug/pprof").
// Make sure to only enable this in development environments as it exposes
// sensitive debugging information.
func ProfilerMiddleware() rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Profiler").
		SetHandler(func(next http.Handler) http.Handler {
			// Create the profiler handler
			profiler := chimiddleware.Profiler()
			
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// If the request path is for the profiler, let the profiler handle it
				if strings.HasPrefix(r.URL.Path, "/debug/pprof/") {
					profiler.ServeHTTP(w, r)
					return
				}
				// Otherwise, pass the request to the next handler
				next.ServeHTTP(w, r)
			})
		})
}
