package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dracory/rtr"
)

func main() {
	// Example: Pure Declarative Configuration
	fmt.Println("=== Declarative Routing System Example ===")
	domain := CreateDeclarativeConfiguration()

	// Display the configuration
	fmt.Printf("Domain: %s\n", domain.Name)
	fmt.Printf("Status: %s\n", domain.Status)
	fmt.Printf("Hosts: %v\n", domain.Hosts)
	fmt.Printf("Items: %d\n", len(domain.Items))
	fmt.Printf("Middlewares: %v\n", domain.Middlewares)

	// Create a HandlerRegistry and register handlers
	registry := rtr.NewHandlerRegistry()
	registerHandlers(registry)
	registerMiddleware(registry)

	fmt.Println("\n=== Configuration Structure ===")
	printDomainStructure(domain)

	fmt.Println("\n=== JSON Serialization Example ===")
	showJSONSerialization(domain)

	fmt.Println("\n=== HandlerRegistry Example ===")
	showHandlerRegistryExample(registry)

	fmt.Println("\n=== Runtime Router Creation ===")
	showRuntimeRouterCreation(domain, registry)
}

// CreateDeclarativeConfiguration demonstrates pure declarative configuration
func CreateDeclarativeConfiguration() *rtr.Domain {
	return &rtr.Domain{
		Name:        "Example API",
		Status:      rtr.StatusEnabled,
		Hosts:       []string{"api.example.com", "*.api.example.com"},
		Middlewares: []string{"cors", "logging"},
		Items: []rtr.ItemInterface{
			// Home route
			&rtr.Route{
				Name:    "Home",
				Status:  rtr.StatusEnabled,
				Method:  rtr.MethodGET,
				Path:    "/",
				Handler: "home",
			},
			// API Group
			&rtr.Group{
				Name:        "API v1",
				Status:      rtr.StatusEnabled,
				Prefix:      "/api/v1",
				Middlewares: []string{"auth", "rate-limit"},
				Routes: []rtr.Route{
					{
						Name:    "List Users",
						Status:  rtr.StatusEnabled,
						Method:  rtr.MethodGET,
						Path:    "/users",
						Handler: "users-list",
					},
					{
						Name:        "Create User",
						Status:      rtr.StatusEnabled,
						Method:      rtr.MethodPOST,
						Path:        "/users",
						Handler:     "users-create",
						Middlewares: []string{"validate-user"},
					},
					{
						Name:    "Get User",
						Status:  rtr.StatusEnabled,
						Method:  rtr.MethodGET,
						Path:    "/users/{id}",
						Handler: "users-get",
					},
				},
			},
			// Admin routes (disabled for maintenance)
			&rtr.Route{
				Name:        "Admin Dashboard",
				Status:      rtr.StatusDisabled,
				Method:      rtr.MethodGET,
				Path:        "/admin",
				Handler:     "admin-dashboard",
				Middlewares: []string{"admin-auth"},
			},
		},
	}
}

// registerHandlers registers all route handlers with the registry
func registerHandlers(registry *rtr.HandlerRegistry) {
	// Home handler using HTMLHandler
	homeRoute := rtr.NewRoute().
		SetName("home").
		SetPath("/").
		SetMethod("GET").
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			return fmt.Sprintf(`<html><body>
				<h1>Welcome to Declarative Router!</h1>
				<p>Current time: %s</p>
				<p>Method: %s</p>
				<p>Path: %s</p>
			</body></html>`, time.Now().Format(time.RFC3339), r.Method, r.URL.Path)
		})
	registry.AddRoute(homeRoute)

	// Users list handler using JSONHandler
	usersListRoute := rtr.NewRoute().
		SetName("users-list").
		SetPath("/users").
		SetMethod("GET").
		SetJSONHandler(func(w http.ResponseWriter, r *http.Request) string {
			users := []map[string]interface{}{
				{"id": 1, "name": "alice", "email": "alice@example.com"},
				{"id": 2, "name": "bob", "email": "bob@example.com"},
				{"id": 3, "name": "charlie", "email": "charlie@example.com"},
			}
			data, _ := json.Marshal(map[string]interface{}{"users": users})
			return string(data)
		})
	registry.AddRoute(usersListRoute)

	// Create user handler using standard Handler
	usersCreateRoute := rtr.NewRoute().
		SetName("users-create").
		SetPath("/users").
		SetMethod("POST").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
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
		})
	registry.AddRoute(usersCreateRoute)

	// Get user handler using JSONHandler
	usersGetRoute := rtr.NewRoute().
		SetName("users-get").
		SetPath("/users/{id}").
		SetMethod("GET").
		SetJSONHandler(func(w http.ResponseWriter, r *http.Request) string {
			// In a real implementation, you'd extract the ID from the path
			userID := "42" // Placeholder
			user := map[string]interface{}{
				"id":    userID,
				"name":  "John Doe",
				"email": "john@example.com",
			}
			data, _ := json.Marshal(user)
			return string(data)
		})
	registry.AddRoute(usersGetRoute)

	// Admin dashboard handler using HTMLHandler
	adminRoute := rtr.NewRoute().
		SetName("admin-dashboard").
		SetPath("/admin").
		SetMethod("GET").
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			return fmt.Sprintf(`<html><body>
				<h1>Admin Dashboard</h1>
				<p>Domain: %s</p>
				<p>Time: %s</p>
				<p><strong>Status:</strong> Maintenance Mode</p>
			</body></html>`, r.Host, time.Now().Format(time.RFC3339))
		})
	registry.AddRoute(adminRoute)
}

