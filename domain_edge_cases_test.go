package rtr_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	rtr "github.com/dracory/rtr"
)

func TestDomainMatching(t *testing.T) {
	tests := []struct {
		name          string
		domainPattern string
		host          string
		shouldMatch   bool
	}{
		{
			name:          "exact match",
			domainPattern: "example.com",
			host:          "example.com",
			shouldMatch:   true,
		},
		{
			name:          "subdomain should not match",
			domainPattern: "example.com",
			host:          "api.example.com",
			shouldMatch:   false,
		},
		{
			name:          "wildcard subdomain match",
			domainPattern: "*.example.com",
			host:          "api.example.com",
			shouldMatch:   true,
		},
		{
			name:          "wildcard subdomain with multiple levels",
			domainPattern: "*.example.com",
			host:          "v1.api.example.com",
			shouldMatch:   true,
		},
		{
			name:          "wildcard subdomain no match",
			domainPattern: "*.example.com",
			host:          "example.com",
			shouldMatch:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := rtr.NewDomain(tc.domainPattern)
			if got := d.Match(tc.host); got != tc.shouldMatch {
				t.Errorf("Match() = %v, want %v for pattern %q and host %q", got, tc.shouldMatch, tc.domainPattern, tc.host)
			}
		})
	}
}

func TestDomainWithPort(t *testing.T) {
	tests := []struct {
		name          string
		domainPattern string
		host          string
		shouldMatch   bool
	}{
		{
			name:          "domain with port match",
			domainPattern: "example.com:8080",
			host:          "example.com:8080",
			shouldMatch:   true,
		},
		{
			name:          "domain with port no match",
			domainPattern: "example.com:8080",
			host:          "example.com:3000",
			shouldMatch:   false,
		},
		{
			name:          "wildcard port match",
			domainPattern: "example.com:*",
			host:          "example.com:8080",
			shouldMatch:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := rtr.NewDomain(tc.domainPattern)
			if got := d.Match(tc.host); got != tc.shouldMatch {
				t.Errorf("Match() = %v, want %v for pattern %q and host %q", got, tc.shouldMatch, tc.domainPattern, tc.host)
			}
		})
	}
}

func TestDomainRoutingEdgeCases(t *testing.T) {
	r := rtr.NewRouter()

	// Create a domain with a route
	domain := rtr.NewDomain("api.example.com")
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("domain handler"))
	}

	domain.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/test").
		SetHandler(handler))

	r.AddDomain(domain)

	// Test matching domain
	t.Run("matching domain", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://api.example.com/test", nil)
		req.Host = "api.example.com"
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
		}

		if rr.Body.String() != "domain handler" {
			t.Errorf("unexpected response body: %s", rr.Body.String())
		}
	})

	// Test non-matching domain
	t.Run("non-matching domain", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://other.com/test", nil)
		req.Host = "other.com"
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
		}
	})
}
