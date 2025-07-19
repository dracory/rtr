package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dracory/rtr"
)

func main() {
	// Example 1: Pure Declarative Approach
	fmt.Println("=== Pure Declarative Router ===")
	declarativeRouter := createDeclarativeRouter()
	declarativeRouter.List()

	// Example 2: Hybrid Approach (Declarative + Imperative)
	fmt.Println("\n=== Hybrid Router (Declarative + Imperative) ===")
	hybridRouter := createHybridRouter()
	hybridRouter.List()

	// Example 3: Converting Imperative to Declarative
	fmt.Println("\n=== Imperative Router ===")
	imperativeRouter := createImperativeRouter()
	imperativeRouter.List()

	// Start server with the declarative router
	port := ":8080"
	fmt.Printf("\nServer running on http://localhost%s\n", port)
	fmt.Println("Try these endpoints:")
	fmt.Println("  GET  /")
	fmt.Println("  GET  /api/users")
	fmt.Println("  POST /api/users")
	fmt.Println("  GET  /api/v1/products")
	fmt.Println("  GET  /admin (with Host: admin.example.com)")
	
	log.Fatal(http.ListenAndServe(port, declarativeRouter))
}

// createDeclarativeRouter demonstrates pure declarative configuration
func createDeclarativeRouter() rtr.RouterInterface {
	// Define middleware
	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("[%s] %s %s\n", r.Method, r.URL.Path, r.RemoteAddr)
			next.ServeHTTP(w, r)
		})
	}

	authMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Simple auth check
			if r.Header.Get("Authorization") == "" {
				w.Header().Set("X-Auth-Required", "true")
			}
			next.ServeHTTP(w, r)
		})
	}

	// Define handlers
	homeHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"message": "Welcome to Declarative Router!", "timestamp": "%s"}`, 
			r.Header.Get("X-Request-Time"))
	}

	usersHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"users": ["alice", "bob", "charlie"]}`)
	}

	createUserHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"message": "User created successfully"}`)
	}

	productsHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"products": ["laptop", "phone", "tablet"]}`)
	}

	adminHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"message": "Admin Dashboard", "domain": "%s"}`, r.Host)
	}

	// Declarative router configuration
	config := rtr.RouterConfig{
		Name: "Declarative API Router",
		BeforeMiddleware: []rtr.Middleware{
			loggingMiddleware,
		},
		Routes: []rtr.RouteConfig{
			rtr.GET("/", homeHandler).
				WithName("Home").
				WithMetadata("description", "Main landing page"),
		},
		Groups: []rtr.GroupConfig{
			rtr.Group("/api",
				// Direct routes in API group
				rtr.GET("/users", usersHandler).
					WithName("List Users").
					WithBeforeMiddleware(authMiddleware),
				rtr.POST("/users", createUserHandler).
					WithName("Create User").
					WithBeforeMiddleware(authMiddleware),
				
				// Nested version groups
				rtr.Group("/v1",
					rtr.GET("/products", productsHandler).
						WithName("List Products V1").
						WithMetadata("version", "1.0"),
				).WithName("API V1"),
				
				rtr.Group("/v2",
					rtr.GET("/products", productsHandler).
						WithName("List Products V2").
						WithMetadata("version", "2.0"),
				).WithName("API V2"),
			).WithName("API Group").
				WithBeforeMiddleware(func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-Type", "application/json")
						next.ServeHTTP(w, r)
					})
				}),
		},
		Domains: []rtr.DomainConfig{
			rtr.Domain([]string{"admin.example.com", "*.admin.example.com"},
				rtr.GET("/", adminHandler).WithName("Admin Home"),
				rtr.Group("/api",
					rtr.GET("/stats", func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprintf(w, `{"stats": {"users": 42, "requests": 1337}}`)
					}).WithName("Admin Stats"),
				).WithName("Admin API"),
			),
		},
	}

	return rtr.NewRouterFromConfig(config)
}

// createHybridRouter demonstrates mixing declarative and imperative approaches
func createHybridRouter() rtr.RouterInterface {
	// Start with declarative base
	config := rtr.RouterConfig{
		Name: "Hybrid Router",
		Routes: []rtr.RouteConfig{
			rtr.GET("/", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "Hybrid Router Home")
			}).WithName("Home"),
		},
		Groups: []rtr.GroupConfig{
			rtr.Group("/api",
				rtr.GET("/users", func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintf(w, `{"users": []}`)
				}).WithName("Users"),
			).WithName("API"),
		},
	}

	router := rtr.NewRouterFromConfig(config)

	// Add imperative routes on top
	router.AddRoute(rtr.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"status": "ok"}`)
	}).SetName("Health Check"))

	// Add imperative group
	adminGroup := rtr.NewGroup().SetPrefix("/admin")
	adminGroup.AddRoute(rtr.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Admin Dashboard")
	}).SetName("Admin Dashboard"))
	router.AddGroup(adminGroup)

	return router
}

// createImperativeRouter demonstrates traditional imperative approach for comparison
func createImperativeRouter() rtr.RouterInterface {
	router := rtr.NewRouter()

	// Add routes imperatively
	router.AddRoute(rtr.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Imperative Router Home")
	}).SetName("Home"))

	// Create API group
	api := rtr.NewGroup().SetPrefix("/api")
	api.AddRoute(rtr.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"users": []}`)
	}).SetName("Users"))

	// Create nested group
	v1 := rtr.NewGroup().SetPrefix("/v1")
	v1.AddRoute(rtr.Get("/products", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"products": []}`)
	}).SetName("Products"))
	api.AddGroup(v1)

	router.AddGroup(api)

	return router
}
