package middlewares

import (
	"net/http"

	"github.com/dracory/rtr"
	"github.com/go-chi/cors"
)

// CORSMiddleware returns a middleware that handles CORS requests.
// It's a thin wrapper around go-chi/cors middleware.
// By default, it allows all origins, methods, and headers.
// Use the options to customize the CORS behavior.
func CORSMiddleware(opts cors.Options) rtr.MiddlewareInterface {
	// Create CORS handler with provided options
	corsHandler := cors.New(opts)

	// Create a middleware with the CORS handler
	return rtr.NewMiddleware(
		rtr.WithName("CORS"),
		rtr.WithHandler(func(next http.Handler) http.Handler {
			return corsHandler.Handler(next)
		}),
	)
}

// DefaultCORSMiddleware returns a CORS middleware with sensible defaults:
// - Allow all origins
// - Allow common HTTP methods
// - Allow common headers
// - Allow credentials
// - Max age: 300 (5 minutes)
func DefaultCORSMiddleware() rtr.MiddlewareInterface {
	return CORSMiddleware(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutes
	})
}
