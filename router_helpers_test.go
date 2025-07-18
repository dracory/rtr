package rtr_test

import (
	"net/http"
	"testing"

	"github.com/dracory/rtr"
)

func TestRouterGettersAndSetters(t *testing.T) {
	t.Run("GetPrefix and SetPrefix", func(t *testing.T) {
		r := rtr.NewRouter()
		
		// Test default prefix (should be empty)
		if got := r.GetPrefix(); got != "" {
			t.Errorf("Expected empty prefix, got %q", got)
		}

		// Set a new prefix
		r.SetPrefix("/api")
		
		// Verify the prefix was set
		if got := r.GetPrefix(); got != "/api" {
			t.Errorf("Expected prefix '/api', got %q", got)
		}
	})

	t.Run("AddGroups and GetGroups", func(t *testing.T) {
		r := rtr.NewRouter()
		
		// Create some test groups
		group1 := rtr.NewGroup().SetPrefix("/v1")
		group2 := rtr.NewGroup().SetPrefix("/v2")
		
		// Add groups individually
		r.AddGroup(group1)
		
		// Add multiple groups at once
		r.AddGroups([]rtr.GroupInterface{group2})
		
		// Get all groups
		groups := r.GetGroups()
		
		// Verify the groups were added
		if len(groups) != 2 {
			t.Fatalf("Expected 2 groups, got %d", len(groups))
		}
		
		// Verify the groups are in the correct order
		if groups[0] != group1 || groups[1] != group2 {
			t.Error("Groups were not added in the expected order")
		}
	})

	t.Run("AddRoutes and GetRoutes", func(t *testing.T) {
		r := rtr.NewRouter()
		
		// Create some test routes
		route1 := rtr.NewRoute().SetMethod(http.MethodGet).SetPath("/test1").SetHandler(nil)
		route2 := rtr.NewRoute().SetMethod(http.MethodPost).SetPath("/test2").SetHandler(nil)
		
		// Add routes individually
		r.AddRoute(route1)
		
		// Add multiple routes at once
		r.AddRoutes([]rtr.RouteInterface{route2})
		
		// Get all routes
		routes := r.GetRoutes()
		
		// Verify the routes were added
		if len(routes) != 2 {
			t.Fatalf("Expected 2 routes, got %d", len(routes))
		}
		
		// Verify the routes are in the correct order
		if routes[0] != route1 || routes[1] != route2 {
			t.Error("Routes were not added in the expected order")
		}
	})

	t.Run("Middleware Getters", func(t *testing.T) {
		r := rtr.NewRouter()
		
		// Get the default middlewares count
		defaultBeforeMWs := r.GetBeforeMiddlewares()
		defaultAfterMWs := r.GetAfterMiddlewares()
		
		// Add some test middlewares
		beforeMW := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { next.ServeHTTP(w, r) })
		}
		
		afterMW := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { next.ServeHTTP(w, r) })
		}
		
		r.AddBeforeMiddlewares([]rtr.Middleware{beforeMW})
		r.AddAfterMiddlewares([]rtr.Middleware{afterMW})
		
		// Test GetBeforeMiddlewares - should include both default and our test middleware
		beforeMWs := r.GetBeforeMiddlewares()
		expectedBeforeCount := len(defaultBeforeMWs) + 1
		if len(beforeMWs) != expectedBeforeCount || beforeMWs[expectedBeforeCount-1] == nil {
			t.Errorf("Expected %d before middlewares (default + 1), got %d", expectedBeforeCount, len(beforeMWs))
		}
		
		// Test GetAfterMiddlewares - should include both default and our test middleware
		afterMWs := r.GetAfterMiddlewares()
		expectedAfterCount := len(defaultAfterMWs) + 1
		if len(afterMWs) != expectedAfterCount || afterMWs[expectedAfterCount-1] == nil {
			t.Errorf("Expected %d after middlewares (default + 1), got %d", expectedAfterCount, len(afterMWs))
		}
	})
}
