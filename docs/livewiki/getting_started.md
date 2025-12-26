---
path: getting_started.md
page-type: tutorial
summary: Complete guide to installing, configuring, and using the RTR router for Go applications.
tags: [tutorial, installation, setup, quickstart]
created: 2025-12-26
updated: 2025-12-26
version: 1.0.0
---

# Getting Started with RTR Router

This guide will help you get up and running with RTR Router in your Go applications.

## Installation

### Requirements

- Go 1.25 or higher
- Compatible with Go modules

### Install via Go Modules

Add RTR to your project using Go modules:

```bash
go get github.com/dracory/rtr
```

### Import in Your Code

```go
import "github.com/dracory/rtr"
```

## Basic Usage

### Simple HTTP Server

Create a basic HTTP server with a single route:

```go
package main

import (
    "net/http"
    "github.com/dracory/rtr"
)

func main() {
    // Create a new router
    router := rtr.NewRouter()
    
    // Add a simple route
    router.AddRoute(rtr.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    }))
    
    // Start the server
    http.ListenAndServe(":8080", router)
}
```

### Run the Example

```bash
go run main.go
```

Now visit `http://localhost:8080` in your browser to see "Hello, World!".

## Route Types

### HTTP Method Shortcuts

RTR provides convenient shortcuts for common HTTP methods:

```go
router := rtr.NewRouter()

// GET route
router.AddRoute(rtr.Get("/users", handleGetUsers))

// POST route
router.AddRoute(rtr.Post("/users", handleCreateUser))

// PUT route
router.AddRoute(rtr.Put("/users/:id", handleUpdateUser))

// DELETE route
router.AddRoute(rtr.Delete("/users/:id", handleDeleteUser))
```

### Method Chaining

Alternatively, use method chaining for more control:

```go
router.AddRoute(rtr.NewRoute().
    SetMethod("GET").
    SetPath("/users").
    SetHandler(handleGetUsers).
    SetName  Name("ENTE  // Optional: name the route
)
```

## Handler Types

### Standard Handler

Full control over the HTTP response:

```go
router.AddRoute(rtr.Get("/standard", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Standard handler response"))
}))
```

### JSON Handler

Automatic JSON content-type handling:

```go
router.AddRoute(rtr.GetJSON("/api/status", func(w http.ResponseWriter, r *http.Request) string {
    return `{"status": "ok", "version": "1.0.0"}`
}))
```

### HTML Handler

Automatic HTML content-type handling:

```go
router.AddRoute(rtr.GetHTML("/page", func(w http.ResponseWriter, r *http.Request) string {
    return `<!DOCTYPE html>
<html>
<head><title>My Page</title></head>
<body><h1>Hello from RTR!</h1></body>
</html>`
}))
```

### Other Handler Types

RTR also supports CSS, XML, Text, and JavaScript handlers:

```go
// CSS Handler
router.AddRoute(rtr.GetCSS("/style.css", func(w http.ResponseWriter, r *http.Request) string {
    return "body { font-family: Arial; }"
}))

// XML Handler
router.AddRoute(rtr.GetXML("/data.xml", func(w http.ResponseWriter, r *http.Request) string {
    return `<?xml version="1.0"?><data><message>Hello XML</message></data>`
}))

// Text Handler
router.AddRoute(rtr.GetText("/robots.txt", func(w http.ResponseWriter, r *http.Request) string {
    return "User-agent: *\nAllow: /"
}))
```

## Path Parameters

### Basic Parameters

Extract values from URL paths using `:param` syntax:

```go
router.AddRoute(rtr.Get("/users/:id", func(w http.ResponseWriter, r *http.Request) {
    userID := rtr.MustGetParam(r, "id")
    w.Write([]byte("User ID: " + userID))
}))
```

### Optional Parameters

Mark parameters as optional with `?`:

```go
// Matches both /articles/tech and /articles/tech/123
router.AddRoute(rtr.Get("/articles/:category/:id?", func(w http.ResponseWriter, r *http.Request) {
    category := rtr.MustGetParam(r, "category")
    id, hasID := rtr.GetParam(r, "id")
    
    if hasID {
        w.Write([]byte("Category: " + category + ", ID: " + id))
    } else {
        w.Write([]byte("Category: " + category))
    }
}))
```

### Greedy Parameters

Capture the remainder of the path with `:param...`:

```go
// Matches /files/images/photo.jpg, /files/user/docs/file.pdf, etc.
router.AddRoute(rtr.Get("/files/:path...", func(w http.ResponseWriter, r *http.Request) {
    filePath := rtr.MustGetParam(r, "path")
    w.Write([]byte("File path: " + filePath))
}))
```

