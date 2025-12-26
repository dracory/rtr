---
path: modules/groups.md
page-type: module
summary: Route grouping system for organizing related routes with shared prefixes and middleware.
tags: [module, groups, organization, prefixes, middleware-inheritance]
created: 2025-12-26
updated: 2025-12-26
version: 1.0.0
---

# Groups Module

The groups module provides a powerful system for organizing related routes with shared prefixes, middleware, and configuration. Groups enable hierarchical route organization and middleware inheritance.

## Overview

Route groups allow you to organize related endpoints under common prefixes while sharing middleware and configuration. This is particularly useful for API versioning, feature modules, and logical route organization.

## Key Features

- **Path Prefixes**: Automatic path prefix application to all group routes
- **Middleware Inheritance**: Shared middleware for all routes in the group
- **Nested Groups**: Support for hierarchical group structures
- **Metadata Support**: Group-level metadata for documentation and tooling
- **Flexible Configuration**: Both imperative and declarative configuration support

## Core Interface

### GroupInterface

```go
type GroupInterface interface {
    // Prefix management
    GetPrefix() string
    SetPrefix(prefix string) GroupInterface
    
    // Route management
    AddRoute(route RouteInterface) GroupInterface
    AddRoutes(routes []RouteInterface) GroupInterface
    GetRoutes() []RouteInterface
    
    // Nested group management
    AddGroup(group GroupInterface) GroupInterface
    AddGroups(groups []GroupInterface) GroupInterface
    GetGroups() []GroupInterface
    
    // Middleware management
    AddBeforeMiddlewares(middlewares []MiddlewareInterface) GroupInterface
    AddAfterMiddlewares(middlewares []MiddlewareInterface) GroupInterface
    GetBeforeMiddlewares() []MiddlewareInterface
    GetAfterMiddlewares() []MiddlewareInterface
    
    // Group metadata
    GetName() string
    SetName(name string) GroupInterface
    GetMetadata() map[string]interface{}
    SetMetadata(metadata map[string]interface{}) GroupInterface
    
    // Utility methods
    String() string
}
```

## Group Creation

### Basic Group Creation

```go
// Create a new group
apiGroup := rtr.NewGroup().SetPrefix("/api")

// Add routes to the group
apiGroup.AddRoute(rtr.Get("/users", usersHandler))
apiGroup.AddRoute(rtr.Post("/users", createUserHandler))

// Add the group to the router
router.AddGroup(apiGroup)
```

### Method Chaining

```go
// Create group with method chaining
apiGroup := rtr.NewGroup().
    SetPrefix("/api/v1").
    SetName("API v1").
    AddBeforeMiddlewares([]rtr.MiddlewareInterface{
        rtr.NewAnonymousMiddleware(apiKeyMiddleware),
        rtr.NewAnonymousMiddleware(rateLimitMiddleware),
    }).
    AddRoute(rtr.Get("/users", usersHandler)).
    AddRoute(rtr.Post("/users", createUserHandler))

router.AddGroup(apiGroup)
```

## Path Prefixes

### Prefix Application

Group prefixes are automatically applied to all routes within the group:

```go
// Group with prefix "/api/v1"
apiGroup := rtr.NewGroup().SetPrefix("/api/v1")

// Route defined as "/users"
apiGroup.AddRoute(rtr.Get("/users", usersHandler))

// Final route path: "/api/v1/users"
```

### Prefix Normalization

- Prefixes are automatically normalized (leading slash ensured)
- Trailing slashes are handled appropriately
- Nested prefixes are combined correctly

```go
// These are equivalent:
group1 := rtr.NewGroup().SetPrefix("api")
group2 := rtr.NewGroup().SetPrefix("/api")
group3 := rtr.NewGroup().SetPrefix("/api/")
```

## Nested Groups

### Basic Nesting

```go
// Parent group
apiGroup := rtr.NewGroup().SetPrefix("/api")

// Child group
v1Group := rtr.NewGroup().
    SetPrefix("/v1").
    AddBeforeMiddlewares([]rtr.MiddlewareInterface{
        rtr.NewAnonymousMiddleware(versionMiddleware),
    })

// Add child to parent
apiGroup.AddGroup(v1Group)
router.AddGroup(apiGroup)

// Routes in v1Group will be accessible at "/api/v1/..."
```

### Complex Nesting

