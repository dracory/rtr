package router_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/dracory/router"
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
			d := router.NewDomain(tt.patterns...)
			got := d.GetPatterns()
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("GetPatterns() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDomain_AddRoute(t *testing.T) {
	d := router.NewDomain("example.com")
	r := router.NewRoute().SetPath("/test")
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
	d := router.NewDomain("example.com")
	r1 := router.NewRoute().SetPath("/test1")
	r2 := router.NewRoute().SetPath("/test2")
	d.AddRoutes([]router.RouteInterface{r1, r2})

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
	d := router.NewDomain("example.com")
	g := router.NewGroup().SetPrefix("/api")
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
	d := router.NewDomain("example.com")
	g1 := router.NewGroup().SetPrefix("/api")
	g2 := router.NewGroup().SetPrefix("/v1")
	d.AddGroups([]router.GroupInterface{g1, g2})

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
	d := router.NewDomain("example.com")

	// Test before middlewares
	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	}
	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	}

	d.AddBeforeMiddlewares([]router.Middleware{mw1, mw2})
	beforeMiddlewares := d.GetBeforeMiddlewares()
	if len(beforeMiddlewares) != 2 {
		t.Errorf("expected 2 before middlewares, got %d", len(beforeMiddlewares))
	}

	// Test after middlewares
	d.AddAfterMiddlewares([]router.Middleware{mw1})
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := router.NewDomain(tt.dPattern)
			got := d.Match(tt.host)
			if got != tt.expected {
				t.Errorf("Match(%q) = %v, want %v", tt.host, got, tt.expected)
			}
		})
	}
}

func TestDomain_Match_WithPort(t *testing.T) {
	d := router.NewDomain("example.com")
	if !d.Match("example.com:8080") {
		t.Error("expected match for domain with port")
	}
}

// Test helper structs

type testRoute struct {
	path string
}

func (r *testRoute) GetPath() string { return r.path }

// Implement other required RouteInterface methods with empty implementations
func (r *testRoute) GetMethod() string                                       { return "" }
func (r *testRoute) SetMethod(method string) router.RouteInterface           { return r }
func (r *testRoute) SetPath(path string) router.RouteInterface               { r.path = path; return r }
func (r *testRoute) GetHandler() router.Handler                              { return nil }
func (r *testRoute) SetHandler(handler router.Handler) router.RouteInterface { return r }
func (r *testRoute) GetName() string                                         { return "" }
func (r *testRoute) SetName(name string) router.RouteInterface               { return r }
func (r *testRoute) AddBeforeMiddlewares(middleware []router.Middleware) router.RouteInterface {
	return r
}
func (r *testRoute) GetBeforeMiddlewares() []router.Middleware { return nil }
func (r *testRoute) AddAfterMiddlewares(middleware []router.Middleware) router.RouteInterface {
	return r
}
func (r *testRoute) GetAfterMiddlewares() []router.Middleware { return nil }
