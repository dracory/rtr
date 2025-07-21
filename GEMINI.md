# GEMINI.md for `dracory/rtr` Project

This file provides comprehensive documentation and guidelines for the `dracory/rtr` Go HTTP router project.

## Project Overview

`dracory/rtr` is a high-performance, feature-rich HTTP router for Go applications. It provides a flexible and intuitive API for building web applications with robust routing capabilities.

## Core Features

- **High Performance**: Optimized for speed with minimal allocations.
- **RESTful Routing**: Intuitive API for defining RESTful endpoints with support for all standard HTTP methods (GET, POST, PUT, DELETE, etc.).
- **Middleware Support**: Flexible middleware chaining with before/after execution at global, domain, group, and route levels. Includes built-in recovery middleware.
- **Route Groups**: Organize routes with shared prefixes and middleware, supporting nested groups for hierarchical routing.
- **Domain-Based Routing**: Handle different domains/subdomains with ease, including wildcard and port matching.
- **Multiple Handler Types**: Support for various response types (JSON, HTML, XML, CSS, Text, JavaScript) with automatic content-type headers.
- **Declarative Configuration**: Define routing structures using configuration objects for better maintainability and tooling support.
- **Path Parameters**: Flexible path parameter extraction with support for required, optional, and wildcard segments.
- **Context Support**: Built-in context support for request-scoped values.
- **Standard Library Compatible**: Implements `http.Handler` for seamless integration with the Go standard library.
- **Comprehensive Testing**: High test coverage with extensive test cases.
- **Route Listing and Debugging**: Built-in functionality to list and visualize the router's configuration for debugging and documentation.

## Detailed Features

### 1. Route Management
The router provides a robust system for managing HTTP routes.
- **HTTP Methods**: Supports `GET`, `POST`, `PUT`, `DELETE`, `PATCH`, `HEAD`, `OPTIONS`, `CONNECT`, `TRACE`.
- **Path Parameters**: Extract values from URL paths using `:param` syntax. Supports optional parameters (`:param?`) and wildcard/catch-all routes (`*filepath`).
- **Flexible Route Matching**: Routes are matched based on exact paths, parameters, and wildcards.
- **Handler Prioritization**: When multiple handlers are set on a single route, they are used in a defined priority order (e.g., `Handler` > `StringHandler` > `HTMLHandler` > `JSONHandler`).

### 2. Middleware System
The middleware system is highly flexible, allowing for powerful request processing pipelines.
- **MiddlewareInterface**: Defines named middleware with metadata support.
- **StdMiddleware**: Standard Go `http.Handler` pattern for middleware.
- **Middleware Chaining and Composition**: Middleware can be chained at global, domain, group, and route levels, executing in a predictable order (before and after the handler).
- **Built-in Middleware**: Includes essential middleware such as Recovery, CORS, Logging, Rate Limiting, Request ID, Secure Headers, and Timeout. For more details, see [Built-in Middleware Documentation](./docs/builtin-middleware.md).
- **Custom Middleware**: Easily create custom middleware by implementing the `Middleware` type.

### 3. Route Organization
The router offers powerful features for organizing routes, improving maintainability and readability.
- **Route Grouping**: Group related routes under a common prefix and apply shared middleware.
- **Nested Groups**: Supports hierarchical routing by allowing groups to be nested within other groups.
- **Domain-Based Routing**: Define routes that only match specific domain names or patterns, including wildcard subdomains and port matching. For more details, see [Domain Routing Documentation](./docs/domains.md).

### 4. Handler Types
The router supports a variety of handler types, automatically setting appropriate HTTP headers for convenience. For complete details, see [Route Handlers Documentation](./docs/route-handlers.md).
- **StdHandler**: Standard `http.HandlerFunc` for full control over the HTTP response.
- **StringHandler**: Returns a plain string response without setting content-type headers.
- **HTMLHandler**: Returns HTML content and automatically sets `Content-Type: text/html; charset=utf-8`.
- **JSONHandler**: Returns JSON content and automatically sets `Content-Type: application/json`.
- **CSSHandler**: Returns CSS content and automatically sets `Content-Type: text/css`.
- **XMLHandler**: Returns XML content and automatically sets `Content-Type: application/xml`.
- **TextHandler**: Returns plain text content and automatically sets `Content-Type: text/plain; charset=utf-8`.
- **JSHandler**: Returns JavaScript content and automatically sets `Content-Type: application/javascript`.
- **ErrorHandler**: Handles errors by returning an error value, allowing for centralized error handling.

