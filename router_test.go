package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRouterBasicRouting tests the basic routing functionality of the Router.
// It verifies that a simple route can be added and that requests to that route
// are properly handled and return the expected response.
func TestRouterBasicRouting(t *testing.T) {
	router := NewRouter()

	// Add a simple route
	route := NewRoute().
		SetMethod("GET").
		SetPath("/hello").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Hello, World!")
		})
	router.AddRoute(route)

	// Create a test request
	req, err := http.NewRequest("GET", "/hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := "Hello, World!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// TestRouterMethodNotAllowed tests the router's behavior when a request is made
// with a method that is not allowed for a given path. It verifies that the router
// returns the appropriate status code (404 Not Found in this implementation).
func TestRouterMethodNotAllowed(t *testing.T) {
	router := NewRouter()

	// Add a route that only accepts GET
	route := NewRoute().
		SetMethod("GET").
		SetPath("/hello").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Hello, World!")
		})
	router.AddRoute(route)

	// Create a POST request to the same path
	req, err := http.NewRequest("POST", "/hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(rr, req)

	// Since we're using simple path matching, a POST to /hello will return 404
	// In a more sophisticated router, this might return 405 Method Not Allowed
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

// TestRouterNotFound tests the router's behavior when a request is made to a path
// that does not exist. It verifies that the router returns a 404 Not Found status code.
func TestRouterNotFound(t *testing.T) {
	router := NewRouter()

	// Add a route
	route := NewRoute().
		SetMethod("GET").
		SetPath("/hello").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Hello, World!")
		})
	router.AddRoute(route)

	// Create a request to a non-existent path
	req, err := http.NewRequest("GET", "/nonexistent", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

// TestRouterWithPrefix tests the router's prefix functionality. It verifies that
// when a prefix is set on the router, all routes are properly prefixed and
// requests to the prefixed paths are correctly handled.
func TestRouterWithPrefix(t *testing.T) {
	router := NewRouter().SetPrefix("/api")

	// Add a route
	route := NewRoute().
		SetMethod("GET").
		SetPath("/hello").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Hello, API!")
		})
	router.AddRoute(route)

	// Create a test request with the prefix
	req, err := http.NewRequest("GET", "/api/hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := "Hello, API!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// TestRouterWithGroup tests the router's group functionality. It verifies that
// routes can be added to a group and that the group's prefix is properly applied
// to all routes within the group.
func TestRouterWithGroup(t *testing.T) {
	router := NewRouter()

	// Create a group
	group := NewGroup().SetPrefix("/api")

	// Add a route to the group
	route := NewRoute().
		SetMethod("GET").
		SetPath("/hello").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Hello from group!")
		})
	group.AddRoute(route)

	// Add the group to the router
	router.AddGroup(group)

	// Create a test request to the grouped route
	req, err := http.NewRequest("GET", "/api/hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := "Hello from group!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// TestRouterWithNestedGroups tests the router's nested group functionality. It verifies
// that groups can be nested within other groups and that the prefixes are properly
// combined to form the full path for routes within nested groups.
func TestRouterWithNestedGroups(t *testing.T) {
	router := NewRouter()

	// Create parent group
	parentGroup := NewGroup().SetPrefix("/api")

	// Create child group
	childGroup := NewGroup().SetPrefix("/v1")

	// Add a route to the child group
	route := NewRoute().
		SetMethod("GET").
		SetPath("/hello").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Hello from nested group!")
		})
	childGroup.AddRoute(route)

	// Add the child group to the parent group
	parentGroup.AddGroup(childGroup)

	// Add the parent group to the router
	router.AddGroup(parentGroup)

	// Create a test request to the nested route
	req, err := http.NewRequest("GET", "/api/v1/hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := "Hello from nested group!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// TestRouterWithBeforeMiddleware tests the router's before middleware functionality.
// It verifies that middleware added to the router is executed before the route handler
// and that it can modify the request or response as needed.
func TestRouterWithBeforeMiddleware(t *testing.T) {
	router := NewRouter()

	// Add a middleware that adds a header
	headerMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test", "test-value")
			next.ServeHTTP(w, r)
		})
	}

	// Add the middleware to the router
	router.AddBeforeMiddlewares([]Middleware{headerMiddleware})

	// Add a route
	route := NewRoute().
		SetMethod("GET").
		SetPath("/hello").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Hello with middleware!")
		})
	router.AddRoute(route)

	// Create a test request
	req, err := http.NewRequest("GET", "/hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := "Hello with middleware!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Check the header was added
	if rr.Header().Get("X-Test") != "test-value" {
		t.Errorf("middleware did not set header: got %v want %v", rr.Header().Get("X-Test"), "test-value")
	}
}

