# Dracory Router vs Gin Router: Routing Capabilities Comparison

## Overview
This document focuses specifically on comparing the routing capabilities of Dracory Router and Gin's router component, ignoring other framework features.

## Feature Comparison Table

| Feature | Dracory Router | Gin Router |
|---------|----------------|------------|
| **Basic Routing** | ✅ All HTTP methods | ✅ All HTTP methods |
| **Path Parameters** | ✅ Supports named parameters with optional segments (`/user/:name`, `/articles/:category/:id?`) | ✅ Supports named parameters (`/user/:name`) |
| **Wildcards** | ✅ Supports catch-all routes (`/static/*filepath`) | ✅ Supports wildcard parameters (`/user/*action`) |
| **Route Groups** | ✅ Full support with nesting | ✅ Full support with nesting |
| **Domain Routing** | ✅ Full support for domain-based routing | ❌ No built-in support |
| **Middleware** | ✅ Before/after middleware with per-route, group, and global support | ✅ Standard middleware support (no built-in after middleware) |
| **Middleware Chaining** | ✅ Supports chaining | ✅ Supports chaining |
| **Performance** | Lightweight, minimal overhead | Highly optimized, uses httprouter |
| **Static Files** | ❌ No built-in support | ✅ Built-in static file serving |
| **Custom Not Found** | ✅ Supported | ✅ Supported |
| **Custom Method Not Allowed** | ❌ Not supported | ✅ Supported |
| **Route Matching** | Supports parameters, optional parameters, and wildcards | Supports parameters and wildcards |
| **Route Naming** | ✅ Supported | ❌ Not supported |
| **Middleware per Route** | ✅ Supported | ✅ Supported |

## Code Examples

### Basic Route Definition

**Dracory Router**
```go
router := rtr.NewRouter()
router.Get("/users", usersHandler)
router.Post("/users", createUserHandler)
```

**Gin**
```go
router := gin.New()
router.GET("/users", usersHandler)
router.POST("/users", createUserHandler)
```

### Route Groups

**Dracory Router**
```go
api := rtr.NewGroup().SetPrefix("/api")
api.AddRoute(rtr.NewRoute().SetMethod("GET").SetPath("/users").SetHandler(usersHandler))
router.AddGroup(api)
```

**Gin**
```go
api := router.Group("/api")
{
    api.GET("/users", usersHandler)
}
```

### Domain-Based Routing

**Dracory Router**
```go
domain := rtr.NewDomain("api.example.com")
domain.AddRoute(rtr.NewRoute().SetMethod("GET").SetPath("/users").SetHandler(apiUsersHandler))
router.AddDomain(domain)
```

**Gin**
```go
// Gin doesn't have built-in domain routing
// You would need to implement this manually or use middleware
```

### Middleware

**Dracory Router**
```go
// Global middleware
router.AddBeforeMiddlewares([]rtr.Middleware{loggerMiddleware})

// Per-route middleware
route := rtr.NewRoute().SetMethod("GET").SetPath("/secure").SetHandler(secureHandler)
route.AddBeforeMiddlewares([]rtr.Middleware{authMiddleware})
```

**Gin**
```go
// Global middleware
router.Use(gin.Logger())

// Per-route middleware
router.GET("/secure", authMiddleware, secureHandler)
```

## Performance Considerations

### Dracory Router
- Uses standard library's `http.ServeMux` for routing
- Minimal overhead due to simple route matching
- No regular expressions or complex path matching

### Gin Router
- Uses a custom version of `httprouter`
- Very fast route matching using a radix tree
- Supports more complex routing patterns with parameters

## When to Choose Which

### Choose Dracory Router if:
- You need domain-based routing
- You prefer explicit route configuration
- You want a lightweight, focused routing solution
- You're building a custom framework

### Choose Gin Router if:
- You need parameterized routes
- You want built-in support for wildcards
- Performance is critical
- You need more advanced routing features

## Conclusion

Both routers have their strengths:
- **Dracory Router** excels in domain-based routing and explicit configuration
- **Gin Router** provides more advanced path matching capabilities and better performance

The choice depends on your specific routing needs and whether you value simplicity and domain routing (Dracory) or advanced path matching and performance (Gin).