### 5. Declarative Configuration
The router supports a declarative approach, allowing you to define your entire routing structure using Go structs or external data sources like JSON, YAML, or a database. This separates routing definition from handler implementation. For more details, see [Declarative Router Configuration](./docs/declarative.md) and [Declarative Routing System from a Database](./docs/declarative-routing-system.md).

### 6. Error Handling
The router provides robust error handling mechanisms.
- **Panic Recovery**: Built-in recovery middleware catches panics and returns 500 Internal Server Error.
- **Not Found Handler**: Customizable handler for 404 Not Found errors.
- **Method Not Allowed Handler**: Customizable handler for 405 Method Not Allowed errors.
- **Global and Route-Level Error Handlers**: Define specific error handling logic at different scopes. For more details, see [Error Handling Documentation](./docs/error-handling.md).

### 7. Route Listing and Debugging
The router includes a built-in `List()` method for debugging and documentation purposes, displaying the router's configuration in formatted tables. This helps visualize the routing structure, identify conflicts, and verify configurations.

### 8. Performance Considerations
The router is designed for performance, with best practices for optimizing applications.
- **Routing Performance**: Discusses linear search, route order, and grouping for efficiency.
- **Middleware Overhead**: Guidance on minimizing overhead by scoping middleware.
- **Concurrency**: Safe for concurrent use with considerations for global state and connection pooling.
- **Memory Usage**: Tips for minimizing memory usage through object reuse and streaming.
- **Caching Strategies**: Example of route caching for high-traffic applications.
- **Benchmarking**: Provides examples for basic and concurrent benchmarking.
- **Production Recommendations**: Includes advice on HTTP/2, reverse proxies, monitoring, and profiling. For more details, see [Performance Considerations Documentation](./docs/performance.md).

### 9. Testing
The router emphasizes comprehensive testing.
- **Unit and Integration Tests**: Examples for testing routes, middleware, and error cases.
- **Mocking Dependencies**: Guidance on using mocking frameworks for isolated testing.
- **Test Helpers**: Provides common test helpers for assertions and setup.
- **Benchmarking**: Examples for benchmarking route matching performance.
- **Best Practices**: Recommendations for table-driven tests, test coverage, parallel testing, and cleanup. For more details, see [Testing Guide Documentation](./docs/testing.md).

## Architecture

### Core Interfaces

1. **RouterInterface**
   - Main router implementation that handles HTTP requests and manages routes and groups.
   - Implements `http.Handler` for seamless integration with the Go standard library.
   - Key methods include: `GetPrefix()`, `SetPrefix()`, `AddGroup()`, `AddGroups()`, `AddRoute()`, `AddRoutes()`, `AddBeforeMiddlewares()`, `AddAfterMiddlewares()`, `ServeHTTP()`.

2. **RouteInterface**
   - Represents individual routes, defining HTTP method, path, and handler associations.
   - Key methods include: `GetMethod()`, `SetMethod()`, `GetPath()`, `SetPath()`, `GetHandler()`, `SetHandler()`, `GetName()`, `SetName()`, `AddBeforeMiddlewares()`, `AddAfterMiddlewares()`.
   - Provides shortcut methods like `Get()`, `Post()`, `Put()`, `Delete()` for concise route creation.

3. **GroupInterface**
   - Manages route groups, allowing for shared prefixes and middleware.
   - Supports nested groups for hierarchical routing.
   - Key methods include: `GetPrefix()`, `SetPrefix()`, `AddRoute()`, `AddRoutes()`, `AddGroup()`, `AddGroups()`, `AddBeforeMiddlewares()`, `AddAfterMiddlewares()`.

4. **DomainInterface**
   - Handles domain-based routing, matching requests against hostnames.
   - Supports exact matches, wildcard subdomains, and port matching.
   - Key methods include: `AddRoute()`, `AddRoutes()`, `AddGroup()`, `AddGroups()`, `AddBeforeMiddlewares()`, `AddAfterMiddlewares()`.

5. **MiddlewareInterface**
   - Defines named middleware, supporting metadata and configuration.
   - Provides a chainable API for building middleware pipelines.

## Middleware System

### Middleware Execution Order

The middleware execution in `dracory/rtr` follows a specific sequence to ensure predictable request processing. The complete order is:

1. **Global Before Middleware**
   - Added via `router.AddBeforeMiddlewares()`
   - Executes first, in the order they were added

2. **Domain Before Middleware**
   - Added via `domain.AddBeforeMiddlewares()`
   - Only executes if the request matches the domain's host pattern
   - Multiple matching domains will execute in the order they were registered

3. **Group Before Middleware**
   - Added via `group.AddBeforeMiddlewares()`
   - Executes from outermost to innermost group
   - For nested groups, parent group middleware runs before child group middleware

