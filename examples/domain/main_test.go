package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dracory/rtr"
)

func TestDomainRouting(t *testing.T) {
	// Create a new router
	r := rtr.NewRouter()

	// Create domains
	apiDomain := rtr.NewDomain("api.example.com", "")
	adminDomain := rtr.NewDomain("admin.example.com", "")

	// Add routes to the API domain
	apiDomain.AddRoute(rtr.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))

	apiDomain.AddRoute(rtr.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`["user1", "user2"]`))
	}))

	// Add routes to the admin domain
	adminDomain.AddRoute(rtr.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(`<title>Admin Panel</title>`))
	}))

	// Add catch-all route for API domain
	apiDomain.AddRoute(rtr.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "Not Found", "message": "The requested resource was not found on this server"}`))
	}))

	// Add catch-all route for Admin domain
	adminDomain.AddRoute(rtr.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`<title>404 Not Found</title>`))
	}))

	// Add domains to the router
	r.AddDomain(apiDomain)
	r.AddDomain(adminDomain)

	// Create test server
	server := httptest.NewServer(r)
	defer server.Close()

	tests := []struct {
		name           string
		url            string
		host           string
		expectedStatus int
		expectedBody   string
		expectedType   string
	}{
		{
			name:           "API status endpoint",
			url:            "/status",
			host:           "api.example.com",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"ok"}`,
			expectedType:   "application/json",
		},
		{
			name:           "API users endpoint",
			url:            "/users",
			host:           "api.example.com",
			expectedStatus: http.StatusOK,
			expectedBody:   `["user1", "user2"]`,
			expectedType:   "application/json",
		},
		{
			name:           "API 404 endpoint",
			url:            "/nonexistent",
			host:           "api.example.com",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `Not Found`,
			expectedType:   "application/json",
		},
		{
			name:           "Admin root endpoint",
			url:            "/",
			host:           "admin.example.com",
			expectedStatus: http.StatusOK,
			expectedBody:   "<title>Admin Panel</title>",
			expectedType:   "text/html",
		},
		{
			name:           "Admin 404 endpoint",
			url:            "/nonexistent",
			host:           "admin.example.com",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "<title>404 Not Found</title>",
			expectedType:   "text/html",
		},
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", server.URL+tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Host = tt.host

			resp, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			// Check status code
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			// Check content type
			contentType := resp.Header.Get("Content-Type")
			if !strings.Contains(contentType, tt.expectedType) {
				t.Errorf("expected content type %s, got %s", tt.expectedType, contentType)
			}

			// Check response body
			buf := make([]byte, 1024)
			n, _ := resp.Body.Read(buf)
			body := string(buf[:n])

			if !strings.Contains(body, tt.expectedBody) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedBody, body)
			}
		})
	}
}
