package router_test

import (
	"net/http"
	"testing"

	"github.com/dracory/router"
)

// TestGroupInterface tests the basic functionality of the GroupInterface implementation.
// It verifies that the group can be created, and that its prefix, routes, and middleware
// can be properly set and retrieved.
func TestGroupInterface(t *testing.T) {
	// Create a new group
	group := router.NewGroup()

	// Test prefix getter/setter
	prefix := "/api"
	group.SetPrefix(prefix)
	if group.GetPrefix() != prefix {
		t.Errorf("Expected prefix %s, got %s", prefix, group.GetPrefix())
	}

	// Test routes getter/setter
	route := router.NewRoute().
		SetMethod("GET").
		SetPath("/test").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {})

	// Add a single route
	group.AddRoute(route)
	if len(group.GetRoutes()) != 1 {
		t.Errorf("Expected 1 route, got %d", len(group.GetRoutes()))
	}

	// Add multiple routes
	route2 := router.NewRoute().SetMethod("POST").SetPath("/test2")
	route3 := router.NewRoute().SetMethod("PUT").SetPath("/test3")
	group.AddRoutes([]router.RouteInterface{route2, route3})

	if len(group.GetRoutes()) != 3 {
		t.Errorf("Expected 3 routes, got %d", len(group.GetRoutes()))
	}

	// Test middleware getters/setters
	middleware1 := func(next http.Handler) http.Handler { return next }
	middleware2 := func(next http.Handler) http.Handler { return next }

	// Test before middlewares
	group.AddBeforeMiddlewares([]router.Middleware{middleware1})
	if len(group.GetBeforeMiddlewares()) != 1 {
		t.Errorf("Expected 1 before middleware, got %d", len(group.GetBeforeMiddlewares()))
	}

	group.AddBeforeMiddlewares([]router.Middleware{middleware2})
	if len(group.GetBeforeMiddlewares()) != 2 {
		t.Errorf("Expected 2 before middlewares, got %d", len(group.GetBeforeMiddlewares()))
	}

	// Test after middlewares
	group.AddAfterMiddlewares([]router.Middleware{middleware1})
	if len(group.GetAfterMiddlewares()) != 1 {
		t.Errorf("Expected 1 after middleware, got %d", len(group.GetAfterMiddlewares()))
	}

	group.AddAfterMiddlewares([]router.Middleware{middleware2})
	if len(group.GetAfterMiddlewares()) != 2 {
		t.Errorf("Expected 2 after middlewares, got %d", len(group.GetAfterMiddlewares()))
	}
}

// TestGroupChaining tests the method chaining functionality of the GroupInterface.
// It verifies that multiple methods can be called in sequence and that the group
// state is correctly updated after each method call.
func TestGroupChaining(t *testing.T) {
	// Test method chaining
	route := router.NewRoute().SetMethod("GET").SetPath("/users")

	group := router.NewGroup().
		SetPrefix("/api").
		AddRoute(route)

	if group.GetPrefix() != "/api" {
		t.Errorf("Expected prefix /api, got %s", group.GetPrefix())
	}

	if len(group.GetRoutes()) != 1 {
		t.Errorf("Expected 1 route, got %d", len(group.GetRoutes()))
	}

	// Test middleware chaining
	middleware := func(next http.Handler) http.Handler { return next }

	group.AddBeforeMiddlewares([]router.Middleware{middleware}).
		AddAfterMiddlewares([]router.Middleware{middleware})

	if len(group.GetBeforeMiddlewares()) != 1 {
		t.Errorf("Expected 1 before middleware, got %d", len(group.GetBeforeMiddlewares()))
	}

	if len(group.GetAfterMiddlewares()) != 1 {
		t.Errorf("Expected 1 after middleware, got %d", len(group.GetAfterMiddlewares()))
	}
}

