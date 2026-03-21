package middlewares

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPSRedirectMiddleware_DefaultConfig(t *testing.T) {
	middleware := NewHTTPSRedirectMiddleware(nil)

	tests := []struct {
		name             string
		host             string
		scheme           string
		expectedStatus   int
		expectedRedirect string
	}{
		{
			name:           "HTTPS request should pass through",
			host:           "example.com",
			scheme:         "https",
			expectedStatus: 200,
		},
		{
			name:           "Localhost should skip redirect",
			host:           "localhost",
			scheme:         "http",
			expectedStatus: 200,
		},
		{
			name:           "127.0.0.1 should skip redirect",
			host:           "127.0.0.1",
			scheme:         "http",
			expectedStatus: 200,
		},
		{
			name:             "HTTP request should redirect",
			host:             "example.com",
			scheme:           "http",
			expectedStatus:   301,
			expectedRedirect: "https://example.com/test?foo=bar",
		},
		{
			name:             "HTTP request with query should redirect with query",
			host:             "example.com",
			scheme:           "http",
			expectedStatus:   301,
			expectedRedirect: "https://example.com/test?foo=bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://example.com/test?foo=bar", nil)
			req.Host = tt.host
			if tt.scheme == "https" {
				req.TLS = &tls.ConnectionState{}
			}

			rr := httptest.NewRecorder()

			handler := middleware.GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			}))

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, status)
			}

			if tt.expectedRedirect != "" {
				location := rr.Header().Get("Location")
				if location != tt.expectedRedirect {
					t.Errorf("Expected redirect to %s, got %s", tt.expectedRedirect, location)
				}
			}
		})
	}
}

func TestHTTPSRedirectMiddleware_CustomConfig(t *testing.T) {
	config := &HTTPSRedirectConfig{
		SkipLocalhost: false,
		CustomSkipFunc: func(r *http.Request) bool {
			return r.Host == "skip.example.com"
		},
	}

	middleware := NewHTTPSRedirectMiddleware(config)

	tests := []struct {
		name           string
		host           string
		expectedStatus int
	}{
		{
			name:           "Custom skip function should work",
			host:           "skip.example.com",
			expectedStatus: 200,
		},
		{
			name:           "Localhost should redirect when SkipLocalhost is false",
			host:           "localhost",
			expectedStatus: 301,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://example.com/test", nil)
			req.Host = tt.host

			rr := httptest.NewRecorder()

			handler := middleware.GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			}))

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, status)
			}
		})
	}
}

func TestIsLocalhost(t *testing.T) {
	tests := []struct {
		host     string
		expected bool
	}{
		{"localhost", true},
		{"127.0.0.1", true},
		{"0.0.0.0", true},
		{"example.local", true},
		{"127.0.0.2", true},
		{"192.168.1.1", true},
		{"10.0.0.1", true},
		{"example.com", false},
		{"google.com", false},
		{"192.169.1.1", false},
		{"10.1.0.1", true},
		{"10.255.255.255", true}, // Edge of 10.0.0.0/8 range
		{"172.16.0.1", false},    // 172.16.0.0/12 range (not implemented)
	}

	for _, tt := range tests {
		t.Run(tt.host, func(t *testing.T) {
			result := isLocalhost(tt.host)
			if result != tt.expected {
				t.Errorf("Expected %v for %s, got %v", tt.expected, tt.host, result)
			}
		})
	}
}
