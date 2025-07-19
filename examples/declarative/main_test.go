package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dracory/rtr"
)

func TestCreateDeclarativeConfiguration(t *testing.T) {
	domain := CreateDeclarativeConfiguration()

	// Test domain properties
	if domain.Name != "Example API" {
		t.Errorf("Expected domain name 'Example API', got %s", domain.Name)
	}

	if domain.Status != rtr.StatusEnabled {
		t.Errorf("Expected domain status %s, got %s", rtr.StatusEnabled, domain.Status)
	}

	if len(domain.Hosts) != 2 {
		t.Errorf("Expected 2 hosts, got %d", len(domain.Hosts))
	}

	expectedHosts := []string{"api.example.com", "*.api.example.com"}
	for i, expectedHost := range expectedHosts {
		if i >= len(domain.Hosts) || domain.Hosts[i] != expectedHost {
			t.Errorf("Expected host[%d] to be %s, got %s", i, expectedHost, domain.Hosts[i])
		}
	}

	// Test domain-level middleware
	if len(domain.Middlewares) != 2 {
		t.Errorf("Expected 2 domain middlewares, got %d", len(domain.Middlewares))
	}

	expectedMiddlewares := []string{"cors", "logging"}
	for i, expectedMw := range expectedMiddlewares {
		if i >= len(domain.Middlewares) || domain.Middlewares[i] != expectedMw {
			t.Errorf("Expected middleware[%d] to be %s, got %s", i, expectedMw, domain.Middlewares[i])
		}
	}

	// Test items count
	if len(domain.Items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(domain.Items))
	}

	// Test home route
	homeRoute, ok := domain.Items[0].(*rtr.Route)
	if !ok {
		t.Errorf("Expected first item to be a Route")
	} else {
		if homeRoute.Name != "Home" {
			t.Errorf("Expected home route name 'Home', got %s", homeRoute.Name)
		}
		if homeRoute.Path != "/" {
			t.Errorf("Expected home route path '/', got %s", homeRoute.Path)
		}
		if homeRoute.Method != rtr.MethodGET {
			t.Errorf("Expected home route method %s, got %s", rtr.MethodGET, homeRoute.Method)
		}
		if homeRoute.Handler != "home" {
			t.Errorf("Expected home route handler 'home', got %s", homeRoute.Handler)
		}
		if homeRoute.Status != rtr.StatusEnabled {
			t.Errorf("Expected home route status %s, got %s", rtr.StatusEnabled, homeRoute.Status)
		}
	}

	// Test API group
	apiGroup, ok := domain.Items[1].(*rtr.Group)
	if !ok {
		t.Errorf("Expected second item to be a Group")
	} else {
		if apiGroup.Name != "API v1" {
			t.Errorf("Expected API group name 'API v1', got %s", apiGroup.Name)
		}
		if apiGroup.Prefix != "/api/v1" {
			t.Errorf("Expected API group prefix '/api/v1', got %s", apiGroup.Prefix)
		}
		if apiGroup.Status != rtr.StatusEnabled {
			t.Errorf("Expected API group status %s, got %s", rtr.StatusEnabled, apiGroup.Status)
		}
		if len(apiGroup.Routes) != 3 {
			t.Errorf("Expected 3 routes in API group, got %d", len(apiGroup.Routes))
		}
		if len(apiGroup.Middlewares) != 2 {
			t.Errorf("Expected 2 middlewares in API group, got %d", len(apiGroup.Middlewares))
		}
	}

	// Test admin route (disabled)
	adminRoute, ok := domain.Items[2].(*rtr.Route)
	if !ok {
		t.Errorf("Expected third item to be a Route")
	} else {
		if adminRoute.Name != "Admin Dashboard" {
			t.Errorf("Expected admin route name 'Admin Dashboard', got %s", adminRoute.Name)
		}
		if adminRoute.Status != rtr.StatusDisabled {
			t.Errorf("Expected admin route status %s, got %s", rtr.StatusDisabled, adminRoute.Status)
		}
		if len(adminRoute.Middlewares) != 1 {
			t.Errorf("Expected 1 middleware for admin route, got %d", len(adminRoute.Middlewares))
		}
	}
}

