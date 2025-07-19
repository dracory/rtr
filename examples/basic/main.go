package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dracory/rtr"
)

func main() {
	// Create a new router
	r := rtr.NewRouter()

	// Root route with web interface listing all available endpoints
	r.AddRoute(rtr.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head>
    <title>Basic Router Example</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; max-width: 800px; margin: 0 auto; padding: 20px; }
        h1 { color: #2c3e50; }
        .endpoint { 
            background: #f8f9fa; 
            border-left: 4px solid #3498db; 
            padding: 10px 15px; 
            margin: 10px 0; 
            border-radius: 0 4px 4px 0;
        }
        .endpoint a { 
            color: #3498db; 
            text-decoration: none; 
            font-weight: bold;
        }
        .endpoint a:hover { text-decoration: underline; }
        .description { color: #7f8c8d; margin: 5px 0 0 0; }
    </style>
</head>
<body>
    <h1>Basic Router Example</h1>
    <p>This example demonstrates basic routing functionality with groups and nested routes.</p>
    
    <div class="endpoint">
        <a href="/hello">GET /hello</a>
        <p class="description">Simple hello world endpoint</p>
    </div>
    
    <div class="endpoint">
        <a href="/api/status">GET /api/status</a>
        <p class="description">API status check endpoint</p>
    </div>
    
    <div class="endpoint">
        <a href="/api/users">GET /api/users</a>
        <p class="description">List of users endpoint</p>
    </div>
    
    <div class="endpoint">
        <a href="/api/users/123">GET /api/users/:id</a>
        <p class="description">Get specific user by ID</p>
    </div>
</body>
</html>`)
	}))

	// Add a simple route using shortcut method
	r.AddRoute(rtr.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	}))

	// Create an API group
	api := rtr.NewGroup().SetPrefix("/api")

	// Add routes to the API group
	api.AddRoute(rtr.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"status": "ok"}`)
	}))

	// Add the API group to the router
	r.AddGroup(api)

	// Create a users group with nested routes
	users := rtr.NewGroup().SetPrefix("/users")

	// Add user routes
	users.AddRoute(rtr.Get("", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "List of users")
	}))

	users.AddRoute(rtr.Get("/:id", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/users/"):]
		fmt.Fprintf(w, "User ID: %s", id)
	}))

	// Add the users group to the API group
	api.AddGroup(users)

	// Start the server
	port := ":8080"
	fmt.Printf("Server running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
