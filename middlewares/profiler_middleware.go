package middlewares

import (
	"net/http"
	"net/http/pprof"
	"strings"

	"github.com/dracory/rtr"
)

// ProfilerMiddleware returns a middleware that serves the Go pprof profiler.
// The profiler is available under /debug/pprof/ and exposes the standard
// net/http/pprof endpoints (index, cmdline, profile, symbol, trace, and the
// individual profile handlers).
// Make sure to only enable this in development environments as it exposes
// sensitive debugging information.
func ProfilerMiddleware() rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Profiler").
		SetHandler(func(next http.Handler) http.Handler {
			// Build a dedicated mux that serves the pprof endpoints so we
			// don't pollute http.DefaultServeMux.
			mux := http.NewServeMux()
			mux.HandleFunc("/debug/pprof/", pprof.Index)
			mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
			mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
			mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
			mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
			mux.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
			mux.Handle("/debug/pprof/block", pprof.Handler("block"))
			mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
			mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
			mux.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
			mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))

			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// If the request path is for the profiler, let the profiler mux handle it
				if strings.HasPrefix(r.URL.Path, "/debug/pprof/") {
					mux.ServeHTTP(w, r)
					return
				}
				// Otherwise, pass the request to the next handler
				next.ServeHTTP(w, r)
			})
		})
}
