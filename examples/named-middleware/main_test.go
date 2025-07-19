package main_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dracory/rtr"
)

func TestNamedMiddlewareEndpoints(t *testing.T) {
	// Create a new router instance with test routes
	r := setupTestRouter()

	tests := []struct {
		name           string
		method         string
		path           string
		headers        map[string]string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "GET /public - no auth required",
			method:         http.MethodGet,
			path:           "/public",
			expectedStatus: http.StatusOK,
			expectedBody:   "This is a public endpoint",
		},
		{
			name:           "GET /protected - no auth header",
			method:         http.MethodGet,
			path:           "/protected",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Unauthorized",
		},
		{
			name:           "GET /protected - with auth header",
			method:         http.MethodGet,
			path:           "/protected",
			headers:        map[string]string{"Authorization": "Bearer token123"},
			expectedStatus: http.StatusOK,
			expectedBody:   "This is a protected endpoint",
		},
		{
			name:           "GET /admin - no auth header",
			method:         http.MethodGet,
			path:           "/admin",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Unauthorized",
		},
		{
			name:           "GET /admin - with auth header",
			method:         http.MethodGet,
			path:           "/admin",
			headers:        map[string]string{"Authorization": "Bearer token123"},
			expectedStatus: http.StatusOK,
			expectedBody:   "This is an admin endpoint",
		},
		{
			name:           "GET /api/users - with logging middleware",
			method:         http.MethodGet,
			path:           "/api/users",
			expectedStatus: http.StatusOK,
			expectedBody:   "Users API endpoint",
		},
		{
			name:           "GET /v1/status - with rate limiting",
			method:         http.MethodGet,
			path:           "/v1/status",
			expectedStatus: http.StatusOK,
			expectedBody:   "API Status: OK",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			
			// Add headers if specified
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}
			
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			body := strings.TrimSpace(w.Body.String())
			if !strings.Contains(body, tt.expectedBody) {
				t.Errorf("Expected body to contain %q, got %q", tt.expectedBody, body)
			}
		})
	}
}

func TestNamedMiddlewareExecution(t *testing.T) {
	// Test that named middleware executes in the correct order
	var executionOrder []string
	
	// Create middleware that records execution order
	middleware1 := rtr.NewMiddleware("First Middleware", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "First")
			next.ServeHTTP(w, r)
		})
	})
	
	middleware2 := rtr.NewMiddleware("Second Middleware", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "Second")
			next.ServeHTTP(w, r)
		})
	})

	// Create a simple route with named middleware
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/test").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "Handler")
			w.WriteHeader(http.StatusOK)
		}).
		AddBeforeMiddlewares([]rtr.MiddlewareInterface{middleware1, middleware2})

	router := rtr.NewRouter().AddRoute(route)

	// Make a request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check execution order
	expectedOrder := []string{"First", "Second", "Handler"}
	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("Expected %d executions, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, expected := range expectedOrder {
		if executionOrder[i] != expected {
			t.Errorf("Expected execution order[%d] to be %q, got %q", i, expected, executionOrder[i])
		}
	}
}

func setupTestRouter() rtr.RouterInterface {
	// Create named middleware
	authMiddleware := rtr.NewMiddleware("User Authentication", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	authzMiddleware := rtr.NewMiddleware("Admin Authorization", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Simple authorization check (in real app, check user roles)
			next.ServeHTTP(w, r)
		})
	})

	loggingMiddleware := rtr.NewMiddleware("Request Logging", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// In a real app, this would log to a proper logger
			next.ServeHTTP(w, r)
		})
	})

	rateLimitMiddleware := rtr.NewMiddleware("Rate Limiting", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Simple rate limiting simulation
			next.ServeHTTP(w, r)
		})
	})

	// Create router
	router := rtr.NewRouter()

	// Add routes with named middleware
	router.AddRoute(rtr.NewRoute().
		SetMethod("GET").
		SetPath("/public").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("This is a public endpoint"))
		}))

	router.AddRoute(rtr.NewRoute().
		SetMethod("GET").
		SetPath("/protected").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("This is a protected endpoint"))
		}).
		AddBeforeMiddlewares([]rtr.MiddlewareInterface{authMiddleware}))

	router.AddRoute(rtr.NewRoute().
		SetMethod("GET").
		SetPath("/admin").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("This is an admin endpoint"))
		}).
		AddBeforeMiddlewares([]rtr.MiddlewareInterface{authMiddleware, authzMiddleware}))

	// Create API group with logging middleware
	apiGroup := rtr.NewGroup().
		SetPrefix("/api").
		AddBeforeMiddlewares([]rtr.MiddlewareInterface{loggingMiddleware})

	apiGroup.AddRoute(rtr.NewRoute().
		SetMethod("GET").
		SetPath("/users").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Users API endpoint"))
		}))

	router.AddGroup(apiGroup)

	// Create versioned group with rate limiting
	v1Group := rtr.NewGroup().
		SetPrefix("/v1").
		AddBeforeMiddlewares([]rtr.MiddlewareInterface{rateLimitMiddleware})

	v1Group.AddRoute(rtr.NewRoute().
		SetMethod("GET").
		SetPath("/status").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("API Status: OK"))
		}))

	router.AddGroup(v1Group)

	return router
}
