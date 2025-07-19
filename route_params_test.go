package rtr_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/dracory/rtr"
)

// panicValue executes a function and returns the recovered value if it panics
func panicValue(fn func()) (val interface{}) {
	defer func() { val = recover() }()
	fn()
	return nil
}

func TestPathParameters(t *testing.T) {
	tests := []struct {
		name           string
		routePath      string
		requestPath    string
		expectedMatch  bool
		expectedParams map[string]string
		handler        http.HandlerFunc
	}{
		{
			name:          "simple parameter",
			routePath:     "/users/:id",
			requestPath:   "/users/123",
			expectedMatch: true,
			expectedParams: map[string]string{
				"id": "123",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				id, exists := rtr.GetParam(r, "id")
				if !exists || id != "123" {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
			},
		},
		{
			name:          "simple parameter",
			routePath:     "/users/:id",
			requestPath:   "/users/123",
			expectedMatch: true,
			expectedParams: map[string]string{
				"id": "123",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				id, exists := rtr.GetParam(r, "id")
				if !exists || id != "123" {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new router
			r := rtr.NewRouter()

			// Create a test handler that captures the request
			var capturedParams map[string]string
			var handler rtr.StdHandler = func(w http.ResponseWriter, r *http.Request) {
				if tc.handler != nil {
					tc.handler(w, r)
				}
				capturedParams = rtr.GetParams(r)
			}

			// Add the route
			r.AddRoute(rtr.NewRoute().
				SetPath(tc.routePath).
				SetMethod(http.MethodGet).
				SetHandler(handler))

			// Create a test request
			req := httptest.NewRequest(http.MethodGet, tc.requestPath, nil)
			w := httptest.NewRecorder()

			// Serve the request
			r.ServeHTTP(w, req)

			// Check if the route matched as expected
			if tc.expectedMatch {
				if w.Code != http.StatusOK {
					t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
				}
				if !reflect.DeepEqual(tc.expectedParams, capturedParams) {
					t.Fatalf("Expected params %v, got %v", tc.expectedParams, capturedParams)
				}
			} else {
				if w.Code != http.StatusNotFound {
					t.Fatalf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
				}
			}
		})
	}
}

func TestMustGetParam(t *testing.T) {
	t.Run("parameter exists", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
		params := map[string]string{"id": "123"}
		req = req.WithContext(context.WithValue(req.Context(), rtr.ParamsKey, params))

		id := rtr.MustGetParam(req, "id")
		if id != "123" {
			t.Fatalf("Expected id '123', got '%s'", id)
		}
	})

	t.Run("parameter does not exist", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
		params := map[string]string{}
		req = req.WithContext(context.WithValue(req.Context(), rtr.ParamsKey, params))

		didPanic := panicValue(func() {
			_ = rtr.MustGetParam(req, "id")
		})
		if didPanic == nil {
			t.Fatal("Expected MustGetParam to panic when parameter doesn't exist")
		}
	})
}

func TestGetParams(t *testing.T) {
	t.Run("no parameters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		params := rtr.GetParams(req)
		if len(params) != 0 {
			t.Fatalf("Expected empty params, got %v", params)
		}
	})

	t.Run("with parameters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/123/orders/456", nil)
		expected := map[string]string{
			"userId":  "123",
			"orderId": "456",
		}
		req = req.WithContext(context.WithValue(req.Context(), rtr.ParamsKey, expected))

		params := rtr.GetParams(req)
		if !reflect.DeepEqual(expected, params) {
			t.Fatalf("Expected params %v, got %v", expected, params)
		}

		// Ensure the returned map is a copy
		params["userId"] = "modified"
		if expected["userId"] != "123" {
			t.Fatal("Original params should not be modified")
		}
	})
}