// TestNestedGroups tests the nested group functionality of the GroupInterface.
// It verifies that groups can be nested within other groups and that the nested
// groups can be properly accessed and manipulated.
func TestNestedGroups(t *testing.T) {
	// Create parent group
	parentGroup := router.NewGroup().SetPrefix("/api")

	// Create child groups
	childGroup1 := router.NewGroup().SetPrefix("/v1")
	childGroup2 := router.NewGroup().SetPrefix("/v2")

	// Add routes to child groups
	childGroup1.AddRoute(router.NewRoute().SetMethod("GET").SetPath("/users"))
	childGroup2.AddRoute(router.NewRoute().SetMethod("GET").SetPath("/products"))

	// Add child groups to parent group
	parentGroup.AddGroup(childGroup1)
	if len(parentGroup.GetGroups()) != 1 {
		t.Errorf("Expected 1 group, got %d", len(parentGroup.GetGroups()))
	}

	parentGroup.AddGroups([]router.GroupInterface{childGroup2})
	if len(parentGroup.GetGroups()) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(parentGroup.GetGroups()))
	}

	// Verify the first child group
	firstChildGroup := parentGroup.GetGroups()[0]
	if firstChildGroup.GetPrefix() != "/v1" {
		t.Errorf("Expected first child group prefix /v1, got %s", firstChildGroup.GetPrefix())
	}

	if len(firstChildGroup.GetRoutes()) != 1 {
		t.Errorf("Expected first child group to have 1 route, got %d", len(firstChildGroup.GetRoutes()))
	}

	// Verify the second child group
	secondChildGroup := parentGroup.GetGroups()[1]
	if secondChildGroup.GetPrefix() != "/v2" {
		t.Errorf("Expected second child group prefix /v2, got %s", secondChildGroup.GetPrefix())
	}

	if len(secondChildGroup.GetRoutes()) != 1 {
		t.Errorf("Expected second child group to have 1 route, got %d", len(secondChildGroup.GetRoutes()))
	}
}

// TestGroupWithMiddlewares tests the middleware functionality of the GroupInterface.
// It verifies that middleware functions can be added to a group and that they are
// properly stored and can be retrieved.
func TestGroupWithMiddlewares(t *testing.T) {
	// Create a group with multiple middlewares
	group := router.NewGroup().SetPrefix("/api")

	// Create middlewares that modify a counter
	counter := 0

	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter += 1
			next.ServeHTTP(w, r)
		})
	}

	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter += 10
			next.ServeHTTP(w, r)
		})
	}

	middleware3 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter += 100
			next.ServeHTTP(w, r)
		})
	}

	// Add middlewares to the group
	group.AddBeforeMiddlewares([]router.Middleware{middleware1, middleware2})
	group.AddAfterMiddlewares([]router.Middleware{middleware3})

	// Check middleware counts
	if len(group.GetBeforeMiddlewares()) != 2 {
		t.Errorf("Expected 2 before middlewares, got %d", len(group.GetBeforeMiddlewares()))
	}

	if len(group.GetAfterMiddlewares()) != 1 {
		t.Errorf("Expected 1 after middleware, got %d", len(group.GetAfterMiddlewares()))
	}
}

