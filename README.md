# RTR Router

A high-performance, flexible, and feature-rich HTTP router for Go applications with support for middleware chains, route grouping, and domain-based routing.

[![Tests Status](https://github.com/dracory/rtr/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/dracory/rtr/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dracory/rtr)](https://goreportcard.com/report/github.com/dracory/rtr)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dracory/rtr)](https://pkg.go.dev/github.com/dracory/rtr)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![GitHub release](https://img.shields.io/github/v/release/dracory/rtr?include_prereleases&style=flat-square)](https://github.com/dracory/rtr/releases)
[![Go version](https://img.shields.io/github/go-mod/go-version/dracory/rtr)](https://github.com/dracory/rtr)
[![codecov](https://codecov.io/gh/dracory/rtr/branch/main/graph/badge.svg)](https://codecov.io/gh/dracory/rtr)

## Features

- **High Performance**: Optimized for speed with minimal allocations
- **RESTful Routing**: Intuitive API for defining RESTful endpoints
- **Middleware Support**: Flexible middleware chaining with before/after execution
- **Route Groups**: Organize routes with shared prefixes and middleware
- **Domain-Based Routing**: Handle different domains/subdomains with ease
- **Multiple Handler Types**: Support for various response types (JSON, HTML, XML, etc.)
- **Context Support**: Built-in context support for request-scoped values
- **Standard Library Compatible**: Implements `http.Handler` for seamless integration
- **Comprehensive Testing**: High test coverage with extensive test cases

## Installation

```bash
go get github.com/dracory/rtr
```

## Quick Start

```go
package main

import (
	"net/http"
	"github.com/dracory/rtr"
)

func main() {
	r := rtr.NewRouter()
	r.AddRoute(rtr.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}))

	http.ListenAndServe(":8080", r)
}
```

## Documentation

For comprehensive documentation, see the [docs/](./docs/) directory in this repository or the package reference on [PkgGoDev](https://pkg.go.dev/github.com/dracory/rtr). Key guides:

- [Middleware Guide](./docs/middleware.md) - Middleware chaining and execution order
- [Handlers Guide](./docs/handlers.md) - Different handler types and usage
- [Domain Routing](./docs/domains.md) - Handle requests based on hostnames
- [Error Handling](./docs/error-handling.md) - Best practices for error handling
- [Testing Guide](./docs/testing.md) - How to test your routes and middleware
- [Performance Guide](./docs/performance.md) - Performance optimization tips

## Examples

Explore complete, runnable examples in the [examples](./examples/) directory.

## Benchmarks

Performance comparison with other popular routers:

```
BenchmarkRouter/Static-8     5000000   300 ns/op   32 B/op   1 allocs/op
BenchmarkRouter/Param-8      3000000   450 ns/op  160 B/op   4 allocs/op
BenchmarkRouter/Regexp-8     2000000   700 ns/op  320 B/op   6 allocs/op
```

## Contributing

We welcome contributions! Please see our [Contributing Guide](./CONTRIBUTING.md) for details.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

- Inspired by [httprouter](https://github.com/julienschmidt/httprouter) and [chi](https://github.com/go-chi/chi)
- Thanks to all [contributors](https://github.com/dracory/router/graphs/contributors)
- Thanks to all [contributors](https://github.com/dracory/rtr/graphs/contributors)
- **Route Groups**: Group related routes with shared prefixes and middleware
- **Middleware Support**: 
  - Pre-route (before) middleware
  - Post-route (after) middleware
  - Support at router, domain, group, and individual route levels
  - Built-in panic recovery middleware
- **Nested Groups**: Create hierarchical route structures with nested groups
- **Flexible API**: Chainable methods for intuitive route and group configuration
- **Standard Interface**: Implements `http.Handler` interface for seamless integration
- **Declarative Configuration**: [Define routes using configuration objects](./docs/declarative.md) for better maintainability and tooling support

## Middleware

### Built-in Middleware

The router comes with several built-in middleware components. For complete documentation on available middleware and usage examples, see [Built-in Middleware Documentation](./docs/builtin-middleware.md).

- **Recovery**: Catches panics and returns 500 errors
- **CORS**: Handles Cross-Origin Resource Sharing
- **Logging**: Request/response logging
- **Rate Limiting**: Request rate limiting
- **Request ID**: Adds unique IDs to requests
- **Security Headers**: Adds security-related HTTP headers
- **Timeouts**: Request timeout handling

Example of adding middleware:

```go
// Add recovery and logging middleware (named middlewares)
r.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    rtr.NewMiddleware(rtr.WithName("Recovery"), rtr.WithHandler(middlewares.RecoveryMiddleware)),
    rtr.NewMiddleware(rtr.WithName("Logger"), rtr.WithHandler(middlewares.Logger())),
})
```

#### Execution order (before/after)

The middleware execution sequence strictly follows:

- globals before → domains before → groups before → routes before
- handler
- routes after → groups after → domains after → globals after

## Core Components

### Router

The main router component that handles HTTP requests and manages routes and groups.

```go
router := rtr.NewRouter()
```

### Routes

Individual route definitions that specify HTTP method, path, and handler.

```go
// Using shortcut methods
route := rtr.Get("/users", handleUsers)         // Exact match: /users
route := rtr.Post("/users", createUser)         // Exact match: /users
route := rtr.Put("/users/123", updateUser)      // Exact match required: /users/123
route := rtr.Delete("/users/123", deleteUser)   // Exact match required: /users/123

// Using method chaining
route := rtr.NewRoute()
    .SetMethod("GET")
    .SetPath("/users")
    .SetHandler(handleUsers)
```

## Handler Types

The router supports multiple handler types that provide different levels of convenience and functionality. Each handler type is designed for specific use cases and automatically handles appropriate HTTP headers.

## Route Handlers

The router supports multiple handler types for different response formats, each automatically handling appropriate HTTP headers. See [Route Handlers Documentation](docs/handlers.md) for complete details.

### Available Handler Types

1. **Standard Handler** - Full HTTP control
2. **String/Text Handlers** - For plain text responses
3. **Web Content Handlers** - HTML, JSON, XML with auto content-type
4. **Asset Handlers** - CSS, JavaScript for static files
5. **Error Handler** - Centralized error handling

### Quick Example

```go
// JSON response
rtr.GetJSON("/api/status", func(w http.ResponseWriter, r *http.Request) string {
    return `{"status":"ok"}`
})

// HTML response with parameters
rtr.GetHTML("/user/:id", func(w http.ResponseWriter, r *http.Request) string {
    userID := rtr.MustGetParam(r, "id")
    return fmt.Sprintf("<h1>User %s</h1>", userID)
})
```

For complete documentation, examples, and best practices, see [Route Handlers Documentation](docs/handlers.md).

### StringHandler

A generic string handler that returns content without setting any headers automatically. Useful when you need full control over headers but want the convenience of returning a string:

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

### HTMLHandler

Returns HTML content and automatically sets `Content-Type: text/html; charset=utf-8`:

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

### JSONHandler

Returns JSON content and automatically sets `Content-Type: application/json`:

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

### CSSHandler

Returns CSS content and automatically sets `Content-Type: text/css`:

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

### XMLHandler

Returns XML content and automatically sets `Content-Type: application/xml`:

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

### TextHandler

Returns plain text content and automatically sets `Content-Type: text/plain; charset=utf-8`:

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

### JSHandler

Returns JavaScript content and automatically sets `Content-Type: application/javascript`:

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

### ErrorHandler

Handles errors by returning an error value. If the error is `nil`, no content is written. If an error is returned, the error message is written to the response:

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

### Handler Combinations

You can set multiple handlers on a single route. The router will use the highest priority handler that is set:

```go
// This route has both HTML and JSON handlers
// HTMLHandler takes priority and will be used
r.AddRoute(rtr.NewRoute().
    SetMethod("GET").
    SetPath("/content").
    SetHTMLHandler(func(w http.ResponseWriter, req *http.Request) string {
        return "<h1>HTML Content</h1>"  // This will be used
    }).
    SetJSONHandler(func(w http.ResponseWriter, req *http.Request) string {
        return `{"message": "JSON Content"}`  // This will be ignored
    }))
```

### Dynamic Content with Parameters

All handler types work seamlessly with path parameters:

```go
// HTML handler with parameters
r.AddRoute(rtr.NewRoute().
    SetMethod("GET").
    SetPath("/user/:id").
    SetHTMLHandler(func(w http.ResponseWriter, req *http.Request) string {
        userID := rtr.MustGetParam(req, "id")
        return fmt.Sprintf(`<h1>User Profile</h1><p>User ID: %s</p>`, userID)
    }))

// JSON handler with parameters
r.AddRoute(rtr.NewRoute().
    SetMethod("GET").
    SetPath("/api/user/:id").
    SetJSONHandler(func(w http.ResponseWriter, req *http.Request) string {
        userID := rtr.MustGetParam(req, "id")
        return fmt.Sprintf(`{"user": {"id": "%s", "name": "User %s"}}`, userID, userID)
    }))
```

### Response Helper Functions

The router provides response helper functions that you can use directly in standard handlers:

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

### Groups

Route groups that share common prefixes and middleware.

```go
group := rtr.NewGroup()
    .SetPrefix("/api")
    .AddRoute(route)
```

## Usage Examples

### Basic Router Setup

```go
r := rtr.NewRouter()

// Add routes using shortcut methods
r.AddRoute(rtr.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
}))

// Add routes using method chaining
r.AddRoute(rtr.NewRoute()
    .SetMethod("GET")
    .SetPath("/users")
    .SetHandler(handleUsers))
```

### Using Route Groups

```go
// Create an API group
apiGroup := rtr.NewGroup().SetPrefix("/api")

// Add routes to the group
apiGroup.AddRoute(rtr.NewRoute()
    .SetMethod("GET")
    .SetPath("/users")
    .SetHandler(handleUsers))

// Add the group to the router
r.AddGroup(apiGroup)
```

### Adding Middleware

```go
// Router-level middleware
r.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    rtr.NewMiddleware(rtr.WithName("Logging"), rtr.WithHandler(loggingMiddleware)),
    rtr.NewMiddleware(rtr.WithName("Auth"), rtr.WithHandler(authenticationMiddleware)),
})

// Group-level middleware
apiGroup.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    rtr.NewMiddleware(rtr.WithName("APIKey"), rtr.WithHandler(apiKeyMiddleware)),
})

// Route-level middleware
route.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    rtr.NewMiddleware(rtr.WithName("SpecificRoute"), rtr.WithHandler(specificRouteMiddleware)),
})
```

## Declarative API

In addition to the imperative API shown above, the router also supports a **declarative configuration approach** that allows you to define your entire routing structure as data structures.

### Basic Declarative Usage

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

### Declarative Route Helpers

```go
// HTTP method helpers
rtr.GET("/users", handler)     // GET route
rtr.POST("/users", handler)    // POST route
rtr.PUT("/users/:id", handler) // PUT route
rtr.DELETE("/users/:id", handler) // DELETE route

// Chainable configuration
rtr.GET("/users", handler).
    WithName("List Users").
    WithBeforeMiddleware(rtr.NewAnonymousMiddleware(authMiddleware)).
    WithMetadata("version", "1.0")
```

### Hybrid Approach

You can mix declarative and imperative approaches:

```go
// Start with declarative configuration
config := rtr.RouterConfig{
    Routes: []rtr.RouteConfig{
        rtr.GET("/", homeHandler).WithName("Home"),
    },
}
router := rtr.NewRouterFromConfig(config)

// Add imperative routes
router.AddRoute(rtr.Get("/health", healthHandler).SetName("Health"))
```

### Benefits of Declarative API

- **Serializable**: Configuration can be exported to JSON/YAML
- **Testable**: Easier to unit test route configurations
- **Readable**: Clear structure and intent
- **Tooling-friendly**: Better IDE support and validation

## Path Parameters

The router supports flexible path parameter extraction with the following features:

### Basic Parameters
Extract values from URL paths using `:param` syntax:

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
```

### Optional Parameters
Mark parameters as optional with `?`:

```go
// Both /articles/tech and /articles/tech/123 will match
r.AddRoute(rtr.NewRoute().
    SetMethod("GET").
    SetPath("/articles/:category/:id?").
    SetHandler(handleArticle))
```

### Greedy Parameters (capture the remainder)
Use `:name...` as the last segment to capture all remaining path segments into a single parameter.
This mirrors the standard library pattern like `/files/{pathname...}`.

```go
// Matches both of the following and captures the remainder into "pathname":
//   /files/images/photo.jpg         -> pathname = "images/photo.jpg"
//   /files/user/albums/2025/img.png -> pathname = "user/albums/2025/img.png"
r.AddRoute(rtr.NewRoute().
    SetMethod("GET").
    SetPath("/files/:pathname...").
    SetHandler(func(w http.ResponseWriter, r *http.Request) {
        pathname := rtr.MustGetParam(r, "pathname")
        _ = pathname // use it
    }))

// Example for thumbnail routes: /th/:extension/:size/:quality/:path...
//   /th/jpg/300x300/80/avatar.png        -> path = "avatar.png"
//   /th/jpg/300x300/80/user/avatar.png   -> path = "user/avatar.png"
r.AddRoute(rtr.NewRoute().
    SetMethod("GET").
    SetPath("/th/:extension/:size/:quality/:path...").
    SetHandler(func(w http.ResponseWriter, r *http.Request) {
        ext := rtr.MustGetParam(r, "extension")
        size := rtr.MustGetParam(r, "size")
        quality := rtr.MustGetParam(r, "quality")
        tail := rtr.MustGetParam(r, "path")
        _, _, _, _ = ext, size, quality, tail
    }))
```

### Wildcard/Catch-all Routes
Use `/*` to allow any suffix after a base path. This is a non-capturing catch-all; if you need the remainder, use a greedy parameter as shown above.

```go
// Matches any path starting with /static/ (non-capturing)
// e.g., /static/js/main.js, /static/css/style.css
r.AddRoute(rtr.NewRoute().
    SetMethod("GET").
    SetPath("/static/*").
    SetHandler(func(w http.ResponseWriter, r *http.Request) {
        // Use a greedy parameter instead if you need the suffix
        w.WriteHeader(http.StatusOK)
    }))
```

### Getting All Parameters
Retrieve all path parameters as a map:

```go
params := rtr.GetParams(r)
// params is a map[string]string of all path parameters
```

## Path Matching Rules

The router uses the following matching rules:
- Paths are matched exactly as defined, with parameter placeholders
- Required parameters must be present in the request path
- Optional parameters can be omitted
- Parameter names must be unique within a route
- A greedy parameter (`:name...`) must be the last segment in the path

### Brace-style parameter aliases
RTR also accepts brace-style parameters as an alias for the colon syntax to stay close to the standard library style. The following are equivalent:

- `{id}` ≡ `:id` (required, one segment)
- `{id?}` ≡ `:id?` (optional)
- `{path...}` ≡ `:path...` (greedy tail, must be last)

Example:

```go
// Using brace aliases
r.AddRoute(rtr.NewRoute().
    SetMethod("GET").
    SetPath("/users/{id}").
    SetHandler(func(w http.ResponseWriter, r *http.Request) {
        id := rtr.MustGetParam(r, "id")
        _ = id
    }))
```

## Domain-based Routing

The router supports domain-based routing, allowing you to define routes that only match specific domain names or patterns.

### Creating a Domain

```go
// Create a domain with exact match
domain := rtr.NewDomain("example.com")

// Create a domain with wildcard subdomain matching
wildcardDomain := rtr.NewDomain("*.example.com")

// Create a domain that matches multiple patterns
multiDomain := rtr.NewDomain("example.com", "api.example.com", "*.example.org")
```

### Adding Routes to a Domain

```go
// Create a new domain
domain := rtr.NewDomain("api.example.com")

// Add routes directly to the domain
domain.AddRoute(rtr.Get("/users", handleUsers))

// Add multiple routes at once
domain.AddRoutes([]rtr.RouteInterface{
    rtr.Get("/users", handleUsers),
    rtr.Post("/users", createUser),
})
```

### Adding Groups to a Domain

```go
// Create a domain
domain := rtr.NewDomain("api.example.com")

// Create an API group
apiGroup := rtr.NewGroup().SetPrefix("/v1")

// Add routes to the group
apiGroup.AddRoute(rtr.Get("/products", handleProducts))

// Add the group to the domain
domain.AddGroup(apiGroup)

// Add the domain to the router
r := rtr.NewRouter()
r.AddDomain(domain)
```

### Domain Matching

Domains are matched against the `Host` header of incoming requests. The matching supports:

#### Basic Domain Matching
- Exact matches (`example.com`)
- Wildcard subdomains (`*.example.com`)
- Multiple patterns per domain

#### Port Matching
- **No port in pattern**: Matches any port on that host
  ```go
  domain := rtr.NewDomain("example.com")  // Matches example.com, example.com:8080, example.com:3000, etc.
  ```

- **Exact port**: Requires exact port match
  ```go
  domain := rtr.NewDomain("example.com:8080")  // Only matches example.com:8080
  ```

- **Wildcard port**: Matches any port on that host
  ```go
  domain := rtr.NewDomain("example.com:*")  // Matches example.com with any port
  ```

- **IPv4 and IPv6 support**:
  ```go
  // IPv4 with port
  ipv4Domain := rtr.NewDomain("127.0.0.1:8080")  // Matches 127.0.0.1:8080
  
  // IPv6 with port (note the square brackets)
  ipv6Domain := rtr.NewDomain("[::1]:8080")  // Matches [::1]:8080
  ```

#### Examples

```go
// Match any port on example.com
anyPort := rtr.NewDomain("example.com")

// Match only port 8080
exactPort := rtr.NewDomain("example.com:8080")

// Match any subdomain on any port
wildcardSubdomain := rtr.NewDomain("*.example.com:*")

// Match localhost on any port
localhost := rtr.NewDomain("localhost:*")

// Match IPv6 localhost on port 3000
ipv6Localhost := rtr.NewDomain("[::1]:3000")
```

### Middleware on Domains

Middleware can be added at the domain level to apply to all routes within that domain:

```go
domain := rtr.NewDomain("admin.example.com")

// Add middleware that will run before all routes in this domain
domain.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    rtr.NewMiddleware(rtr.WithName("AdminAuth"), rtr.WithHandler(adminAuthMiddleware)),
    rtr.NewMiddleware(rtr.WithName("Logging"), rtr.WithHandler(loggingMiddleware)),
})

// Add middleware that will run after all routes in this domain
domain.AddAfterMiddlewares([]rtr.MiddlewareInterface{
    rtr.NewMiddleware(rtr.WithName("ResponseTime"), rtr.WithHandler(responseTimeMiddleware)),
})
```

## Interfaces

### RouterInterface

The main router interface that provides methods for managing routes and groups:

- `GetPrefix()` / `SetPrefix()`: Manage router prefix
- `AddGroup()` / `AddGroups()`: Add route groups
- `AddRoute()` / `AddRoutes()`: Add individual routes
- `AddBeforeMiddlewares()` / `AddAfterMiddlewares()`: Add middleware chains
- `ServeHTTP()`: Handle HTTP requests

### GroupInterface

Interface for managing route groups:

- `GetPrefix()` / `SetPrefix()`: Manage group prefix
- `AddRoute()` / `AddRoutes()`: Add routes to the group
- `AddGroup()` / `AddGroups()`: Add nested groups
- `AddBeforeMiddlewares()` / `AddAfterMiddlewares()`: Add group-level middleware

### RouteInterface

Interface for configuring individual routes:

- `GetMethod()` / `SetMethod()`: HTTP method configuration
- `GetPath()` / `SetPath()`: URL path configuration
- `GetHandler()` / `SetHandler()`: Route handler configuration
- `GetName()` / `SetName()`: Route naming
- `AddBeforeMiddlewares()` / `AddAfterMiddlewares()`: Route-specific middleware

#### Shortcut Methods

The package provides shortcut methods for common HTTP methods:

- `Get(path string, handler Handler) RouteInterface` - Creates a GET route
- `Post(path string, handler Handler) RouteInterface` - Creates a POST route
- `Put(path string, handler Handler) RouteInterface` - Creates a PUT route
- `Delete(path string, handler Handler) RouteInterface` - Creates a DELETE route

These methods automatically set the HTTP method, path, and handler, making route creation more concise.

## Route Listing and Debugging

The router provides a built-in `List()` method for debugging and documentation purposes. This method displays the router's configuration in formatted tables, making it easy to visualize your routing structure.

### Using the List Method

```go
router := rtr.NewRouter()

// Add some routes and middleware
router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{rtr.NewMiddleware(rtr.WithName("Logging"), rtr.WithHandler(loggingMiddleware))})
router.AddRoute(rtr.Get("/", homeHandler).SetName("Home"))

// Create a group
apiGroup := rtr.NewGroup().SetPrefix("/api")
apiGroup.AddRoute(rtr.Get("/users", usersHandler).SetName("List Users"))
router.AddGroup(apiGroup)

// Display the router configuration
router.List()
```

### Output Format

The `List()` method displays:

1. **Global Middleware Table**: Shows before and after middleware applied at the router level
2. **Domain Routes Tables**: Shows routes organized by domain (if using domain-based routing)
3. **Direct Routes Table**: Shows routes added directly to the router
4. **Group Routes Tables**: Shows routes organized by groups with their prefixes

#### Example Output

```
+------------------------------------+
| GLOBAL BEFORE MIDDLEWARE LIST (TOTAL: 2) |
+---+--------------------------------+------+
| # | MIDDLEWARE NAME                | TYPE |
+---+--------------------------------+------+
| 1 | RecoveryMiddleware             | Before |
| 2 | LoggingMiddleware              | Before |
+---+--------------------------------+------+

+---------------------------------------------------------------+
| DIRECT ROUTES LIST (TOTAL: 1)                                |
+---+------------+--------+------------+---------------------+
| # | ROUTE PATH | METHOD | ROUTE NAME | MIDDLEWARE LIST     |
+---+------------+--------+------------+---------------------+
| 1 | /          | GET    | Home       | none                |
+---+------------+--------+------------+---------------------+

+---------------------------------------------------------------+
| GROUP ROUTES [/api] (TOTAL: 1)                               |
+---+------------+--------+------------+---------------------+
| # | ROUTE PATH | METHOD | ROUTE NAME | MIDDLEWARE LIST     |
+---+------------+--------+------------+---------------------+
| 1 | /api/users | GET    | List Users | none                |
+---+------------+--------+------------+---------------------+
```

### Middleware Name Detection

The List method attempts to extract meaningful names from middleware functions using reflection:

- **Named functions**: Shows the actual function name (e.g., `RecoveryMiddleware`)
- **Anonymous functions**: Shows `anonymous` or attempts to extract from closure context
- **Method receivers**: Shows the method name when middleware is defined on a struct

### Use Cases

- **Development**: Quickly verify your routing configuration
- **Debugging**: Identify routing conflicts or missing routes
- **Documentation**: Generate route documentation for your API
- **Testing**: Validate that routes are configured as expected

## Testing

The package includes comprehensive test coverage:

- `router_test.go`: Core router functionality tests
- `router_integration_test.go`: Integration tests
- `route_test.go`: Route-specific tests
- `group_test.go`: Group functionality tests
- `examples/basic/`: Complete example with tests

Run tests using:

```bash
# From the root directory
go test .

# Or to run all tests including examples
go test ./...
```
