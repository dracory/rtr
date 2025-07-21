package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/rtr/middlewares"
)

func TestLoggerMiddleware(t *testing.T) {
	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "http://example.com/test?q=1", nil)
	w := httptest.NewRecorder()

	// Create middleware and apply to handler
	middleware := middlewares.LoggerMiddleware()
	middleware.GetHandler()(handler).ServeHTTP(w, req)

	// Check that the response is OK
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}
