package middlewares_test

import (
	"bufio"
	"bytes"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/rtr/middlewares"
)

// hijackableResponseWriter is a test ResponseWriter that implements http.Hijacker
// to verify that the logger middleware preserves Hijacker support.
type hijackableResponseWriter struct {
	*httptest.ResponseRecorder
}

func (h *hijackableResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, nil
}

func TestLoggerMiddleware(t *testing.T) {
	t.Run("response is correct for successful request", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "http://example.com/test?q=1", nil)
		w := httptest.NewRecorder()

		middleware := middlewares.LoggerMiddleware()
		middleware.GetHandler()(handler).ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("response body is preserved", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("response body content"))
		})

		req := httptest.NewRequest("GET", "http://example.com/test", nil)
		w := httptest.NewRecorder()

		middleware := middlewares.LoggerMiddleware()
		middleware.GetHandler()(handler).ServeHTTP(w, req)

		if w.Body.String() != "response body content" {
			t.Errorf("Expected body %q, got %q", "response body content", w.Body.String())
		}
	})

	t.Run("handles different status codes", func(t *testing.T) {
		codes := []int{
			http.StatusOK,
			http.StatusNotFound,
			http.StatusInternalServerError,
			http.StatusMovedPermanently,
			http.StatusUnauthorized,
		}

		for _, code := range codes {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(code)
			})

			req := httptest.NewRequest("GET", "http://example.com/test", nil)
			w := httptest.NewRecorder()

			middleware := middlewares.LoggerMiddleware()
			middleware.GetHandler()(handler).ServeHTTP(w, req)

			if w.Code != code {
				t.Errorf("Expected status code %d, got %d", code, w.Code)
			}
		}
	})

	t.Run("handles POST requests", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
		})

		req := httptest.NewRequest("POST", "http://example.com/api", bytes.NewBufferString(`{"key":"value"}`))
		w := httptest.NewRecorder()

		middleware := middlewares.LoggerMiddleware()
		middleware.GetHandler()(handler).ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
		}
	})

	t.Run("preserves response headers", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Custom-Header", "custom-value")
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "http://example.com/test", nil)
		w := httptest.NewRecorder()

		middleware := middlewares.LoggerMiddleware()
		middleware.GetHandler()(handler).ServeHTTP(w, req)

		if h := w.Header().Get("X-Custom-Header"); h != "custom-value" {
			t.Errorf("Expected X-Custom-Header %q, got %q", "custom-value", h)
		}
	})

	t.Run("preserves Hijacker interface from underlying writer", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// The handler should be able to type-assert to http.Hijacker
			// because the underlying ResponseWriter implements it.
			hj, ok := w.(http.Hijacker)
			if !ok {
				http.Error(w, "Hijacker not available", http.StatusInternalServerError)
				return
			}
			_, _, err := hj.Hijack()
			if err != nil {
				http.Error(w, "Hijack failed", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "http://example.com/ws", nil)
		w := &hijackableResponseWriter{ResponseRecorder: httptest.NewRecorder()}

		middleware := middlewares.LoggerMiddleware()
		middleware.GetHandler()(handler).ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d (Hijacker interface was not preserved)", http.StatusOK, w.Code)
		}
	})
}