4. **Route Before Middleware**
   - Added via `route.AddBeforeMiddlewares()`
   - Executes in the order they were added to the specific route

5. **Route Handler**
   - The actual route handler function processes the request

6. **Route After Middleware**
   - Added via `route.AddAfterMiddlewares()`
   - Executes in reverse order (last added, first executed)

7. **Group After Middleware**
   - Added via `group.AddAfterMiddlewares()`
   - Executes from innermost to outermost group
   - For nested groups, child group middleware runs before parent group middleware

8. **Domain After Middleware**
   - Added via `domain.AddAfterMiddlewares()`
   - Only executes if the request matches the domain's host pattern
   - Multiple matching domains will execute in reverse order of registration

9. **Global After Middleware**
   - Added via `router.AddAfterMiddlewares()`
   - Executes last, in reverse order they were added

### Visual Representation

```
Incoming Request
        ↓
┌───────────────────────┐
│  Global Before (1..N) │
└──────────┬────────────┘
           ↓
┌───────────────────────┐
│  Domain Before (1..N) │
└──────────┬────────────┘
           ↓
┌───────────────────────┐
│  Group Before (Outer) │
└──────────┬────────────┘
           ↓
┌───────────────────────┐
│  Group Before (Inner) │
└──────────┬────────────┘
           ↓
┌───────────────────────┐
│  Route Before (1..N)  │
└──────────┬────────────┘
           ↓
┌───────────────────────┐
│     Route Handler     │
└──────────┬────────────┘
           ↓
┌───────────────────────┐
│  Route After (N..1)   │
└──────────┬────────────┘
           ↓
┌───────────────────────┐
│  Group After (Inner)  │
└──────────┬────────────┘
           ↓
┌───────────────────────┐
│  Group After (Outer)  │
└──────────┬────────────┘
           ↓
┌───────────────────────┐
│  Domain After (N..1)  │
└──────────┬────────────┘
           ↓
┌───────────────────────┐
│  Global After (N..1)  │
└──────────┬────────────┘
           ↓
     Response Sent
```

### Important Notes

- Middleware at the same level executes in the order it was added
- The execution order is always from global → domain → group → route → handler → route → group → domain → global
- After middleware executes in reverse order (last added, first executed) within their respective scopes
- If any middleware in the "before" chain writes a response and doesn't call `next.ServeHTTP()`, the remaining middleware in that chain and the route handler will be skipped

### Middleware Types
1. **Before Middleware**: Executed before the route handler
2. **After Middleware**: Executed after the route handler
3. **Recovery Middleware**: Handles panics in handlers

## Usage Examples

### Quick Start
A minimal example to get started with `dracory/rtr`.
```go
package main

import (
	"net/http"
	"github.com/dracory/router"
)

func main() {
	r := router.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	http.ListenAndServe(":8080", r)
}
```

### Basic Router Setup
```go
r := router.NewRouter()

// Add routes using shortcut methods
r.AddRoute(router.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
}))

// Add routes using method chaining
r.AddRoute(router.NewRoute()
    .SetMethod("GET")
    .SetPath("/users")
    .SetHandler(handleUsers))
```

### Using Route Groups
Organize related routes with shared prefixes and middleware.
```go
// Create an API group
apiGroup := router.NewGroup().SetPrefix("/api")

// Add routes to the group
apiGroup.AddRoute(router.NewRoute()
    .SetMethod("GET")
    .SetPath("/users")
    .SetHandler(handleUsers))

// Add the group to the router
r.AddGroup(apiGroup)
```

### Adding Middleware
Middleware can be added at router, group, or route levels.
```go
// Router-level middleware
r.AddBeforeMiddlewares([]router.Middleware{
    loggingMiddleware,
    authenticationMiddleware,
})

// Group-level middleware
apiGroup.AddBeforeMiddlewares([]router.Middleware{
    apiKeyMiddleware,
})

// Route-level middleware
route.AddBeforeMiddlewares([]router.Middleware{
    specificRouteMiddleware,
})
```

### Declarative API
Define your routing structure using Go structs for better maintainability.
```go
config := rtr.RouterConfig{
    Name: "My API",
    Routes: []rtr.RouteConfig{
        rtr.GET("/", homeHandler).WithName("Home"),
        rtr.POST("/users", createUserHandler).WithName("Create User"),
    },
    Groups: []rtr.GroupConfig{
        rtr.Group("/api",
            rtr.GET("/users", usersHandler).WithName("List Users"),
            rtr.GET("/products", productsHandler).WithName("List Products"),
        ).WithName("API Group"),
    },
}

router := rtr.NewRouterFromConfig(config)
```

