package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/rtr/middlewares"
	"github.com/go-chi/cors"
)

func TestCORSMiddleware(t *testing.T) {
	t.Run("adds CORS headers to response", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		// Create middleware with default CORS options
		middleware := middlewares.DefaultCORSMiddleware()

		req := httptest.NewRequest("GET", "http://example.com/api", nil)
		req.Header.Set("Origin", "http://example.com")
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		// Check CORS headers
		if origin := resp.Header.Get("Access-Control-Allow-Origin"); origin != "*" {
			t.Errorf("Expected Access-Control-Allow-Origin: %q, got %q", "*", origin)
		}
	})

	t.Run("handles preflight requests", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Handler should not be called for OPTIONS requests")
		})

		// Create middleware with default CORS options
		middleware := middlewares.DefaultCORSMiddleware()

		req := httptest.NewRequest("OPTIONS", "http://example.com/api", nil)
		req.Header.Set("Origin", "http://example.com")
		req.Header.Set("Access-Control-Request-Method", "POST")
		req.Header.Set("Access-Control-Request-Headers", "Content-Type") // Request headers
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		// Should return 200 for preflight
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		// Should have CORS headers
		if origin := resp.Header.Get("Access-Control-Allow-Origin"); origin != "*" {
			t.Errorf("Expected Access-Control-Allow-Origin: *, got %q", origin)
		}

		headers := map[string]string{
			"Access-Control-Allow-Methods": "POST",
			"Access-Control-Allow-Headers": "Content-Type",
		}

		for header, expected := range headers {
			if got := resp.Header.Get(header); got != expected {
				t.Errorf("Expected header %s: %q, got %q", header, expected, got)
			}
		}
	})

	t.Run("respects custom CORS options", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		// Create middleware with custom CORS options
		middleware := middlewares.CORSMiddleware(
			cors.Options{
				AllowedOrigins:   []string{"https://example.com"},
				AllowedMethods:   []string{"GET", "POST"},
				AllowedHeaders:   []string{"Content-Type"},
				AllowCredentials: true,
			},
		)

		t.Run("allows request from allowed origin", func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://api.example.com/resource", nil)
			req.Header.Set("Origin", "https://example.com")
			w := httptest.NewRecorder()

			middleware.GetHandler()(handler).ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if origin := resp.Header.Get("Access-Control-Allow-Origin"); origin != "https://example.com" {
				t.Errorf("Expected Access-Control-Allow-Origin: https://example.com, got %q", origin)
			}
		})

		t.Run("rejects request from disallowed origin", func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://api.example.com/resource", nil)
			req.Header.Set("Origin", "https://malicious.com")
			w := httptest.NewRecorder()

			middleware.GetHandler()(handler).ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			// Should not include CORS headers for disallowed origins
			if origin := resp.Header.Get("Access-Control-Allow-Origin"); origin != "" {
				t.Errorf("Expected no Access-Control-Allow-Origin, got %q", origin)
			}
		})
	})
}
