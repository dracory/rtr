package router

import (
	"log"
	"net/http"
)

// RecoveryMiddleware creates a new middleware that recovers from panics.
// It logs the panic details and returns a 500 Internal Server Error response.
// This should typically be added as one of the first middlewares in the chain.
func RecoveryMiddleware(next http.Handler) http.Handler {
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

// DefaultMiddlewares returns a slice of default middlewares that should be used with the router.
// Currently, it only includes the RecoveryMiddleware.
func DefaultMiddlewares() []Middleware {
	return []Middleware{
		RecoveryMiddleware,
	}
}