func TestHandlerRegistry(t *testing.T) {
	registry := rtr.NewHandlerRegistry()
	registerHandlers(registry)

	// Test that handlers were registered
	handlerNames := []string{"home", "users-list", "users-create", "users-get", "admin-dashboard"}
	for _, name := range handlerNames {
		handler := registry.FindRoute(name)
		if handler == nil {
			t.Errorf("Expected to find handler '%s', but it was not registered", name)
		}
	}

	// Test specific handler properties
	homeHandler := registry.FindRoute("home")
	if homeHandler != nil {
		if homeHandler.GetName() != "home" {
			t.Errorf("Expected home handler name 'home', got %s", homeHandler.GetName())
		}
		if homeHandler.GetPath() != "/" {
			t.Errorf("Expected home handler path '/', got %s", homeHandler.GetPath())
		}
		if homeHandler.GetMethod() != "GET" {
			t.Errorf("Expected home handler method 'GET', got %s", homeHandler.GetMethod())
		}
	}

	usersListHandler := registry.FindRoute("users-list")
	if usersListHandler != nil {
		if usersListHandler.GetName() != "users-list" {
			t.Errorf("Expected users-list handler name 'users-list', got %s", usersListHandler.GetName())
		}
		if usersListHandler.GetPath() != "/users" {
			t.Errorf("Expected users-list handler path '/users', got %s", usersListHandler.GetPath())
		}
		if usersListHandler.GetMethod() != "GET" {
			t.Errorf("Expected users-list handler method 'GET', got %s", usersListHandler.GetMethod())
		}
	}
}

func TestMiddlewareRegistry(t *testing.T) {
	registry := rtr.NewHandlerRegistry()
	registerMiddleware(registry)

	// Test that middleware was registered
	middlewareNames := []string{"cors", "logging", "auth", "rate-limit", "validate-user", "admin-auth"}
	for _, name := range middlewareNames {
		middleware := registry.FindMiddleware(name)
		if middleware == nil {
			t.Errorf("Expected to find middleware '%s', but it was not registered", name)
		} else {
			if middleware.GetName() != name {
				t.Errorf("Expected middleware name '%s', got %s", name, middleware.GetName())
			}
		}
	}
}

func TestHandlerFunctionality(t *testing.T) {
	registry := rtr.NewHandlerRegistry()
	registerHandlers(registry)

	// Test home handler (HTMLHandler)
	homeHandler := registry.FindRoute("home")
	if homeHandler != nil {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		
		handlerFunc := homeHandler.GetHandler()
		http.HandlerFunc(handlerFunc).ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
		
		contentType := w.Header().Get("Content-Type")
		if !strings.Contains(contentType, "text/html") {
			t.Errorf("Expected HTML content type, got %s", contentType)
		}
		
		body := w.Body.String()
		if !strings.Contains(body, "Welcome to Declarative Router!") {
			t.Errorf("Expected body to contain welcome message, got %s", body)
		}
	}

	// Test users-list handler (JSONHandler)
	usersListHandler := registry.FindRoute("users-list")
	if usersListHandler != nil {
		req := httptest.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()
		
		handlerFunc := usersListHandler.GetHandler()
		http.HandlerFunc(handlerFunc).ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
		
		contentType := w.Header().Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			t.Errorf("Expected JSON content type, got %s", contentType)
		}
		
		body := w.Body.String()
		var result map[string]interface{}
		err := json.Unmarshal([]byte(body), &result)
		if err != nil {
			t.Errorf("Expected valid JSON response, got error: %v", err)
		}
		
		if users, ok := result["users"]; ok {
			if usersList, ok := users.([]interface{}); ok {
				if len(usersList) != 3 {
					t.Errorf("Expected 3 users, got %d", len(usersList))
				}
			} else {
				t.Errorf("Expected users to be an array")
			}
		} else {
			t.Errorf("Expected 'users' field in response")
		}
	}

	// Test users-create handler (standard Handler)
	usersCreateHandler := registry.FindRoute("users-create")
	if usersCreateHandler != nil {
		req := httptest.NewRequest("POST", "/users", nil)
		w := httptest.NewRecorder()
		
		handlerFunc := usersCreateHandler.GetHandler()
		http.HandlerFunc(handlerFunc).ServeHTTP(w, req)
		
		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", w.Code)
		}
		
		contentType := w.Header().Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			t.Errorf("Expected JSON content type, got %s", contentType)
		}
		
		body := w.Body.String()
		if !strings.Contains(body, "User created successfully") {
			t.Errorf("Expected success message in response, got %s", body)
		}
	}
}

func TestMiddlewareFunctionality(t *testing.T) {
	registry := rtr.NewHandlerRegistry()
	registerMiddleware(registry)

	// Test CORS middleware
	corsMiddleware := registry.FindMiddleware("cors")
	if corsMiddleware != nil {
		// Create a simple handler to wrap
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		
		// Wrap with CORS middleware
		wrappedHandler := corsMiddleware.GetHandler()(testHandler)
		
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		
		wrappedHandler.ServeHTTP(w, req)
		
		// Check CORS headers
		if w.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Errorf("Expected CORS origin header to be '*', got %s", w.Header().Get("Access-Control-Allow-Origin"))
		}
		
		if !strings.Contains(w.Header().Get("Access-Control-Allow-Methods"), "GET") {
			t.Errorf("Expected CORS methods to contain 'GET', got %s", w.Header().Get("Access-Control-Allow-Methods"))
		}
	}

	// Test auth middleware
	authMiddleware := registry.FindMiddleware("auth")
	if authMiddleware != nil {
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		
		wrappedHandler := authMiddleware.GetHandler()(testHandler)
		
		// Test without authorization header
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		
		wrappedHandler.ServeHTTP(w, req)
		
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401 without auth, got %d", w.Code)
		}
		
		// Test with authorization header
		req = httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer token")
		w = httptest.NewRecorder()
		
		wrappedHandler.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 with auth, got %d", w.Code)
		}
	}

	// Test rate-limit middleware
	rateLimitMiddleware := registry.FindMiddleware("rate-limit")
	if rateLimitMiddleware != nil {
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		
		wrappedHandler := rateLimitMiddleware.GetHandler()(testHandler)
		
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		
		wrappedHandler.ServeHTTP(w, req)
		
		// Check rate limit headers
		if w.Header().Get("X-RateLimit-Limit") != "100" {
			t.Errorf("Expected rate limit to be '100', got %s", w.Header().Get("X-RateLimit-Limit"))
		}
		
		if w.Header().Get("X-RateLimit-Remaining") != "99" {
			t.Errorf("Expected rate limit remaining to be '99', got %s", w.Header().Get("X-RateLimit-Remaining"))
		}
	}
}

