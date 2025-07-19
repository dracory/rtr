# Declarative Router Configuration

The rtr router supports a powerful declarative configuration approach that allows you to define your entire routing structure using configuration objects. This approach is particularly useful for:

- Large applications with complex routing requirements
- Applications where routes need to be loaded from configuration files
- Cases where you want to separate route definition from implementation
- Tools that generate or modify routing configuration

## Table of Contents

- [Basic Concepts](#basic-concepts)
- [Configuration Structure](#configuration-structure)
  - [Route](#route-configuration)
  - [Group](#group-configuration)
  - [Domain](#domain-configuration)
- [Handler Registry](#handler-registry)
- [Middleware](#middleware)
- [Complete Example](#complete-example)
- [Best Practices](#best-practices)

## Basic Concepts

1. **Configuration as Code**: Define your routes, groups, and domains using Go structs
2. **Separation of Concerns**: Keep route definitions separate from handler implementations
3. **Reusability**: Share common configuration across different environments
4. **Runtime Flexibility**: Enable/disable routes and middleware at runtime

## Configuration Structure

### Route Configuration

Routes define individual HTTP endpoints with their handlers and middleware:

```go
type Route struct {
    Status       string   // Status of the route (enabled/disabled)
    Path         string   // URL path pattern
    Method       string   // HTTP method (GET, POST, etc.)
    Name         string   // Route name for reference
    Middlewares  []string // List of middleware names to apply
    
    // Handler references (only one should be set per route)
    Handler      string   // Standard HTTP handler reference
    HTMLHandler  string   // HTML handler reference
    JSONHandler  string   // JSON handler reference
    CSSHandler   string   // CSS handler reference
    XMLHandler   string   // XML handler reference
    TextHandler  string   // Plain text handler reference
    JSHandler    string   // JavaScript handler reference
    ErrorHandler string   // Error handler reference
}
```

### Group Configuration

Groups allow you to apply common configuration to multiple routes:

```go
type Group struct {
    Status      string   // Status of the group (enabled/disabled)
    Prefix      string   // URL prefix for all routes in this group
    Name        string   // Group name for reference
    Middlewares []string // Middleware to apply to all routes in this group
    Routes      []Route  // Nested routes
}
```

### Domain Configuration

Domains allow you to handle different hostnames with separate routing rules:

```go
type Domain struct {
    Status      string          // Status of the domain (enabled/disabled)
    Hosts       []string        // List of hostnames this domain handles
    Name        string          // Domain name for reference
    Middlewares []string        // Middleware to apply to all routes in this domain
    Items       []ItemInterface // Can contain both Groups and Routes
}
```

## Handler Registry

The `HandlerRegistry` is used to register and look up handlers by name:

```go
// Create a new registry
registry := rtr.NewHandlerRegistry()

// Register a standard HTTP handler
registry.AddRoute(rtr.NewRoute().
    SetName("home").
    SetHandler(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    }))

// Register an HTML handler (returns string)
registry.AddRoute(rtr.NewRoute().
    SetName("about").
    SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
        return "<h1>About Us</h1><p>Welcome to our website!</p>"
    }))

// Register a JSON handler (returns string)
registry.AddRoute(rtr.NewRoute().
    SetName("api_status").
    SetJSONHandler(func(w http.ResponseWriter, r *http.Request) string {
        return `{"status":"ok","version":"1.0.0"}`
    }))
```

## Middleware

Middleware can be registered and applied at different levels (domain, group, or route):

```go
// Register middleware
registry.AddMiddleware(rtr.NewMiddleware("logger", func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
    })
}))

// Register authentication middleware
registry.AddMiddleware(rtr.NewMiddleware("auth", func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !isAuthenticated(r) {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}))
```

## Complete Example

Here's a complete example demonstrating the declarative configuration approach:

```go
package main

import (
	"fmt"
	"net/http"
	"github.com/dracory/rtr"
)

func main() {
	// Create registry
	registry := rtr.NewHandlerRegistry()

	// Register handlers
	handlers := map[string]rtr.StdHandler{
		"home": func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("<h1>Welcome Home</h1>"))
		},
		"about": func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("<h1>About Us</h1>"))
		},
		"api_status": func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status":"ok"}`))
		},
	}

	// Register all handlers
	for name, handler := range handlers {
		registry.AddRoute(rtr.NewRoute().
			SetName(name).
			SetHandler(handler))
	}

	// Register middleware
	registry.AddMiddleware(rtr.NewMiddleware("logger", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("Request: %s %s\n", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}))

	// Define domain configuration
	domain := &rtr.Domain{
		Name:   "Example Domain",
		Status: rtr.StatusEnabled,
		Hosts:  []string{"example.com"},
		Middlewares: []string{"logger"},
		Items: []rtr.ItemInterface{
			// Direct route
			&rtr.Route{
				Name:    "Home",
				Method:  "GET",
				Path:    "/",
				Handler: "home",
				Status:  rtr.StatusEnabled,
			},
			// Group with common prefix and middleware
			&rtr.Group{
				Name:   "API",
				Prefix: "/api",
				Status: rtr.StatusEnabled,
				Middlewares: []string{"auth"},
				Routes: []rtr.Route{
					{
						Name:    "API Status",
						Method:  "GET",
						Path:    "/status",
						Handler: "api_status",
						Status:  rtr.StatusEnabled,
					},
				},
			},
		},
	}

	// In a real application, you would use this configuration to build your router
	// and start the HTTP server
	fmt.Printf("Domain configuration created: %s\n", domain.Name)
	fmt.Printf("Registered %d routes and %d middleware\n", 
		len(registry.FindAllRoutes()),
		len(registry.FindAllMiddlewares()),
	)
}
```

## Best Practices

1. **Organize by Feature**: Group related routes together using the `Group` type
2. **Use Meaningful Names**: Choose clear, descriptive names for routes and middleware
3. **Leverage Middleware**: Use middleware for cross-cutting concerns like logging and auth
4. **Enable/Disable Features**: Use the `Status` field to toggle features without code changes
5. **Validate Early**: Validate your configuration before using it in production
6. **Document Thoroughly**: Document your routes, groups, and middleware for future maintainers
7. **Test Configurations**: Write tests that verify your routing configuration works as expected

For more examples and advanced usage, see the [declarative example](../examples/declarative) in the examples directory.
