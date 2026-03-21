package middlewares

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJailBotsMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
		config         JailBotsConfig
	}{
		{
			name:           "Safe path should pass",
			path:           "/safe/path",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
			config:         JailBotsConfig{},
		},
		{
			name:           "Blacklisted path should be jailed",
			path:           "/wp-admin",
			expectedStatus: http.StatusForbidden,
			expectedBody:   "malicious access not allowed (jb)",
			config:         JailBotsConfig{},
		},
		{
			name:           "Excluded path should pass even if blacklisted",
			path:           "/wp-admin",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
			config:         JailBotsConfig{ExcludePaths: []string{"/wp*"}},
		},
		{
			name:           "Excluded item should pass",
			path:           "/wp-admin",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
			config:         JailBotsConfig{Exclude: []string{"wp"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test handler
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			})

			// Create middleware
			middleware := JailBotsMiddleware(tt.config)
			handler := middleware.GetHandler()(testHandler)

			// Create request
			req := httptest.NewRequest("GET", tt.path, nil)
			req.RemoteAddr = "192.168.1.1:12345"

			// Create response recorder
			rr := httptest.NewRecorder()

			// Serve the request
			handler.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Check body
			body := strings.TrimSpace(rr.Body.String())
			if body != tt.expectedBody {
				t.Errorf("Expected body %q, got %q", tt.expectedBody, body)
			}
		})
	}
}

func TestJailBotsMiddleware_IPJailing(t *testing.T) {
	config := JailBotsConfig{}
	middleware := JailBotsMiddleware(config)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := middleware.GetHandler()(testHandler)

	// First request to a blacklisted path should jail the IP
	req := httptest.NewRequest("GET", "/wp-admin", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("Expected first request to be jailed, got status %d", rr.Code)
	}

	// Second request from the same IP to a safe path should still be jailed
	req = httptest.NewRequest("GET", "/safe/path", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("Expected second request from jailed IP to be blocked, got status %d", rr.Code)
	}

	// Request from different IP should pass
	req = httptest.NewRequest("GET", "/safe/path", nil)
	req.RemoteAddr = "192.168.1.2:12345"
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected request from different IP to pass, got status %d", rr.Code)
	}
}

func TestJailBotsMiddleware_XForwardedFor(t *testing.T) {
	config := JailBotsConfig{}
	middleware := JailBotsMiddleware(config)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := middleware.GetHandler()(testHandler)

	// Request with X-Forwarded-For header
	req := httptest.NewRequest("GET", "/wp-admin", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.1")
	req.RemoteAddr = "192.168.1.1:12345"
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("Expected request with X-Forwarded-For to be jailed, got status %d", rr.Code)
	}

	// Second request from same forwarded IP to safe path should still be jailed
	req = httptest.NewRequest("GET", "/safe/path", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.1")
	req.RemoteAddr = "192.168.1.2:12345"
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("Expected second request from same forwarded IP to be jailed, got status %d", rr.Code)
	}
}

func TestGetIP(t *testing.T) {
	tests := []struct {
		name       string
		request    *http.Request
		expectedIP string
	}{
		{
			name: "X-Forwarded-For single IP",
			request: &http.Request{
				Header:     http.Header{"X-Forwarded-For": []string{"203.0.113.1"}},
				RemoteAddr: "192.168.1.1:12345",
			},
			expectedIP: "203.0.113.1",
		},
		{
			name: "X-Forwarded-For multiple IPs",
			request: &http.Request{
				Header:     http.Header{"X-Forwarded-For": []string{"203.0.113.1, 198.51.100.1"}},
				RemoteAddr: "192.168.1.1:12345",
			},
			expectedIP: "203.0.113.1",
		},
		{
			name: "X-Real-IP",
			request: func() *http.Request {
				req := &http.Request{
					Header:     http.Header{},
					RemoteAddr: "192.168.1.1:12345",
				}
				req.Header.Set("X-Real-IP", "203.0.113.1")
				return req
			}(),
			expectedIP: "203.0.113.1",
		},
		{
			name: "RemoteAddr only",
			request: &http.Request{
				RemoteAddr: "192.168.1.1:12345",
			},
			expectedIP: "192.168.1.1",
		},
		{
			name: "RemoteAddr without port",
			request: &http.Request{
				RemoteAddr: "192.168.1.1",
			},
			expectedIP: "192.168.1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := getIP(tt.request)
			if ip != tt.expectedIP {
				t.Errorf("Expected IP %q, got %q", tt.expectedIP, ip)
			}
		})
	}
}

func TestJailBotsMiddleware_Integration(t *testing.T) {
	// Test middleware chaining and configuration
	config := JailBotsConfig{
		ExcludePaths: []string{"/api*"},
		Exclude:      []string{"admin"},
	}

	middleware := JailBotsMiddleware(config)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := middleware.GetHandler()(testHandler)

	tests := []struct {
		path           string
		expectedStatus int
	}{
		{"/safe", http.StatusOK},
		{"/wp-admin", http.StatusForbidden},
		{"/api/status", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			req.RemoteAddr = "192.168.1.1:12345"
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d for %s, got %d", tt.expectedStatus, tt.path, rr.Code)
			}
		})
	}
}
