package rtr_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/dracory/rtr"
)

func TestNewDomain(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		expected []string
	}{
		{
			name:     "single pattern",
			patterns: []string{"example.com"},
			expected: []string{"example.com"},
		},
		{
			name:     "multiple patterns",
			patterns: []string{"example.com", "*.example.com"},
			expected: []string{"example.com", "*.example.com"},
		},
		{
			name:     "no patterns",
			patterns: []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := rtr.NewDomain(tt.patterns...)
			got := d.GetPatterns()
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("GetPatterns() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDomain_AddRoute(t *testing.T) {
	d := rtr.NewDomain("example.com")
	r := rtr.NewRoute().SetPath("/test")
	d.AddRoute(r)

	routes := d.GetRoutes()
	if len(routes) != 1 {
		t.Fatalf("expected 1 route, got %d", len(routes))
	}
	if routes[0].GetPath() != "/test" {
		t.Errorf("expected path '/test', got '%s'", routes[0].GetPath())
	}
}

func TestDomain_AddRoutes(t *testing.T) {
	d := rtr.NewDomain("example.com")
	r1 := rtr.NewRoute().SetPath("/test1")
	r2 := rtr.NewRoute().SetPath("/test2")
	d.AddRoutes([]rtr.RouteInterface{r1, r2})

	routes := d.GetRoutes()
	if len(routes) != 2 {
		t.Fatalf("expected 2 routes, got %d", len(routes))
	}
	if routes[0].GetPath() != "/test1" {
		t.Errorf("first route: expected path '/test1', got '%s'", routes[0].GetPath())
	}
	if routes[1].GetPath() != "/test2" {
		t.Errorf("second route: expected path '/test2', got '%s'", routes[1].GetPath())
	}
}

func TestDomain_AddGroup(t *testing.T) {
	d := rtr.NewDomain("example.com")
	g := rtr.NewGroup().SetPrefix("/api")
	d.AddGroup(g)

	groups := d.GetGroups()
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if got := groups[0].GetPrefix(); got != "/api" {
		t.Errorf("expected prefix '/api', got '%s'", got)
	}
}

func TestDomain_AddGroups(t *testing.T) {
	d := rtr.NewDomain("example.com")
	g1 := rtr.NewGroup().SetPrefix("/api")
	g2 := rtr.NewGroup().SetPrefix("/v1")
	d.AddGroups([]rtr.GroupInterface{g1, g2})

	groups := d.GetGroups()
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if got := groups[0].GetPrefix(); got != "/api" {
		t.Errorf("first group: expected prefix '/api', got '%s'", got)
	}
	if got := groups[1].GetPrefix(); got != "/v1" {
		t.Errorf("second group: expected prefix '/v1', got '%s'", got)
	}
}

func TestDomain_Middleware(t *testing.T) {
	d := rtr.NewDomain("example.com")

	// Test before middlewares
	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	}
	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	}

	d.AddBeforeMiddlewares(rtr.MiddlewaresToInterfaces([]rtr.StdMiddleware{mw1, mw2}))
	beforeMiddlewares := d.GetBeforeMiddlewares()
	if len(beforeMiddlewares) != 2 {
		t.Errorf("expected 2 before middlewares, got %d", len(beforeMiddlewares))
	}

	// Test after middlewares
	d.AddAfterMiddlewares(rtr.MiddlewaresToInterfaces([]rtr.StdMiddleware{mw1}))
	afterMiddlewares := d.GetAfterMiddlewares()
	if len(afterMiddlewares) != 1 {
		t.Errorf("expected 1 after middleware, got %d", len(afterMiddlewares))
	}
}

