package middlewares_test

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/rtr/middlewares"
)

func TestCompressMiddleware(t *testing.T) {
	t.Run("compresses response with gzip when accepted", func(t *testing.T) {
		// Create a test handler that returns some content
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			_, _ = w.Write([]byte("test response content"))
		})

		// Create middleware with default compression level
		middleware := middlewares.CompressMiddleware(gzip.DefaultCompression)

		// Create a request that accepts gzip encoding
		req := httptest.NewRequest("GET", "http://example.com/test", nil)
		req.Header.Set("Accept-Encoding", "gzip")

		w := httptest.NewRecorder()

		// Apply middleware and serve the request
		middleware.GetHandler()(handler).ServeHTTP(w, req)

		// Check response
		resp := w.Result()
		defer resp.Body.Close()

		// Verify Content-Encoding header
		if ce := resp.Header.Get("Content-Encoding"); ce != "gzip" {
			t.Errorf("Expected Content-Encoding: gzip, got %s", ce)
		}

		// Decompress and verify content
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			t.Fatalf("Failed to create gzip reader: %v", err)
		}
		defer reader.Close()

		content, err := io.ReadAll(reader)
		if err != nil {
			t.Fatalf("Failed to read decompressed content: %v", err)
		}

		if string(content) != "test response content" {
			t.Errorf("Unexpected response content: %s", string(content))
		}
	})

	t.Run("does not compress when client does not accept gzip", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("test response content"))
		})

		// Create middleware with default compression level
		middleware := middlewares.CompressMiddleware(gzip.DefaultCompression)

		// Request without Accept-Encoding header
		req := httptest.NewRequest("GET", "http://example.com/test", nil)
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		// Should not have Content-Encoding header
		if ce := resp.Header.Get("Content-Encoding"); ce != "" {
			t.Errorf("Expected no Content-Encoding, got %s", ce)
		}

		// Verify content is not compressed
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		if string(content) != "test response content" {
			t.Errorf("Unexpected response content: %s", string(content))
		}
	})

	t.Run("respects content types parameter", func(t *testing.T) {
		// Test handler that sets content type
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"message":"test"}`))
		})

		// Create middleware that only compresses JSON
		middleware := middlewares.CompressMiddleware(gzip.DefaultCompression, "application/json")

		t.Run("compresses when content type matches", func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://example.com/api", nil)
			req.Header.Set("Accept-Encoding", "gzip")
			w := httptest.NewRecorder()

			middleware.GetHandler()(handler).ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if ce := resp.Header.Get("Content-Encoding"); ce != "gzip" {
				t.Errorf("Expected Content-Encoding: gzip, got %s", ce)
			}
		})

		t.Run("does not compress when content type does not match", func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				_, _ = w.Write([]byte("test"))
			})

			req := httptest.NewRequest("GET", "http://example.com/api", nil)
			req.Header.Set("Accept-Encoding", "gzip")
			w := httptest.NewRecorder()

			middleware.GetHandler()(handler).ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if ce := resp.Header.Get("Content-Encoding"); ce != "" {
				t.Errorf("Expected no Content-Encoding, got %s", ce)
			}
		})
	})

	t.Run("compresses with deflate when only deflate accepted", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			_, _ = w.Write([]byte("test response content"))
		})

		middleware := middlewares.CompressMiddleware(flate.DefaultCompression)

		req := httptest.NewRequest("GET", "http://example.com/test", nil)
		req.Header.Set("Accept-Encoding", "deflate")
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if ce := resp.Header.Get("Content-Encoding"); ce != "deflate" {
			t.Errorf("Expected Content-Encoding: deflate, got %s", ce)
		}
	})

	t.Run("prefers gzip over deflate when both accepted", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			_, _ = w.Write([]byte("test response content"))
		})

		middleware := middlewares.CompressMiddleware(gzip.DefaultCompression)

		req := httptest.NewRequest("GET", "http://example.com/test", nil)
		req.Header.Set("Accept-Encoding", "deflate, gzip")
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if ce := resp.Header.Get("Content-Encoding"); ce != "gzip" {
			t.Errorf("Expected Content-Encoding: gzip (preferred), got %s", ce)
		}
	})

	t.Run("adds Vary header when compressing", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			_, _ = w.Write([]byte("test"))
		})

		middleware := middlewares.CompressMiddleware(gzip.DefaultCompression)

		req := httptest.NewRequest("GET", "http://example.com/test", nil)
		req.Header.Set("Accept-Encoding", "gzip")
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if vary := resp.Header.Get("Vary"); vary != "Accept-Encoding" {
			t.Errorf("Expected Vary: Accept-Encoding, got %q", vary)
		}
	})

	t.Run("does not compress when Content-Encoding already set", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Content-Encoding", "gzip")
			_, _ = w.Write([]byte("already compressed"))
		})

		middleware := middlewares.CompressMiddleware(gzip.DefaultCompression)

		req := httptest.NewRequest("GET", "http://example.com/test", nil)
		req.Header.Set("Accept-Encoding", "gzip")
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		// Should not double-compress; the existing Content-Encoding should be preserved
		body, _ := io.ReadAll(resp.Body)
		if string(body) != "already compressed" {
			t.Errorf("Expected uncompressed body, got %q", string(body))
		}
	})

	t.Run("compresses default content types when no types specified", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write([]byte("<html><body>test</body></html>"))
		})

		middleware := middlewares.CompressMiddleware(gzip.DefaultCompression)

		req := httptest.NewRequest("GET", "http://example.com/test", nil)
		req.Header.Set("Accept-Encoding", "gzip")
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if ce := resp.Header.Get("Content-Encoding"); ce != "gzip" {
			t.Errorf("Expected Content-Encoding: gzip for text/html, got %s", ce)
		}
	})

	t.Run("supports wildcard content types", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/custom")
			_, _ = w.Write([]byte("test"))
		})

		middleware := middlewares.CompressMiddleware(gzip.DefaultCompression, "text/*")

		req := httptest.NewRequest("GET", "http://example.com/test", nil)
		req.Header.Set("Accept-Encoding", "gzip")
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if ce := resp.Header.Get("Content-Encoding"); ce != "gzip" {
			t.Errorf("Expected Content-Encoding: gzip for text/custom (wildcard), got %s", ce)
		}
	})

	t.Run("strips content type parameters before checking", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, _ = w.Write([]byte("test"))
		})

		middleware := middlewares.CompressMiddleware(gzip.DefaultCompression)

		req := httptest.NewRequest("GET", "http://example.com/test", nil)
		req.Header.Set("Accept-Encoding", "gzip")
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if ce := resp.Header.Get("Content-Encoding"); ce != "gzip" {
			t.Errorf("Expected Content-Encoding: gzip for text/plain; charset=utf-8, got %s", ce)
		}
	})

	t.Run("compresses with case-insensitive content type matching", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "TEXT/HTML")
			_, _ = w.Write([]byte("<html>test</html>"))
		})

		middleware := middlewares.CompressMiddleware(gzip.DefaultCompression)

		req := httptest.NewRequest("GET", "http://example.com/test", nil)
		req.Header.Set("Accept-Encoding", "gzip")
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if ce := resp.Header.Get("Content-Encoding"); ce != "gzip" {
			t.Errorf("Expected Content-Encoding: gzip for TEXT/HTML, got %s", ce)
		}
	})

	t.Run("trims whitespace in content type before matching", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html ; charset=utf-8")
			_, _ = w.Write([]byte("<html>test</html>"))
		})

		middleware := middlewares.CompressMiddleware(gzip.DefaultCompression)

		req := httptest.NewRequest("GET", "http://example.com/test", nil)
		req.Header.Set("Accept-Encoding", "gzip")
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if ce := resp.Header.Get("Content-Encoding"); ce != "gzip" {
			t.Errorf("Expected Content-Encoding: gzip for 'text/html ; charset=utf-8', got %s", ce)
		}
	})

	t.Run("respects q=0 as explicit refusal", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			_, _ = w.Write([]byte("test response content"))
		})

		middleware := middlewares.CompressMiddleware(gzip.DefaultCompression)

		req := httptest.NewRequest("GET", "http://example.com/test", nil)
		req.Header.Set("Accept-Encoding", "gzip;q=0")
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		// q=0 means "explicitly not acceptable" per RFC 7231
		if ce := resp.Header.Get("Content-Encoding"); ce != "" {
			t.Errorf("Expected no Content-Encoding for gzip;q=0, got %s", ce)
		}
	})

	t.Run("does not match encoding substrings", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			_, _ = w.Write([]byte("test response content"))
		})

		middleware := middlewares.CompressMiddleware(gzip.DefaultCompression)

		req := httptest.NewRequest("GET", "http://example.com/test", nil)
		// "xgzip" should not match "gzip"
		req.Header.Set("Accept-Encoding", "xgzip")
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if ce := resp.Header.Get("Content-Encoding"); ce != "" {
			t.Errorf("Expected no Content-Encoding for xgzip, got %s", ce)
		}
	})

	t.Run("panics on invalid compression level", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for invalid compression level")
			}
		}()

		// Level 10 is invalid for gzip (valid range: -2 to 9)
		middlewares.CompressMiddleware(10)
	})

	t.Run("accepts level 0 (NoCompression)", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Level 0 should be valid, but panicked: %v", r)
			}
		}()

		// Level 0 (gzip.NoCompression) is valid — it writes a valid gzip
		// stream with no actual compression. Useful for testing and proxies.
		middleware := middlewares.CompressMiddleware(0)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			_, _ = w.Write([]byte("test"))
		})

		req := httptest.NewRequest("GET", "http://example.com/test", nil)
		req.Header.Set("Accept-Encoding", "gzip")
		w := httptest.NewRecorder()

		middleware.GetHandler()(handler).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if ce := resp.Header.Get("Content-Encoding"); ce != "gzip" {
			t.Errorf("Expected Content-Encoding: gzip for level 0, got %s", ce)
		}
	})
}
