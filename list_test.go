package rtr

import (
	"net/http"
	"strings"
	"testing"
)

// TestList verifies that the List method works without panicking
func TestList(t *testing.T) {
	// Create a router with various configurations
	router := NewRouter()

	// Add some middleware
	router.AddBeforeMiddlewares([]Middleware{
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		},
	})

	// Add direct routes
	router.AddRoute(NewRoute().
		SetMethod("GET").
		SetPath("/").
		SetName("Home").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Home"))
		}))

	router.AddRoute(NewRoute().
		SetMethod("POST").
		SetPath("/users").
		SetName("Create User").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("User created"))
		}))

	// Add a route group
	apiGroup := NewGroup().SetPrefix("/api")
	apiGroup.AddRoute(NewRoute().
		SetMethod("GET").
		SetPath("/users").
		SetName("List Users").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Users list"))
		}))

	// Add middleware to the group
	apiGroup.AddBeforeMiddlewares([]Middleware{
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		},
	})

	router.AddGroup(apiGroup)

	// Add a domain
	domain := NewDomain("example.com")
	domain.AddRoute(NewRoute().
		SetMethod("GET").
		SetPath("/admin").
		SetName("Admin Panel").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Admin"))
		}))

	router.AddDomain(domain)

	// Test that List doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("List method panicked: %v", r)
		}
	}()

	// Call List - this should print formatted tables to stdout
	router.List()

	// If we reach here without panicking, the test passes
	t.Log("List method executed successfully")
}

// TestGetMiddlewareName tests the middleware name extraction
func TestGetMiddlewareName(t *testing.T) {
	// Test with RecoveryMiddleware
	name := getMiddlewareName(RecoveryMiddleware)
	if name == "" {
		t.Error("Expected non-empty middleware name")
	}
	t.Logf("RecoveryMiddleware name: %s", name)

	// Test with anonymous middleware
	anonymousMiddleware := func(next http.Handler) http.Handler {
		return next
	}
	name = getMiddlewareName(anonymousMiddleware)
	if name == "" {
		t.Error("Expected non-empty middleware name for anonymous middleware")
	}
	t.Logf("Anonymous middleware name: %s", name)

	// Test with nil middleware
	name = getMiddlewareName(nil)
	if name != "nil" {
		t.Errorf("Expected 'nil' for nil middleware, got: %s", name)
	}
}

// TestRouteMiddlewareNames tests route middleware name extraction
func TestRouteMiddlewareNames(t *testing.T) {
	route := NewRoute()
	
	// Test route with no middleware
	names := getRouteMiddlewareNames(route)
	if len(names) != 1 || names[0] != "none" {
		t.Errorf("Expected ['none'] for route with no middleware, got: %v", names)
	}

	// Add some middleware
	route.AddBeforeMiddlewares([]Middleware{
		func(next http.Handler) http.Handler { return next },
	})
	route.AddAfterMiddlewares([]Middleware{
		func(next http.Handler) http.Handler { return next },
	})

	names = getRouteMiddlewareNames(route)
	if len(names) != 2 {
		t.Errorf("Expected 2 middleware names, got: %d", len(names))
	}

	// Check that names contain expected suffixes
	hasBeforeMiddleware := false
	hasAfterMiddleware := false
	for _, name := range names {
		if strings.Contains(name, "(before)") {
			hasBeforeMiddleware = true
		}
		if strings.Contains(name, "(after)") {
			hasAfterMiddleware = true
		}
	}

	if !hasBeforeMiddleware {
		t.Error("Expected to find before middleware in names")
	}
	if !hasAfterMiddleware {
		t.Error("Expected to find after middleware in names")
	}
}
