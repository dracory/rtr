package rtr_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dracory/rtr"
)

// !!! Important !!!
// Expected order: globals before → domains before → groups before → routes before → handler → routes after → groups after → domains after → globals after
// For routes at the same level as groups, they should only get global and domain middlewares

// testMiddlewareSetup creates a new router with test middlewares configured
// Returns the router, domain, and a trace middleware factory function
func testMiddlewareSetup(t *testing.T) (rtr.RouterInterface, rtr.DomainInterface, func(string) rtr.MiddlewareInterface, func(string) rtr.MiddlewareInterface) {
	// Create middleware factory functions
	traceBeforeMiddleware := func(name string) rtr.MiddlewareInterface {
		return rtr.NewAnonymousMiddleware(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Get or initialize the execution order slice
				var executionOrder *[]string
				if val := r.Context().Value(rtr.ExecutionSequenceKey); val != nil {
					executionOrder = val.(*[]string)
				} else {
					executionOrder = &[]string{}
					r = r.WithContext(context.WithValue(r.Context(), rtr.ExecutionSequenceKey, executionOrder))
				}

				// Log entry point
				// t.Logf("ENTERING BEFORE MIDDLEWARE: %s", name)
				// *executionOrder = append(*executionOrder, name+"_enter")

				// Log execution point (before calling next)
				t.Logf("EXECUTING BEFORE MIDDLEWARE: %s", name)
				*executionOrder = append(*executionOrder, name+"_execute")

				// Create a response recorder to capture the response
				rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

				// Call the next handler
				next.ServeHTTP(rw, r)

				// Log exit point
				// t.Logf("EXITING BEFORE MIDDLEWARE: %s", name)
				// *executionOrder = append(*executionOrder, name+"_exit")
			})
		})
	}

	traceAfterMiddleware := func(name string) rtr.MiddlewareInterface {
		return rtr.NewAnonymousMiddleware(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Get or initialize the execution order slice
				var executionOrder *[]string
				if val := r.Context().Value(rtr.ExecutionSequenceKey); val != nil {
					executionOrder = val.(*[]string)
				} else {
					executionOrder = &[]string{}
					r = r.WithContext(context.WithValue(r.Context(), rtr.ExecutionSequenceKey, executionOrder))
				}

				// Log entry point
				// t.Logf("ENTERING AFTER MIDDLEWARE: %s", name)
				// *executionOrder = append(*executionOrder, name+"_enter")

				// Create a response recorder to capture the response
				rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

				// Call the next handler
				next.ServeHTTP(rw, r)

				// Log execution point (after next returns, for AFTER middlewares)
				t.Logf("EXECUTING AFTER MIDDLEWARE: %s", name)
				*executionOrder = append(*executionOrder, name+"_execute")

				// Log exit point
				// t.Logf("EXITING AFTER MIDDLEWARE: %s", name)
				// *executionOrder = append(*executionOrder, name+"_exit")
			})
		})
	}

	// Create a new router
	r := rtr.NewRouter()

	// We'll modify the route handler to record its own execution
	// This will be used in the test routes

	// Global middlewares should wrap everything (added first)
	r.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
		traceBeforeMiddleware("global_before_1"), // should be executed first
		traceBeforeMiddleware("global_before_2"), // should be executed second
	})

	// Domain middlewares wrap domain-specific routes
	domain := rtr.NewDomain("example.com")
	domain.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
		traceBeforeMiddleware("domain_before_1"),
		traceBeforeMiddleware("domain_before_2"),
	})
	domain.AddAfterMiddlewares([]rtr.MiddlewareInterface{
		traceAfterMiddleware("domain_after_1"), // should be executed first
		traceAfterMiddleware("domain_after_2"), // should be executed second
	})
	r.AddDomain(domain)

	// Global after middlewares wrap everything (added last)
	r.AddAfterMiddlewares([]rtr.MiddlewareInterface{
		traceAfterMiddleware("global_after_1"), // should be executed first
		traceAfterMiddleware("global_after_2"), // should be executed second
	})

	return r, domain, traceBeforeMiddleware, traceAfterMiddleware
}

