package main_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	main "github.com/dracory/rtr/examples/declarative"
)

func TestDeclarativeRouter(t *testing.T) {
	router := main.CreateDeclarativeRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	tests := []struct {
		name           string
		method         string
		url            string
		host           string
		headers        map[string]string
		expectedStatus int
		expectedBody   string
		expectedType   string
	}{
		{
			name:           "Home endpoint",
			method:         "GET",
			url:            "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "Welcome to Declarative Router!",
			expectedType:   "application/json",
		},
		{
			name:           "API users list without auth",
			method:         "GET",
			url:            "/api/users",
			expectedStatus: http.StatusOK,
			expectedBody:   "alice",
			expectedType:   "application/json",
		},
		{
			name:           "API users list with auth",
			method:         "GET",
			url:            "/api/users",
			headers:        map[string]string{"Authorization": "Bearer token"},
			expectedStatus: http.StatusOK,
			expectedBody:   "alice",
			expectedType:   "application/json",
		},
		{
			name:           "API create user",
			method:         "POST",
			url:            "/api/users",
			expectedStatus: http.StatusCreated,
			expectedBody:   "User created successfully",
			expectedType:   "application/json",
		},
		{
			name:           "API v1 products",
			method:         "GET",
			url:            "/api/v1/products",
			expectedStatus: http.StatusOK,
			expectedBody:   "laptop",
			expectedType:   "application/json",
		},
		{
			name:           "API v2 products",
			method:         "GET",
			url:            "/api/v2/products",
			expectedStatus: http.StatusOK,
			expectedBody:   "laptop",
			expectedType:   "application/json",
		},
		{
			name:           "Admin domain home",
			method:         "GET",
			url:            "/",
			host:           "admin.example.com",
			expectedStatus: http.StatusOK,
			expectedBody:   "Admin Dashboard",
			expectedType:   "application/json",
		},
		{
			name:           "Admin domain stats",
			method:         "GET",
			url:            "/api/stats",
			host:           "admin.example.com",
			expectedStatus: http.StatusOK,
			expectedBody:   "users",
			expectedType:   "application/json",
		},
		{
			name:           "Wildcard admin domain",
			method:         "GET",
			url:            "/",
			host:           "sub.admin.example.com",
			expectedStatus: http.StatusOK,
			expectedBody:   "Admin Dashboard",
			expectedType:   "application/json",
		},
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, server.URL+tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			if tt.host != "" {
				req.Host = tt.host
			}

			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			// Check status code
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			// Check content type if specified
			if tt.expectedType != "" {
				contentType := resp.Header.Get("Content-Type")
				if !strings.Contains(contentType, tt.expectedType) {
					t.Errorf("expected content type to contain %s, got %s", tt.expectedType, contentType)
				}
			}

			// Check response body
			buf := make([]byte, 1024)
			n, _ := resp.Body.Read(buf)
			body := string(buf[:n])

			if tt.expectedBody != "" && !strings.Contains(body, tt.expectedBody) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedBody, body)
			}
		})
	}
}

func TestHybridRouter(t *testing.T) {
	router := main.CreateHybridRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Declarative home",
			method:         "GET",
			url:            "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "Hybrid Router Home",
		},
		{
			name:           "Declarative API users",
			method:         "GET",
			url:            "/api/users",
			expectedStatus: http.StatusOK,
			expectedBody:   "users",
		},
		{
			name:           "Imperative health check",
			method:         "GET",
			url:            "/health",
			expectedStatus: http.StatusOK,
			expectedBody:   "ok",
		},
		{
			name:           "Imperative admin dashboard",
			method:         "GET",
			url:            "/admin/dashboard",
			expectedStatus: http.StatusOK,
			expectedBody:   "Admin Dashboard",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, server.URL+tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			buf := make([]byte, 1024)
			n, _ := resp.Body.Read(buf)
			body := string(buf[:n])

			if !strings.Contains(body, tt.expectedBody) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedBody, body)
			}
		})
	}
}

func TestImperativeRouter(t *testing.T) {
	router := main.CreateImperativeRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Home endpoint",
			method:         "GET",
			url:            "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "Imperative Router Home",
		},
		{
			name:           "API users",
			method:         "GET",
			url:            "/api/users",
			expectedStatus: http.StatusOK,
			expectedBody:   "users",
		},
		{
			name:           "API v1 products",
			method:         "GET",
			url:            "/api/v1/products",
			expectedStatus: http.StatusOK,
			expectedBody:   "products",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, server.URL+tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			buf := make([]byte, 1024)
			n, _ := resp.Body.Read(buf)
			body := string(buf[:n])

			if !strings.Contains(body, tt.expectedBody) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedBody, body)
			}
		})
	}
}

func TestRouterConfiguration(t *testing.T) {
	router := main.CreateDeclarativeRouter()

	// Test that router was created successfully
	if router == nil {
		t.Fatal("router should not be nil")
	}

	// Test that routes were registered (this would depend on router implementation)
	// For now, we'll just verify the router responds to requests
	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestMiddlewareExecution(t *testing.T) {
	router := main.CreateDeclarativeRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	// Test that auth middleware sets header when no authorization
	req, err := http.NewRequest("GET", server.URL+"/api/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Check if auth middleware set the header
	authRequired := resp.Header.Get("X-Auth-Required")
	if authRequired != "true" {
		t.Errorf("expected X-Auth-Required header to be 'true', got %q", authRequired)
	}

	// Test with authorization header
	req, err = http.NewRequest("GET", server.URL+"/api/users", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer token")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Check that auth middleware didn't set the header
	authRequired = resp.Header.Get("X-Auth-Required")
	if authRequired == "true" {
		t.Error("X-Auth-Required header should not be set when Authorization header is present")
	}
}
