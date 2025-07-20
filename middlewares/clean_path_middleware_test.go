package middlewares_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/rtr/middlewares"
)

func TestCleanPath_NoRedirect(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expectPath  string
		setRawPath  bool
	}{
		{
			name:       "clean path stays the same",
			path:       "/users/1",
			expectPath: "/users/1",
		},
		{
			name:       "trailing slash is preserved",
			path:       "/users/1/",
			expectPath: "/users/1/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			if tt.setRawPath {
				req.URL.RawPath = req.URL.Path
			}
			rec := httptest.NewRecorder()

			handler := middlewares.CleanPathMiddleware().GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(r.URL.Path))
			}))

			handler.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if status := res.StatusCode; status != http.StatusOK {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}

			body, _ := io.ReadAll(res.Body)
			if got := string(body); got != tt.expectPath {
				t.Errorf("handler returned unexpected body: got %v want %v", got, tt.expectPath)
			}

			// Verify no redirect
			if location := res.Header.Get("Location"); location != "" {
				t.Errorf("unexpected redirect to: %s", location)
			}
		})
	}
}

func TestCleanPath_Redirect(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expectPath  string
		setRawPath  bool
	}{
		{
			name:       "double slashes are cleaned",
			path:       "/users//1",
			expectPath: "/users/1",
		},
		{
			name:       "multiple slashes are cleaned",
			path:       "//users////1",
			expectPath: "/users/1",
		},
		{
			name:       "root path with multiple slashes",
			path:       "///",
			expectPath: "/",
		},
		{
			name:       "preserve query parameters",
			path:       "/users//1?name=test",
			expectPath: "/users/1?name=test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			rec := httptest.NewRecorder()

			handler := middlewares.CleanPathMiddleware().GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Error("handler should not be called for redirects")
			}))

			handler.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if status := res.StatusCode; status != http.StatusMovedPermanently {
				t.Fatalf("handler returned wrong status code: got %v want %v", 
					status, http.StatusMovedPermanently)
			}

			location, err := res.Location()
			if err != nil {
				t.Fatalf("error getting location: %v", err)
			}
			if got := location.String(); got != tt.expectPath {
				t.Errorf("handler returned wrong location header: got %v want %v",
					got, tt.expectPath)
			}
		})
	}
}

func TestCleanPath_NilSafety(t *testing.T) {
	t.Run("nil handler", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rec := httptest.NewRecorder()

		// Test with nil handler
		handler := middlewares.CleanPathMiddleware().GetHandler()(nil)
		handler.ServeHTTP(rec, req)

		// Should not panic, but may not do anything useful
		if status := rec.Code; status != http.StatusOK {
			t.Errorf("unexpected status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("nil request", func(t *testing.T) {
		rec := httptest.NewRecorder()

		// Test with nil request
		handler := middlewares.CleanPathMiddleware().GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("handler should not be called with nil request")
		}))

		handler.ServeHTTP(rec, nil)
	})

	t.Run("nil URL", func(t *testing.T) {
		req := &http.Request{} // No URL set
		rec := httptest.NewRecorder()

		handler := middlewares.CleanPathMiddleware().GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("handler should not be called with nil URL")
		}))

		handler.ServeHTTP(rec, req)
	})
}

func TestCleanPath_URLEncoded(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expectPath  string
		setRawPath  bool
	}{
		{
			name:       "handle URL-encoded paths",
			path:       "/users%2F1//profile",
			expectPath: "/users/1/profile",
			setRawPath: true,
		},
		{
			name:       "handle URL-encoded paths with query",
			path:       "/users%2F1//profile?name=test%20user",
			expectPath: "/users/1/profile?name=test%20user",
			setRawPath: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			if tt.setRawPath {
				req.URL.RawPath = req.URL.Path
			}
			rec := httptest.NewRecorder()

			handler := middlewares.CleanPathMiddleware().GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Error("handler should not be called for redirects")
			}))

			handler.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if status := res.StatusCode; status != http.StatusMovedPermanently {
				t.Fatalf("handler returned wrong status code: got %v want %v", 
					status, http.StatusMovedPermanently)
			}

			location, err := res.Location()
			if err != nil {
				t.Fatalf("error getting location: %v", err)
			}
			if got := location.String(); got != tt.expectPath {
				t.Errorf("handler returned wrong location header: got %v want %v",
					got, tt.expectPath)
			}
		})
	}
}