// responseWriter is a wrapper around http.ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// TestDirectRouteMiddlewareOrder tests middleware execution order for direct routes
func TestDirectRouteMiddlewareOrder(t *testing.T) {
	r, domain, traceBeforeMiddleware, traceAfterMiddleware := testMiddlewareSetup(t)

	// Create a slice to hold the execution order
	executionOrder := []string{}

	// Create a direct route with middlewares
	directRoute := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/direct").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			// Record handler execution in the context
			if val := r.Context().Value(rtr.ExecutionSequenceKey); val != nil {
				execOrder := val.(*[]string)
				t.Logf("HANDLER EXECUTION")
				*execOrder = append(*execOrder, "handler")
			}
			w.WriteHeader(http.StatusOK)
		}).
		AddBeforeMiddlewares([]rtr.MiddlewareInterface{
			traceBeforeMiddleware("route_before_1"), // should be executed first
			traceBeforeMiddleware("route_before_2"), // should be executed second
		}).
		AddAfterMiddlewares([]rtr.MiddlewareInterface{
			traceAfterMiddleware("route_after_1"), // should be executed first
			traceAfterMiddleware("route_after_2"), // should be executed second
		})

	// Add the route to the domain instead of the router
	domain.AddRoute(directRoute)

	req := httptest.NewRequest("GET", "http://example.com/direct", nil)
	req.Host = "example.com" // Set the Host header to match the domain
	w := httptest.NewRecorder()

	// Create a request with a context that has our execution order slice
	req = req.WithContext(context.WithValue(req.Context(), rtr.ExecutionSequenceKey, &executionOrder))
	r.ServeHTTP(w, req)

	// Verify the execution order matches the defined middleware order
	// The expected order shows the flow of execution:
	// 1. Before middlewares execute in order (global → domain → route)
	// 2. Handler executes
	// 3. After middlewares execute in reverse order (route → domain → global)
	// Each middleware has an _enter (when it starts) and _exit (when it completes) phase
	expectedOrder := []string{
		// Before middlewares enter (in definition order)
		"global_before_1_execute",
		"global_before_2_execute",
		"domain_before_1_execute",
		"domain_before_2_execute",
		"route_before_1_execute",
		"route_before_2_execute",

		// Handler executes here (no exit, just the handler itself)
		"handler",

		// After middlewares execute in reverse order (route → domain → global)
		"route_after_1_execute",
		"route_after_2_execute",
		"domain_after_1_execute",
		"domain_after_2_execute",
		"global_after_1_execute",
		"global_after_2_execute",
	}

	assertMiddlewareOrder(t, executionOrder, expectedOrder)
}

// TestGroupMiddlewareOrder tests middleware execution order for routes within groups
func TestGroupMiddlewareOrder(t *testing.T) {
	r, domain, traceBeforeMiddleware, traceAfterMiddleware := testMiddlewareSetup(t)
	var executionOrder []string

	// Create a group with middlewares
	group := rtr.NewGroup().
		SetPrefix("/api").
		AddBeforeMiddlewares([]rtr.MiddlewareInterface{
			traceBeforeMiddleware("group_before_1"),
			traceBeforeMiddleware("group_before_2"),
		}).
		AddAfterMiddlewares([]rtr.MiddlewareInterface{
			traceAfterMiddleware("group_after_1"),
			traceAfterMiddleware("group_after_2"),
		})

	// Add a route to the group
	groupRoute := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/users").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			// Record handler execution in the context
			if val := r.Context().Value(rtr.ExecutionSequenceKey); val != nil {
				execOrder := val.(*[]string)
				t.Logf("HANDLER EXECUTION")
				*execOrder = append(*execOrder, "handler")
			}
			w.WriteHeader(http.StatusOK)
		}).
		AddBeforeMiddlewares([]rtr.MiddlewareInterface{
			traceBeforeMiddleware("route_before_1"),
		}).
		AddAfterMiddlewares([]rtr.MiddlewareInterface{
			traceAfterMiddleware("route_after_1"),
		})

	group.AddRoute(groupRoute)
	domain.AddGroup(group)

	req := httptest.NewRequest("GET", "/api/users", nil)
	req.Host = "example.com"
	w := httptest.NewRecorder()

	// Create a request with a context that has our execution order slice
	req = req.WithContext(context.WithValue(req.Context(), rtr.ExecutionSequenceKey, &executionOrder))
	r.ServeHTTP(w, req)

	expectedOrder := []string{
		// Before middlewares (outer to inner)
		"global_before_1_execute",
		"global_before_2_execute",
		"domain_before_1_execute",
		"domain_before_2_execute",
		"group_before_1_execute",
		"group_before_2_execute",
		"route_before_1_execute",
		"handler",
		// After middlewares (inner to outer)
		"route_after_1_execute",
		"group_after_1_execute",
		"group_after_2_execute",
		"domain_after_1_execute",
		"domain_after_2_execute",
		"global_after_1_execute",
		"global_after_2_execute",
	}

	assertMiddlewareOrder(t, executionOrder, expectedOrder)
}

