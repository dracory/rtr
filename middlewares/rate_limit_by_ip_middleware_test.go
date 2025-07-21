package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dracory/rtr/middlewares"
)

func TestRateLimitByIPMiddleware(t *testing.T) {
	maxRequests := 2
	seconds := 1
	middleware := middlewares.RateLimitByIPMiddleware(maxRequests, seconds)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("rate limits requests from the same IP", func(t *testing.T) {
		for i := 0; i < maxRequests; i++ {
			req := httptest.NewRequest("GET", "/", nil)
			req.RemoteAddr = "192.0.2.1:1234"
			w := httptest.NewRecorder()
			middleware.GetHandler()(handler).ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status code %d, got %d on request %d", http.StatusOK, w.Code, i+1)
			}
		}

		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "192.0.2.1:1234"
		w := httptest.NewRecorder()
		middleware.GetHandler()(handler).ServeHTTP(w, req)

		if w.Code != http.StatusTooManyRequests {
			t.Errorf("Expected status code %d, got %d", http.StatusTooManyRequests, w.Code)
		}
	})

	t.Run("does not rate limit requests from different IPs", func(t *testing.T) {
		req1 := httptest.NewRequest("GET", "/", nil)
		req1.RemoteAddr = "192.0.2.2:1234"
		w1 := httptest.NewRecorder()
		middleware.GetHandler()(handler).ServeHTTP(w1, req1)

		if w1.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w1.Code)
		}

		req2 := httptest.NewRequest("GET", "/", nil)
		req2.RemoteAddr = "192.0.2.3:1234"
		w2 := httptest.NewRecorder()
		middleware.GetHandler()(handler).ServeHTTP(w2, req2)

		if w2.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w2.Code)
		}
	})

	t.Run("rate limit resets after the specified time", func(t *testing.T) {
		for i := 0; i < maxRequests; i++ {
			req := httptest.NewRequest("GET", "/", nil)
			req.RemoteAddr = "192.0.2.4:1234"
			w := httptest.NewRecorder()
			middleware.GetHandler()(handler).ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status code %d, got %d on request %d", http.StatusOK, w.Code, i+1)
			}
		}

		time.Sleep(time.Duration(seconds) * time.Second)

		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "192.0.2.4:1234"
		w := httptest.NewRecorder()
		middleware.GetHandler()(handler).ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d after waiting, got %d", http.StatusOK, w.Code)
		}
	})
}
