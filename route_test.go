package rtr_test

import (
	"net/http"
	"testing"

	"github.com/dracory/rtr"
)

// TestRouteInterface tests the basic functionality of the RouteInterface implementation.
// It verifies that a route can be created, and that its method, path, name, handler, and middleware
// can be properly set and retrieved.
func TestRouteInterface(t *testing.T) {
	// Create a new route
	route := rtr.NewRoute()

	// Test method getter/setter
	method := "GET"
	route.SetMethod(method)
	if route.GetMethod() != method {
		t.Errorf("Expected method %s, got %s", method, route.GetMethod())
	}

	// Test path getter/setter
	path := "/test/path"
	route.SetPath(path)
	if route.GetPath() != path {
		t.Errorf("Expected path %s, got %s", path, route.GetPath())
	}

	// Test name getter/setter
	name := "test-route"
	route.SetName(name)
	if route.GetName() != name {
		t.Errorf("Expected name %s, got %s", name, route.GetName())
	}

	// Test handler getter/setter
	handler := func(w http.ResponseWriter, r *http.Request) {}
	route.SetHandler(handler)
	if route.GetHandler() == nil {
		t.Errorf("Expected handler to be set, got nil")
	}

	// Test middleware getters/setters
	middleware1 := rtr.Middleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	})
	middleware2 := rtr.Middleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	})

	// Test before middlewares
	route.AddBeforeMiddlewares([]rtr.Middleware{middleware1})
	if len(route.GetBeforeMiddlewares()) != 1 {
		t.Errorf("Expected 1 before middleware, got %d", len(route.GetBeforeMiddlewares()))
	}

	route.AddBeforeMiddlewares([]rtr.Middleware{middleware2})
	if len(route.GetBeforeMiddlewares()) != 2 {
		t.Errorf("Expected 2 before middlewares, got %d", len(route.GetBeforeMiddlewares()))
	}

	// Test after middlewares
	route.AddAfterMiddlewares([]rtr.Middleware{middleware1})
	if len(route.GetAfterMiddlewares()) != 1 {
		t.Errorf("Expected 1 after middleware, got %d", len(route.GetAfterMiddlewares()))
	}

	route.AddAfterMiddlewares([]rtr.Middleware{middleware2})
	if len(route.GetAfterMiddlewares()) != 2 {
		t.Errorf("Expected 2 after middlewares, got %d", len(route.GetAfterMiddlewares()))
	}
}

// TestRouteChaining tests the method chaining functionality of the RouteInterface.
// It verifies that multiple methods can be called in sequence and that the route
// state is correctly updated after each method call.
func TestRouteChaining(t *testing.T) {
	// Test method chaining
	route := rtr.NewRoute().
		SetMethod("POST").
		SetPath("/api/users").
		SetName("create-user").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {})

	if route.GetMethod() != "POST" {
		t.Errorf("Expected method POST, got %s", route.GetMethod())
	}

	if route.GetPath() != "/api/users" {
		t.Errorf("Expected path /api/users, got %s", route.GetPath())
	}

	if route.GetName() != "create-user" {
		t.Errorf("Expected name create-user, got %s", route.GetName())
	}

	if route.GetHandler() == nil {
		t.Errorf("Expected handler to be set, got nil")
	}

	// Test middleware chaining
	middleware := func(next http.Handler) http.Handler { return next }

	route.AddBeforeMiddlewares([]rtr.Middleware{middleware}).
		AddAfterMiddlewares([]rtr.Middleware{middleware})

	if len(route.GetBeforeMiddlewares()) != 1 {
		t.Errorf("Expected 1 before middleware, got %d", len(route.GetBeforeMiddlewares()))
	}

	if len(route.GetAfterMiddlewares()) != 1 {
		t.Errorf("Expected 1 after middleware, got %d", len(route.GetAfterMiddlewares()))
	}
}

