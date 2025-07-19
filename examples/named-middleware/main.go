package main

import (
	"fmt"
	"net/http"

	"github.com/dracory/rtr"
)

func main() {
	// Create named middleware similar to the old gouniverse/router pattern
	authMiddleware := rtr.NewMiddleware("User Authentication", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check for authentication token
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			
			fmt.Printf("Auth middleware '%s' executed for %s\n", "User Authentication", r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})

	loggingMiddleware := rtr.NewMiddleware("Request Logger", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("Logging middleware '%s' - %s %s\n", "Request Logger", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})

	corsMiddleware := rtr.NewMiddleware("CORS Handler", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			fmt.Printf("CORS middleware '%s' executed\n", "CORS Handler")
			next.ServeHTTP(w, r)
		})
	})

	// Create router with global named middleware
	router := rtr.NewRouter()

	// Add global middleware using the new named middleware approach
	globalMiddlewares := []rtr.MiddlewareInterface{
		loggingMiddleware,
		corsMiddleware,
	}

	// Add named middleware to router
	router.AddBeforeMiddlewares(globalMiddlewares)

	// Example 1: Route with named middleware (similar to old gouniverse/router)
	protectedRoute := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/api/protected").
		SetName("Protected API Endpoint").
		AddBeforeMiddlewares([]rtr.MiddlewareInterface{authMiddleware}).
		SetJSONHandler(func(w http.ResponseWriter, r *http.Request) string {
			return `{"message": "This is a protected endpoint", "user": "authenticated"}`
		})

	// Example 2: Public route without authentication
	publicRoute := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/api/public").
		SetName("Public API Endpoint").
		SetJSONHandler(func(w http.ResponseWriter, r *http.Request) string {
			return `{"message": "This is a public endpoint", "status": "ok"}`
		})

	// Example 3: Route with multiple named middleware
	adminRoute := rtr.NewRoute().
		SetMethod("POST").
		SetPath("/api/admin").
		SetName("Admin API Endpoint").
		AddBeforeMiddlewares([]rtr.MiddlewareInterface{
			authMiddleware,
			rtr.NewMiddleware("Admin Check", func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Check if user is admin
					role := r.Header.Get("X-User-Role")
					if role != "admin" {
						http.Error(w, "Forbidden: Admin access required", http.StatusForbidden)
						return
					}
					fmt.Printf("Admin middleware executed for %s\n", r.URL.Path)
					next.ServeHTTP(w, r)
				})
			}),
		}).
		SetJSONHandler(func(w http.ResponseWriter, r *http.Request) string {
			return `{"message": "Admin endpoint accessed successfully"}`
		})

	// Add routes to router
	router.AddRoute(protectedRoute)
	router.AddRoute(publicRoute)
	router.AddRoute(adminRoute)

	// Example 4: Using RouteConfig with named middleware (declarative approach)
	configRoute := &rtr.RouteConfig{
		Name:   "Config Route",
		Method: "GET",
		Path:   "/api/config",
		JSONHandler: func(w http.ResponseWriter, r *http.Request) string {
			return `{"config": "loaded", "version": "1.0"}`
		},
	}

	// Add named middleware to config route
	configRoute.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
		rtr.NewMiddleware("Config Validator", func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Printf("Config validation middleware executed\n")
				next.ServeHTTP(w, r)
			})
		}),
	})

	router.AddRoute(configRoute)

	// Example 5: Group with named middleware
	apiGroup := rtr.NewGroup().
		SetPrefix("/v1").
		AddBeforeMiddlewares([]rtr.MiddlewareInterface{
			rtr.NewMiddleware("API Version Check", func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("X-API-Version", "v1")
					fmt.Printf("API version middleware executed\n")
					next.ServeHTTP(w, r)
				})
			}),
		})

	// Add routes to the group
	apiGroup.AddRoute(rtr.NewRoute().
		SetMethod("GET").
		SetPath("/users").
		SetJSONHandler(func(w http.ResponseWriter, r *http.Request) string {
			return `{"users": []}`
		}))

	router.AddGroup(apiGroup)

	// Display middleware information
	fmt.Println("=== Named Middleware Example ===")
	fmt.Printf("Auth Middleware: %s\n", authMiddleware.GetName())
	fmt.Printf("Logging Middleware: %s\n", loggingMiddleware.GetName())
	fmt.Printf("CORS Middleware: %s\n", corsMiddleware.GetName())
	fmt.Println()

	// Display route information with their named middleware
	fmt.Println("Routes with Named Middleware:")
	fmt.Printf("- %s: %d named before middlewares\n", 
		protectedRoute.GetName(), len(protectedRoute.GetBeforeMiddlewares()))
	fmt.Printf("- %s: %d named before middlewares\n", 
		adminRoute.GetName(), len(adminRoute.GetBeforeMiddlewares()))
	fmt.Printf("- %s: %d named before middlewares\n", 
		configRoute.GetName(), len(configRoute.GetBeforeMiddlewares()))
	fmt.Println()

	// List all middleware names for the admin route
	fmt.Println("Admin route middleware chain:")
	for i, mw := range adminRoute.GetBeforeMiddlewares() {
		fmt.Printf("  %d. %s\n", i+1, mw.GetName())
	}
	fmt.Println()

	fmt.Println("Starting server on :8080...")
	fmt.Println("Try these endpoints:")
	fmt.Println("  GET  /api/public                    (no auth required)")
	fmt.Println("  GET  /api/protected                 (requires Authorization header)")
	fmt.Println("  POST /api/admin                     (requires Authorization + X-User-Role: admin)")
	fmt.Println("  GET  /api/config                    (no auth required)")
	fmt.Println("  GET  /v1/users                      (API v1 endpoint)")
	fmt.Println()
	fmt.Println("Example requests:")
	fmt.Println("  curl http://localhost:8080/api/public")
	fmt.Println("  curl -H 'Authorization: Bearer token123' http://localhost:8080/api/protected")
	fmt.Println("  curl -X POST -H 'Authorization: Bearer token123' -H 'X-User-Role: admin' http://localhost:8080/api/admin")

	// Start the server
	http.ListenAndServe(":8080", router)
}
