# dracory/rtr vs gouniverse/router: Technical Comparison

## 1. Overview

### dracory/rtr
A lightweight, flexible HTTP router for Go with a focus on declarative configuration and middleware chaining. It provides a clean API for building web applications with support for domain-based routing, middleware pipelines, and various response types.

### gouniverse/router
A minimal, fast HTTP router designed for the Gouniverse framework. It emphasizes simplicity and performance while providing essential routing capabilities for web applications.

### Project Status
- **dracory/rtr**: Actively maintained, with recent updates focusing on middleware improvements and type safety.
- **gouniverse/router**: Part of the larger Gouniverse ecosystem, with stable but less frequent updates.

## 2. Feature Comparison Table

| Feature | dracory/rtr | gouniverse/router |
|---------|-------------|-------------------|
| **Basic Routing** | ✅ All HTTP methods with exact path matching | ✅ All HTTP methods with exact path matching |
| **Path Parameters** | ✅ Named parameters (`/users/:id`) with optional segments | ✅ Basic parameter support (`/users/:id`) |
| **Wildcard Routes** | ✅ Full support (`/static/*filepath`) | ✅ Basic support (`/*`) |
| **Route Groups** | ✅ Full support with nesting and shared middleware | ✅ Basic grouping with path prefixing |
| **Domain Routing** | ✅ Full support with wildcard subdomains | ❌ Not supported |
| **Middleware** | ✅ Before/after middleware at all levels | ✅ Global and per-route middleware |
| **Middleware Types** | ✅ Supports both `StdMiddleware` and `MiddlewareInterface` | ✅ Standard middleware pattern |
| **Response Helpers** | ✅ Built-in for HTML, JSON, CSS, XML, Text, JS | ✅ Automatic wrapping for HTML/JSON |
| **Performance** | ⚡ Lightweight, optimized for speed | ⚡ Minimal overhead, fast routing |
| **Testing** | ✅ Standard Go testing patterns | ✅ Easy testing with string returns |
| **Documentation** | 📚 Comprehensive with examples | 📚 Basic documentation |
| **Community** | 👥 Growing community | 👥 Part of Gouniverse ecosystem |

## 3. Code Examples

### Basic Route Definition

**dracory/rtr**
```go
router := rtr.NewRouter()
router.AddRoute(rtr.Get("/users", usersHandler))
router.AddRoute(rtr.Post("/users", createUserHandler))

// Using specialized handlers
router.AddRoute(rtr.GetHTML("/page", func(w http.ResponseWriter, r *http.Request) string {
    return "<h1>Hello World</h1>"
}))
```

**gouniverse/router**
```go
r := router.New()
r.GET("/users", usersHandler)
r.POST("/users", createUserHandler)

// Using handler with string return
r.HTML("/page", func(w http.ResponseWriter, r *http.Request) string {
    return "<h1>Hello World</h1>"
})
```

### Middleware Usage

**dracory/rtr**
```go
// Create middleware
authMiddleware := rtr.NewAnonymousMiddleware(func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Authentication logic
        next.ServeHTTP(w, r)
    })
})

// Apply to route
router.AddRoute(rtr.Get("/profile", profileHandler).WithBeforeMiddleware(authMiddleware))

// Or to group
group := rtr.NewGroup("/api")
group.WithBeforeMiddleware(authMiddleware)
group.AddRoute(rtr.Get("/users", usersHandler))
```

**gouniverse/router**
```go
// Middleware function
authMiddleware := func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Authentication logic
        next.ServeHTTP(w, r)
    })
}

// Apply to route
r.GET("/profile", authMiddleware(profileHandler))

// Apply to group
api := r.Group("/api")
api.Use(authMiddleware)
api.GET("/users", usersHandler)
```

### Route Groups

**dracory/rtr**
```go
api := rtr.NewGroup("/api")
api.WithBeforeMiddleware(loggingMiddleware)

// Nested groups
v1 := api.Group("/v1")
v1.AddRoute(rtr.Get("/users", getUsersV1))

// Domain-based routing
domain := rtr.NewDomain("api.example.com")
domain.AddRoute(rtr.Get("/status", statusHandler))
router.AddDomain(domain)
```

**gouniverse/router**
```go
api := r.Group("/api")
api.Use(loggingMiddleware)

// Nested groups
v1 := api.Group("/v1")
v1.GET("/users", getUsersV1)

// No built-in domain support
// Would need to use separate router instances
```

## 4. Performance Comparison

### dracory/rtr
- **Memory Usage**: Low memory footprint
- **Request Handling**: Optimized for fast route matching
- **Benchmarks**: Slightly faster in benchmarks due to simpler routing algorithm

### gouniverse/router
- **Memory Usage**: Minimal overhead
- **Request Handling**: Fast, but may slow down with many routes
- **Benchmarks**: Slightly slower than dracory/rtr in direct comparisons

## 5. Use Cases

### Choose dracory/rtr when:
- You need domain-based routing
- You want a more flexible middleware system
- You need advanced routing features like optional segments
- You prefer a more type-safe approach

### Choose gouniverse/router when:
- You want a simple, minimal router
- You're already using other Gouniverse components
- You don't need domain routing
- You prefer a more traditional middleware approach

## 6. Migration Guide

### From gouniverse/router to dracory/rtr:
1. **Update Imports**: Change import paths
2. **Route Registration**: Update to use `AddRoute` with `rtr.Get()`, `rtr.Post()`, etc.
3. **Middleware**: Convert middleware to use `rtr.NewAnonymousMiddleware`
4. **Response Helpers**: Update to use dracory's response helpers
5. **Groups**: Update group creation and middleware application

### Key Differences to Watch For:
- Different method chaining patterns
- Middleware signature differences
- Response handling variations

## 7. Conclusion

### dracory/rtr
**Strengths**:
- More feature-rich with domain routing
- Flexible middleware system
- Better type safety
- Active development

**Weaknesses**:
- Slightly steeper learning curve
- More complex API

### gouniverse/router
**Strengths**:
- Simpler API
- Easier to get started
- Good for basic routing needs

**Weaknesses**:
- Lacks advanced features like domain routing
- Less flexible middleware system

### Final Recommendation
- For complex applications needing advanced routing features, choose **dracory/rtr**
- For simpler applications or when already using Gouniverse, **gouniverse/router** is a good choice