// TestMultipleRoutes tests the creation and manipulation of multiple routes.
// It verifies that multiple routes can be created independently and that modifications
// to one route do not affect other routes.
func TestMultipleRoutes(t *testing.T) {
	// Create multiple routes and ensure they don't interfere with each other
	route1 := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/route1").
		SetName("route1")

	route2 := rtr.NewRoute().
		SetMethod("POST").
		SetPath("/route2").
		SetName("route2")

	// Verify route1 properties
	if route1.GetMethod() != "GET" {
		t.Errorf("Expected route1 method GET, got %s", route1.GetMethod())
	}

	if route1.GetPath() != "/route1" {
		t.Errorf("Expected route1 path /route1, got %s", route1.GetPath())
	}

	if route1.GetName() != "route1" {
		t.Errorf("Expected route1 name route1, got %s", route1.GetName())
	}

	// Verify route2 properties
	if route2.GetMethod() != "POST" {
		t.Errorf("Expected route2 method POST, got %s", route2.GetMethod())
	}

	if route2.GetPath() != "/route2" {
		t.Errorf("Expected route2 path /route2, got %s", route2.GetPath())
	}

	if route2.GetName() != "route2" {
		t.Errorf("Expected route2 name route2, got %s", route2.GetName())
	}

	// Modify route1 and ensure route2 is unaffected
	route1.SetMethod("PUT")

	if route1.GetMethod() != "PUT" {
		t.Errorf("Expected route1 method PUT, got %s", route1.GetMethod())
	}

	if route2.GetMethod() != "POST" {
		t.Errorf("Expected route2 method to remain POST, got %s", route2.GetMethod())
	}
}

// TestRouteWithMiddlewares tests the middleware functionality of the RouteInterface.
// It verifies that middleware functions can be added to a route and that they are
// properly stored and can be retrieved.
func TestRouteWithMiddlewares(t *testing.T) {
	// Create a route with multiple middlewares
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/middleware-test")

	// Create middlewares that modify a counter
	counter := 0

	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter += 1
			next.ServeHTTP(w, r)
		})
	}

	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter += 10
			next.ServeHTTP(w, r)
		})
	}

	middleware3 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter += 100
			next.ServeHTTP(w, r)
		})
	}

	// Add middlewares to the route
	route.AddBeforeMiddlewares([]rtr.Middleware{
		rtr.Middleware(middleware1),
		rtr.Middleware(middleware2),
	})
	route.AddAfterMiddlewares([]rtr.Middleware{
		rtr.Middleware(middleware3),
	})

	// Check middleware counts
	if len(route.GetBeforeMiddlewares()) != 2 {
		t.Errorf("Expected 2 before middlewares, got %d", len(route.GetBeforeMiddlewares()))
	}

	if len(route.GetAfterMiddlewares()) != 1 {
		t.Errorf("Expected 1 after middleware, got %d", len(route.GetAfterMiddlewares()))
	}
}

// TestRouteShortcuts tests the shortcut methods for creating routes with different HTTP methods.
// It verifies that each shortcut method correctly sets the method, path, and handler.
func TestRouteShortcuts(t *testing.T) {
	// Test handler for all routes
	handler := func(w http.ResponseWriter, r *http.Request) {}
	path := "/test"

	// Test GET shortcut
	getRoute := rtr.NewRoute().SetMethod(http.MethodGet).SetPath(path).SetHandler(handler)
	if getRoute.GetMethod() != http.MethodGet {
		t.Errorf("Expected GET method, got %s", getRoute.GetMethod())
	}
	if getRoute.GetPath() != path {
		t.Errorf("Expected path %s, got %s", path, getRoute.GetPath())
	}
	if getRoute.GetHandler() == nil {
		t.Error("Expected handler to be set, got nil")
	}

	// Test POST shortcut
	postRoute := rtr.NewRoute().SetMethod(http.MethodPost).SetPath(path).SetHandler(handler)
	if postRoute.GetMethod() != http.MethodPost {
		t.Errorf("Expected POST method, got %s", postRoute.GetMethod())
	}
	if postRoute.GetPath() != path {
		t.Errorf("Expected path %s, got %s", path, postRoute.GetPath())
	}
	if postRoute.GetHandler() == nil {
		t.Error("Expected handler to be set, got nil")
	}

	// Test PUT shortcut
	putRoute := rtr.NewRoute().SetMethod(http.MethodPut).SetPath(path).SetHandler(handler)
	if putRoute.GetMethod() != http.MethodPut {
		t.Errorf("Expected PUT method, got %s", putRoute.GetMethod())
	}
	if putRoute.GetPath() != path {
		t.Errorf("Expected path %s, got %s", path, putRoute.GetPath())
	}
	if putRoute.GetHandler() == nil {
		t.Error("Expected handler to be set, got nil")
	}

	// Test DELETE shortcut
	deleteRoute := rtr.NewRoute().SetMethod(http.MethodDelete).SetPath(path).SetHandler(handler)
	if deleteRoute.GetMethod() != http.MethodDelete {
		t.Errorf("Expected DELETE method, got %s", deleteRoute.GetMethod())
	}
	if deleteRoute.GetPath() != path {
		t.Errorf("Expected path %s, got %s", path, deleteRoute.GetPath())
	}
	if deleteRoute.GetHandler() == nil {
		t.Error("Expected handler to be set, got nil")
	}
}
