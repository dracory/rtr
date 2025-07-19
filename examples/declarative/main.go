package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dracory/rtr"
)

func main() {
	// Example 1: Pure Declarative Approach
	fmt.Println("=== Pure Declarative Router ===")
	declarativeRouter := CreateDeclarativeRouter()
	declarativeRouter.List()

	// Example 2: Hybrid Approach (Declarative + Imperative)
	fmt.Println("\n=== Hybrid Router (Declarative + Imperative) ===")
	hybridRouter := CreateHybridRouter()
	hybridRouter.List()

	// Example 3: Converting Imperative to Declarative
	fmt.Println("\n=== Imperative Router ===")
	imperativeRouter := CreateImperativeRouter()
	imperativeRouter.List()

	// Start server with the declarative router
	port := ":8080"
	fmt.Printf("\nServer running on http://localhost%s\n", port)
	fmt.Println("Try these endpoints:")
	fmt.Println("\n=== Main Endpoints ===")
	fmt.Println("  GET  /                    - Home page")
	fmt.Println("  GET  /health              - Health check")
	fmt.Println("  GET  /error               - Error demo")
	fmt.Println("\n=== API Endpoints ===")
	fmt.Println("  GET  /api/users           - List users (requires auth header)")
	fmt.Println("  POST /api/users           - Create user (requires auth header)")
	fmt.Println("  GET  /api/v1/products     - List products v1")
	fmt.Println("  GET  /api/v2/products     - List products v2")
	fmt.Println("\n=== Admin Endpoints (Host: admin.example.com) ===")
	fmt.Println("  GET  /                    - Admin dashboard")
	fmt.Println("  GET  /api/stats           - System stats")
	fmt.Println("  GET  /api/users           - User stats")
	fmt.Println("\n=== Testing Commands ===")
	fmt.Println("  curl http://localhost:8080/")
	fmt.Println("  curl -H 'Authorization: Bearer token' http://localhost:8080/api/users")
	fmt.Println("  curl -H 'Host: admin.example.com' http://localhost:8080/")
	fmt.Println("  curl -H 'Accept: text/html' -H 'Host: admin.example.com' http://localhost:8080/")
	
	log.Fatal(http.ListenAndServe(port, declarativeRouter))
}

