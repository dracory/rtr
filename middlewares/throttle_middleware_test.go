package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dracory/rtr/middlewares"
)

func TestThrottleMiddleware(t *testing.T) {
	t.Run("allows requests within limit", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		// Allow 2 requests per second
		middleware := middlewares.ThrottleMiddleware(2, time.Second)

		req1 := httptest.NewRequest("GET", "http://example.com/test", nil)
		req1.RemoteAddr = "192.168.1.1:12345"
		w1 := httptest.NewRecorder()

		// First request should succeed
		middleware.GetHandler()(handler).ServeHTTP(w1, req1)
		if w1.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, w1.Code)
		}

		req2 := httptest.NewRequest("GET", "http://example.com/test", nil)
		req2.RemoteAddr = "192.168.1.1:12345"
		w2 := httptest.NewRecorder()

		// Second request should also succeed
		middleware.GetHandler()(handler).ServeHTTP(w2, req2)
		if w2.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, w2.Code)
		}
	})

	t.Run("rejects requests over limit", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		// Allow only 1 request per second
		middleware := middlewares.ThrottleMiddleware(1, time.Second)

		req1 := httptest.NewRequest("GET", "http://example.com/test", nil)
		req1.RemoteAddr = "192.168.1.1:12345"
		w1 := httptest.NewRecorder()

		// First request should succeed
		middleware.GetHandler()(handler).ServeHTTP(w1, req1)
		if w1.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, w1.Code)
		}

		req2 := httptest.NewRequest("GET", "http://example.com/test", nil)
		req2.RemoteAddr = "192.168.1.1:12345"
		w2 := httptest.NewRecorder()

		// Second request should be rate limited
		middleware.GetHandler()(handler).ServeHTTP(w2, req2)
		if w2.Code != http.StatusTooManyRequests {
			t.Fatalf("Expected status code %d, got %d", http.StatusTooManyRequests, w2.Code)
		}
	})

	t.Run("allows requests from different IPs", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		// Allow only 1 request per second per IP
		middleware := middlewares.ThrottleMiddleware(1, time.Second)

		req1 := httptest.NewRequest("GET", "http://example.com/test", nil)
		req1.RemoteAddr = "192.168.1.1:12345"
		w1 := httptest.NewRecorder()

		// First IP's request should succeed
		middleware.GetHandler()(handler).ServeHTTP(w1, req1)
		if w1.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, w1.Code)
		}

		req2 := httptest.NewRequest("GET", "http://example.com/test", nil)
		req2.RemoteAddr = "192.168.1.2:12345" // Different IP
		w2 := httptest.NewRecorder()

		// Second IP's request should also succeed
		middleware.GetHandler()(handler).ServeHTTP(w2, req2)
		if w2.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, w2.Code)
		}
	})
}
