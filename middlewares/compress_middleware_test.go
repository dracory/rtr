package middlewares_test

import (
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
			w.Write([]byte("test response content"))
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
			w.Write([]byte("test response content"))
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
			w.Write([]byte(`{"message":"test"}`))
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
				w.Write([]byte("test"))
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
}