### Path Parameters
Extract values from URL paths, including optional and wildcard segments.
```go
// Define a route with parameters
r.AddRoute(rtr.NewRoute().
    SetMethod("GET").
    SetPath("/users/:id").
    SetHandler(func(w http.ResponseWriter, r *http.Request) {
        // Get a required parameter
        id := rtr.MustGetParam(r, "id")
        
        // Or safely get an optional parameter
        if name, exists := rtr.GetParam(r, "name"); exists {
            // Parameter exists
        }
    }))

// Wildcard/Catch-all Routes
// Matches /static/js/main.js, /static/css/style.css, etc.
r.AddRoute(rtr.NewRoute().
    SetMethod("GET").
    SetPath("/static/*filepath").
    SetHandler(serveStaticFile))
```

### Domain-based Routing
Handle requests based on the `Host` header, supporting wildcards and port matching.
```go
// Create a domain for api.example.com
apiDomain := rtr.NewDomain("api.example.com")

// Add routes to the domain
apiDomain.AddRoute(rtr.NewRoute().
    SetMethod("GET").
    SetPath("/users").
    SetHandler(apiUsersHandler))

// Add the domain to the router
router.AddDomain(apiDomain)

// Wildcard subdomain matching
wildcardDomain := rtr.NewDomain("*.example.com")

// IPv6 with port
ipv6Localhost := rtr.NewDomain("[::1]:3000")
```

### Handler Types
Examples for various built-in handler types.

#### StringHandler
```go
r.AddRoute(rtr.NewRoute().
    SetMethod("GET").
    SetPath("/custom").
    SetStringHandler(func(w http.ResponseWriter, req *http.Request) string {
        w.Header().Set("Content-Type", "text/custom")
        w.Header().Set("X-Custom-Header", "value")
        return "Custom content with custom headers"
    }))
```

#### HTMLHandler
```go
r.AddRoute(rtr.NewRoute().
    SetMethod("GET").
    SetPath("/page").
    SetHTMLHandler(func(w http.ResponseWriter, req *http.Request) string {
        return `<!DOCTYPE html>
<html>
<head><title>My Page</title></head>
<body><h1>Hello World!</h1></body>
</html>`
    }))
```

#### JSONHandler
```go
r.AddRoute(rtr.NewRoute().
    SetMethod("GET").
    SetPath("/api/users").
    SetJSONHandler(func(w http.ResponseWriter, req *http.Request) string {
        return `{
    "users": [
        {"id": 1, "name": "Alice"},
        {"id": 2, "name": "Bob"}
    ]
}`
    }))
```

#### CSSHandler
```go
r.AddRoute(rtr.NewRoute().
    SetMethod("GET").
    SetPath("/styles.css").
    SetCSSHandler(func(w http.ResponseWriter, req *http.Request) string {
        return `body {
    font-family: Arial, sans-serif;
    background-color: #f0f0f0;
}

h1 {
    color: #333;
    border-bottom: 2px solid #007acc;
}`
    }))
```

#### XMLHandler
```go
r.AddRoute(rtr.NewRoute().
    SetMethod("GET").
    SetPath("/api/data.xml").
    SetXMLHandler(func(w http.ResponseWriter, req *http.Request) string {
        return `<?xml version="1.0" encoding="UTF-8"?>
<users>
    <user id="1">
        <name>Alice</name>
        <email>alice@example.com</email>
    </user>
</users>`
    }))
```

#### TextHandler
```go
r.AddRoute(rtr.NewRoute().
    SetMethod("GET").
    SetPath("/robots.txt").
    SetTextHandler(func(w http.ResponseWriter, req *http.Request) string {
        return `User-agent: *
Disallow: /admin/
Allow: /

Sitemap: https://example.com/sitemap.xml`
    }))
```

#### JSHandler
```go
r.AddRoute(rtr.NewRoute().
    SetMethod("GET").
    SetPath("/script.js").
    SetJSHandler(func(w http.ResponseWriter, req *http.Request) string {
        return `console.log("Hello from RTR Router!");

function initApp() {
    document.addEventListener('DOMContentLoaded', function() {
        console.log('App initialized');
    });
}

initApp();`
    }))
```

#### ErrorHandler
```go
r.AddRoute(rtr.NewRoute().
    SetMethod("GET").
    SetPath("/might-fail").
    SetErrorHandler(func(w http.ResponseWriter, req *http.Request) error {
        // Some logic that might fail
        if someCondition {
            return errors.New("something went wrong")
        }
        // Success case - no error, no output
        return nil
    }))
