package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/rtr/middlewares"
)

func TestRequestIDMiddleware(t *testing.T) {
	// Create a simple test handler that checks for the request ID in the context
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Print all context keys for debugging
		ctx := r.Context()
		t.Logf("Context values: %+v", ctx)

		// Try to get the request ID using the package's getter
		reqID := middlewares.GetRequestID(ctx)
		if reqID == "" {
			http.Error(w, "Request ID not found in context", http.StatusInternalServerError)
			return
		}
		t.Logf("Found request ID in context: %s", reqID)
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	// Create middleware and apply to handler
	middleware := middlewares.RequestIDMiddleware()
	middleware.GetHandler()(handler).ServeHTTP(w, req)

	// Check that the response has a request ID header
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	header := w.Header().Get("X-Request-Id")
	if header == "" {
		t.Fatal("Expected X-Request-Id header to be set")
	}
}
