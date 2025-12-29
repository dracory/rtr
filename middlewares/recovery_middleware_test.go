package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/rtr"
	"github.com/dracory/rtr/middlewares"
)

func TestRecoveryMiddleware(t *testing.T) {
	t.Run("recovers from panic", func(t *testing.T) {
		// Create a test handler that panics
		panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("test panic")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()

		// Create and execute middleware
		recoveryMw := middlewares.RecoveryMiddleware()
		handler := recoveryMw.Execute(panicHandler)

		// This should not panic
		handler.ServeHTTP(rr, req)

		// Check the response
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, status)
		}

		expectedBody := "Internal Server Error\n"
		if rr.Body.String() != expectedBody {
			t.Errorf("expected body %q, got %q", expectedBody, rr.Body.String())
		}

		// Check middleware name
		if recoveryMw.GetName() != "Recovery Middleware" {
			t.Errorf("expected name %q, got %q", "Recovery Middleware", recoveryMw.GetName())
		}

		// Check handler is set
		if recoveryMw.GetHandler() == nil {
			t.Error("expected handler to be set")
		}
	})

	t.Run("with nil next handler", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()

		recoveryMw := middlewares.RecoveryMiddleware()
		handler := recoveryMw.Execute(nil) // Pass nil handler

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, status)
		}
	})

	t.Run("with nil request and response", func(t *testing.T) {
		recoveryMw := middlewares.RecoveryMiddleware()
		handler := recoveryMw.Execute(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("should not be called with nil request/response")
		}))

		// This should not panic
		handler.ServeHTTP(nil, nil)
	})
}

func TestRecoveryMiddleware_EdgeCases(t *testing.T) {
	t.Run("with already written response", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()

		// Write to response before panic
		panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Partial"))
			panic("test panic after write")
		})

		recoveryMw := middlewares.RecoveryMiddleware()
		handler := recoveryMw.Execute(panicHandler)

		handler.ServeHTTP(rr, req)

		// Should not override the existing response
		if rr.Body.String() != "Partial" {
			t.Errorf("expected body to be \"Partial\", got %q", rr.Body.String())
		}
	})

	t.Run("with custom handler", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()

		// Create a custom recovery handler
		customHandler := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					if r := recover(); r != nil {
						w.WriteHeader(http.StatusBadGateway)
						_, _ = w.Write([]byte("Custom Error"))
					}
				}()
				next.ServeHTTP(w, r)
			})
		}

		// Create middleware with custom handler
		recoveryMw := rtr.NewMiddleware().
			SetName("Custom Recovery").
			SetHandler(customHandler)

		// Test handler that panics
		panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("test panic")
		})

		handler := recoveryMw.Execute(panicHandler)
		handler.ServeHTTP(rr, req)

		// Should use custom error handling
		if status := rr.Code; status != http.StatusBadGateway {
			t.Errorf("expected status %d, got %d", http.StatusBadGateway, status)
		}

		expectedBody := "Custom Error"
		if rr.Body.String() != expectedBody {
			t.Errorf("expected body %q, got %q", expectedBody, rr.Body.String())
		}
	})
}
