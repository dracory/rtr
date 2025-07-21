package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dracory/rtr"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// RequestIDKey is the key used to store the request ID in the context.
// This matches Chi's RequestIDKey for consistency.
const RequestIDKey = chimiddleware.RequestIDKey

// RequestIDMiddleware returns a middleware that adds a unique request ID to the context and response headers.
// The request ID can be retrieved using GetRequestID(ctx).
// This is a thin wrapper around Chi's RequestID middleware for consistency with the project.
func RequestIDMiddleware() rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Request ID").
		SetHandler(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Generate a new request ID and convert to string
				reqID := fmt.Sprintf("%d", chimiddleware.NextRequestID())
				
				// Create a new context with the request ID
				ctx := context.WithValue(r.Context(), chimiddleware.RequestIDKey, reqID)
				
				// Set the request ID in the response header
				w.Header().Set("X-Request-Id", reqID)
				
				// Call the next handler with the new context
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		})
}

// GetRequestID retrieves the request ID from the context.
// Returns an empty string if no request ID is found.
func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(chimiddleware.RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}
