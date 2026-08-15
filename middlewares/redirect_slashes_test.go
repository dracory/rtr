package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/rtr/middlewares"
)

func TestRedirectSlashesMiddleware(t *testing.T) {
	middleware := middlewares.RedirectSlashesMiddleware()

	t.Run("redirects trailing slash", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Next handler should not be called")
		})

		req := httptest.NewRequest("GET", "/api/v1/", nil)
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMovedPermanently {
			t.Errorf("Expected status code %d, got %d", http.StatusMovedPermanently, resp.StatusCode)
		}

		if location := resp.Header.Get("Location"); location != "/api/v1" {
			t.Errorf("Expected location %q, got %q", "/api/v1", location)
		}
	})

	t.Run("does not redirect internal double slashes without trailing slash", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "/api//v1", nil)
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if !called {
			t.Error("Next handler was not called")
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("redirects internal double slashes with trailing slash", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Next handler should not be called")
		})

		req := httptest.NewRequest("GET", "/api/v1", nil)
		req.URL.Path = "/api//v1/"
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMovedPermanently {
			t.Fatalf("Expected status code %d, got %d", http.StatusMovedPermanently, resp.StatusCode)
		}

		// http.Redirect cleans the path, so internal double slashes are
		// collapsed in the Location header. This matches chi's behaviour
		// since chi also uses http.Redirect.
		location := resp.Header.Get("Location")
		if location != "/api/v1" {
			t.Errorf("Expected location %q, got %q", "/api/v1", location)
		}
	})

	t.Run("does not redirect valid path", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "/api/v1", nil)
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if !called {
			t.Error("Next handler was not called")
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("does not redirect root path", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if !called {
			t.Error("Next handler was not called for root path")
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d for root, got %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("preserves query string on redirect", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Next handler should not be called")
		})

		req := httptest.NewRequest("GET", "/api/v1/?foo=bar&baz=qux", nil)
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMovedPermanently {
			t.Fatalf("Expected status code %d, got %d", http.StatusMovedPermanently, resp.StatusCode)
		}

		location := resp.Header.Get("Location")
		if location != "/api/v1?foo=bar&baz=qux" {
			t.Errorf("Expected location %q, got %q", "/api/v1?foo=bar&baz=qux", location)
		}
	})

	t.Run("normalizes backslashes in path", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Next handler should not be called")
		})

		// Path with trailing slash and backslash that should be normalized
		req := httptest.NewRequest("GET", "/api\\v1/", nil)
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMovedPermanently {
			t.Fatalf("Expected status code %d, got %d", http.StatusMovedPermanently, resp.StatusCode)
		}

		// Backslash should be normalized to forward slash
		location := resp.Header.Get("Location")
		if location != "/api/v1" {
			t.Errorf("Expected location %q, got %q", "/api/v1", location)
		}
	})

	t.Run("collapses multiple trailing slashes", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Next handler should not be called")
		})

		// httptest.NewRequest normalizes the URL, so we set the path manually
		req := httptest.NewRequest("GET", "/api/v1", nil)
		req.URL.Path = "/api/v1///"
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMovedPermanently {
			t.Fatalf("Expected status code %d, got %d", http.StatusMovedPermanently, resp.StatusCode)
		}

		location := resp.Header.Get("Location")
		if location != "/api/v1" {
			t.Errorf("Expected location %q, got %q", "/api/v1", location)
		}
	})

	t.Run("redirects slash-only paths to root", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Next handler should not be called")
		})

		// Paths like "//" or "///" should redirect to "/"
		req := httptest.NewRequest("GET", "/x", nil)
		req.URL.Path = "//"
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMovedPermanently {
			t.Fatalf("Expected status code %d, got %d", http.StatusMovedPermanently, resp.StatusCode)
		}

		location := resp.Header.Get("Location")
		if location != "/" {
			t.Errorf("Expected location %q, got %q", "/", location)
		}
	})
}