// // TestNestedGroupMiddlewareOrder tests middleware execution order for nested groups
// func TestNestedGroupMiddlewareOrder(t *testing.T) {
// 	r, domain, traceBeforeMiddleware, traceAfterMiddleware := testMiddlewareSetup(t)
// 	var executionOrder []string

// 	// Create parent group
// 	parentGroup := rtr.NewGroup().
// 		SetPrefix("/api").
// 		AddBeforeMiddlewares([]rtr.MiddlewareInterface{
// 			traceBeforeMiddleware("parent_group_before_1"),
// 		}).
// 		AddAfterMiddlewares([]rtr.MiddlewareInterface{
// 			traceAfterMiddleware("parent_group_after_1"),
// 		})

// 	// Create child group
// 	childGroup := rtr.NewGroup().
// 		SetPrefix("/v1").
// 		AddBeforeMiddlewares([]rtr.MiddlewareInterface{
// 			traceBeforeMiddleware("child_group_before_1"),
// 		}).
// 		AddAfterMiddlewares([]rtr.MiddlewareInterface{
// 			traceAfterMiddleware("child_group_after_1"),
// 		})

// 	// Add route to child group
// 	nestedRoute := rtr.NewRoute().
// 		SetMethod("GET").
// 		SetPath("/users").
// 		SetHandler(func(w http.ResponseWriter, r *http.Request) {
// 			if val := r.Context().Value(rtr.ExecutionSequenceKey); val != nil {
// 				executionOrder = *val.(*[]string)
// 			}
// 			executionOrder = append(executionOrder, "handler")
// 			w.WriteHeader(http.StatusOK)
// 		}).
// 		AddBeforeMiddlewares([]rtr.MiddlewareInterface{
// 			traceBeforeMiddleware("route_before_1"),
// 		})

// 	childGroup.AddRoute(nestedRoute)
// 	parentGroup.AddGroup(childGroup)
// 	domain.AddGroup(parentGroup)

// 	req := httptest.NewRequest("GET", "/api/v1/users", nil)
// 	req.Host = "example.com"
// 	w := httptest.NewRecorder()

// 	r.ServeHTTP(w, req)

// 	expectedOrder := []string{
// 		// Before middlewares (outer to inner)
// 		"global_before_1",
// 		"global_before_2",
// 		"domain_before_1",
// 		"domain_before_2",
// 		"parent_group_before_1",
// 		"child_group_before_1",
// 		"route_before_1",

// 		"handler",

// 		// After middlewares (inner to outer)
// 		"child_group_after_1",
// 		"parent_group_after_1",
// 		"domain_after_1",
// 		"domain_after_2",
// 		"global_after_1",
// 		"global_after_2",
// 	}

// 	assertMiddlewareOrder(t, executionOrder, expectedOrder)
// }

// // TestSiblingRoutesAndGroups tests that routes at the same level as groups only get global and domain middlewares
// func TestSiblingRoutesAndGroups(t *testing.T) {
// 	r, domain, traceBeforeMiddleware, _ := testMiddlewareSetup(t)
// 	var executionOrder []string

// 	// Create a group with middlewares
// 	group := rtr.NewGroup().
// 		SetPrefix("/api").
// 		AddBeforeMiddlewares([]rtr.MiddlewareInterface{
// 			traceBeforeMiddleware("group_before_1"),
// 		})

