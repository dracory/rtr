# Declarative Router Example

This example demonstrates the powerful declarative configuration approach of the rtr router, showcasing how to build complex routing structures using configuration objects instead of imperative code.

## Features

- **Pure Declarative Configuration** - Define entire router structure using `RouterConfig`
- **Hybrid Approach** - Mix declarative and imperative patterns
- **Nested Route Groups** - Organize routes with hierarchical grouping
- **Domain-Based Routing** - Handle multiple domains with different route sets
- **Middleware Integration** - Apply middleware at router, group, and route levels
- **Route Metadata** - Attach custom metadata to routes for documentation and tooling
- **Multiple Handler Types** - Support for standard HTTP handlers
- **Comprehensive Testing** - Full test coverage demonstrating all features

## Quick Start

1. **Run the example:**
   ```bash
   go run main.go
   ```

2. **Test the endpoints:**
   ```bash
   # Main API endpoints
   curl http://localhost:8080/
   curl http://localhost:8080/api/users
   curl -X POST http://localhost:8080/api/users
   curl http://localhost:8080/api/v1/products
   curl http://localhost:8080/api/v2/products
   
   # Domain-specific endpoints (add to /etc/hosts first)
   curl -H "Host: admin.example.com" http://localhost:8080/
   curl -H "Host: admin.example.com" http://localhost:8080/api/stats
   ```

3. **Run tests:**
   ```bash
   go test -v
   ```

## Configuration Structure

### RouterConfig

The main configuration object that defines the entire router:

```go
config := rtr.RouterConfig{
    Name: "My API Router",
    BeforeMiddleware: []rtr.Middleware{
        loggingMiddleware,
    },
    Routes: []rtr.RouteConfig{
        // Direct routes
    },
    Groups: []rtr.GroupConfig{
        // Route groups
    },
    Domains: []rtr.DomainConfig{
        // Domain-specific routing
    },
}
```

### Route Configuration

Individual routes can be configured with:

```go
rtr.GET("/users", handler).
    WithName("List Users").
    WithMetadata("version", "1.0").
    WithBeforeMiddleware(authMiddleware)
```

### Group Configuration

Organize related routes into groups:

```go
rtr.Group("/api",
    rtr.GET("/users", usersHandler),
    rtr.POST("/users", createUserHandler),
    
    // Nested groups
    rtr.Group("/v1",
        rtr.GET("/products", productsHandler),
    ).WithName("API V1"),
).WithName("API Group").
  WithBeforeMiddleware(jsonMiddleware)
```

### Domain Configuration

Handle multiple domains with different route sets:

```go
rtr.Domain([]string{"admin.example.com", "*.admin.example.com"},
    rtr.GET("/", adminHandler),
    rtr.Group("/api",
        rtr.GET("/stats", statsHandler),
    ),
)
```

## Examples Included

### 1. Pure Declarative Router

Demonstrates building a complete router using only declarative configuration:

- Nested route groups (`/api`, `/api/v1`, `/api/v2`)
- Multiple middleware layers
- Domain-based routing
- Route metadata and naming

### 2. Hybrid Router

Shows how to combine declarative configuration with imperative route additions:

- Start with declarative base configuration
- Add additional routes imperatively
- Mix configuration styles as needed

### 3. Imperative Router (for comparison)

Traditional imperative approach for comparison:

- Manual route and group creation
- Explicit router assembly
- Shows the difference in code organization

## Middleware Examples

The example includes several middleware patterns:

### Logging Middleware
```go
loggingMiddleware := func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("[%s] %s %s\n", r.Method, r.URL.Path, r.RemoteAddr)
        next.ServeHTTP(w, r)
    })
}
```

### Authentication Middleware
```go
authMiddleware := func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("Authorization") == "" {
            w.Header().Set("X-Auth-Required", "true")
        }
        next.ServeHTTP(w, r)
    })
}
```

### Content-Type Middleware
```go
jsonMiddleware := func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        next.ServeHTTP(w, r)
    })
}
```

