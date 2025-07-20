package middlewares

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/dracory/rtr"
)

// RecoveryMiddleware creates a new middleware that recovers from panics.
// It logs the panic details and returns a 500 Internal Server Error response.
// This should typically be added as one of the first middlewares in the chain.
func RecoveryMiddleware() rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Recovery Middleware").
		SetHandler(recoveryHandler())
}

// recoveryHandler returns the default recovery handler
func recoveryHandler() rtr.StdMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if w == nil || r == nil {
				// If we don't have valid request/response objects, we can't do much
				return
			}

			// Create a response recorder to capture the response
			rw := newResponseRecorder(w)

			defer func() {
				if err := recover(); err != nil {
					// Log the error with stack trace
					log.Printf("Recovered from panic: %v\n%s", err, string(debug.Stack()))

					// Only write error response if nothing was written yet
					if !rw.Written() {
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					}
				}
			}()

			next.ServeHTTP(rw, r)

			// If the handler didn't write anything, ensure we have a status code
			if !rw.Written() && rw.status == 0 {
				rw.WriteHeader(http.StatusOK)
			}

			// Copy the response to the original writer if not already done
			rw.WriteTo(w)
		})
	}
}

// responseRecorder is a wrapper around http.ResponseWriter that records the response
// and allows checking if anything was written
type responseRecorder struct {
	http.ResponseWriter
	status int
	size   int
	header http.Header
}

func newResponseRecorder(w http.ResponseWriter) *responseRecorder {
	return &responseRecorder{
		ResponseWriter: w,
		header:         w.Header().Clone(),
	}
}

func (r *responseRecorder) WriteHeader(code int) {
	if !r.Written() {
		r.status = code
		r.ResponseWriter.WriteHeader(code)
	}
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.size += len(b)
	if !r.Written() {
		r.status = http.StatusOK // Default status code if none was set
	}
	return r.ResponseWriter.Write(b)
}

func (r *responseRecorder) Written() bool {
	return r.status != 0
}

func (r *responseRecorder) WriteTo(w http.ResponseWriter) {
	// If we already wrote to the original writer, don't write again
	if r.Written() {
		return
	}

	// Copy headers
	headers := w.Header()
	for k, v := range r.header {
		headers[k] = v
	}

	// Set status code if it was set
	if r.status != 0 {
		w.WriteHeader(r.status)
	}
}