```

### Response Helper Functions
Use these helpers in standard `http.HandlerFunc` for consistent responses.
```go
// Using response helpers in a standard handler
r.AddRoute(rtr.NewRoute().
    SetMethod("GET").
    SetPath("/manual").
    SetHandler(func(w http.ResponseWriter, req *http.Request) {
        // These functions set appropriate headers and write content
        rtr.JSONResponse(w, req, `{"message": "Hello JSON"}`)
        // or
        rtr.HTMLResponse(w, req, "<h1>Hello HTML</h1>")
        // or
        rtr.CSSResponse(w, req, "body { color: red; }")
        // or
        rtr.XMLResponse(w, req, "<?xml version='1.0'?><root></root>")
        // or
        rtr.TextResponse(w, req, "Hello Text")
        // or
        rtr.JSResponse(w, req, "console.log('Hello JS');")
    }))
```

## Benchmarks

Performance comparison with other popular routers:

```
BenchmarkRouter/Static-8     5000000   300 ns/op   32 B/op   1 allocs/op
BenchmarkRouter/Param-8      3000000   450 ns/op  160 B/op   4 allocs/op
BenchmarkRouter/Regexp-8     2000000   700 ns/op  320 B/op   6 allocs/op
```

For more detailed benchmarking information and best practices, refer to the [Performance Guide](./docs/performance.md).

## Development Guidelines

### Code Style
- Follow standard Go formatting (`go fmt`)
- Use `goimports` for import organization
- Write clear, descriptive variable and function names
- Document all exported functions and types
- Keep functions small and focused

### Error Handling
- Always check and handle errors
- Provide meaningful error messages
- Use custom error types when appropriate

### Testing
- Write comprehensive tests for all public APIs, routes, and middleware.
- Use table-driven tests where applicable for multiple scenarios.
- Test edge cases and error conditions, including 404 Not Found and 405 Method Not Allowed.
- Utilize Go's `httptest` package for efficient HTTP handler testing.
- Employ mocking frameworks (e.g., Testify Mock) for isolating dependencies.
- Include benchmarks for performance-critical paths.
- Maintain good test coverage.

For detailed testing strategies, examples, and best practices, refer to the [Testing Guide](./docs/testing.md).

## Best Practices

1. **Middleware**
   - Keep middleware focused on a single responsibility
   - Use named middleware for better debugging
   - Document middleware dependencies

2. **Routing**
   - Group related routes together
   - Use meaningful route prefixes
   - Document route requirements and behaviors

3. **Performance**
   - Minimize middleware overhead
   - Use sync.Pool for frequently allocated objects
   - Avoid expensive operations in hot paths

4. **Security**
   - Validate all inputs
   - Use context for request-scoped values
   - Implement proper CORS and CSRF protection

## Project Structure

The project is organized into a clear and modular structure:

- `constants.go`: Defines various constants used throughout the router.
- `interfaces.go`: Contains core interface definitions for `RouterInterface`, `RouteInterface`, `GroupInterface`, `DomainInterface`, and `MiddlewareInterface`.
- `router_implementation.go`: Implements the main `Router` logic.
- `route_implementation.go`: Handles individual route definitions and matching.
- `group_implementation.go`: Implements route grouping functionality.
- `middleware_implementation.go`: Core middleware system implementation.
- `middleware_chain.go`: Manages the chaining and execution of middleware.
- `domain.go`: Implements domain-based routing logic.
- `declarative.go`: Provides the declarative API for defining routes from configuration.
- `list.go`: Contains the `List()` method for debugging and visualizing router configuration.
- `params.go`: Handles path parameter extraction.
- `responses.go`: Provides helper functions for various HTTP response types.
- `functions.go`: Contains general utility functions.
- `middlewares/`: Directory containing various built-in middleware implementations (e.g., `recovery_middleware.go`, `cors_middleware.go`, `logger_middleware.go`).
- `examples/`: Contains runnable examples demonstrating different features of the router.
- `docs/`: Comprehensive documentation, including detailed guides on middleware, handlers, domains, error handling, performance, and testing. Also contains AI prompts for comparisons.
- `README.md`: The main project overview and quick start guide.
- `GEMINI.md`: This file, providing detailed documentation and guidelines for the project.
- `CONTRIBUTING.md`: Guidelines for contributing to the project.
- `LICENSE`: Project license information.
- `go.mod` and `go.sum`: Go module dependency files.
- `.github/workflows/tests.yml`: GitHub Actions workflow for running tests.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run tests and linters
6. Submit a pull request

## License

MIT License - See LICENSE file for details
