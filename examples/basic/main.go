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
