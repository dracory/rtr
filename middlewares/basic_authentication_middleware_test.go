package middlewares_test

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dracory/rtr/middlewares"
)

func TestBasicAuthenticationMiddleware(t *testing.T) {
	// Test credentials
	username := "testuser"
	password := "testpass123"

	tests := []struct {
		name           string
		setAuth        bool
		username       string
		password       string
		expectedStatus int
		authHeader     string
	}{
		{
			name:           "valid credentials",
			setAuth:        true,
			username:       username,
			password:       password,
			expectedStatus: http.StatusOK,
			authHeader:     "Basic " + basicAuth(username, password),
		},
		{
			name:           "invalid username",
			setAuth:        true,
			username:       "wronguser",
			password:       password,
			expectedStatus: http.StatusUnauthorized,
			authHeader:     "Basic " + basicAuth("wronguser", password),
		},
		{
			name:           "invalid password",
			setAuth:        true,
			username:       username,
			password:       "wrongpass",
			expectedStatus: http.StatusUnauthorized,
			authHeader:     "Basic " + basicAuth(username, "wrongpass"),
		},
		{
			name:           "missing credentials",
			setAuth:        false,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "malformed auth header",
			setAuth:        true,
			expectedStatus: http.StatusUnauthorized,
			authHeader:     "Basic invalidbase64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.setAuth {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rec := httptest.NewRecorder()

			// Create middleware with test credentials
			middleware := middlewares.BasicAuthenticationMiddleware(username, password)
			handler := middleware.GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("authenticated"))
			}))

			handler.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			// Check status code
			if res.StatusCode != tt.expectedStatus {
				t.Errorf("status code mismatch: got %v, want %v", res.StatusCode, tt.expectedStatus)
			}

			// If unauthorized, check for WWW-Authenticate header
			if res.StatusCode == http.StatusUnauthorized {
				authHeader := res.Header.Get("WWW-Authenticate")
				if !strings.Contains(authHeader, `Basic realm="restricted"`) {
					t.Error("missing WWW-Authenticate header with Basic realm")
				}
			}
		})
	}
}

// Helper function to create basic auth header value
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func TestBasicAuthenticationMiddleware_NilNextHandler(t *testing.T) {
	username := "testuser"
	password := "testpass123"

	req := httptest.NewRequest("GET", "/", nil)
	req.SetBasicAuth(username, password)

	rec := httptest.NewRecorder()

	middleware := middlewares.BasicAuthenticationMiddleware(username, password)
	handler := middleware.GetHandler()(nil) // Pass nil as next handler

	handler.ServeHTTP(rec, req)

	// Should return 404 since we use http.NotFoundHandler as fallback for nil next handler
	if status := rec.Result().StatusCode; status != http.StatusNotFound {
		t.Errorf("status code should be 404 with nil next handler, got %d", status)
	}
}

func TestBasicAuthenticationMiddleware_NilRequest(t *testing.T) {
	username := "testuser"
	password := "testpass123"

	rec := httptest.NewRecorder()

	middleware := middlewares.BasicAuthenticationMiddleware(username, password)
	handler := middleware.GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// This should not panic
	handler.ServeHTTP(rec, nil)

	if status := rec.Result().StatusCode; status != http.StatusBadRequest {
		t.Errorf("status code should be 400 with nil request, got %d", status)
	}
}

func TestBasicAuthenticationMiddleware_NilResponseWriter(t *testing.T) {
	username := "testuser"
	password := "testpass123"

	req := httptest.NewRequest("GET", "/", nil)
	req.SetBasicAuth(username, password)

	middleware := middlewares.BasicAuthenticationMiddleware(username, password)
	handler := middleware.GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("next handler should not be called with nil response writer")
	}))

	// This should not panic
	handler.ServeHTTP(nil, req)
}

func TestBasicAuthenticationMiddleware_EmptyCredentials(t *testing.T) {
	username := "testuser"
	password := "testpass123"

	req := httptest.NewRequest("GET", "/", nil)
	req.SetBasicAuth(username, password)

	rec := httptest.NewRecorder()

	// Create middleware with empty credentials
	middleware := middlewares.BasicAuthenticationMiddleware("", "")
	handler := middleware.GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("next handler should not be called with invalid configuration")
	}))

	handler.ServeHTTP(rec, req)

	if status := rec.Result().StatusCode; status != http.StatusBadRequest {
		t.Errorf("status code should be 400 with empty credentials, got %d", status)
	}
}

func TestBasicAuthenticationMiddleware_WithNextHandler(t *testing.T) {
	username := "testuser"
	password := "testpass123"

	t.Run("calls next handler when authenticated", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.SetBasicAuth(username, password)

		rec := httptest.NewRecorder()

		var nextCalled bool
		middleware := middlewares.BasicAuthenticationMiddleware(username, password)
		handler := middleware.GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
			w.WriteHeader(http.StatusOK)
		}))

		handler.ServeHTTP(rec, req)

		if !nextCalled {
			t.Error("next handler should be called")
		}
		if status := rec.Result().StatusCode; status != http.StatusOK {
			t.Errorf("status code should be 200, got %d", status)
		}
	})

	t.Run("does not call next handler when not authenticated", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.SetBasicAuth("wronguser", "wrongpass")

		rec := httptest.NewRecorder()

		var nextCalled bool
		middleware := middlewares.BasicAuthenticationMiddleware(username, password)
		handler := middleware.GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
			w.WriteHeader(http.StatusOK)
		}))

		handler.ServeHTTP(rec, req)

		if nextCalled {
			t.Error("next handler should not be called")
		}
		if status := rec.Result().StatusCode; status != http.StatusUnauthorized {
			t.Errorf("status code should be 401, got %d", status)
		}
	})
}
