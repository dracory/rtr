package router_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/router"

	_ "github.com/mattn/go-sqlite3"
)

// TestRouterWithDatabase tests the router with a real SQLite database connection
func TestRouterWithDatabase(t *testing.T) {
	// Set up in-memory SQLite database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create a test table
	_, err = db.Exec(`CREATE TABLE users (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL
	)`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Insert test data
	_, err = db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", "Test User", "test@example.com")
	if err != nil {
		t.Fatalf("Failed to insert data: %v", err)
	}

	// Create a router
	r := router.NewRouter()

	// Add a route that uses the database
	route := router.NewRoute().
		SetMethod("GET").
		SetPath("/users").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			// Query the database
			rows, err := db.Query("SELECT id, name, email FROM users")
			if err != nil {
				http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			// Build a response
			var response string
			for rows.Next() {
				var id int
				var name, email string
				if err := rows.Scan(&id, &name, &email); err != nil {
					http.Error(w, fmt.Sprintf("Row scan error: %v", err), http.StatusInternalServerError)
					return
				}
				response += fmt.Sprintf("User %d: %s (%s)\n", id, name, email)
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
	expected := "User 1: Test User (test@example.com)\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// TestRouterWithDatabaseMiddleware tests middleware that injects a database connection
func TestRouterWithDatabaseMiddleware(t *testing.T) {
	// Set up in-memory SQLite database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create a test table
	_, err = db.Exec(`CREATE TABLE items (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL
	)`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Insert test data
	_, err = db.Exec("INSERT INTO items (name) VALUES (?)", "Test Item")
	if err != nil {
		t.Fatalf("Failed to insert data: %v", err)
	}

	// Create a router
	r := router.NewRouter()

	// Create a middleware that adds the database connection to the request context
	dbMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// In a real application, you would use context to pass the db
			// For this test, we're using a closure to access the db variable
			next.ServeHTTP(w, r)
		})
	}

	// Add the middleware to the router
	r.AddBeforeMiddlewares([]router.Middleware{dbMiddleware})

	// Add a route that uses the database from the middleware
	route := router.NewRoute().
		SetMethod("GET").
		SetPath("/items").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			// Query the database (accessing it from the closure)
			rows, err := db.Query("SELECT id, name FROM items")
			if err != nil {
				http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			// Build a response
			var response string
			for rows.Next() {
				var id int
				var name string
				if err := rows.Scan(&id, &name); err != nil {
					http.Error(w, fmt.Sprintf("Row scan error: %v", err), http.StatusInternalServerError)
					return
				}
				response += fmt.Sprintf("Item %d: %s\n", id, name)
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
	expected := "Item 1: Test Item\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// TestRouterWithMultipleDatabaseOperations tests a more complex scenario with multiple database operations
func TestRouterWithMultipleDatabaseOperations(t *testing.T) {
	// Set up in-memory SQLite database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create a test table
	_, err = db.Exec(`CREATE TABLE products (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		price REAL NOT NULL
	)`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Insert test data
	_, err = db.Exec("INSERT INTO products (name, price) VALUES (?, ?)", "Product 1", 10.99)
	if err != nil {
		t.Fatalf("Failed to insert data: %v", err)
	}
	_, err = db.Exec("INSERT INTO products (name, price) VALUES (?, ?)", "Product 2", 20.49)
	if err != nil {
		t.Fatalf("Failed to insert data: %v", err)
	}

	// Create a router
	r := router.NewRouter()

	// Create a group for product-related routes
	productGroup := router.NewGroup().SetPrefix("/products")

	// Add routes to the group
	productGroup.AddRoute(router.NewRoute().SetMethod("GET").SetPath("").SetHandler(func(w http.ResponseWriter, r *http.Request) {
		// List all products
		rows, err := db.Query("SELECT id, name, price FROM products")
		if err != nil {
			http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var response string
		for rows.Next() {
			var id int
			var name string
			var price float64
			if err := rows.Scan(&id, &name, &price); err != nil {
				http.Error(w, fmt.Sprintf("Row scan error: %v", err), http.StatusInternalServerError)
				return
			}
			response += fmt.Sprintf("Product %d: %s ($%.2f)\n", id, name, price)
		}

		fmt.Fprint(w, response)
	}))

	productGroup.AddRoute(router.NewRoute().SetMethod("GET").SetPath("/total").SetHandler(func(w http.ResponseWriter, r *http.Request) {
		// Calculate total price of all products
		var total float64
		err := db.QueryRow("SELECT SUM(price) FROM products").Scan(&total)
		if err != nil {
			http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
			return
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
