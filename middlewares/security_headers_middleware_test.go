package middlewares

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSecurityHeadersMiddleware_DefaultConfig(t *testing.T) {
	middleware := NewSecurityHeadersMiddleware(nil)

	req := httptest.NewRequest("GET", "http://example.com/test", nil)
	rr := httptest.NewRecorder()

	handler := middleware.GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	handler.ServeHTTP(rr, req)

	// Check default security headers
	expectedHeaders := map[string]string{
		"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
		"X-Frame-Options":           "DENY",
		"X-Content-Type-Options":    "nosniff",
		"X-XSS-Protection":          "1; mode=block",
		"Referrer-Policy":           "strict-origin-when-cross-origin",
	}

	for header, expectedValue := range expectedHeaders {
		actualValue := rr.Header().Get(header)
		if actualValue != expectedValue {
			t.Errorf("Expected %s header to be %s, got %s", header, expectedValue, actualValue)
		}
	}

	// Check Permissions-Policy header (order may vary due to map iteration)
	permissionsPolicy := rr.Header().Get("Permissions-Policy")
	expectedPermissions := []string{"geolocation=()", "microphone=()", "camera=()"}

	// Split and check each expected permission is present
	actualPermissions := strings.Split(permissionsPolicy, ", ")
	for _, expected := range expectedPermissions {
		found := false
		for _, actual := range actualPermissions {
			if strings.TrimSpace(actual) == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected Permissions-Policy to contain %s, got %s", expected, permissionsPolicy)
		}
	}

	// Check CSP header
	cspHeader := rr.Header().Get("Content-Security-Policy")
	if cspHeader == "" {
		t.Error("Expected Content-Security-Policy header to be set")
	}

	// Check that CSP contains expected directives
	expectedCSPParts := []string{
		"default-src 'self'",
		"script-src 'self' 'unsafe-inline'",
		"style-src 'self' 'unsafe-inline'",
		"font-src 'self'",
		"img-src 'self' data:",
		"upgrade-insecure-requests",
	}

	for _, part := range expectedCSPParts {
		if !contains(cspHeader, part) {
			t.Errorf("Expected CSP to contain %s, got %s", part, cspHeader)
		}
	}
}

func TestSecurityHeadersMiddleware_CustomConfig(t *testing.T) {
	config := &SecurityHeadersConfig{
		CSP: &CSPConfig{
			Enabled:    true,
			DefaultSrc: []string{"'self'", "https://trusted.com"},
			ScriptSrc:  []string{"'self'"},
			StyleSrc:   []string{"'self'", "https://styles.com"},
		},
		HSTS: &HSTSConfig{
			Enabled:           true,
			MaxAge:            86400,
			IncludeSubDomains: false,
			Preload:           true,
		},
		FrameOptions: &FrameOptionsConfig{
			Enabled: true,
			Option:  "SAMEORIGIN",
		},
		ContentTypeNosniff: false,
		XSSProtection: &XSSProtectionConfig{
			Enabled: false,
		},
		ReferrerPolicy: "no-referrer",
		PermissionsPolicy: map[string][]string{
			"geolocation": {"self"},
		},
		CustomHeaders: map[string]string{
			"X-Custom-Security": "custom-value",
		},
	}

	middleware := NewSecurityHeadersMiddleware(config)

	req := httptest.NewRequest("GET", "http://example.com/test", nil)
	rr := httptest.NewRecorder()

	handler := middleware.GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	handler.ServeHTTP(rr, req)

	// Check custom HSTS configuration
	hstsHeader := rr.Header().Get("Strict-Transport-Security")
	expectedHSTS := "max-age=86400; preload"
	if hstsHeader != expectedHSTS {
		t.Errorf("Expected HSTS to be %s, got %s", expectedHSTS, hstsHeader)
	}

	// Check custom frame options
	frameOptions := rr.Header().Get("X-Frame-Options")
	if frameOptions != "SAMEORIGIN" {
		t.Errorf("Expected X-Frame-Options to be SAMEORIGIN, got %s", frameOptions)
	}

	// Check that disabled headers are not set
	if rr.Header().Get("X-Content-Type-Options") != "" {
		t.Error("Expected X-Content-Type-Options to not be set when disabled")
	}

	if rr.Header().Get("X-XSS-Protection") != "" {
		t.Error("Expected X-XSS-Protection to not be set when disabled")
	}

	// Check custom referrer policy
	referrerPolicy := rr.Header().Get("Referrer-Policy")
	if referrerPolicy != "no-referrer" {
		t.Errorf("Expected Referrer-Policy to be no-referrer, got %s", referrerPolicy)
	}

	// Check custom permissions policy
	permissionsPolicy := rr.Header().Get("Permissions-Policy")
	expectedPermissions := "geolocation=self"
	if permissionsPolicy != expectedPermissions {
		t.Errorf("Expected Permissions-Policy to be %s, got %s", expectedPermissions, permissionsPolicy)
	}

	// Check custom headers
	customHeader := rr.Header().Get("X-Custom-Security")
	if customHeader != "custom-value" {
		t.Errorf("Expected X-Custom-Security to be custom-value, got %s", customHeader)
	}

	// Check custom CSP
	cspHeader := rr.Header().Get("Content-Security-Policy")
	if !contains(cspHeader, "default-src 'self' https://trusted.com") {
		t.Errorf("Expected CSP to contain custom default-src, got %s", cspHeader)
	}
	if !contains(cspHeader, "style-src 'self' https://styles.com") {
		t.Errorf("Expected CSP to contain custom style-src, got %s", cspHeader)
	}
}

func TestSecurityHeadersMiddleware_DisabledCSP(t *testing.T) {
	config := &SecurityHeadersConfig{
		CSP: &CSPConfig{
			Enabled: false,
		},
	}

	middleware := NewSecurityHeadersMiddleware(config)

	req := httptest.NewRequest("GET", "http://example.com/test", nil)
	rr := httptest.NewRecorder()

	handler := middleware.GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	handler.ServeHTTP(rr, req)

	// CSP should not be set when disabled
	cspHeader := rr.Header().Get("Content-Security-Policy")
	if cspHeader != "" {
		t.Error("Expected Content-Security-Policy to not be set when disabled")
	}
}

func TestBuildCSPValue(t *testing.T) {
	config := &CSPConfig{
		Enabled:                 true,
		DefaultSrc:              []string{"'self'"},
		ScriptSrc:               []string{"'self'", "https://cdn.com"},
		StyleSrc:                []string{"'self'"},
		UpgradeInsecureRequests: true,
	}

	result := buildCSPValue(config)
	expected := "default-src 'self'; script-src 'self' https://cdn.com; style-src 'self'; upgrade-insecure-requests"

	if result != expected {
		t.Errorf("Expected CSP value %s, got %s", expected, result)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				indexOf(s, substr) >= 0)))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
