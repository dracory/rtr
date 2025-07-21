package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/rtr/middlewares"
)

func TestProfilerMiddleware(t *testing.T) {
	t.Run("handles requests to profiler paths", func(t *testing.T) {
		handlerCalled := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			t.Error("Handler should not be called for profiler paths")
		})

		// Create the profiler middleware
		middleware := middlewares.ProfilerMiddleware()

		// Request the profiler path
		req := httptest.NewRequest("GET", "http://example.com/debug/pprof/", nil)
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		// Should not call the next handler
		if handlerCalled {
			t.Error("Expected handler not to be called for profiler paths")
		}
	})

	t.Run("does not call next handler for any request", func(t *testing.T) {
		handlerCalled := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
		})

		// Create the profiler middleware
		middleware := middlewares.ProfilerMiddleware()

		// Request a non-profiler path
		req := httptest.NewRequest("GET", "http://example.com/api/status", nil)
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		// Should call the next handler
		if !handlerCalled {
			t.Error("Expected handler to be called for non-profiler requests")
		}

		// Should return 200 OK from the handler
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})
}
