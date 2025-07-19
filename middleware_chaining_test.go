package rtr_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	rtr "github.com/dracory/rtr"
)

func TestMiddlewareChaining(t *testing.T) {
	tests := []struct {
		name            string
		middlewareOrder []string
		expectedHeaders map[string]string
	}{
		{
			name:            "single middleware",
			middlewareOrder: []string{"first"},
			expectedHeaders: map[string]string{
				"X-First": "true",
			},
		},
		{
			name:            "multiple middlewares in order",
			middlewareOrder: []string{"first", "second", "third"},
			expectedHeaders: map[string]string{
				"X-First":  "true",
				"X-Second": "true",
				"X-Third":  "true",
			},
		},
		{
			name:            "reverse order middlewares",
			middlewareOrder: []string{"third", "second", "first"},
			expectedHeaders: map[string]string{
				"X-First":  "true",
				"X-Second": "true",
				"X-Third":  "true",
			},
		},
	}

	// Define test middlewares
	middlewares := map[string]rtr.StdMiddleware{
		"first": func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-First", "true")
				next.ServeHTTP(w, r)
			})
		},
		"second": func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Second", "true")
				next.ServeHTTP(w, r)
			})
		},
		"third": func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Third", "true")
				next.ServeHTTP(w, r)
			})
		},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := rtr.NewRouter()

			// Create a route with the specified middlewares
			route := rtr.NewRoute().
				SetMethod(http.MethodGet).
				SetPath("/test").
				SetHandler(handler)

			// Add middlewares in the specified order
			for _, mwName := range tc.middlewareOrder {
				mw := middlewares[mwName]
				route.AddBeforeMiddlewares(rtr.MiddlewaresToInterfaces([]rtr.StdMiddleware{mw}))
			}

			r.AddRoute(route)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)

			// Check that all expected headers are set
			for key, value := range tc.expectedHeaders {
				if got := rr.Header().Get(key); got != value {
					t.Errorf("expected header %s=%s, got %s", key, value, got)
				}
			}

			// Verify no unexpected headers were set
			expectedHeaderCount := len(tc.expectedHeaders)
			if len(rr.Header()) != expectedHeaderCount {
				t.Errorf("expected %d headers, got %d", expectedHeaderCount, len(rr.Header()))
			}
		})
	}
}

func TestMiddlewareAbort(t *testing.T) {
	r := rtr.NewRouter()

	// Track which middlewares were called
	var middlewareCalls []string

	// Middleware that aborts the request chain
	abortMiddleware := func(name string) rtr.StdMiddleware {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				middlewareCalls = append(middlewareCalls, name)
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("access denied"))
				// Don't call next.ServeHTTP to abort the chain
			})
		}
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}

	r.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/secure").
		AddBeforeMiddlewares(rtr.MiddlewaresToInterfaces([]rtr.StdMiddleware{
			abortMiddleware("abort"),
		})).
		SetHandler(handler))

	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	// Check that the abort middleware was called
	if len(middlewareCalls) == 0 || middlewareCalls[0] != "abort" {
		t.Fatalf("expected abort middleware to be called first, got calls: %v", middlewareCalls)
	}

	// Check response
	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, rr.Code)
	}

	if rr.Body.String() != "access denied" {
		t.Errorf("unexpected response body: %s", rr.Body.String())
	}
}

func TestAfterMiddleware(t *testing.T) {
	r := rtr.NewRouter()

	// Before middleware that adds a header
	beforeMW := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Before", "true")
			next.ServeHTTP(w, r)
		})
	}

	// After middleware that modifies the response
	afterMW := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Call the next handler first
			next.ServeHTTP(w, r)
			// Then modify the response
			w.Header().Set("X-After", "true")
			// This should not affect the status code already set by the handler
			w.WriteHeader(http.StatusTeapot) // This should be ignored
		})
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	}

	r.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/test").
		AddBeforeMiddlewares(rtr.MiddlewaresToInterfaces([]rtr.StdMiddleware{beforeMW})).
		AddAfterMiddlewares(rtr.MiddlewaresToInterfaces([]rtr.StdMiddleware{afterMW})).
		SetHandler(handler))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	// Check status code (should be from handler, not after middleware)
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	// Check that both before and after middleware headers are set
	headers := map[string]string{
		"X-Before": "true",
		"X-After":  "true",
	}

	for key, value := range headers {
		if got := rr.Header().Get(key); got != value {
			t.Errorf("expected header %s=%s, got %s", key, value, got)
		}
	}
}
