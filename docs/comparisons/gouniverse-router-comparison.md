# Dracory Router vs Gouniverse Router: Routing Capabilities Comparison

## Overview
This document compares the routing capabilities of Dracory Router and [Gouniverse Router](https://github.com/gouniverse/router), focusing on their design philosophies and practical differences.

## Feature Comparison Table

| Feature | Dracory Router | Gouniverse Router |
|---------|----------------|-------------------|
| **Basic Routing** | ✅ All HTTP methods with exact path matching | ✅ All HTTP methods (implicit support for all verbs unless specified) |
| **Path Parameters** | ✅ Named parameters with optional segments (`/user/:id`, `/articles/:category/:id?`) | ✅ Basic parameter support (limited documentation) |
| **Wildcards** | ✅ Catch-all routes (`/static/*filepath`) | ✅ Catch-all routes (`/*`) |
| **Route Groups** | ✅ Full support with nesting and shared middleware | ✅ Helper functions for route grouping (`RoutesPrependPath`) |
| **Domain Routing** | ✅ Full support with wildcard subdomains and port matching | ❌ No built-in support |
| **Middleware** | ✅ Before/after middleware at router, group, and route levels | ✅ Global and per-route middleware support |
| **Middleware Chaining** | ✅ Hierarchical chaining | ✅ Explicit middleware definition |
| **Performance** | Lightweight with standard `http.Handler` interface | Lightweight, designed to be fast |
| **Handler Types** | ✅ Standard `http.HandlerFunc` | ✅ Multiple types: HTML, JSON, and Idiomatic handlers |
| **Route Listing** | ❌ No built-in support | ✅ Built-in route listing and visualization |
| **Controller Support** | ✅ Through standard patterns | ✅ Built-in MVC controller interfaces |
| **Custom Not Found** | ✅ Supported through middleware | ✅ Catch-all route support |
| **Route Naming** | ✅ Supported | ✅ Required for each route |
| **Testing** | ✅ Standard Go testing patterns | ✅ Easy testing with string return values |
| **Integration** | ✅ Standard `http.Handler` interface | ✅ Chi router integration built-in |
| **Learning Curve** | Standard Go HTTP patterns | Declarative, opinionated approach |

## Code Examples

### Basic Route Definition

**Dracory Router**
```go
router := rtr.NewRouter()
router.AddRoute(rtr.Get("/users", usersHandler))
router.AddRoute(rtr.Post("/users", createUserHandler))
```

**Gouniverse Router**
```go
routes := []router.RouteInterface{
    &router.Route{
        Name: "Users List",
        Path: "/users",
        HTMLHandler: func(w http.ResponseWriter, r *http.Request) string {
            return "Users list"
        },
    },
    &router.Route{
        Name: "Create User",
        Path: "/users",
        Methods: []string{http.MethodPost},
        JSONHandler: func(w http.ResponseWriter, r *http.Request) string {
            return api.Success("User created")
        },
    },
}
```

### Route Groups

**Dracory Router**
```go
api := rtr.NewGroup().SetPrefix("/api")
api.AddRoute(rtr.Get("/users", usersHandler))
api.AddRoute(rtr.Post("/users", createUserHandler))
router.AddGroup(api)
```

**Gouniverse Router**
```go
// Using helper function to prepend path
userRoutes := []router.RouteInterface{
    &router.Route{Name: "Users", Path: "/users", HTMLHandler: usersHandler},
    &router.Route{Name: "Create User", Path: "/users", Methods: []string{"POST"}, JSONHandler: createUserHandler},
}
apiRoutes := router.RoutesPrependPath(userRoutes, "/api")
```

### Middleware

**Dracory Router**
```go
// Router-level middleware
router.AddBeforeMiddlewares([]rtr.Middleware{loggingMiddleware})

// Group-level middleware
api.AddBeforeMiddlewares([]rtr.Middleware{authMiddleware})

// Route-level middleware
route.AddBeforeMiddlewares([]rtr.Middleware{specificMiddleware})
```

**Gouniverse Router**
```go
// Global middleware
globalMiddlewares := []router.Middleware{
    {Name: "Logger", Handler: loggingMiddleware},
    {Name: "Auth", Handler: authMiddleware},
}

// Per-route middleware
&router.Route{
    Name: "Protected Route",
    Path: "/dashboard",
    Middlewares: []router.Middleware{
        {Name: "Check Auth", Handler: checkAuthMiddleware},
    },
    HTMLHandler: dashboardHandler,
}
```

### Domain-based Routing

**Dracory Router**
```go
// Create domain with wildcard support
domain := rtr.NewDomain("*.example.com")
domain.AddRoute(rtr.Get("/api/users", apiUsersHandler))
router.AddDomain(domain)

// Port-specific matching
localDomain := rtr.NewDomain("localhost:8080")
localDomain.AddRoute(rtr.Get("/debug", debugHandler))
router.AddDomain(localDomain)
```

**Gouniverse Router**
```go
// No built-in domain routing support
// Would need to implement custom logic in handlers
```

### Path Parameters

**Dracory Router**
```go
// Required and optional parameters
router.AddRoute(rtr.Get("/users/:id", userHandler))
router.AddRoute(rtr.Get("/articles/:category/:id?", articleHandler))
router.AddRoute(rtr.Get("/static/*filepath", staticHandler))

// In handler
func userHandler(w http.ResponseWriter, r *http.Request) {
    id := rtr.MustGetParam(r, "id")
    // Handle user with id
}
```

**Gouniverse Router**
```go
// Basic parameter support (documentation limited)
&router.Route{
    Name: "User Profile",
    Path: "/users/:id",
    HTMLHandler: func(w http.ResponseWriter, r *http.Request) string {
        // Parameter extraction method not clearly documented
        return "User profile"
    },
}
```

### Route Listing

**Dracory Router**
```go
// No built-in route listing
// Would need custom implementation
```

**Gouniverse Router**
```go
// Built-in route listing
router.List(globalMiddlewares, routes)
// Outputs formatted table with routes, methods, and middleware
```

## Design Philosophy Comparison

### Dracory Router
- **Traditional Approach**: Follows standard Go HTTP patterns
- **Flexibility First**: Designed for complex routing scenarios
- **Standard Interface**: Uses `http.Handler` for maximum compatibility
- **Hierarchical**: Supports nested groups and middleware inheritance

### Gouniverse Router
- **Declarative Approach**: Routes defined as data structures
- **Simplicity First**: Prioritizes clear, explicit configuration
- **String-based Handlers**: Returns strings for easier testing
- **Opinionated**: Makes specific choices to reduce cognitive load

## Performance Considerations

| Aspect | Dracory Router | Gouniverse Router |
|--------|----------------|-------------------|
| **Memory Usage** | Standard Go patterns | String returns may use more memory |
| **Route Matching** | Efficient exact matching | Standard matching performance |
| **Middleware Overhead** | Minimal with standard interface | Explicit middleware definition |
| **Startup Time** | Fast initialization | Route listing adds minimal overhead |

## Use Case Recommendations

### Choose Dracory Router When:
- Building multi-domain applications
- Need complex routing hierarchies
- Require fine-grained middleware control
- Working with existing Go HTTP ecosystem
- Building high-performance APIs
- Need IPv4/IPv6 with port-specific routing

### Choose Gouniverse Router When:
- Building simple to medium complexity applications
- Prefer explicit, declarative configuration
- Want easy route debugging and visualization
- Building MVC-style applications
- Need multiple handler types (HTML/JSON)
- Want simplified testing with string returns

## Conclusion

Both routers serve different needs:

- **Dracory Router** excels in complex, enterprise-level applications requiring maximum flexibility and standard Go patterns.
- **Gouniverse Router** shines in applications where simplicity, explicitness, and easy debugging are prioritized over flexibility.

The choice depends on your project's complexity, team preferences, and specific routing requirements.
