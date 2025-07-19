# HTTP Router Package

A flexible and feature-rich HTTP router implementation for Go applications that supports route grouping, middleware chains, and nested routing structures.

<img src="https://opengraph.githubassets.com/5b92c81c05d64a82c3fb4ba95739403a2d38cbad61f260a0701b3366b3d10327/dracory/router" width="300" />

[![Tests Status](https://github.com/dracory/router/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/dracory/router/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dracory/router)](https://goreportcard.com/report/github.com/dracory/router)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dracory/router)](https://pkg.go.dev/github.com/dracory/router)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

## Features

- **Route Management**: Define and manage HTTP routes with support for all standard HTTP methods using exact path matching
- **Route Groups**: Group related routes with shared prefixes and middleware
- **Middleware Support**: 
  - Pre-route (before) middleware
  - Post-route (after) middleware
  - Support at router, group, and individual route levels
  - Built-in panic recovery middleware
- **Nested Groups**: Create hierarchical route structures with nested groups
- **Flexible API**: Chainable methods for intuitive route and group configuration
- **Standard Interface**: Implements `http.Handler` interface for seamless integration

## Middleware

### Built-in Middleware

#### Recovery Middleware
The router includes a built-in recovery middleware that catches panics in your handlers and returns a 500 Internal Server Error response instead of crashing the server. This middleware is added by default when you create a new router with `NewRouter()`.

```go
// This is automatically added when you create a new router
router := router.NewRouter()

// But you can also add it manually if needed
router.AddBeforeMiddlewares([]router.Middleware{router.RecoveryMiddleware})
```

#### Custom Middleware
You can create your own middleware by implementing the `Middleware` type:

```go
func myMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Do something before the handler runs
        log.Println("Before handler")
        
        // Call the next handler
        next.ServeHTTP(w, r)
        
        // Do something after the handler runs
        log.Println("After handler")
    })
}

// Add it to your router
router.AddBeforeMiddlewares([]router.Middleware{myMiddleware})
```

## Core Components

### Router

The main router component that handles HTTP requests and manages routes and groups.

```go
router := router.NewRouter()
```

### Routes

Individual route definitions that specify HTTP method, path, and handler.

```go
// Using shortcut methods
route := router.Get("/users", handleUsers)      // Exact match: /users
route := router.Post("/users", createUser)     // Exact match: /users
route := router.Put("/users/123", updateUser)  // Exact match required: /users/123
route := router.Delete("/users/123", deleteUser) // Exact match required: /users/123

// Using method chaining
route := router.NewRoute()
    .SetMethod("GET")
    .SetPath("/users")
    .SetHandler(handleUsers)
```

### Groups

Route groups that share common prefixes and middleware.

```go
group := router.NewGroup()
    .SetPrefix("/api")
    .AddRoute(route)
```

## Usage Examples

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
    WithBeforeMiddleware(authMiddleware).
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

### Wildcard/Catch-all Routes
Use `*` to match all remaining path segments:

```go
// Matches /static/js/main.js, /static/css/style.css, etc.
r.AddRoute(rtr.NewRoute().
    SetMethod("GET").
    SetPath("/static/*filepath").
    SetHandler(serveStaticFile))
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
- The wildcard parameter must be the last segment in the path

## Domain-based Routing

The router supports domain-based routing, allowing you to define routes that only match specific domain names or patterns.

### Creating a Domain

```go
// Create a domain with exact match
domain := router.NewDomain("example.com")

// Create a domain with wildcard subdomain matching
wildcardDomain := router.NewDomain("*.example.com")

// Create a domain that matches multiple patterns
multiDomain := router.NewDomain("example.com", "api.example.com", "*.example.org")
```

### Adding Routes to a Domain

```go
// Create a new domain
domain := router.NewDomain("api.example.com")

// Add routes directly to the domain
domain.AddRoute(router.Get("/users", handleUsers))

// Add multiple routes at once
domain.AddRoutes([]router.RouteInterface{
    router.Get("/users", handleUsers),
    router.Post("/users", createUser),
})
```

### Adding Groups to a Domain

```go
// Create a domain
domain := router.NewDomain("api.example.com")

// Create an API group
apiGroup := router.NewGroup().SetPrefix("/v1")

// Add routes to the group
apiGroup.AddRoute(router.Get("/products", handleProducts))

// Add the group to the domain
domain.AddGroup(apiGroup)

// Add the domain to the router
router.AddDomain(domain)
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
  domain := router.NewDomain("example.com")  // Matches example.com, example.com:8080, example.com:3000, etc.
  ```

- **Exact port**: Requires exact port match
  ```go
  domain := router.NewDomain("example.com:8080")  // Only matches example.com:8080
  ```

- **Wildcard port**: Matches any port on that host
  ```go
  domain := router.NewDomain("example.com:*")  // Matches example.com with any port
  ```

- **IPv4 and IPv6 support**:
  ```go
  // IPv4 with port
  ipv4Domain := router.NewDomain("127.0.0.1:8080")  // Matches 127.0.0.1:8080
  
  // IPv6 with port (note the square brackets)
  ipv6Domain := router.NewDomain("[::1]:8080")  // Matches [::1]:8080
  ```

#### Examples

```go
// Match any port on example.com
anyPort := router.NewDomain("example.com")

// Match only port 8080
exactPort := router.NewDomain("example.com:8080")

// Match any subdomain on any port
wildcardSubdomain := router.NewDomain("*.example.com:*")

// Match localhost on any port
localhost := router.NewDomain("localhost:*")

// Match IPv6 localhost on port 3000
ipv6Localhost := router.NewDomain("[::1]:3000")
```

### Middleware on Domains

Middleware can be added at the domain level to apply to all routes within that domain:

```go
domain := router.NewDomain("admin.example.com")

// Add middleware that will run before all routes in this domain
domain.AddBeforeMiddlewares([]router.Middleware{
    adminAuthMiddleware,
    loggingMiddleware,
})

// Add middleware that will run after all routes in this domain
domain.AddAfterMiddlewares([]router.Middleware{
    responseTimeMiddleware,
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
router.AddBeforeMiddlewares([]rtr.Middleware{loggingMiddleware})
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