func TestDomain_Match(t *testing.T) {
	tests := []struct {
		name     string
		dPattern string
		host     string
		expected bool
	}{
		// Basic domain matching
		{
			name:     "exact match",
			dPattern: "example.com",
			host:     "example.com",
			expected: true,
		},
		{
			name:     "wildcard subdomain match",
			dPattern: "*.example.com",
			host:     "api.example.com",
			expected: true,
		},
		{
			name:     "wildcard subdomain match with multiple subdomains",
			dPattern: "*.example.com",
			host:     "v1.api.example.com",
			expected: true,
		},
		{
			name:     "no match",
			dPattern: "example.com",
			host:     "api.example.com",
			expected: false,
		},
		{
			name:     "wildcard subdomain no match",
			dPattern: "*.example.com",
			host:     "example.com",
			expected: false,
		},

		// Port matching tests
		{
			name:     "domain with port matches any port when no port in pattern",
			dPattern: "example.com",
			host:     "example.com:8080",
			expected: true,
		},
		{
			name:     "exact port match",
			dPattern: "example.com:8080",
			host:     "example.com:8080",
			expected: true,
		},
		{
			name:     "port mismatch",
			dPattern: "example.com:8080",
			host:     "example.com:8081",
			expected: false,
		},
		{
			name:     "wildcard port matches any port",
			dPattern: "example.com:*",
			host:     "example.com:8080",
			expected: true,
		},
		{
			name:     "wildcard port with subdomain",
			dPattern: "*.example.com:*",
			host:     "api.example.com:8080",
			expected: true,
		},
		{
			name:     "IPv4 with port",
			dPattern: "127.0.0.1:8080",
			host:     "127.0.0.1:8080",
			expected: true,
		},
		{
			name:     "IPv6 with port",
			dPattern: "[::1]:8080",
			host:     "[::1]:8080",
			expected: true,
		},
		{
			name:     "IPv6 with different port",
			dPattern: "[::1]:8080",
			host:     "[::1]:8081",
			expected: false,
		},
		{
			name:     "IPv6 with wildcard port",
			dPattern: "[::1]:*",
			host:     "[::1]:8080",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := rtr.NewDomain(tt.dPattern)
			got := d.Match(tt.host)
			if got != tt.expected {
				t.Errorf("Match(%q) = %v, want %v", tt.host, got, tt.expected)
			}
		})
	}
}

func TestDomain_Match_WithPort(t *testing.T) {
	tests := []struct {
		name     string
		dPattern string
		host     string
		expected bool
	}{
		{
			name:     "domain without port matches any port",
			dPattern: "example.com",
			host:     "example.com:8080",
			expected: true,
		},
		{
			name:     "domain with specific port matches exact port",
			dPattern: "example.com:3000",
			host:     "example.com:3000",
			expected: true,
		},
		{
			name:     "domain with specific port doesn't match different port",
			dPattern: "example.com:3000",
			host:     "example.com:3001",
			expected: false,
		},
		{
			name:     "domain with wildcard port matches any port",
			dPattern: "example.com:*",
			host:     "example.com:8080",
			expected: true,
		},
		{
			name:     "multiple patterns with different ports",
			dPattern: "example.com:8080,example.com:8081,example.com:8082",
			host:     "example.com:8081",
			expected: true,
		},
		{
			name:     "multiple patterns with wildcard port",
			dPattern: "example.com:8080,example.com:*,example.com:8082",
			host:     "example.com:9999",
			expected: true,
		},
		{
			name:     "multiple patterns with no match",
			dPattern: "example.com:8080,example.com:8081",
			host:     "example.com:9999",
			expected: false,
		},
		{
			name:     "multiple domains with different ports",
			dPattern: "api.example.com:8080,admin.example.com:8081",
			host:     "admin.example.com:8081",
			expected: true,
		},
		{
			name:     "multiple domains with port mismatch",
			dPattern: "api.example.com:8080,admin.example.com:8081",
			host:     "api.example.com:8081",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := rtr.NewDomain(tt.dPattern)
			got := d.Match(tt.host)
			if got != tt.expected {
				t.Errorf("Match(%q) with pattern %q = %v, want %v", tt.host, tt.dPattern, got, tt.expected)
			}
		})
	}
}