```go
// Create hierarchical structure
router := rtr.NewRouter()

// API group
apiGroup := rtr.NewGroup().
    SetPrefix("/api").
    SetName("API").
    AddBeforeMiddlewares([]rtr.MiddlewareInterface{
        rtr.NewAnonymousMiddleware(corsMiddleware),
        rtr.NewAnonymousMiddleware(loggingMiddleware),
    })

// Version groups
v1Group := rtr.NewGroup().
    SetPrefix("/v1").
    SetName("API v1").
    AddBeforeMiddlewares([]rtr.MiddlewareInterface{
        rtr.NewAnonymousMiddleware(apiKeyMiddleware),
    })

v2Group := rtr.NewGroup().
    SetPrefix("/v2").
    SetName("API v2").
    AddBeforeMiddlewares([]rtr.MiddlewareInterface{
        rtr.NewAnonymousMiddleware(oauthMiddleware),
    })

// Feature groups within versions
usersGroup := rtr.NewGroup().
    SetPrefix("/users").
    SetName("Users").
    AddRoute(rtr.Get("/", listUsersHandler)).
    AddRoute(rtr.Post("/", createUserHandler))

// Build hierarchy
v1Group.AddGroup(usersGroup)
apiGroup.AddGroups([]rtr.GroupInterface{v1Group, v2Group})
router.AddGroup(apiGroup)

// Resulting paths:
// /api/v1/users/
// /api/v2/users/
```

## Middleware Inheritance

### Group-Level Middleware

```go
// Group with middleware
apiGroup := rtr.NewGroup().
    SetPrefix("/api").
    AddBeforeMiddlewares([]rtr.MiddlewareInterface{
        rtr.NewMiddleware("APIAuth", apiAuthMiddleware),
        rtr.NewMiddleware("RateLimit", rateLimitMiddleware),
    }).
    AddAfterMiddlewares([]rtr.MiddlewareInterface{
        rtr.NewMiddleware("APILogging", apiLoggingMiddleware),
    })

// All routes in this group inherit the middleware
apiGroup.AddRoute(rtr.Get("/users", usersHandler))
apiGroup.AddRoute(rtr.Get("/products", productsHandler))
```

### Middleware Composition

```mermaid
graph TD
    A[Global Middleware] --> B B[Group Group Before Middleware]
    B --> C[Nested Group Before Middleware]
    C --> D[Route Before Middleware]
    D --> E[Handler]
    E --> F[Route After Middleware]
    F --> G[NNested Group After Middleware]
    G --> H[Group After Middleware]
    H --> I[Global After Middleware]
```

### Middleware Override

Routes can add their own middleware in addition to inherited middleware:

```go
apiGroup := rtr.NewGroup().
    SetPrefix("/api").
    AddBeforeMiddlewares([]rtr.MiddlewareInterface{
        rtr.NewAnonymousMiddleware(apiKeyMiddleware),
    })

// Route with additional middleware
apiGroup.AddRoute(rtr.Get("/admin", adminHandler).
    AddBeforeMiddlewares([]rtr.MiddlewareInterface{
        rtr.NewAnonymousMiddleware(adminAuthMiddleware),
    }))

// Execution order:
// 1. Group middleware (apiKeyMiddleware)
// 2. Route middleware (adminAuthMiddleware)
// 3. Handler
```

## Configuration Support

### Declarative Group Configuration

```go
config := rtr.GroupConfig{
    Prefix: "/api/v1",
    Name: "API v1",
    BeforeMiddleware: []rtr.MiddlewareConfig{
        rtr.NewMiddlewareConfig("APIKey", apiKeyMiddleware),
        rtr.NewMiddlewareConfig("RateLimit", rateLimitMiddleware),
    },
    AfterMiddleware: []rtr.MiddlewareConfig{
        rtr.NewMiddlewareConfig("APILogging", apiLoggingMiddleware),
    },
    Routes: []rtr.RouteConfig{
        rtr.GET("/users", usersHandler).WithName("List Users"),
        rtr.POST("/users", createUserHandler).WithName("Create User"),
    },
    Groups: []rtr.GroupConfig{
        rtr.Group("/admin",
            rtr.GET("/dashboard", adminDashboardHandler),
        ).WithName("Admin").
        WithBeforeMiddleware(
            rtr.NewAnonymousMiddleware(adminAuthMiddleware),
        ),
    },
    Metadata: map[string]interface{}{
        "version": "1.0",
        "deprecated": false,
    },
}
```

### Group Helper Function