## Testing

The example includes comprehensive tests covering:

- **Route Resolution** - Verify all routes respond correctly
- **Middleware Execution** - Test middleware behavior
- **Domain Routing** - Validate domain-based routing
- **HTTP Methods** - Test GET, POST, and other methods
- **Error Handling** - Verify proper error responses
- **Header Handling** - Test custom headers and content types

Run tests with:
```bash
go test -v                    # Run all tests
go test -run TestDeclarative  # Run specific test
go test -cover               # Run with coverage
```

## Domain Setup

To test domain-based routing, add these entries to your hosts file:

**Unix/Mac (`/etc/hosts`):**
```
127.0.0.1 admin.example.com
```

**Windows (`C:\Windows\System32\drivers\etc\hosts`):**
```
127.0.0.1 admin.example.com
```

## Advanced Features

### Route Metadata

Attach custom metadata to routes for documentation, versioning, or tooling:

```go
rtr.GET("/products", handler).
    WithMetadata("version", "2.0").
    WithMetadata("deprecated", "false").
    WithMetadata("rateLimit", "100/hour")
```

### Wildcard Domains

Support wildcard domain matching:

```go
rtr.Domain([]string{"*.admin.example.com"}, ...)
```

### Middleware Ordering

Control middleware execution order:

```go
config := rtr.RouterConfig{
    BeforeMiddleware: []rtr.Middleware{
        loggingMiddleware,    // Executes first
        authMiddleware,      // Executes second
    },
    // ...
}
```

## Benefits of Declarative Configuration

1. **Readability** - Router structure is immediately visible
2. **Maintainability** - Easy to modify and extend
3. **Testing** - Configuration can be easily tested
4. **Documentation** - Self-documenting code structure
5. **Tooling** - Configuration can be analyzed by external tools
6. **Serialization** - Configuration can be loaded from files
7. **Validation** - Configuration can be validated before use

## Comparison with Imperative Approach

| Aspect | Declarative | Imperative |
|--------|-------------|------------|
| Code Organization | Structured, hierarchical | Linear, procedural |
| Readability | High - structure is visible | Medium - requires reading through code |
| Maintainability | High - easy to modify | Medium - requires careful editing |
| Testing | Easy - test configuration | Medium - test assembled router |
| Flexibility | High - can be loaded from files | Medium - hardcoded in Go |
| Learning Curve | Low - follows familiar patterns | Medium - requires understanding API |

## Best Practices

1. **Group Related Routes** - Use groups to organize related functionality
2. **Apply Middleware Strategically** - Use appropriate middleware at the right level
3. **Use Meaningful Names** - Name routes and groups for better debugging
4. **Add Metadata** - Include version, description, and other useful metadata
5. **Test Thoroughly** - Write comprehensive tests for all routes and middleware
6. **Document Domains** - Clearly document domain requirements and setup
7. **Validate Configuration** - Check configuration before creating router

## Troubleshooting

### Common Issues

1. **Domain routing not working**
   - Check hosts file configuration
   - Verify Host header in requests
   - Ensure domain patterns match exactly

2. **Middleware not executing**
   - Check middleware order
   - Verify middleware is properly attached
   - Ensure middleware calls `next.ServeHTTP()`

3. **Routes not matching**
   - Verify route patterns
   - Check HTTP method matching
   - Ensure no conflicting routes

### Debug Tips

1. **Enable Logging** - Use logging middleware to trace requests
2. **List Routes** - Call `router.List()` to see registered routes
3. **Test Incrementally** - Build configuration step by step
4. **Use Tests** - Write tests to verify expected behavior

## Related Examples

- [Domain Example](../domain/) - Focus on domain-based routing
- [Handlers Example](../handlers/) - Different handler types (HTML, JSON, etc.)
- [Middleware Example](../middleware/) - Advanced middleware patterns

## Contributing

When contributing to this example:

1. Add tests for new features
2. Update documentation
3. Follow Go conventions
4. Ensure backward compatibility
5. Add meaningful examples
