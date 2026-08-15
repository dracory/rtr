package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/dracory/rtr"
)

// requestIDKey is the context key used to store the request ID.
// It is an unexported type to avoid key collisions with other packages.
type requestIDKey struct{}

// RequestIDKey is exported for consumers that need to read the request ID
// directly from the context. Prefer GetRequestID(ctx) where possible.
var RequestIDKey = requestIDKey(struct{}{})

// requestIDCounter is a process-wide monotonically increasing counter used to
// generate unique request IDs without external dependencies.
var requestIDCounter uint64

// RequestIDMiddleware returns a middleware that adds a unique request ID to the
// context and response headers. The request ID can be retrieved using
// GetRequestID(ctx).
func RequestIDMiddleware() rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Request ID").
		SetHandler(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Generate a new request ID from the atomic counter
				reqID := fmt.Sprintf("%d", atomic.AddUint64(&requestIDCounter, 1))

				// Create a new context with the request ID
				ctx := context.WithValue(r.Context(), RequestIDKey, reqID)

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
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}
