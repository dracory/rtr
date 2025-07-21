package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/rtr/middlewares"
)

func TestRealIPMiddleware(t *testing.T) {
	// Create a test handler that checks for the real IP in the context
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The RealIP middleware should have set the RemoteAddr to the X-Forwarded-For header
		if r.RemoteAddr != "192.168.1.1" {
			http.Error(w, "Unexpected RemoteAddr", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.1")
	req.RemoteAddr = "10.0.0.1:12345" // This should be overridden by the middleware

	w := httptest.NewRecorder()

	// Create middleware and apply to handler
	middleware := middlewares.RealIPMiddleware()
	middleware.GetHandler()(handler).ServeHTTP(w, req)

	// Check that the response is OK
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}