// TestComplexGroupStructure tests the creation and manipulation of a complex nested group structure.
// It verifies that groups can be nested at multiple levels and that the structure
// can be properly traversed and verified.
func TestComplexGroupStructure(t *testing.T) {
	// Create a complex group structure with nested groups and routes

	// Main group
	apiGroup := router.NewGroup().SetPrefix("/api")

	// Version groups
	v1Group := router.NewGroup().SetPrefix("/v1")
	v2Group := router.NewGroup().SetPrefix("/v2")

	// Resource groups
	usersGroup := router.NewGroup().SetPrefix("/users")
	productsGroup := router.NewGroup().SetPrefix("/products")

	// Add routes to resource groups
	usersGroup.AddRoute(router.NewRoute().SetMethod("GET").SetPath(""))
	usersGroup.AddRoute(router.NewRoute().SetMethod("GET").SetPath("/:id"))
	usersGroup.AddRoute(router.NewRoute().SetMethod("POST").SetPath(""))

	productsGroup.AddRoute(router.NewRoute().SetMethod("GET").SetPath(""))
	productsGroup.AddRoute(router.NewRoute().SetMethod("GET").SetPath("/:id"))

	// Add resource groups to version groups
	v1Group.AddGroup(usersGroup)
	v1Group.AddGroup(productsGroup)

	v2Group.AddGroup(usersGroup)

	// Add version groups to API group
	apiGroup.AddGroup(v1Group)
	apiGroup.AddGroup(v2Group)

	// Verify structure
	if len(apiGroup.GetGroups()) != 2 {
		t.Errorf("Expected API group to have 2 subgroups, got %d", len(apiGroup.GetGroups()))
	}

	v1 := apiGroup.GetGroups()[0]
	if v1.GetPrefix() != "/v1" {
		t.Errorf("Expected first subgroup prefix /v1, got %s", v1.GetPrefix())
	}

	if len(v1.GetGroups()) != 2 {
		t.Errorf("Expected v1 group to have 2 subgroups, got %d", len(v1.GetGroups()))
	}

	v2 := apiGroup.GetGroups()[1]
	if v2.GetPrefix() != "/v2" {
		t.Errorf("Expected second subgroup prefix /v2, got %s", v2.GetPrefix())
	}

	if len(v2.GetGroups()) != 1 {
		t.Errorf("Expected v2 group to have 1 subgroup, got %d", len(v2.GetGroups()))
	}

	// Check users group in v1
	usersInV1 := v1.GetGroups()[0]
	if usersInV1.GetPrefix() != "/users" {
		t.Errorf("Expected users group prefix /users, got %s", usersInV1.GetPrefix())
	}

	if len(usersInV1.GetRoutes()) != 3 {
		t.Errorf("Expected users group to have 3 routes, got %d", len(usersInV1.GetRoutes()))
	}

	// Check products group in v1
	productsInV1 := v1.GetGroups()[1]
	if productsInV1.GetPrefix() != "/products" {
		t.Errorf("Expected products group prefix /products, got %s", productsInV1.GetPrefix())
	}

	if len(productsInV1.GetRoutes()) != 2 {
		t.Errorf("Expected products group to have 2 routes, got %d", len(productsInV1.GetRoutes()))
	}
}

// TestGroupWithDatabaseIntegration tests the integration of the GroupInterface with a simulated database.
// It verifies that routes can be added to a group with handlers that would interact with a database,
// and that the routes can be properly accessed and verified.
func TestGroupWithDatabaseIntegration(t *testing.T) {
	// Create a group for database-related routes
	dbGroup := router.NewGroup().SetPrefix("/db")

	// Add routes for database operations
	dbGroup.AddRoute(router.NewRoute().
		SetMethod("GET").
		SetPath("/query").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			// In a real test, this would query a database
		}))

	dbGroup.AddRoute(router.NewRoute().
		SetMethod("POST").
		SetPath("/execute").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			// In a real test, this would execute a database command
		}))

	// Verify the group structure
	if dbGroup.GetPrefix() != "/db" {
		t.Errorf("Expected group prefix /db, got %s", dbGroup.GetPrefix())
	}

	if len(dbGroup.GetRoutes()) != 2 {
		t.Errorf("Expected 2 routes, got %d", len(dbGroup.GetRoutes()))
	}

	// Verify the first route
	queryRoute := dbGroup.GetRoutes()[0]
	if queryRoute.GetMethod() != "GET" {
		t.Errorf("Expected query route method GET, got %s", queryRoute.GetMethod())
	}

	if queryRoute.GetPath() != "/query" {
		t.Errorf("Expected query route path /query, got %s", queryRoute.GetPath())
	}

	// Verify the second route
	executeRoute := dbGroup.GetRoutes()[1]
	if executeRoute.GetMethod() != "POST" {
		t.Errorf("Expected execute route method POST, got %s", executeRoute.GetMethod())
	}

	if executeRoute.GetPath() != "/execute" {
		t.Errorf("Expected execute route path /execute, got %s", executeRoute.GetPath())
	}
}