// CreateDeclarativeRouter demonstrates pure declarative configuration
func CreateDeclarativeRouter() rtr.RouterInterface {
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

	// Define different types of handlers to showcase router capabilities
	
	// Standard HTTP handler
	homeHandler := func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"message":   "Welcome to Declarative Router!",
			"timestamp": time.Now().Format(time.RFC3339),
			"method":    r.Method,
			"path":      r.URL.Path,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}

	// JSON handler - returns JSON string (would use JSONHandler if available)
	usersHandler := func(w http.ResponseWriter, r *http.Request) {
		users := []map[string]interface{}{
			{"id": 1, "name": "alice", "email": "alice@example.com"},
			{"id": 2, "name": "bob", "email": "bob@example.com"},
			{"id": 3, "name": "charlie", "email": "charlie@example.com"},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"users": users})
	}

	// Create user handler with proper error handling
	createUserHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
			return
		}
		
		response := map[string]interface{}{
			"message": "User created successfully",
			"id":      42,
			"created": time.Now().Format(time.RFC3339),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}

	// Products handler with version-aware responses
	productsHandler := func(w http.ResponseWriter, r *http.Request) {
		products := []map[string]interface{}{
			{"id": 1, "name": "laptop", "price": 999.99, "category": "electronics"},
			{"id": 2, "name": "phone", "price": 699.99, "category": "electronics"},
			{"id": 3, "name": "tablet", "price": 399.99, "category": "electronics"},
		}
		
		response := map[string]interface{}{
			"products": products,
			"total":    len(products),
			"version":  r.Header.Get("X-API-Version"),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}

	// Admin handler with HTML response
	adminHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") == "text/html" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintf(w, `<!DOCTYPE html>
			<html><head><title>Admin Dashboard</title></head>
			<body><h1>Admin Dashboard</h1><p>Domain: %s</p></body></html>`, r.Host)
		} else {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Admin Dashboard",
				"domain":  r.Host,
				"time":    time.Now().Format(time.RFC3339),
			})
		}
	}

	// Health check handler
	healthHandler := func(w http.ResponseWriter, r *http.Request) {
		health := map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
			"uptime":    "24h", // This would be calculated in real app
			"version":   "1.0.0",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(health)
	}

	// Error handler for demonstration
	errorHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Internal server error",
			"code":  "500",
		})
	}

	// Declarative router configuration
	config := rtr.RouterConfig{
		Name: "Comprehensive Declarative API Router",
		BeforeMiddleware: []rtr.Middleware{
			loggingMiddleware,
		},
		Routes: []rtr.RouteConfig{
			rtr.GET("/", homeHandler).
				WithName("Home").
				WithMetadata("description", "Main landing page").
				WithMetadata("public", "true"),
			
			rtr.GET("/health", healthHandler).
				WithName("Health Check").
				WithMetadata("description", "Service health status").
				WithMetadata("monitoring", "true"),
				
			rtr.GET("/error", errorHandler).
				WithName("Error Demo").
				WithMetadata("description", "Demonstrates error handling"),
		},
		Groups: []rtr.GroupConfig{
			rtr.Group("/api",
				// Direct routes in API group
				rtr.GET("/users", usersHandler).
					WithName("List Users").
					WithMetadata("description", "Get all users").
					WithMetadata("auth_required", "true").
					WithBeforeMiddleware(authMiddleware),
					
				rtr.POST("/users", createUserHandler).
					WithName("Create User").
					WithMetadata("description", "Create a new user").
					WithMetadata("auth_required", "true").
					WithBeforeMiddleware(authMiddleware),
				
				// Nested version groups
				rtr.Group("/v1",
					rtr.GET("/products", productsHandler).
						WithName("List Products V1").
						WithMetadata("version", "1.0").
						WithMetadata("description", "Get products (v1)").
						WithBeforeMiddleware(func(next http.Handler) http.Handler {
							return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
								r.Header.Set("X-API-Version", "1.0")
								next.ServeHTTP(w, r)
							})
						}),
				).WithName("API V1"),
				
				rtr.Group("/v2",
					rtr.GET("/products", productsHandler).
						WithName("List Products V2").
						WithMetadata("version", "2.0").
						WithMetadata("description", "Get products (v2)").
						WithBeforeMiddleware(func(next http.Handler) http.Handler {
							return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
								r.Header.Set("X-API-Version", "2.0")
								next.ServeHTTP(w, r)
							})
						}),
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
				rtr.GET("/", adminHandler).
					WithName("Admin Home").
					WithMetadata("description", "Admin dashboard home").
					WithMetadata("auth_level", "admin"),
					
				rtr.Group("/api",
					rtr.GET("/stats", func(w http.ResponseWriter, r *http.Request) {
						stats := map[string]interface{}{
							"users":         42,
							"requests":      1337,
							"uptime":        "24h",
							"last_updated": time.Now().Format(time.RFC3339),
						}
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(map[string]interface{}{"stats": stats})
					}).
						WithName("Admin Stats").
						WithMetadata("description", "System statistics").
						WithMetadata("auth_level", "admin"),
						
					rtr.GET("/users", func(w http.ResponseWriter, r *http.Request) {
						users := map[string]interface{}{
							"total":        42,
							"active":       38,
							"last_24h":     5,
							"last_updated": time.Now().Format(time.RFC3339),
						}
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(map[string]interface{}{"users": users})
					}).
						WithName("Admin User Stats").
						WithMetadata("description", "User statistics").
						WithMetadata("auth_level", "admin"),
				).WithName("Admin API").
				  WithBeforeMiddleware(func(next http.Handler) http.Handler {
					  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						  w.Header().Set("X-Admin-API", "true")
						  next.ServeHTTP(w, r)
					  })
				  }),
			),
		},
	}

	return rtr.NewRouterFromConfig(config)
}

// CreateHybridRouter demonstrates mixing declarative and imperative approaches
func CreateHybridRouter() rtr.RouterInterface {
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

// CreateImperativeRouter demonstrates traditional imperative approach for comparison
func CreateImperativeRouter() rtr.RouterInterface {
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
