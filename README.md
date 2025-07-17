# HTTP Router Package

A flexible and feature-rich HTTP router implementation for Go applications that supports route grouping, middleware chains, and nested routing structures.

## Features

- **Route Management**: Define and manage HTTP routes with support for all standard HTTP methods using exact path matching
- **Route Groups**: Group related routes with shared prefixes and middleware
- **Middleware Support**: 
  - Pre-route (before) middleware
  - Post-route (after) middleware
  - Support at router, group, and individual route levels
- **Nested Groups**: Create hierarchical route structures with nested groups
- **Flexible API**: Chainable methods for intuitive route and group configuration
- **Standard Interface**: Implements `http.Handler` interface for seamless integration

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

## Path Matching

This router uses **exact path matching**:
- Paths must match exactly as defined
- No built-in path parameters (e.g., `/users/:id`)
- No automatic trailing slash redirection

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