```go
// Using the Group helper function
apiGroup := rtr.Group("/api/v1",
    rtr.GET("/users", usersHandler).WithName("List Users"),
    rtr.POST("/users", createUserHandler).WithName("Create User"),
    rtr.GET("/products", productsHandler).WithName("List Products"),
).
WithName("API v1").
WithBeforeMiddleware(
    rtr.NewAnonymousMiddleware(apiKeyMiddleware),
).
WithMetadata("version", "1.0")
```

## Use Cases

### API Versioning

```go
// Version 1 API
v1Group := rtr.NewGroup().
    SetPrefix("/v1").
    SetName("API v1").
    AddRoute(rtr.Get("/users", v1UsersHandler))

// Version 2 API
v2Group := rtr.NewGroup().
    SetPrefix("/v2").
    SetName("API v2").
    AddRoute(rtr.Get("/users", v2UsersHandler))

apiGroup := rtr.NewGroup().SetPrefix("/api")
apiGroup.AddGroups([]rtr.GroupInterface{v1Group, v2Group})
router.AddGroup(apiGroup)

// Routes:
// /api/v1/users
// /api/v2/users
```

### Feature Modules

```go
// Users module
usersGroup := rtr.NewGroup().
    SetPrefix("/users").
    SetName("Users").
    AddBeforeMiddlewares([]rtr.MiddlewareInterface{
        rtr.NewAnonymousMiddleware(userAuthMiddleware),
    }).
    AddRoute(rtr.Get("/", listUsersHandler)).
    AddRoute(rtr.Post("/", createUserHandlerHandler)).
    AddRoute(rtr.Get("/:id", getUserHandler))

// Products module
productsGroup := rtr.NewGroup().
    SetPrefix("/products").
    SetName("Products").
    AddBeforeMiddlewares([]rtr.MiddlewareInterface{
        rtr.NewAnonymousMiddleware(productAuthMiddleware),
    }).
    AddRoute(rtr.Get("/", listProductsHandler)).
    AddRoute(rtr.Post("/", createProductHandler))

router.AddGroups([]rtr.GroupInterface{usersGroup, productsGroup})
```

### Administrative Interfaces

```go
// Admin group with authentication
adminGroup := rtr.NewGroup().
    SetPrefix("/admin").
    SetName("Admin").
    AddBeforeMiddlewares([]rtr.MiddlewareInterface{
        rtr.NewAnonymousMiddleware(adminAuthMiddleware),
        rtr.NewAnonymousMiddleware(adminLoggingMiddleware),
    }).
    AddRoute(rtr.Get("/dashboard", adminDashboardHandler)).
    AddRoute(rtr.Get("/users", adminUsersHandler))

// Nested admin modules
usersAdminGroup := rtr.NewGroup().
    SetPrefix("/users").
    SetName("User Admin").
    AddRoute(rtr.Get("/", adminListUsersHandler)).
    AddRoute(rtr.Post("/:id/ban", adminBanUserHandler))

adminGroup.AddGroup(usersAdminGroup)
router.AddGroup(adminGroup)
```

## Group Metadata

### Adding Metadata

```go
apiGroup := rtr.NewGroup().
    SetPrefix("/api/v1").
    SetName("API v1").
    SetMetadata(map[string]interface{}{
        "version": "1.0",
        "deprecated": false,
        "description": "Version 1 of the API",
        "maintainer": "api-team@example.com",
    })
```

### Using Metadata

```go
// Access group metadata
metadata := apiGroup.GetMetadata()
version := metadata["version"].(string)
deprecated := metadata["deprecated"].(bool)
```

## Performance Considerations

### Group Matching Efficiency

- Groups are matched in order of addition
- More specific groups should be added first
- Consider the depth of nesting for performance

```go
// Good: Specific groups first
router.AddGroup(specificGroup)    // /api/v1/users
router.AddGroup(generalGroup)      // /api/v1

// Avoid: General groups first (less efficient)
router.AddGroup(generalGroup)      // /api/v1
router.AddGroup(specificGroup)    // /api/v1/users
```

### Middleware Optimization

- Place frequently used middleware higher in the hierarchy
- Avoid duplicating middleware across nested groups
- Consider the cost of middleware when designing group structure

## Testing Groups

### Unit Testing Groups