// registerMiddleware registers all middleware with the registry
func registerMiddleware(registry *rtr.HandlerRegistry) {
	// CORS middleware
	corsMiddleware := rtr.NewMiddleware(
		rtr.WithName("cors"),
		rtr.WithHandler(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				next.ServeHTTP(w, r)
			})
		}))
	registry.AddMiddleware(corsMiddleware)

	// Logging middleware
	loggingMiddleware := rtr.NewMiddleware(
		rtr.WithName("logging"),
		rtr.WithHandler(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()
				fmt.Printf("[%s] %s %s\n", start.Format(time.RFC3339), r.Method, r.URL.Path)
				next.ServeHTTP(w, r)
				fmt.Printf("[%s] Completed in %v\n", time.Now().Format(time.RFC3339), time.Since(start))
			})
		}),
	)
	registry.AddMiddleware(loggingMiddleware)

	// Auth middleware
	authMiddleware := rtr.NewMiddleware(
		rtr.WithName("auth"),
		rtr.WithHandler(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				auth := r.Header.Get("Authorization")
				if auth == "" {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(map[string]string{"error": "Authorization required"})
					return
				}
				next.ServeHTTP(w, r)
			})
		}),
	)
	registry.AddMiddleware(authMiddleware)

	// Rate limiting middleware
	rateLimitMiddleware := rtr.NewMiddleware(
		rtr.WithName("rate-limit"),
		rtr.WithHandler(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Simple rate limiting simulation
				w.Header().Set("X-RateLimit-Limit", "100")
				w.Header().Set("X-RateLimit-Remaining", "99")
				next.ServeHTTP(w, r)
			})
		}),
	)
	registry.AddMiddleware(rateLimitMiddleware)

	// User validation middleware
	validateUserMiddleware := rtr.NewMiddleware(
		rtr.WithName("validate-user"),
		rtr.WithHandler(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Content-Type") != "application/json" {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]string{"error": "Content-Type must be application/json"})
					return
				}
				next.ServeHTTP(w, r)
			})
		}),
	)
	registry.AddMiddleware(validateUserMiddleware)

	// Admin auth middleware
	adminAuthMiddleware := rtr.NewMiddleware(
		rtr.WithName("admin-auth"),
		rtr.WithHandler(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				adminToken := r.Header.Get("X-Admin-Token")
				if adminToken != "admin-secret-token" {
					w.WriteHeader(http.StatusForbidden)
					json.NewEncoder(w).Encode(map[string]string{"error": "Admin access required"})
					return
				}
				next.ServeHTTP(w, r)
			})
		}),
	)
	registry.AddMiddleware(adminAuthMiddleware)
}

// printDomainStructure prints the domain configuration structure
func printDomainStructure(domain *rtr.Domain) {
	fmt.Printf("Domain: %s (%s)\n", domain.Name, domain.Status)
	fmt.Printf("  Hosts: %v\n", domain.Hosts)
	fmt.Printf("  Middlewares: %v\n", domain.Middlewares)
	fmt.Printf("  Items: %d\n", len(domain.Items))

	for i, item := range domain.Items {
		switch v := item.(type) {
		case *rtr.Route:
			fmt.Printf("  [%d] Route: %s (%s) %s %s -> %s\n", i, v.Name, v.Status, v.Method, v.Path, v.Handler)
			if len(v.Middlewares) > 0 {
				fmt.Printf("      Middlewares: %v\n", v.Middlewares)
			}
		case *rtr.Group:
			fmt.Printf("  [%d] Group: %s (%s) %s\n", i, v.Name, v.Status, v.Prefix)
			if len(v.Middlewares) > 0 {
				fmt.Printf("      Middlewares: %v\n", v.Middlewares)
			}
			fmt.Printf("      Routes: %d\n", len(v.Routes))
			for j, route := range v.Routes {
				fmt.Printf("        [%d] %s (%s) %s %s -> %s\n", j, route.Name, route.Status, route.Method, route.Path, route.Handler)
				if len(route.Middlewares) > 0 {
					fmt.Printf("            Middlewares: %v\n", route.Middlewares)
				}
			}
		}
	}
}

