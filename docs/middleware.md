# Middleware Guide

## Overview

Middleware in the router provides a way to intercept and process HTTP requests and responses. This guide covers the middleware system in detail, including execution order, built-in middleware, and creating custom middleware.

## Table of Contents
- [Execution Order](#execution-order)
- [Built-in Middleware](#built-in-middleware)
- [Custom Middleware](#custom-middleware)
- [Middleware Chaining](#middleware-chaining)
- [Best Practices](#best-practices)

## Execution Order

Middleware execution follows a specific order based on its scope:

1. **Before Middleware** (in order of addition):
   - Global (Router-level) middleware
   - Domain-level middleware (if domain matches)
   - Group-level middleware (for matching groups)
   - Route-level middleware

2. **Handler Execution**

3. **After Middleware** (in order of addition):
   - Route-level middleware
   - Group-level middleware
   - Domain-level middleware
   - Global middleware

### Example Execution Flow

```
Global Middleware Before 1
 Global Middleware Before 2
  Domain Middleware Before 1
   Domain Middleware Before 2
    Group Middleware Before 1
     Group Middleware Before 2
      Route Middleware Before 1
       Route Middleware Before 2
        Handler
       Route Middleware After 1
      Route Middleware After 2
     Group Middleware After 1
    Group Middleware After 2
   Domain Middleware After 1
  Domain Middleware After 2
 Global Middleware After 1
Global Middleware After 2
```

## Built-in Middleware

### Recovery Middleware

Automatically recovers from panics in your handlers and returns a 500 Internal Server Error response.

```go
// Automatically added when creating a new router
router := router.NewRouter()

// Can be added manually if needed
router.AddBeforeMiddlewares([]router.Middleware{router.RecoveryMiddleware})
```

## Custom Middleware

Create custom middleware by implementing the `Middleware` type:

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Before handler execution
        log.Printf("Started %s %s", r.Method, r.URL.Path)
        
        // Call the next handler
        next.ServeHTTP(w, r)
        
        // After handler execution
        log.Printf("Completed %s %s in %v", 
            r.Method, r.URL.Path, time.Since(start))
    })
}
```

## Middleware Chaining

Middleware can be chained together using the router's methods:

```go
router := NewRouter()

// Global middleware
router.AddBeforeMiddlewares([]Middleware{
    loggingMiddleware,
    authMiddleware,
})

// Group with specific middleware
group := NewGroup("/api")
group.AddBeforeMiddlewares([]Middleware{apiAuthMiddleware})

// Route with specific middleware
route := NewRoute("/admin", adminHandler)
route.AddBeforeMiddlewares([]Middleware{adminOnlyMiddleware})
```

## Best Practices

1. **Keep Middleware Focused**: Each middleware should do one thing well.
2. **Order Matters**: Place middleware in the correct order (e.g., authentication before authorization).
3. **Error Handling**: Always handle errors and don't let panics propagate.
4. **Performance**: Be mindful of middleware that performs expensive operations.
5. **Logging**: Include request/response logging in appropriate middleware.
6. **Context**: Use `context.Context` to pass values between middleware and handlers.

For more advanced middleware patterns and examples, see [Advanced Middleware Patterns](./advanced-middleware.md).
