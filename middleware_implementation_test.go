package rtr

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestMiddleware tests the basic functionality of middleware
func TestMiddleware(t *testing.T) {
	// Test creating a new middleware with name and handler
	handlerFunc := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test-Middleware", "test")
			next.ServeHTTP(w, r)
		})
	}

	mw := NewMiddleware(
		WithName("Test Middleware"),
		WithHandler(handlerFunc),
	)

	// Verify the middleware was created correctly
	if mw.GetName() != "Test Middleware" {
		t.Errorf("Expected middleware name to be 'Test Middleware', got '%s'", mw.GetName())
	}

	// Test the middleware execution
	handler := mw.Execute(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test"))
	}))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Verify the middleware was executed
	if w.Header().Get("X-Test-Middleware") != "test" {
		t.Error("Expected X-Test-Middleware header to be set")
	}
}

// TestMiddlewareSetName tests setting middleware name
func TestMiddlewareSetName(t *testing.T) {
	// Create middleware with initial name
	mw := NewMiddleware(
		WithName("Original Name"),
		WithHandler(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		}),
	)

	if mw.GetName() != "Original Name" {
		t.Errorf("Expected 'Original Name', got '%s'", mw.GetName())
	}

	// Change the name
	mw.SetName("New Name")

	if mw.GetName() != "New Name" {
		t.Errorf("Expected 'New Name', got '%s'", mw.GetName())
	}
}

// TestAnonymousMiddleware tests middleware without names
func TestAnonymousMiddleware(t *testing.T) {
	// Create anonymous middleware
	anonMiddleware := NewAnonymousMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Anon-Middleware", "executed")
			next.ServeHTTP(w, r)
		})
	})

	if anonMiddleware.GetName() != "" {
		t.Errorf("Expected empty name for anonymous middleware, got '%s'", anonMiddleware.GetName())
	}
}

// TestMiddlewareConversion tests conversion between middleware types
func TestMiddlewareConversion(t *testing.T) {
	// Create a middleware with a name and handler
	handlerFunc := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Middleware", "1")
			next.ServeHTTP(w, r)
		})
	}

	mw1 := NewMiddleware()
	mw1.SetName("First")
	mw1.SetHandler(handlerFunc)

	// Convert to StdMiddleware and back
	stdMW := mw1.GetHandler()
	mw2 := NewAnonymousMiddleware(stdMW)

	// Verify the conversion
	if mw2.GetName() != "" {
		t.Error("Expected anonymous middleware to have empty name")
	}

	// Test the converted middleware
	handler := mw2.Execute(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test"))
	}))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Header().Get("X-Middleware") != "1" {
		t.Error("Expected X-Middleware header to be set")
	}
}

// TestMiddlewaresConversion tests conversion between middleware slices
func TestMiddlewaresConversion(t *testing.T) {
	// Create a slice to collect middleware execution order
	var executionOrder []string

	handler1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "1")
			next.ServeHTTP(w, r)
		})
	}

	handler2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "2")
			next.ServeHTTP(w, r)
		})
	}

	mw1 := NewMiddleware()
	mw1.SetName("First")
	mw1.SetHandler(handler1)

	mw2 := NewMiddleware()
	mw2.SetName("Second")
	mw2.SetHandler(handler2)

	// Test conversion from MiddlewareInterface to StdMiddleware
	middlewareFuncs := InterfacesToMiddlewares([]MiddlewareInterface{mw1, mw2})
	if len(middlewareFuncs) != 2 {
		t.Fatalf("Expected 2 middleware functions, got %d", len(middlewareFuncs))
	}

	// Test execution of converted middlewares
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test"))
	})

	// Apply middlewares in reverse order to test chaining
	for i := len(middlewareFuncs) - 1; i >= 0; i-- {
		handler = http.HandlerFunc(middlewareFuncs[i](handler).ServeHTTP)
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Reset execution order before serving
	executionOrder = nil
	handler.ServeHTTP(w, req)

	// Check execution order instead of headers
	if len(executionOrder) != 2 {
		t.Fatalf("Expected 2 middleware executions, got %d", len(executionOrder))
	}

	// Middleware should execute in the order they were defined (1 then 2)
	// because the second middleware wraps the first
	if executionOrder[0] != "1" || executionOrder[1] != "2" {
		t.Errorf("Expected execution order [1, 2], got %v", executionOrder)
	}

	// Also verify the response is OK
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

// TestMiddlewareChain tests middleware chaining
func TestMiddlewareChain(t *testing.T) {
	// Create a chain of middlewares
	handler1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headers := w.Header()
			headers["X-Middleware"] = append(headers["X-Middleware"], "1")
			next.ServeHTTP(w, r)
		})
	}

	handler2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headers := w.Header()
			headers["X-Middleware"] = append(headers["X-Middleware"], "2")
			next.ServeHTTP(w, r)
		})
	}

	handler3 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headers := w.Header()
			headers["X-Middleware"] = append(headers["X-Middleware"], "3")
			next.ServeHTTP(w, r)
		})
	}

	mw1 := NewMiddleware()
	mw1.SetName("First")
	mw1.SetHandler(handler1)

	mw2 := NewMiddleware()
	mw2.SetName("Second")
	mw2.SetHandler(handler2)

	mw3 := NewMiddleware()
	mw3.SetName("Third")
	mw3.SetHandler(handler3)

	// Create a chain: mw1 -> mw2 -> mw3 -> final handler
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test"))
	})

	// Chain the middlewares in reverse order
	handler := mw1.Execute(mw2.Execute(mw3.Execute(finalHandler)))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Check execution order
	headers := w.Header()["X-Middleware"]
	expectedOrder := []string{"1", "2", "3"}

	if len(headers) != len(expectedOrder) {
		t.Fatalf("Expected %d middleware executions, got %d", len(expectedOrder), len(headers))
	}

	for i, v := range expectedOrder {
		if headers[i] != v {
			t.Errorf("Expected header[%d] to be '%s', got '%s'", i, v, headers[i])
		}
	}
}

// BenchmarkMiddleware benchmarks middleware performance
func BenchmarkMiddleware(b *testing.B) {
	// Create a simple middleware
	mwHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}

	mw := NewMiddleware()
	mw.SetName("Benchmark Middleware")
	mw.SetHandler(mwHandler)

	handler := mw.Execute(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(w, req)
		_ = w.Result().Body.Close()
	}
}

// BenchmarkMiddlewareChain benchmarks middleware chain execution
func BenchmarkMiddlewareChain(b *testing.B) {
	middlewares := make([]MiddlewareInterface, 5)
	for i := 0; i < 5; i++ {
		middlewares[i] = NewMiddleware(
			WithName("Middleware"),
			WithHandler(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					next.ServeHTTP(w, r)
				})
			}),
		)
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
