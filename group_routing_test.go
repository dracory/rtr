package rtr_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/rtr"
)

func TestGroupRouting(t *testing.T) {
	r := rtr.NewRouter()

	// Create a group with a prefix and middleware
	group := rtr.NewGroup().
		SetPrefix("/api/v1").
		AddBeforeMiddlewares(rtr.MiddlewaresToInterfaces([]rtr.StdMiddleware{
			func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("X-API-Version", "v1")
					next.ServeHTTP(w, r)
				})
			},
		}))

	// Add routes to the group
	group.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/users").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("users list"))
		}))

	group.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/users/:id").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			// Use the helper function to get the parameter
			id, exists := rtr.GetParam(r, "id")
			if !exists {
				http.Error(w, "ID parameter not found", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("user " + id))
		}))

	r.AddGroup(group)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
		expectedHeader string
	}{
		{
			name:           "group route with prefix",
			method:         http.MethodGet,
			path:           "/api/v1/users",
			expectedStatus: http.StatusOK,
			expectedBody:   "users list",
			expectedHeader: "v1",
		},
		{
			name:           "group route with params",
			method:         http.MethodGet,
			path:           "/api/v1/users/123",
			expectedStatus: http.StatusOK,
			expectedBody:   "user 123",
			expectedHeader: "v1",
		},
		{
			name:           "non-existent route",
			method:         http.MethodGet,
			path:           "/api/v1/nonexistent",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404 page not found\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

			if tc.expectedBody != "" && rr.Body.String() != tc.expectedBody {
				t.Errorf("expected body %q, got %q", tc.expectedBody, rr.Body.String())
			}

			if tc.expectedHeader != "" {
				if got := rr.Header().Get("X-API-Version"); got != tc.expectedHeader {
					t.Errorf("expected header X-API-Version=%s, got %s", tc.expectedHeader, got)
				}
			}
		})
	}

}

func TestNestedGroupRouting(t *testing.T) {
	r := rtr.NewRouter()

	// Track middleware execution
	var parentMiddlewareCalled, childMiddlewareCalled bool
	var parentCtx, childCtx context.Context

	// Create a parent group
	parentGroup := rtr.NewGroup().
		SetPrefix("/parent").
		AddBeforeMiddlewares(rtr.MiddlewaresToInterfaces([]rtr.StdMiddleware{
			func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					parentMiddlewareCalled = true
					parentCtx = r.Context()
					w.Header().Set("X-Parent-Middleware", "executed")
					next.ServeHTTP(w, r)
				})
			},
		}))

	// Create a child group
	childGroup := rtr.NewGroup().
		SetPrefix("/child").
		AddBeforeMiddlewares(rtr.MiddlewaresToInterfaces([]rtr.StdMiddleware{
			func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					childMiddlewareCalled = true
					childCtx = r.Context()
					w.Header().Set("X-Child-Middleware", "executed")
					next.ServeHTTP(w, r)
				})
			},
		}))

	// Add a route to the child group
	childGroup.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/test").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("nested group test"))
		}))

	// Add child group to parent group
	parentGroup.AddGroup(childGroup)

	// Add parent group to router
	r.AddGroup(parentGroup)

	// Test the nested route
	req := httptest.NewRequest(http.MethodGet, "/parent/child/test", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if rr.Body.String() != "nested group test" {
		t.Errorf("unexpected response body: %s", rr.Body.String())
	}

	// Check response
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if rr.Body.String() != "nested group test" {
		t.Errorf("unexpected response body: %s", rr.Body.String())
	}

	// Check that both middlewares were executed
	if !parentMiddlewareCalled {
		t.Error("parent middleware was not called")
	} else if parentCtx == nil {
		t.Error("parent middleware did not set context")
	}

	if !childMiddlewareCalled {
		t.Error("child middleware was not called")
	} else if childCtx == nil {
		t.Error("child middleware did not set context")
	}

	// Also check headers as a secondary verification
	if got := rr.Header().Get("X-Parent-Middleware"); got != "executed" {
		t.Error("parent middleware did not set header")
	}

	if got := rr.Header().Get("X-Child-Middleware"); got != "executed" {
		t.Error("child middleware did not set header")
	}
}
