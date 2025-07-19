package rtr_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	rtr "github.com/dracory/rtr"
)

func TestDomainRouting(t *testing.T) {
	tests := []struct {
		name           string
		domainPattern  string
		requestHost    string
		expectedMatch  bool
		expectedStatus int
	}{
		{
			name:           "exact domain match",
			domainPattern:  "example.com",
			requestHost:    "example.com",
			expectedMatch:  true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "subdomain not matching",
			domainPattern:  "example.com",
			requestHost:    "api.example.com",
			expectedMatch:  false,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "wildcard subdomain match",
			domainPattern:  "*.example.com",
			requestHost:    "api.example.com",
			expectedMatch:  true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "wildcard subdomain with multiple levels",
			domainPattern:  "*.example.com",
			requestHost:    "v1.api.example.com",
			expectedMatch:  true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "wildcard subdomain no match",
			domainPattern:  "*.example.com",
			requestHost:    "example.com",
			expectedMatch:  false,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "port in pattern",
			domainPattern:  "example.com:8080",
			requestHost:    "example.com:8080",
			expectedMatch:  true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "port mismatch",
			domainPattern:  "example.com:8080",
			requestHost:    "example.com:3000",
			expectedMatch:  false,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "wildcard port",
			domainPattern:  "example.com:*",
			requestHost:    "example.com:8080",
			expectedMatch:  true,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new router
			r := rtr.NewRouter()

			// Create a domain with the test pattern
			domain := rtr.NewDomain(tc.domainPattern)

			// Add a test route to the domain
			handler := func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			}

			domain.AddRoute(rtr.NewRoute().
				SetMethod(http.MethodGet).
				SetPath("/test").
				SetHandler(handler))

			// Add the domain to the router
			r.AddDomain(domain)

			// Create a test request
			req := httptest.NewRequest(http.MethodGet, "http://"+tc.requestHost+"/test", nil)
			req.Host = tc.requestHost // Important: Set the Host header

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Serve the request
			r.ServeHTTP(rr, req)

			// Check the response status code
			if tc.expectedMatch && rr.Code != http.StatusOK {
				t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
			} else if !tc.expectedMatch && rr.Code != http.StatusNotFound {
				t.Errorf("expected status %d for non-matching domain, got %d", http.StatusNotFound, rr.Code)
			}
		})
	}
}

func TestDomainRoutingWithGroups(t *testing.T) {
	r := rtr.NewRouter()

	// Create an API domain with versioned groups
	apiDomain := rtr.NewDomain("api.example.com")

	// Create version 1 group
	v1Group := rtr.NewGroup().SetPrefix("/v1")
	v1Group.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/users").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("v1 users"))
		}))

	// Create version 2 group
	v2Group := rtr.NewGroup().SetPrefix("/v2")
	v2Group.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/users").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("v2 users"))
		}))

	// Add groups to domain
	apiDomain.AddGroup(v1Group)
	apiDomain.AddGroup(v2Group)

	// Add domain to router
	r.AddDomain(apiDomain)

	// Test v1 endpoint
	t.Run("v1 API endpoint", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://api.example.com/v1/users", nil)
		req.Host = "api.example.com"
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK || rr.Body.String() != "v1 users" {
			t.Errorf("v1 endpoint failed: status=%d, body=%s", rr.Code, rr.Body.String())
		}
	})

	// Test v2 endpoint
	t.Run("v2 API endpoint", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://api.example.com/v2/users", nil)
		req.Host = "api.example.com"
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK || rr.Body.String() != "v2 users" {
			t.Errorf("v2 endpoint failed: status=%d, body=%s", rr.Code, rr.Body.String())
		}
	})

	// Test non-matching domain
	t.Run("non-matching domain", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://other.com/v1/users", nil)
		req.Host = "other.com"
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected 404 for non-matching domain, got %d", rr.Code)
		}
	})
}

func TestDomainRoutingWithMiddleware(t *testing.T) {
	r := rtr.NewRouter()

	// Create a test domain with middleware
	domain := rtr.NewDomain("auth.example.com")

	// Add middleware that adds a custom header
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test-Middleware", "executed")
			next.ServeHTTP(w, r)
		})
	}

	domain.AddBeforeMiddlewares(rtr.MiddlewaresToInterfaces([]rtr.Middleware{middleware}))

	// Add a test route
	domain.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/secure").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("secure content"))
		}))

	r.AddDomain(domain)

	// Test the endpoint
	req := httptest.NewRequest(http.MethodGet, "http://auth.example.com/secure", nil)
	req.Host = "auth.example.com"
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	if rr.Header().Get("X-Test-Middleware") != "executed" {
		t.Error("domain middleware was not executed")
	}
}