// 	// Add a route to the group
// 	groupRoute := rtr.NewRoute().
// 		SetMethod("GET").
// 		SetPath("/users").
// 		SetHandler(func(w http.ResponseWriter, r *http.Request) {
// 			executionOrder = *r.Context().Value("executionOrder").(*[]string)
// 			executionOrder = append(executionOrder, "route_handler")
// 			w.WriteHeader(http.StatusOK)
// 		})

// 	group.AddRoute(groupRoute)
// 	domain.AddGroup(group)

// 	// Add a route at the same level as the group
// 	siblingRoute := rtr.NewRoute().
// 		SetMethod("GET").
// 		SetPath("/sibling").
// 		SetHandler(func(w http.ResponseWriter, r *http.Request) {
// 			executionOrder = *r.Context().Value("executionOrder").(*[]string)
// 			executionOrder = append(executionOrder, "sibling_route_handler")
// 			w.WriteHeader(http.StatusOK)
// 		}).
// 		AddBeforeMiddlewares([]rtr.MiddlewareInterface{
// 			traceBeforeMiddleware("sibling_before_1"),
// 		})

// 	domain.AddRoute(siblingRoute)

// 	// Test the sibling route (should only get global and domain middlewares)
// 	req := httptest.NewRequest("GET", "/sibling", nil)
// 	req.Host = "example.com"
// 	w := httptest.NewRecorder()
// 	executionOrder = nil
// 	r.ServeHTTP(w, req)

// 	expectedSiblingOrder := []string{
// 		// Before middlewares (only global and domain)
// 		"global_before_1 before",
// 		"global_before_2 before",
// 		"domain_before_1 before",
// 		"domain_before_2 before",
// 		"sibling_before_1 before",

// 		"sibling_route_handler",

// 		// After middlewares (only domain and global)
// 		"domain_after_1 after",
// 		"domain_after_2 after",
// 		"global_after_1 after",
// 		"global_after_2 after",
// 	}

// 	assertMiddlewareOrder(t, executionOrder, expectedSiblingOrder)

// 	// Test the group route (should get all middlewares)
// 	req = httptest.NewRequest("GET", "/api/users", nil)
// 	req.Host = "example.com"
// 	w = httptest.NewRecorder()
// 	executionOrder = nil
// 	r.ServeHTTP(w, req)

// 	expectedGroupOrder := []string{
// 		// Before middlewares (all levels)
// 		"global_before_1 before",
// 		"global_before_2 before",
// 		"domain_before_1 before",
// 		"domain_before_2 before",
// 		"group_before_1 before",
// 		"route_handler",
// 		// After middlewares (all levels, reversed)
// 		"domain_after_1 after",
// 		"domain_after_2 after",
// 		"group_after_1 after",
// 		"group_after_2 after",
// 		"route_after_1 after",
// 		"global_after_1 after",
// 		"global_after_2 after",
// 	}

// 	assertMiddlewareOrder(t, executionOrder, expectedGroupOrder)
// }

// assertMiddlewareOrder is a helper function to assert the middleware execution order
// It verifies that the actual middleware execution order matches the expected order
func assertMiddlewareOrder(t *testing.T, actual, expected []string) {
	t.Helper()

	t.Logf("\n=== Expected Middleware Order ===")
	for i, v := range expected {
		t.Logf("  %2d: %s", i, v)
	}

	t.Logf("\n=== Actual Middleware Order ===")
	for i, v := range actual {
		t.Logf("  %2d: %s", i, v)
	}

	if len(actual) != len(expected) {
		t.Errorf("\nMiddleware execution length mismatch.\n\tExpected (%d): %v\n\tGot (%d): %v",
			len(expected), strings.Join(expected, " -> "),
			len(actual), strings.Join(actual, " -> "))
		return
	}

	for i := range actual {
		if actual[i] != expected[i] {
			t.Errorf("\nMiddleware execution order mismatch at index %d.\n\tExpected: %s\n\tGot: %s",
				i, expected[i], actual[i])
			t.Logf("Full expected order: %v", expected)
			t.Logf("Full actual order: %v", actual)
			break
		}
	}
}
