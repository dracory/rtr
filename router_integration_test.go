package rtr_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/rtr"
)

// TestRouterWithDataHandling tests the router with data handling scenarios
func TestRouterWithDataHandling(t *testing.T) {
	// Simulate some in-memory data
	users := []map[string]interface{}{
		{"id": 1, "name": "Test User", "email": "test@example.com"},
		{"id": 2, "name": "Another User", "email": "another@example.com"},
	}

	// Create a router
	r := rtr.NewRouter()

	// Add a route that returns user data
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/users").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			// Build a response from in-memory data
			var response string
			for _, user := range users {
				response += fmt.Sprintf("User %v: %v (%v)\n", user["id"], user["name"], user["email"])
			}
			fmt.Fprint(w, response)
		})
	r.AddRoute(route)

	// Create a test request
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response contains the expected data
	expected := "User 1: Test User (test@example.com)\nUser 2: Another User (another@example.com)\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// TestRouterWithDataMiddleware tests middleware that injects data handling
func TestRouterWithDataMiddleware(t *testing.T) {
	// Simulate some in-memory data
	items := []map[string]interface{}{
		{"id": 1, "name": "Test Item"},
		{"id": 2, "name": "Another Item"},
	}

	// Create a router
	r := rtr.NewRouter()

	// Create a middleware that adds logging or data validation
	dataMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add some middleware logic (e.g., logging, validation)
			w.Header().Set("X-Middleware", "data-middleware")
			next.ServeHTTP(w, r)
		})
	}

	// Add the middleware to the router
	r.AddBeforeMiddlewares(rtr.MiddlewaresToInterfaces([]rtr.StdMiddleware{dataMiddleware}))

	// Add a route that uses the in-memory data
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/items").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			// Build a response from in-memory data
			var response string
			for _, item := range items {
				response += fmt.Sprintf("Item %v: %v\n", item["id"], item["name"])
			}
			fmt.Fprint(w, response)
		})
	r.AddRoute(route)

	// Create a test request
	req, err := http.NewRequest("GET", "/items", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response contains the expected data
	expected := "Item 1: Test Item\nItem 2: Another Item\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Check that middleware header was set
	if middleware := rr.Header().Get("X-Middleware"); middleware != "data-middleware" {
		t.Errorf("middleware header not set correctly: got %v want %v", middleware, "data-middleware")
	}
}

// TestRouterWithMultipleRouteOperations tests a more complex scenario with multiple route operations
func TestRouterWithMultipleRouteOperations(t *testing.T) {
	// Simulate some in-memory data
	products := []map[string]interface{}{
		{"id": 1, "name": "Product 1", "price": 10.99},
		{"id": 2, "name": "Product 2", "price": 20.49},
	}

	// Create a router
	r := rtr.NewRouter()

	// Create a group for product-related routes
	productGroup := rtr.NewGroup().SetPrefix("/products")

	// Add routes to the group
	productGroup.AddRoute(rtr.NewRoute().SetMethod("GET").SetPath("").SetHandler(func(w http.ResponseWriter, r *http.Request) {
		// List all products
		var response string
		for _, product := range products {
			response += fmt.Sprintf("Product %v: %v ($%.2f)\n", product["id"], product["name"], product["price"])
		}
		fmt.Fprint(w, response)
	}))

	productGroup.AddRoute(rtr.NewRoute().SetMethod("GET").SetPath("/total").SetHandler(func(w http.ResponseWriter, r *http.Request) {
		// Calculate total price of all products
		var total float64
		for _, product := range products {
			if price, ok := product["price"].(float64); ok {
				total += price
			}
		}
		fmt.Fprintf(w, "Total price: $%.2f", total)
	}))

	// Add the group to the router
	r.AddGroup(productGroup)

	// Test the list products route
	req1, err := http.NewRequest("GET", "/products", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr1 := httptest.NewRecorder()
	r.ServeHTTP(rr1, req1)

	if status := rr1.Code; status != http.StatusOK {
		t.Errorf("list handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected1 := "Product 1: Product 1 ($10.99)\nProduct 2: Product 2 ($20.49)\n"
	if rr1.Body.String() != expected1 {
		t.Errorf("list handler returned unexpected body: got %v want %v", rr1.Body.String(), expected1)
	}

	// Test the total price route
	req2, err := http.NewRequest("GET", "/products/total", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr2 := httptest.NewRecorder()
	r.ServeHTTP(rr2, req2)

	if status := rr2.Code; status != http.StatusOK {
		t.Errorf("total handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected2 := "Total price: $31.48"
	if rr2.Body.String() != expected2 {
		t.Errorf("total handler returned unexpected body: got %v want %v", rr2.Body.String(), expected2)
	}
}
