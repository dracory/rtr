package middlewares

import (
	"log"
	"net/http"

	"github.com/dracory/rtr"
)

// Ensure recoveryMiddleware implements MiddlewareInterface at compile time
var _ rtr.MiddlewareInterface = (*recoveryMiddleware)(nil)

// recoveryMiddleware implements the MiddlewareInterface for panic recovery.
type recoveryMiddleware struct {
	name    string
	handler rtr.StdMiddleware
}

// RecoveryMiddleware creates a new middleware that recovers from panics.
// It logs the panic details and returns a 500 Internal Server Error response.
// This should typically be added as one of the first middlewares in the chain.
func RecoveryMiddleware() rtr.MiddlewareInterface {
	// Create our recovery middleware with the configured values
	rm := &recoveryMiddleware{
		name:    "Recovery Middleware",
		handler: defaultRecoveryHandler(),
	}

	return rm
}

// GetName returns the name of the middleware
func (m *recoveryMiddleware) GetName() string {
	return m.name
}

// SetName sets the name of the middleware and returns the middleware for chaining
func (m *recoveryMiddleware) SetName(name string) rtr.MiddlewareInterface {
	m.name = name
	return m
}

// GetHandler returns the handler function
func (m *recoveryMiddleware) GetHandler() rtr.StdMiddleware {
	return m.handler
}

// SetHandler sets the handler function and returns the middleware for chaining
func (m *recoveryMiddleware) SetHandler(handler rtr.StdMiddleware) rtr.MiddlewareInterface {
	m.handler = handler
	return m
}

// Execute implements the middleware interface
func (m *recoveryMiddleware) Execute(next http.Handler) http.Handler {
	if m.handler == nil {
		m.handler = defaultRecoveryHandler()
	}
	return m.handler(next)
}

// defaultRecoveryHandler returns the default recovery handler
func defaultRecoveryHandler() rtr.StdMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("Recovered from panic: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