## Route Groups

Organize related routes with shared prefixes and middleware:

```go
// Create an API group
apiGroup := rtr.NewGroup().SetPrefix("/api/v1")

// Add routes to the group
apiGroup.AddRoute(rtr.Get("/users", handleUsers))
apiGroup.AddRoute(rtr.Post("/users", handleCreateUser))
apiGroup.AddRoute(rtr.Get("/products", handleProducts))

// Add the group to the router
router.AddGroup(apiGroup)
```

## Middleware

### Global Middleware

Add middleware that runs for all routes:

```go
// Add logging middleware
router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    rtr.NewAnonymousMiddleware(loggingMiddleware),
})

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("Request: %s %s\n", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}
```

### Group-Level Middleware

Add middleware to specific groups:

```go
apiGroup.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    rtr.NewAnonymousMiddleware(apiKeyMiddleware),
})
```

### Route-Level Middleware

Add middleware to individual routes:

```go
router.AddRoute(rtr.Get("/admin", handleAdmin).
    AddBeforeMiddlewares([]rtr.MiddlewareInterface{
        rtr.NewAnonymousMiddleware(authMiddleware),
    }))
```

## Built-in Middleware

RTR provides several built-in middleware components:

```go
import "github.com/dracory/rtr/middlewares"

// Recovery middleware
router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    rtr.NewAnonymousMiddleware(middlewares.RecoveryMiddleware),
})

// CORS middleware
router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    rtr.NewAnonymousMiddleware(middlewares.CorsMiddleware),
})

// Rate limiting middleware
router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    rtr.NewAnonymousMiddleware(middlewares.RateLimitByIPMiddleware),
})
```

## Declarative Configuration

For complex applications, use the declarative API:

```go
config := rtr.RouterConfig{
    Name: "My API",
    BeforeMiddleware: []rtr.MiddlewareConfig{
        rtr.NewMiddlewareConfig("Recovery", middlewares.RecoveryMiddleware),
    },
    Routes: []rtr.RouteConfig{
        rtr.GET("/", homeHandler).WithName("Home"),
        rtr.GET("/health", healthHandler).WithName("Health Check"),
    },
    Groups: []rtr.GroupConfig{
        rtr.Group("/api/v1",
            rtr.GET("/users", usersHandler).WithName("List Users"),
            rtr.POST("/users", createUserHandler).WithName("Create User"),
        ).WithName("API Group"),
    },
}

router := rtr.NewRouterFromConfig(config)
```

## Domain-Based Routing

Handle different domains with specific routes:

```go
// Create a domain for API requests
apiDomain := rtr.NewDomain("api.example.com")
apiDomain.AddRoute(rtr.Get("/users", handleAPIUsers))

// Create a domain for web requests
webDomain := rtr.NewDomain("www.example.com")
webDomain.AddRoute(rtr.Get("/", handleWebHome))

// Add domains to router
router.AddDomain(apiDomain)
router.AddDomain(webDomain)
```

## Testing Your Routes

Use the standard Go testing framework:

```go
func TestHelloHandler(t *testing.T) {
    router := rtr.NewRouter()
    router.AddRoute(rtr.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    }))

    req := httptest.NewRequest("GET", "/hello", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }

    if w.Body.String() != "Hello, World!" {
        t.Errorf("Expected 'Hello, World!', got '%s'", w.Body.String())
    }
}
```

## Next Steps

- [Architecture Documentation](architecture.md) - Understand the system design
- [API Reference](api_reference.md) - Complete API documentation
- [Middleware Guide](modules/middleware.md) - Learn about middleware patterns
- [Examples](../examples/) - Explore practical examples

## Common Issues

### Route Not Matching

Ensure your route patterns are correct and that you're using the right HTTP method:

```go
// This will NOT match POST requests
router.AddRoute(rtr.Get("/users", handler))

// Use this for POST requests
router.AddRoute(rtr.Post("/users", handler))
```

### Middleware Not Running

Check the middleware execution order:

1. Global before middleware
2. Domain before middleware
3. Group before middleware
4. Route before middleware
5. Handler
6. Route after middleware
7. Group after middleware
8. Domain after middleware
9. Global after middleware

### Parameters Not Found

Use `MustGetParam()` for required parameters and `GetParam()` for optional ones:

```go
// Required parameter - panics if not found
id := rtr.MustGetParam(r, "id")

// Optional parameter - returns bool indicating existence
id, exists := rtr.GetParam(r, "id")
```