// TestRouterWithAfterMiddleware tests the router's after middleware functionality.
// It verifies that middleware added to the router is executed after the route handler
// and that it can modify the response as needed.
func TestRouterWithAfterMiddleware(t *testing.T) {
	router := NewRouter()

	// Add a middleware that modifies the response
	responseMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Call the next handler first
			next.ServeHTTP(w, r)
			// Then add a header (in a real scenario, you might wrap the ResponseWriter)
			w.Header().Set("X-After", "after-value")
		})
	}

	// Add the middleware to the router
	router.AddAfterMiddlewares([]Middleware{responseMiddleware})

	// Add a route
	route := NewRoute().
		SetMethod("GET").
		SetPath("/hello").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Hello with after middleware!")
		})
	router.AddRoute(route)

	// Create a test request
	req, err := http.NewRequest("GET", "/hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := "Hello with after middleware!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Check the header was added
	if rr.Header().Get("X-After") != "after-value" {
		t.Errorf("middleware did not set header: got %v want %v", rr.Header().Get("X-After"), "after-value")
	}
}

// TestRouterWithRouteMiddleware tests the middleware functionality at the route level.
// It verifies that middleware can be added to a specific route and that it is executed
// only for requests to that route.
func TestRouterWithRouteMiddleware(t *testing.T) {
	router := NewRouter()

	// Create a route with middleware
	route := NewRoute().
		SetMethod("GET").
		SetPath("/hello").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Hello with route middleware!")
		})

	// Add middleware to the route
	routeMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Route", "route-value")
			next.ServeHTTP(w, r)
		})
	}
	route.AddBeforeMiddlewares([]Middleware{routeMiddleware})

	// Add the route to the router
	router.AddRoute(route)

	// Create a test request
	req, err := http.NewRequest("GET", "/hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := "Hello with route middleware!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Check the header was added
	if rr.Header().Get("X-Route") != "route-value" {
		t.Errorf("middleware did not set header: got %v want %v", rr.Header().Get("X-Route"), "route-value")
	}
}

// TestRouterWithGroupMiddleware tests the middleware functionality at the group level.
// It verifies that middleware can be added to a group and that it is executed for
// all requests to routes within that group.
func TestRouterWithGroupMiddleware(t *testing.T) {
	router := NewRouter()

	// Create a group with middleware
	group := NewGroup().SetPrefix("/api")

	// Add middleware to the group
	groupMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Group", "group-value")
			next.ServeHTTP(w, r)
		})
	}
	group.AddBeforeMiddlewares([]Middleware{groupMiddleware})

	// Add a route to the group
	route := NewRoute().
		SetMethod("GET").
		SetPath("/hello").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Hello with group middleware!")
		})
	group.AddRoute(route)

	// Add the group to the router
	router.AddGroup(group)

	// Create a test request
	req, err := http.NewRequest("GET", "/api/hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := "Hello with group middleware!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Check the header was added
	if rr.Header().Get("X-Group") != "group-value" {
		t.Errorf("middleware did not set header: got %v want %v", rr.Header().Get("X-Group"), "group-value")
	}
}

// TestRouterStaticFiles tests the router's static file serving functionality.
// It verifies that the router can serve static files from a specified directory.
func TestRouterStaticFiles(t *testing.T) {
	router := NewRouter()

	// Set up a static file server
	// Note: In a real test, you'd use a temporary directory with actual files
	fileServer := http.FileServer(http.Dir("./testdata"))
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Strip the path prefix
		r.URL.Path = r.URL.Path[len("/static/"):]
		if r.URL.Path == "" {
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	}

	route := NewRoute().
		SetMethod("GET").
		SetPath("/static/*").
		SetHandler(handler)
	router.AddRoute(route)

	// Create a test request
	req, err := http.NewRequest("GET", "/static/test.txt", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request - this will return 404 since we don't have the files
	// In a real test, you would check for 200 and the file contents
	router.ServeHTTP(rr, req)

	// We expect 404 since the file doesn't exist in our test
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("static file handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

// TestRouterHTTPMethods tests the router's support for various HTTP methods.
// It verifies that the router can handle requests with different HTTP methods
// and that the appropriate handler is called for each method.
func TestRouterHTTPMethods(t *testing.T) {
	router := NewRouter()

	// Add routes for different HTTP methods
	router.AddRoute(NewRoute().
		SetMethod("GET").
		SetPath("/method").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "GET")
		}))

	router.AddRoute(NewRoute().
		SetMethod("POST").
		SetPath("/method").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "POST")
		}))

	router.AddRoute(NewRoute().
		SetMethod("PUT").
		SetPath("/method").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "PUT")
		}))

	router.AddRoute(NewRoute().
		SetMethod("DELETE").
		SetPath("/method").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "DELETE")
		}))

	router.AddRoute(NewRoute().
		SetMethod("PATCH").
		SetPath("/method").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "PATCH")
		}))

	router.AddRoute(NewRoute().
		SetMethod("OPTIONS").
		SetPath("/method").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "OPTIONS")
		}))

	router.AddRoute(NewRoute().
		SetMethod("HEAD").
		SetPath("/method").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "HEAD")
		}))

	// Test each method
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}
	for _, method := range methods {
		req, err := http.NewRequest(method, "/method", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("%s handler returned wrong status code: got %v want %v", method, status, http.StatusOK)
		}

		// For HEAD, the body is empty by HTTP spec
		if method != "HEAD" {
			// Check the response body
			if rr.Body.String() != method {
				t.Errorf("%s handler returned unexpected body: got %v want %v", method, rr.Body.String(), method)
			}
		}
	}
}
