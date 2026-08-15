package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/rtr/middlewares"
)

func TestRealIPMiddleware(t *testing.T) {
	t.Run("uses X-Forwarded-For header", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.RemoteAddr != "192.168.1.1" {
				http.Error(w, "Unexpected RemoteAddr", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		req.Header.Set("X-Forwarded-For", "192.168.1.1")
		req.RemoteAddr = "10.0.0.1:12345"

		w := httptest.NewRecorder()
		middleware := middlewares.RealIPMiddleware()
		middleware.GetHandler()(handler).ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("uses True-Client-IP header (highest priority)", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.RemoteAddr != "10.0.0.99" {
				http.Error(w, "Unexpected RemoteAddr", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		req.Header.Set("True-Client-IP", "10.0.0.99")
		req.Header.Set("X-Real-IP", "10.0.0.50")
		req.Header.Set("X-Forwarded-For", "10.0.0.1")
		req.RemoteAddr = "10.0.0.1:12345"

		w := httptest.NewRecorder()
		middleware := middlewares.RealIPMiddleware()
		middleware.GetHandler()(handler).ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("uses X-Real-IP header (second priority)", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.RemoteAddr != "10.0.0.50" {
				http.Error(w, "Unexpected RemoteAddr", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		req.Header.Set("X-Real-IP", "10.0.0.50")
		req.Header.Set("X-Forwarded-For", "10.0.0.1")
		req.RemoteAddr = "10.0.0.1:12345"

		w := httptest.NewRecorder()
		middleware := middlewares.RealIPMiddleware()
		middleware.GetHandler()(handler).ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("takes first IP from X-Forwarded-For with multiple IPs", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.RemoteAddr != "192.168.1.1" {
				http.Error(w, "Unexpected RemoteAddr", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		req.Header.Set("X-Forwarded-For", "192.168.1.1, 10.0.0.1, 10.0.0.2")
		req.RemoteAddr = "10.0.0.1:12345"

		w := httptest.NewRecorder()
		middleware := middlewares.RealIPMiddleware()
		middleware.GetHandler()(handler).ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("does not change RemoteAddr for invalid IP", func(t *testing.T) {
		originalAddr := "10.0.0.1:12345"
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.RemoteAddr != originalAddr {
				http.Error(w, "RemoteAddr was changed", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		req.Header.Set("X-Forwarded-For", "not-an-ip")
		req.RemoteAddr = originalAddr

		w := httptest.NewRecorder()
		middleware := middlewares.RealIPMiddleware()
		middleware.GetHandler()(handler).ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("does not change RemoteAddr when no headers present", func(t *testing.T) {
		originalAddr := "10.0.0.1:12345"
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.RemoteAddr != originalAddr {
				http.Error(w, "RemoteAddr was changed", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		req.RemoteAddr = originalAddr

		w := httptest.NewRecorder()
		middleware := middlewares.RealIPMiddleware()
		middleware.GetHandler()(handler).ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("supports IPv6 addresses", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.RemoteAddr != "2001:db8::1" {
				http.Error(w, "Unexpected RemoteAddr", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		req.Header.Set("X-Forwarded-For", "2001:db8::1")
		req.RemoteAddr = "[::1]:12345"

		w := httptest.NewRecorder()
		middleware := middlewares.RealIPMiddleware()
		middleware.GetHandler()(handler).ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})
}
