package middlewares

import (
	"log"
	"net/http"

	"github.com/dracory/rtr"
)

// recoveryMiddleware implements the MiddlewareInterface for panic recovery.
type recoveryMiddleware struct {
	name    string
	handler rtr.StdMiddleware
}

// RecoveryMiddleware creates a new middleware that recovers from panics.
// It logs the panic details and returns a 500 Internal Server Error response.
// This should typically be added as one of the first middlewares in the chain.
func RecoveryMiddleware() rtr.MiddlewareInterface {
	return &recoveryMiddleware{
		name:    "Recovery Middleware",
		handler: defaultRecoveryHandler(),
	}
}

// GetName returns the name identifier associated with this middleware.
func (rm *recoveryMiddleware) GetName() string {
	return rm.name
}

// SetName sets the name identifier for this middleware and returns the middleware for method chaining.
func (rm *recoveryMiddleware) SetName(name string) rtr.MiddlewareInterface {
	rm.name = name
	return rm
}

// defaultRecoveryHandler returns the default recovery middleware handler.
func defaultRecoveryHandler() rtr.StdMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Log the panic
					log.Printf("Recovered from panic: %v", err)

					// Return 500 Internal Server Error
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// GetHandler returns the underlying middleware function.
func (rm *recoveryMiddleware) GetHandler() rtr.StdMiddleware {
	return rm.handler
}

// SetHandler sets the middleware function and returns the middleware for method chaining.
// This allows for custom recovery behavior, which is useful for testing or specialized handling.
func (rm *recoveryMiddleware) SetHandler(handler rtr.StdMiddleware) rtr.MiddlewareInterface {
	rm.handler = handler
	return rm
}

// Execute applies the middleware to the given handler.
func (rm *recoveryMiddleware) Execute(next http.Handler) http.Handler {
	return rm.GetHandler()(next)
}