// showJSONSerialization demonstrates JSON serialization of the configuration
func showJSONSerialization(domain *rtr.Domain) {
	data, err := json.MarshalIndent(domain, "", "  ")
	if err != nil {
		fmt.Printf("Error serializing domain: %v\n", err)
		return
	}
	fmt.Printf("JSON Configuration:\n%s\n", string(data))
}

// showHandlerRegistryExample demonstrates the HandlerRegistry functionality
func showHandlerRegistryExample(registry *rtr.HandlerRegistry) {
	fmt.Println("Handler Registry Contents:")

	// Try to find registered handlers
	handlerNames := []string{"home", "users-list", "users-create", "users-get", "admin-dashboard"}
	for _, name := range handlerNames {
		handler := registry.FindRoute(name)
		if handler != nil {
			fmt.Printf("  ✓ Route '%s': %s %s\n", name, handler.GetMethod(), handler.GetPath())
		} else {
			fmt.Printf("  ✗ Route '%s': not found\n", name)
		}
	}

	// Try to find registered middleware
	middlewareNames := []string{"cors", "logging", "auth", "rate-limit", "validate-user", "admin-auth"}
	for _, name := range middlewareNames {
		middleware := registry.FindMiddleware(name)
		if middleware != nil {
			fmt.Printf("  ✓ Middleware '%s': %s\n", name, middleware.GetName())
		} else {
			fmt.Printf("  ✗ Middleware '%s': not found\n", name)
		}
	}
}

// showRuntimeRouterCreation demonstrates how to build a runtime router from declarative config
func showRuntimeRouterCreation(domain *rtr.Domain, registry *rtr.HandlerRegistry) {
	fmt.Println("Runtime Router Creation Process:")
	fmt.Printf("1. Domain configuration loaded: %s\n", domain.Name)
	fmt.Printf("2. Handler registry populated with %d routes and middleware\n", len(domain.Items))
	fmt.Println("3. Building runtime router...")

	// In a real implementation, you would:
	// - Create a new router instance
	// - Iterate through domain.Items
	// - For each route, find the handler in the registry and add it to the router
	// - For each group, create a group and add its routes
	// - Apply middleware based on the middleware names

	fmt.Println("\nSteps to build runtime router:")
	fmt.Println("  a) Create new router instance")
	fmt.Println("  b) Apply domain-level middleware:")
	for _, mw := range domain.Middlewares {
		fmt.Printf("     - Apply middleware: %s\n", mw)
	}

	fmt.Println("  c) Process domain items:")
	for i, item := range domain.Items {
		switch v := item.(type) {
		case *rtr.Route:
			if v.Status == rtr.StatusEnabled {
				fmt.Printf("     [%d] Add route: %s %s -> %s\n", i, v.Method, v.Path, v.Handler)
				for _, mw := range v.Middlewares {
					fmt.Printf("         - Apply middleware: %s\n", mw)
				}
			} else {
				fmt.Printf("     [%d] Skip disabled route: %s\n", i, v.Name)
			}
		case *rtr.Group:
			if v.Status == rtr.StatusEnabled {
				fmt.Printf("     [%d] Create group: %s with prefix %s\n", i, v.Name, v.Prefix)
				for _, mw := range v.Middlewares {
					fmt.Printf("         - Apply group middleware: %s\n", mw)
				}
				for j, route := range v.Routes {
					if route.Status == rtr.StatusEnabled {
						fmt.Printf("         [%d] Add group route: %s %s -> %s\n", j, route.Method, route.Path, route.Handler)
						for _, mw := range route.Middlewares {
							fmt.Printf("             - Apply route middleware: %s\n", mw)
						}
					} else {
						fmt.Printf("         [%d] Skip disabled group route: %s\n", j, route.Name)
					}
				}
			} else {
				fmt.Printf("     [%d] Skip disabled group: %s\n", i, v.Name)
			}
		}
	}

	fmt.Println("  d) Router ready for HTTP server")
	fmt.Println("\nNote: This is a demonstration. In a real implementation,")
	fmt.Println("you would use the registry to resolve handler names to actual functions.")
}
