package rtr_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	rtr "github.com/dracory/rtr"
)

func TestDeclarativeAPI(t *testing.T) {
	// Test handlers
	homeHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Home"))
	}
	usersHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Users"))
	}
	createUserHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("User created"))
	}
	adminHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Admin"))
	}

	// Test middleware
	testMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test", "middleware")
			next.ServeHTTP(w, r)
		})
	}

	// Create router using declarative configuration
	config := rtr.RouterConfig{
		Name:             "Test Router",
		Prefix:           "",
		BeforeMiddleware: []rtr.Middleware{testMiddleware},
		Routes: []rtr.RouteConfig{
			rtr.GET("/", homeHandler).WithName("Home"),
		},
		Groups: []rtr.GroupConfig{
			rtr.Group("/api",
				rtr.GET("/users", usersHandler).WithName("List Users"),
				rtr.POST("/users", createUserHandler).WithName("Create User"),
				rtr.Group("/admin",
					rtr.GET("/dashboard", adminHandler).WithName("Admin Dashboard"),
				).WithName("Admin Group"),
			).WithName("API Group"),
		},
		Domains: []rtr.DomainConfig{
			rtr.Domain([]string{"admin.example.com"},
				rtr.GET("/", adminHandler).WithName("Admin Home"),
			),
		},
	}

	router := rtr.NewRouterFromConfig(config)

	// Test the routes
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
		expectedHeader string
	}{
		{
			name:           "Home route",
			method:         "GET",
			path:           "/",
			expectedStatus: 200,
			expectedBody:   "Home",
			expectedHeader: "middleware",
		},
		{
			name:           "API Users route",
			method:         "GET",
			path:           "/api/users",
			expectedStatus: 200,
			expectedBody:   "Users",
			expectedHeader: "middleware",
		},
		{
			name:           "API Create User route",
			method:         "POST",
			path:           "/api/users",
			expectedStatus: 200,
			expectedBody:   "User created",
			expectedHeader: "middleware",
		},
		{
			name:           "Admin Dashboard route",
			method:         "GET",
			path:           "/api/admin/dashboard",
			expectedStatus: 200,
			expectedBody:   "Admin",
			expectedHeader: "middleware",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if w.Body.String() != tt.expectedBody {
				t.Errorf("Expected body %q, got %q", tt.expectedBody, w.Body.String())
			}

			if w.Header().Get("X-Test") != tt.expectedHeader {
				t.Errorf("Expected header %q, got %q", tt.expectedHeader, w.Header().Get("X-Test"))
			}
		})
	}
}

func TestDeclarativeHelperFunctions(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {}

	// Test route helper functions
	getRoute := rtr.GET("/test", handler)
	if getRoute.Method != "GET" || getRoute.Path != "/test" {
		t.Error("GET helper function failed")
	}

	postRoute := rtr.POST("/test", handler)
	if postRoute.Method != "POST" || postRoute.Path != "/test" {
		t.Error("POST helper function failed")
	}

	putRoute := rtr.PUT("/test", handler)
	if putRoute.Method != "PUT" || putRoute.Path != "/test" {
		t.Error("PUT helper function failed")
	}

	deleteRoute := rtr.DELETE("/test", handler)
	if deleteRoute.Method != "DELETE" || deleteRoute.Path != "/test" {
		t.Error("DELETE helper function failed")
	}

	patchRoute := rtr.PATCH("/test", handler)
	if patchRoute.Method != "PATCH" || patchRoute.Path != "/test" {
		t.Error("PATCH helper function failed")
	}

	optionsRoute := rtr.OPTIONS("/test", handler)
	if optionsRoute.Method != "OPTIONS" || optionsRoute.Path != "/test" {
		t.Error("OPTIONS helper function failed")
	}
}

func TestRouteConfigChaining(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {}
	middleware := func(next http.Handler) http.Handler { return next }

	route := rtr.GET("/test", handler).
		WithName("Test Route").
		WithBeforeMiddleware(middleware).
		WithAfterMiddleware(middleware).
		WithMetadata("version", "1.0")

	if route.Name != "Test Route" {
		t.Error("WithName chaining failed")
	}

	if len(route.BeforeMiddleware) != 1 {
		t.Error("WithBeforeMiddleware chaining failed")
	}

	if len(route.AfterMiddleware) != 1 {
		t.Error("WithAfterMiddleware chaining failed")
	}

	if route.Metadata["version"] != "1.0" {
		t.Error("WithMetadata chaining failed")
	}
}

func TestGroupConfigChaining(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {}
	middleware := func(next http.Handler) http.Handler { return next }

	group := rtr.Group("/api",
		rtr.GET("/users", handler),
		rtr.POST("/users", handler),
	).WithName("API Group").
		WithBeforeMiddleware(middleware).
		WithAfterMiddleware(middleware)

	if group.Name != "API Group" {
		t.Error("Group WithName chaining failed")
	}

	if len(group.BeforeMiddleware) != 1 {
		t.Error("Group WithBeforeMiddleware chaining failed")
	}

	if len(group.AfterMiddleware) != 1 {
		t.Error("Group WithAfterMiddleware chaining failed")
	}

	if len(group.Routes) != 2 {
		t.Error("Group route addition failed")
	}
}

func TestDomainConfig(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Admin"))
	}

	config := rtr.RouterConfig{
		Domains: []rtr.DomainConfig{
			rtr.Domain([]string{"admin.example.com", "*.admin.example.com"},
				rtr.GET("/", handler).WithName("Admin Home"),
				rtr.Group("/api",
					rtr.GET("/users", handler).WithName("Admin API Users"),
				),
			),
		},
	}

	router := rtr.NewRouterFromConfig(config)

	// Test domain route with correct host
	req := httptest.NewRequest("GET", "/", nil)
	req.Host = "admin.example.com"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Domain route failed, status: %d", w.Code)
	}

	if w.Body.String() != "Admin" {
		t.Errorf("Expected 'Admin', got %q", w.Body.String())
	}
}
