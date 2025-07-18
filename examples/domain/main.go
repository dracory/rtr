package main

import (
	"fmt"
	"log"
	"net/http"

	rtr "github.com/dracory/rtr"
)

func main() {
	// Create a new router
	r := rtr.NewRouter()

	// Create domains
	apiDomain := rtr.NewDomain("api.example.com", "localhost:8080")
	adminDomain := rtr.NewDomain("admin.example.com", "localhost:8081")

	// Add routes to the API domain
	apiDomain.AddRoute(rtr.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	}))

	apiDomain.AddRoute(rtr.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`["user1", "user2"]`))
	}))

	// Add routes to the admin domain
	adminDomain.AddRoute(rtr.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Admin Panel</title>
		</head>
		<body>
			<h1>Welcome to Admin Panel</h1>
			<p>This is the admin interface for example.com</p>
		</body>
		</html>
		`)
	}))

	// Add catch-all route for API domain
	apiDomain.AddRoute(rtr.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Not Found", "message": "The requested resource was not found on this server"}`))
	}))

	// Add catch-all route for Admin domain
	adminDomain.AddRoute(rtr.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>404 Not Found</title>
		</head>
		<body>
			<h1>404 Not Found</h1>
			<p>The requested page was not found on this server.</p>
			<p><a href="/">Return to Admin Panel</a></p>
		</body>
		</html>
		`)
	}))

	// Add domains to the router
	r.AddDomain(apiDomain)
	r.AddDomain(adminDomain)

	// Start the server
	fmt.Println("Starting server on :8080 (api.example.com) and :8081 (admin.example.com)")
	fmt.Println("To test, add the following to your /etc/hosts or C:\\Windows\\System32\\drivers\\etc\\hosts file:")
	fmt.Println("127.0.0.1 api.example.com admin.example.com")
	fmt.Println("Then visit:")
	fmt.Println("- http://api.example.com:8080/status")
	fmt.Println("- http://api.example.com:8080/users")
	fmt.Println("- http://admin.example.com:8081/")

	// Start the server on multiple ports
	go func() {
		log.Fatal(http.ListenAndServe(":8080", r))
	}()

	log.Fatal(http.ListenAndServe(":8081", r))
}
