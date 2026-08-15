package middlewares_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dracory/rtr/middlewares"
)

func TestTimeoutMiddleware(t *testing.T) {
	t.Run("request completes before timeout", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "http://example.com/test", nil)
		w := httptest.NewRecorder()

		// Create middleware with a timeout longer than the handler takes
		middleware := middlewares.TimeoutMiddleware(100 * time.Millisecond)
		middleware.GetHandler()(handler).ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("request exceeds timeout", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Simulate a long-running request that respects context cancellation
			select {
			case <-time.After(200 * time.Millisecond):
				w.WriteHeader(http.StatusOK)
			case <-r.Context().Done():
				// Context was canceled by the timeout middleware
				// The middleware should have already written the 504 response
				return
			}
		})

		req := httptest.NewRequest("GET", "http://example.com/test", nil)
		w := httptest.NewRecorder()

		// Create middleware with a timeout shorter than the handler takes
		middleware := middlewares.TimeoutMiddleware(50 * time.Millisecond)
		middleware.GetHandler()(handler).ServeHTTP(w, req)

		if w.Code != http.StatusGatewayTimeout {
			t.Fatalf("Expected status code %d, got %d", http.StatusGatewayTimeout, w.Code)
		}
	})

	t.Run("context is canceled when timeout is reached", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Wait for the context to be done
			<-r.Context().Done()
			if r.Context().Err() == context.DeadlineExceeded {
				w.WriteHeader(http.StatusRequestTimeout)
				return
			}
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "http://example.com/test", nil)
		w := httptest.NewRecorder()

		// Create middleware with a very short timeout
		middleware := middlewares.TimeoutMiddleware(1 * time.Millisecond)
		middleware.GetHandler()(handler).ServeHTTP(w, req)

		if w.Code != http.StatusRequestTimeout {
			t.Fatalf("Expected status code %d, got %d", http.StatusRequestTimeout, w.Code)
		}
	})

	t.Run("handler can write response body before timeout", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("response body"))
		})

		req := httptest.NewRequest("GET", "http://example.com/test", nil)
		w := httptest.NewRecorder()

		middleware := middlewares.TimeoutMiddleware(100 * time.Millisecond)
		middleware.GetHandler()(handler).ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		if w.Body.String() != "response body" {
			t.Errorf("Expected body %q, got %q", "response body", w.Body.String())
		}
	})

	t.Run("propagates context cancellation to handler", func(t *testing.T) {
		canceled := make(chan struct{})
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			<-r.Context().Done()
			close(canceled)
		})

		req := httptest.NewRequest("GET", "http://example.com/test", nil)
		w := httptest.NewRecorder()

		middleware := middlewares.TimeoutMiddleware(10 * time.Millisecond)
		middleware.GetHandler()(handler).ServeHTTP(w, req)

		select {
		case <-canceled:
			// Good, handler detected cancellation
		case <-time.After(1 * time.Second):
			t.Fatal("Handler did not detect context cancellation")
		}
	})
}
