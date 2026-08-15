package middlewares_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/rtr/middlewares"
)

func TestHeartbeatMiddleware(t *testing.T) {
	t.Run("returns 200 OK when path matches", func(t *testing.T) {
		endpoint := "/ping"
		middleware := middlewares.HeartbeatMiddleware(endpoint)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Next handler should not be called")
		})

		req := httptest.NewRequest("GET", endpoint, nil)
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		if string(body) != "." {
			t.Errorf("Expected body '.', got %q", string(body))
		}
	})

	t.Run("calls next handler when path does not match", func(t *testing.T) {
		endpoint := "/ping"
		middleware := middlewares.HeartbeatMiddleware(endpoint)

		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusAccepted)
		})

		req := httptest.NewRequest("GET", "/api", nil)
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if !called {
			t.Error("Next handler was not called")
		}

		if resp.StatusCode != http.StatusAccepted {
			t.Errorf("Expected status code %d, got %d", http.StatusAccepted, resp.StatusCode)
		}
	})

	t.Run("responds to HEAD method", func(t *testing.T) {
		endpoint := "/ping"
		middleware := middlewares.HeartbeatMiddleware(endpoint)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Next handler should not be called")
		})

		req := httptest.NewRequest("HEAD", endpoint, nil)
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d for HEAD, got %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("does not respond to POST method", func(t *testing.T) {
		endpoint := "/ping"
		middleware := middlewares.HeartbeatMiddleware(endpoint)

		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("POST", endpoint, nil)
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if !called {
			t.Error("Next handler should be called for POST")
		}
	})

	t.Run("matches path case-insensitively", func(t *testing.T) {
		endpoint := "/ping"
		middleware := middlewares.HeartbeatMiddleware(endpoint)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Next handler should not be called")
		})

		req := httptest.NewRequest("GET", "/PING", nil)
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d for case-insensitive match, got %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("does not match partial path prefix", func(t *testing.T) {
		endpoint := "/ping"
		middleware := middlewares.HeartbeatMiddleware(endpoint)

		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "/pingpong", nil)
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if !called {
			t.Error("Next handler should be called for /pingpong")
		}
	})

	t.Run("sets content type to text/plain", func(t *testing.T) {
		endpoint := "/ping"
		middleware := middlewares.HeartbeatMiddleware(endpoint)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Next handler should not be called")
		})

		req := httptest.NewRequest("GET", endpoint, nil)
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if ct := resp.Header.Get("Content-Type"); ct != "text/plain" {
			t.Errorf("Expected Content-Type %q, got %q", "text/plain", ct)
		}
	})
}