```go
func TestGroupCreation(t *testing.T) {
    group := rtr.NewGroup().
        SetPrefix("/api").
        SetName("API").
        AddRoute(rtr.Get("/users", usersHandler))
    
    assert.Equal(t, "/api", group.GetPrefix())
    assert.Equal(t, "API", group.GetName())
    assert.Len(t, group.GetRoutes(), 1)
}

func TestGroupMiddleware(t *testing.T) {
    var executionOrder []string
    
    group := rtr.NewGroup().
        SetPrefix("/api").
        AddBeforeMiddlewares([]rtr.MiddlewareInterface{
            rtr.NewAnonymousMiddleware(func(next http.Handler) http.Handler {
                return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                    executionOrder = append(executionOrder, "group-before")
                    next.ServeHTTP(w, r)
                })
            }),
        }).
        AddRoute(rtr.Get("/test", func(w http.ResponseWriter, r *http.Request) {
            executionOrder = append(executionOrder, "handler")
            w.WriteHeader(http.StatusOK)
        }))
    
    router := rtr.NewRouter()
    router.AddGroup(group)
    
    req := httptest.NewRequest("GET", "/api/test", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    expected := []string{"group-before", "handler"}
    assert.Equal(t, expected, executionOrder)
}
```

### Integration Testing Groups

```go
func TestGroupIntegration(t *testing.T) {
    router := rtr.NewRouter()
    
    apiGroup := rtr.NewGroup().
        SetPrefix("/api/v1").
        AddBeforeMiddlewares([]rtr.MiddlewareInterface{
            rtr.NewAnonymousMiddleware(func(next http.Handler) http.Handler {
                return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                    w.Header().Set("X-API-Version", "v1")
                    next.ServeHTTP(w, r)
                })
            }),
        }).
        AddRoute(rtr.Get("/users", func(w http.ResponseWriter, r *http.Request) {
            w.Write([]byte("users"))
        }))
    
    router.AddGroup(apiGroup)
    
    req := httptest.NewRequest("GET", "/api/v1/users", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    assert.Equal(t, "v1", w.Header().Get("X-API-Version"))
    assert.Equal(t, "users", w.Body.String())
}
```

## Best Practices

### 1. Organize by Feature

```go
// Good: Feature-based organization
usersGroup := rtr.NewGroup().SetPrefix("/users")
productsGroup := rtr.NewGroup().SetPrefix("/products")
ordersGroup := rtr.NewGroup().SetPrefix("/orders")

// Avoid: Mixed organization
mixedGroup := rtr.NewGroup().SetPrefix("/api")
// Contains users, products, and orders mixed together
```

### 2. Use Consistent Prefixes

```go
// Good: Consistent API prefixing
apiGroup := rtr.NewGroup().SetPrefix("/api")
v1Group := rtr.NewGroup().SetPrefix("/v1")
usersGroup := rtr.NewGroup().SetPrefix("/users")

// Result: /api/v1/users
```

### 3. Leverage Middleware Inheritance

```go
// Good: Shared middleware at group level
apiGroup := rtr.NewGroup().
    SetPrefix("/api").
    AddBeforeMiddlewares([]rtr.MiddlewareInterface{
        rtr.NewAnonymousMiddleware(corsMiddleware),
        rtr.NewAnonymousMiddleware(rateLimitMiddleware),
    })

// Avoid: Duplicating middleware on each route
// apiGroup.AddRoute(rtr.Get("/users", handler).
//     AddBeforeMiddlewares([]rtr.MiddlewareInterface{corsMiddleware, rateLimitMiddleware}))
```

### 4. Use Descriptive Names

```go
// Good: Descriptive group names
apiGroup := rtr.NewGroup().SetName("API v1")
adminGroup := rtr.NewGroup().SetName("Admin Panel")
usersGroup := rtr.NewGroup().SetName("User Management")

// Include metadata for documentation
apiGroup.SetMetadata(map[string]interface{}{
    "description": "Version 1 of the REST API",
    "maintainer": "api-team@example.com",
})
```

### 5. Plan for Versioning

```go
// Good: Structure for versioning
apiGroup := rtr.NewGroup().SetPrefix("/api")
v1Group := rtr.NewGroup().SetPrefix("/v1")
v2Group := rtr.NewGroup().SetPrefix("/v2")

// Easy to add new versions
v3Group := rtr.NewGroup().SetPrefix("/v3")
apiGroup.AddGroup(v3Group)
```

## See Also

- [Router Core Module](router_core.md) - Main router component
- [Routes Module](routes.md) - Route management and matching
- [Middleware Module](middleware.md) - Middleware system
- [Domains Module](domains.md) - Domain-based routing
- [API Reference](../api_reference.md) - Complete API documentation