func TestJSONSerialization(t *testing.T) {
	domain := CreateDeclarativeConfiguration()
	
	// Test that domain can be serialized to JSON
	data, err := json.Marshal(domain)
	if err != nil {
		t.Errorf("Expected domain to be serializable, got error: %v", err)
	}
	
	// Test that JSON contains expected fields
	jsonStr := string(data)
	if !strings.Contains(jsonStr, "Example API") {
		t.Errorf("Expected JSON to contain domain name, got %s", jsonStr)
	}
	
	if !strings.Contains(jsonStr, "api.example.com") {
		t.Errorf("Expected JSON to contain host, got %s", jsonStr)
	}
	
	if !strings.Contains(jsonStr, "Home") {
		t.Errorf("Expected JSON to contain route name, got %s", jsonStr)
	}
	
	// Test that JSON can be deserialized back
	var deserializedDomain rtr.Domain
	err = json.Unmarshal(data, &deserializedDomain)
	if err != nil {
		t.Errorf("Expected JSON to be deserializable, got error: %v", err)
	}
	
	// Verify deserialized data
	if deserializedDomain.Name != domain.Name {
		t.Errorf("Expected deserialized name to match original, got %s", deserializedDomain.Name)
	}
	
	if len(deserializedDomain.Hosts) != len(domain.Hosts) {
		t.Errorf("Expected deserialized hosts count to match original, got %d", len(deserializedDomain.Hosts))
	}
}

func TestIntegration(t *testing.T) {
	// Test the complete flow: create config, register handlers/middleware, verify everything works together
	domain := CreateDeclarativeConfiguration()
	registry := rtr.NewHandlerRegistry()
	registerHandlers(registry)
	registerMiddleware(registry)
	
	// Verify that all handlers referenced in the config are registered
	for _, item := range domain.Items {
		switch v := item.(type) {
		case *rtr.Route:
			if v.Status == rtr.StatusEnabled {
				handler := registry.FindRoute(v.Handler)
				if handler == nil {
					t.Errorf("Handler '%s' referenced in config but not found in registry", v.Handler)
				}
			}
		case *rtr.Group:
			if v.Status == rtr.StatusEnabled {
				for _, route := range v.Routes {
					if route.Status == rtr.StatusEnabled {
						handler := registry.FindRoute(route.Handler)
						if handler == nil {
							t.Errorf("Handler '%s' referenced in group route but not found in registry", route.Handler)
						}
					}
				}
			}
		}
	}
	
	// Verify that all middleware referenced in the config are registered
	allMiddleware := make(map[string]bool)
	
	// Collect all middleware names from domain
	for _, mw := range domain.Middlewares {
		allMiddleware[mw] = true
	}
	
	// Collect middleware from items
	for _, item := range domain.Items {
		switch v := item.(type) {
		case *rtr.Route:
			for _, mw := range v.Middlewares {
				allMiddleware[mw] = true
			}
		case *rtr.Group:
			for _, mw := range v.Middlewares {
				allMiddleware[mw] = true
			}
			for _, route := range v.Routes {
				for _, mw := range route.Middlewares {
					allMiddleware[mw] = true
				}
			}
		}
	}
	
	// Verify all middleware are registered
	for mwName := range allMiddleware {
		middleware := registry.FindMiddleware(mwName)
		if middleware == nil {
			t.Errorf("Middleware '%s' referenced in config but not found in registry", mwName)
		}
	}
	
	// Test that disabled routes are properly marked
	adminRoute, ok := domain.Items[2].(*rtr.Route)
	if ok && adminRoute.Status == rtr.StatusDisabled {
		// This is expected - admin route should be disabled
		if adminRoute.Name != "Admin Dashboard" {
			t.Errorf("Expected disabled route to be admin dashboard")
		}
	} else {
		t.Errorf("Expected third item to be a disabled admin route")
	}
}