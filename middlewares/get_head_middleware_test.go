package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/rtr"
	"github.com/dracory/rtr/middlewares"
)

func TestGetHeadMiddleware(t *testing.T) {
	// Test case 1: HEAD request should return headers and status code from GET handler
	t.Run("HEAD request to GET handler", func(t *testing.T) {
		// Create a test handler that sets some headers and returns a body
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Test", "test-value")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("response body"))
		})

		req := httptest.NewRequest("HEAD", "/test", nil)
		rr := httptest.NewRecorder()

		// Create and execute middleware
		mw := middlewares.GetHead()
		handler = mw.Execute(handler).(http.HandlerFunc)
		handler.ServeHTTP(rr, req)

		// Verify status code is the same as GET handler
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Verify headers are copied from GET handler
		expectedHeaders := map[string]string{
			"Content-Type": "application/json",
			"X-Test":       "test-value",
		}

		for k, v := range expectedHeaders {
			if got := rr.Header().Get(k); got != v {
				t.Errorf("header %s: got %v want %v", k, got, v)
			}
		}

		// Verify no body is returned for HEAD request
		if body := rr.Body.String(); body != "" {
			t.Errorf("handler returned unexpected body: got %v want empty", body)
		}
	})

	// Test case 2: GET request should pass through unchanged
	t.Run("GET request passes through", func(t *testing.T) {
		var called bool
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			if r.Method != "GET" {
				t.Errorf("expected method GET, got %s", r.Method)
			}
		})

		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()

		mw := middlewares.GetHead()
		handler = mw.Execute(handler).(http.HandlerFunc)
		handler.ServeHTTP(rr, req)

		if !called {
			t.Error("handler was not called")
		}
	})

	// Test case 3: POST request should pass through unchanged
	t.Run("POST request passes through", func(t *testing.T) {
		var called bool
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			if r.Method != "POST" {
				t.Errorf("expected method POST, got %s", r.Method)
			}
		})

		req := httptest.NewRequest("POST", "/test", nil)
		rr := httptest.NewRecorder()

		mw := middlewares.GetHead()
		handler = mw.Execute(handler).(http.HandlerFunc)
		handler.ServeHTTP(rr, req)

		if !called {
			t.Error("handler was not called")
		}
	})

	// Test case 4: HEAD request to non-existent route should return 404
	t.Run("HEAD request to non-existent route", func(t *testing.T) {
		// Create a router with a GET handler for /exists but not for /not-exists
		router := rtr.NewRouter()
		router.AddRoute(rtr.NewRoute().
			SetMethod("GET").
			SetPath("/exists").
			SetHandler(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

		req := httptest.NewRequest("HEAD", "/not-exists", nil)
		rr := httptest.NewRecorder()

		// Wrap the router with the middleware
		handler := middlewares.GetHead().Execute(router)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})

	// Test case 5: HEAD request should work with middleware chain
	t.Run("HEAD request with middleware chain", func(t *testing.T) {
		var middlewareCalled bool
		var handlerCalled bool

		// Create a test handler
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
		})

		// Create a test middleware
		testMiddleware := rtr.NewMiddleware().
			SetName("Test Middleware").
			SetHandler(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					middlewareCalled = true
					next.ServeHTTP(w, r)
				})
			})

		req := httptest.NewRequest("HEAD", "/test", nil)
		rr := httptest.NewRecorder()

		// Create middleware chain: GetHead -> Test Middleware -> Handler
		chain := testMiddleware.Execute(handler)
		chain = middlewares.GetHead().Execute(chain)
		chain.ServeHTTP(rr, req)

		if !middlewareCalled {
			t.Error("middleware was not called")
		}
		if !handlerCalled {
			t.Error("handler was not called")
		}
	})
}
