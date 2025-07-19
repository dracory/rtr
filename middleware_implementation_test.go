package rtr

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestMiddleware tests the basic functionality of  middleware
func TestMiddleware(t *testing.T) {
	// Create a  middleware
	authMiddleware := NewMiddleware("Authentication Check", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add a header to indicate this middleware ran
			w.Header().Set("X-Auth-Middleware", "executed")
			next.ServeHTTP(w, r)
		})
	})

	// Test middleware properties
	if authMiddleware.GetName() != "Authentication Check" {
		t.Errorf("Expected middleware name 'Authentication Check', got '%s'", authMiddleware.GetName())
	}

	if authMiddleware.GetHandler() == nil {
		t.Error("Expected middleware handler to be set")
	}

	// Test middleware execution
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	wrappedHandler := authMiddleware.Execute(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(w, req)

	if w.Header().Get("X-Auth-Middleware") != "executed" {
		t.Error("Expected middleware to set X-Auth-Middleware header")
	}

	if w.Body.String() != "test response" {
		t.Errorf("Expected 'test response', got '%s'", w.Body.String())
	}
}

// TestAnonymousMiddleware tests middleware without names
func TestAnonymousMiddleware(t *testing.T) {
	middleware := NewAnonymousMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Anonymous", "true")
			next.ServeHTTP(w, r)
		})
	})

	if middleware.GetName() != "" {
		t.Errorf("Expected empty name for anonymous middleware, got '%s'", middleware.GetName())
	}
}

// TestMiddlewareConversion tests conversion between middleware types
func TestMiddlewareConversion(t *testing.T) {
	// Create regular middleware functions
	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Middleware-1", "executed")
			next.ServeHTTP(w, r)
		})
	}

	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Middleware-2", "executed")
			next.ServeHTTP(w, r)
		})
	}

	// Convert to MiddlewareInterface slice
	middlewares := []StdMiddleware{middleware1, middleware2}
	interfaceMiddlewares := MiddlewaresToInterfaces(middlewares)

	if len(interfaceMiddlewares) != 2 {
		t.Errorf("Expected 2 middleware interfaces, got %d", len(interfaceMiddlewares))
	}

	// Convert back to Middleware slice
	backToMiddlewares := InterfacesToMiddlewares(interfaceMiddlewares)

	if len(backToMiddlewares) != 2 {
		t.Errorf("Expected 2 middleware functions, got %d", len(backToMiddlewares))
	}
}

// TestRouteWithMiddleware tests adding  middleware to routes
func TestRouteWithMiddleware(t *testing.T) {
	// Create  middleware
	logMiddleware := NewMiddleware("Request Logger", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Logged", "true")
			next.ServeHTTP(w, r)
		})
	})

	authMiddleware := NewMiddleware("Auth Check", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Authenticated", "true")
			next.ServeHTTP(w, r)
		})
	})

	// Create a route with  middleware
	route := NewRoute().
		SetMethod("GET").
		SetPath("/protected").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("protected content"))
		}).
		AddBeforeMiddlewares([]MiddlewareInterface{logMiddleware, authMiddleware})

	// Test that middleware was added
	Middlewares := route.GetBeforeMiddlewares()
	if len(Middlewares) != 2 {
		t.Errorf("Expected 2  middlewares, got %d", len(Middlewares))
	}

	if Middlewares[0].GetName() != "Request Logger" {
		t.Errorf("Expected first middleware name 'Request Logger', got '%s'", Middlewares[0].GetName())
	}

	if Middlewares[1].GetName() != "Auth Check" {
		t.Errorf("Expected second middleware name 'Auth Check', got '%s'", Middlewares[1].GetName())
	}
}

// TestExecuteMiddlewareChain tests the middleware chain execution
func TestExecuteMiddlewareChain(t *testing.T) {
	var executionOrder []string

	// Create middleware that tracks execution order
	middleware1 := NewMiddleware("First", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "First")
			next.ServeHTTP(w, r)
		})
	})

	middleware2 := NewMiddleware("Second", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "Second")
			next.ServeHTTP(w, r)
		})
	})

	middleware3 := NewMiddleware("Third", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "Third")
			next.ServeHTTP(w, r)
		})
	})

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		executionOrder = append(executionOrder, "Handler")
		w.WriteHeader(http.StatusOK)
	})

	middlewares := []MiddlewareInterface{middleware1, middleware2, middleware3}
	chainedHandler := ExecuteMiddlewareChain(middlewares, finalHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	chainedHandler.ServeHTTP(w, req)

	expectedOrder := []string{"First", "Second", "Third", "Handler"}
	if len(executionOrder) != len(expectedOrder) {
		t.Errorf("Expected %d execution steps, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, expected := range expectedOrder {
		if i >= len(executionOrder) || executionOrder[i] != expected {
			t.Errorf("Expected execution order %v, got %v", expectedOrder, executionOrder)
			break
		}
	}
}

// TestMiddlewareSetName tests changing middleware names
func TestMiddlewareSetName(t *testing.T) {
	middleware := NewMiddleware("Original Name", func(next http.Handler) http.Handler {
		return next
	})

	if middleware.GetName() != "Original Name" {
		t.Errorf("Expected 'Original Name', got '%s'", middleware.GetName())
	}

	// Change the name
	middleware.SetName("New Name")

	if middleware.GetName() != "New Name" {
		t.Errorf("Expected 'New Name', got '%s'", middleware.GetName())
	}
}

// BenchmarkMiddleware benchmarks  middleware performance
func BenchmarkMiddleware(b *testing.B) {
	middleware := NewMiddleware("Benchmark Middleware", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrappedHandler := middleware.Execute(handler)
	req := httptest.NewRequest("GET", "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(w, req)
	}
}

// BenchmarkMiddlewareChain benchmarks middleware chain execution
func BenchmarkMiddlewareChain(b *testing.B) {
	middlewares := make([]MiddlewareInterface, 5)
	for i := 0; i < 5; i++ {
		middlewares[i] = NewMiddleware("Middleware", func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		})
	}

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	chainedHandler := ExecuteMiddlewareChain(middlewares, finalHandler)
	req := httptest.NewRequest("GET", "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		chainedHandler.ServeHTTP(w, req)
	}
}
