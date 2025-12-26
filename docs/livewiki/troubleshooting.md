---
path: troubleshooting.md
page-type: tutorial
summary: Common issues, debugging techniques, and solutions for RTR router problems.
tags: [troubleshooting, debugging, issues, solutions, faq]
created: 2025-12-26
updated: 2025-12-26
version: 1.0.0
---

# RTR Router Troubleshooting Guide

This guide covers common issues, debugging techniques, and solutions for problems you might encounter while using the RTR router.

## Common Issues

### Route Not Matching

#### Problem
Routes are not being matched and returning 404 errors.

#### Common Causes

**1. Method Mismatch**
```go
// Route defined as GET
router.AddRoute(rtr.Get("/users", handler))

// But request is POST
// This will NOT match
```

**Solution:**
```go
// Add POST route or use proper method
router.AddRoute(rtr.Post("/users", postHandler))
```

**2. Path Pattern Issues**
```go
// Route with trailing slash
router.AddRoute(rtr.Get("/users/", handler))

// Request without trailing slash
// GET /users  -> 404
```

**Solution:**
```go
// Use consistent paths or add both routes
router.AddRoute(rtr.Get("/users", handler))
router.AddRoute(rtr.Get("/users/", handler))

// Or use redirect middleware
router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    rtr.NewAnonymousMiddleware(middlewares.RedirectSlashes),
})
```

**3. Parameter Extraction Errors**
```go
// Route with required parameter
router.AddRoute(rtr.Get("/users/:id", handler))

// Request without parameter
// GET /users  -> 404
```

**Solution:**
```go
// Use optional parameter if needed
router.AddRoute(rtr.Get("/users/:id?", handler))

// Or add separate route
router.AddRoute(rtr.Get("/users", listUsersHandler))
router.AddRoute(rtr.Get("/users/:id", getUserHandler))
```

#### Debugging Route Matching

```go
// Add debugging middleware to see what routes are being checked
router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    rtr.NewAnonymousMiddleware(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            fmt.Printf("Request: %s %s\n", r.Method, r.URL.Path)
            next.ServeHTTP(w, r)
        })
    }),
})

// List all routes to verify configuration
router.List()
```

### Middleware Not Executing

#### Problem
Middleware is not running as expected.

#### Common Causes

**1. Wrong Middleware Type**
```go
// Trying to use StdMiddleware where MiddlewareInterface is expected
router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    loggingMiddleware, // This is StdMiddleware, not MiddlewareInterface
})
```

**Solution:**
```go
// Convert to MiddlewareInterface
router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    rtr.NewAnonymousMiddleware(loggingMiddleware),
    // or
    rtr.NewMiddleware("Logger", loggingMiddleware),
})
```

**2. Middleware Order Issues**
```go
// After middleware added before before middleware
router.AddAfterMiddlewares([]rtr.MiddlewareInterface{
    loggingMiddleware, // This runs AFTER the handler
})
```

**Solution:**
```go
// Use before middleware for preprocessing
router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    loggingMiddleware,
})
```

**3. Group-Level Middleware Not Inherited**
```go
// Middleware on group doesn't apply to direct routes
apiGroup := rtr.NewGroup().SetPrefix("/api")
apiGroup.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    authMiddleware,
})

// This route won't get auth middleware
router.AddRoute(rtr.Get("/api/health", healthHandler))
```

**Solution:**
```go
// Add route to group, not router
apiGroup.AddRoute(rtr.Get("/health", healthHandler))
router.AddGroup(apiGroup)
```

#### Debugging Middleware

```go
// Add debugging middleware to trace execution
func debugMiddleware(name string) rtr.MiddlewareInterface {
    return rtr.NewAnonymousMiddleware(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            fmt.Printf("Entering middleware: %s\n", name)
            start := time.Now()
            next.ServeHTTP(w, r)
            duration := time.Since(start)
            fmt.Printf("Exiting middleware: %s (took %v)\n", name, duration)
        })
    })
}

// Add at different levels
router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    debugMiddleware("global-before"),
})

group.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    debugMiddleware("group-before"),
})
```

### Parameter Extraction Issues

#### Problem
Path parameters are not being extracted correctly.

#### Common Causes

**1. Parameter Name Mismatch**
```go
// Route defined with :id
router.AddRoute(rtr.Get("/users/:id", handler))

// But trying to get :user_id
userID := rtr.MustGetParam(r, "user_id") // Panic!
```

**Solution:**
```go
// Use correct parameter name
userID := rtr.MustGetParam(r, "id")
```

**2. Optional Parameter Issues**
```go
// Route with optional parameter
router.AddRoute(rtr.Get("/users/:id?", handler))

// MustGetParam on optional parameter
id := rtr.MustGetParam(r, "id") // Panic if not provided!
```

**Solution:**
```go
// Use GetParam for optional parameters
if id, exists := rtr.GetParam(r, "id"); exists {
    // Use the parameter
}
```

**3. Greedy Parameter Issues**
```go
// Greedy parameter not at end
router.AddRoute(rtr.Get("/files/:path.../download", handler)) // Invalid!
```

**Solution:**
```go
// Greedy parameter must be last
router.AddRoute(rtr.Get("/files/:path...", handler))
```

#### Debugging Parameters

```go
router.AddRoute(rtr.Get("/users/:id/posts/:post_id", func(w http.ResponseWriter, r *http.Request) {
    // Log all parameters
    params := rtr.GetParams(r)
    fmt.Printf("All parameters: %+v\n", params)
    
    // Log individual parameters
    id, hasID := rtr.GetParam(r, "id")
    postID, hasPostID := rtr.GetParam(r, "post_id")
    fmt.Printf("ID: %s (exists: %v), PostID: %s (exists: %v)\n", 
        id, hasID, postID, hasPostID)
    
    // Continue with handler logic
}))
```

### Handler Not Working

#### Problem
Handler is not executing or not producing expected output.

#### Common Causes

**1. Handler Priority Issues**
```go
// Multiple handlers set
route := rtr.NewRoute().
    SetHandler(standardHandler).
    SetJSONHandler(jsonHandler) // This takes priority!

// Only JSONHandler will execute
```

**Solution:**
```go
// Set only one handler type
route := rtr.NewRoute().
    SetJSONHandler(jsonHandler)
```

**2. Response Already Written**
```go
// Middleware writes response but doesn't stop chain
router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    rtr.NewAnonymousMiddleware(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte("Unauthorized"))
            next.ServeHTTP(w, r) // Handler still runs!
        })
    }),
})
```

**Solution:**
```go
// Stop chain after writing response
router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    rtr.NewAnonymousMiddleware(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !isAuthorized(r) {
                w.WriteHeader(http.StatusUnauthorized)
                w.Write([]byte("Unauthorized"))
                return // Don't call next.ServeHTTP
            }
            next.ServeHTTP(w, r)
        })
    }),
})
```

#### Debugging Handlers

```go
// Add debugging to handler
router.AddRoute(rtr.Get("/debug", func(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("Handler executing for: %s %s\n", r.Method, r.URL.Path)
    
    // Check if response writer is already committed
    if rw, ok := w.(interface{ Status() int }); ok {
        fmt.Printf("Response status: %d\n", rw.Status())
    }
    
    w.Write([]byte("Debug info logged"))
}))
```

### Domain Routing Issues

#### Problem
Domain-based routing is not working as expected.

#### Common Causes

**1. Host Header Issues**
```go
// Domain configured for example.com
domain := rtr.NewDomain("example.com")

// But request comes to localhost:8080
// Won't match
```

**Solution:**
```go
// Configure for localhost during development
domain := rtr.NewDomain("localhost", "example.com")

// Or use wildcard
domain := rtr.NewDomain("*")
```

**2. Port Mismatch**
```go
// Domain configured without port
domain := rtr.NewDomain("example.com")

// Request comes to example.com:8080
// Will match (port is ignored when not specified)
```

**Solution:**
```go
// Be explicit about ports if needed
domain := rtr.NewDomain("example.com:8080")
// Or use wildcard port
domain := rtr.NewDomain("example.com:*")
```

#### Debugging Domain Matching

```go
// Add middleware to debug domain matching
router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    rtr.NewAnonymousMiddleware(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            host := r.Host
            fmt.Printf("Request host: %s\n", host)
            
            // Check which domain would match
            for _, domain := range router.GetDomains() {
                if domain.Matches(host) {
                    fmt.Printf("Matches domain: %s\n", domain.GetName())
                    break
                }
            }
            
            next.ServeHTTP(w, r)
        })
    }),
})
```

## Performance Issues

### Slow Route Matching

#### Problem
Router is slow, especially with many routes.

#### Solutions

**1. Optimize Route Order**
```go
// Put most specific routes first
router.AddRoute(rtr.Get("/api/v1/users/123", specificUserHandler))
router.AddRoute(rtr.Get("/api/v1/users/:id", userHandler))
router.AddRoute(rtr.Get("/api/v1/users", usersHandler))
```

**2. Use Groups for Organization**
```go
// Instead of many routes with same prefix
apiGroup := rtr.NewGroup().SetPrefix("/api/v1")
apiGroup.AddRoute(rtr.Get("/users", usersHandler))
apiGroup.AddRoute(rtr.Get("/users/:id", userHandler))
router.AddGroup(apiGroup)
```

**3. Profile with Benchmarks**
```go
func BenchmarkRouter(b *testing.B) {
    router := setupLargeRouter()
    req := httptest.NewRequest("GET", "/api/v1/users/123", nil)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)
    }
}
```

### Memory Leaks

#### Problem
Memory usage increases over time.

#### Common Causes

**1. Middleware Holding References**
```go
// Middleware storing request data in global variable
var requests []*http.Request

router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    rtr.NewAnonymousMiddleware(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            requests = append(requests, r) // Memory leak!
            next.ServeHTTP(w, r)
        })
    }),
})
```

**Solution:**
```go
// Use bounded storage or cleanup
var requestPool = sync.Pool{
    New: func() interface{} {
        return make([]*http.Request, 0, 100)
    },
}

router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    rtr.NewAnonymousMiddleware(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            requests := requestPool.Get().([]*http.Request)
            requests = append(requests, r)
            defer func() {
                // Clear and return to pool
                requests = requests[:0]
                requestPool.Put(requests)
            }()
            next.ServeHTTP(w, r)
        })
    }),
})
```

## Testing Issues

### Test Failures

#### Problem
Tests are failing intermittently or unexpectedly.

#### Common Causes

**1. Race Conditions**
```go
// Test modifying shared state
var counter int

func TestConcurrentRequests(t *testing.T) {
    router := setupRouter()
    
    for i := 0; i < 100; i++ {
        go func() {
            router.AddRoute(rtr.Get(fmt.Sprintf("/route%d", i), handler))
        }()
    }
    // Race condition!
}
```

**Solution:**
```go
// Use proper synchronization
var mu sync.Mutex
var counter int

func TestConcurrentRequests(t *testing.T) {
    router := setupRouter()
    
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            mu.Lock()
            router.AddRoute(rtr.Get(fmt.Sprintf("/route%d", id), handler))
            mu.Unlock()
        }(i)
    }
    wg.Wait()
}
```

**2. Test Isolation**
```go
// Tests affecting each other
func TestRouteA(t *testing.T) {
    router = rtr.NewRouter() // Global variable!
    router.AddRoute(rtr.Get("/test", handler))
}

func TestRouteB(t *testing.T) {
    // Router still has routes from TestRouteA
    w := executeRequest(router, "GET", "/test") // Unexpected success!
}
```

**Solution:**
```go
// Create fresh router for each test
func TestRouteA(t *testing.T) {
    router := rtr.NewRouter() // Local variable
    router.AddRoute(rtr.Get("/test", handler))
}

func TestRouteB(t *testing.T) {
    router := rtr.NewRouter() // Fresh router
    // Test logic
}
```

## Debugging Tools

### Built-in Debugging

#### Route Listing
```go
router := setupRouter()
router.List() // Prints detailed configuration
```

#### String Representation
```go
fmt.Println(router.String()) // Get router as string
```

### Custom Debugging

#### Request Tracing
```go
func requestTracingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Capture request details
        fmt.Printf("=== Request Start ===\n")
        fmt.Printf("Method: %s\n", r.Method)
        fmt.Printf("Path: %s\n", r.URL.Path)
        fmt.Printf("Host: %s\n", r.Host)
        fmt.Printf("Headers: %+v\n", r.Header)
        
        // Wrap response writer to capture response
        rw := &responseWriter{ResponseWriter: w, statusCode: 200}
        
        next.ServeHTTP(rw, r)
        
        duration := time.Since(start)
        fmt.Printf("Status: %d\n", rw.statusCode)
        fmt.Printf("Duration: %v\n", duration)
        fmt.Printf("=== Request End ===\n\n")
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
```

#### Middleware Chain Tracing
```go
func chainTracingMiddleware(name string) rtr.MiddlewareInterface {
    return rtr.NewAnonymousMiddleware(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            fmt.Printf("→ Entering %s\n", name)
            defer fmt.Printf("← Exiting %s\n", name)
            next.ServeHTTP(w, r)
        })
    })
}

// Usage
router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    chainTracingMiddleware("global-middleware-1"),
    chainTracingMiddleware("global-middleware-2"),
})
```

## Environment-Specific Issues

### Development vs Production

#### Problem
Code works in development but fails in production.

#### Common Causes

**1. Different Host Headers**
```go
// Development: localhost:8080
// Production: api.example.com
```

**Solution:**
```go
// Configure for both environments
domains := []string{"localhost", "api.example.com"}
if port := os.Getenv("PORT"); port != "" {
    domains = append(domains, "localhost:"+port)
}
domain := rtr.NewDomain(domains...)
```

**2. Different Paths**
```go
// Development: /
// Production: /app/
```

**Solution:**
```go
// Use configurable prefix
prefix := os.Getenv("APP_PREFIX")
if prefix == "" {
    prefix = "/"
}
router.SetPrefix(prefix)
```

### Container Environments

#### Problem
Router doesn't work in Docker/Kubernetes.

#### Solutions

**1. Health Checks**
```go
router.AddRoute(rtr.Get("/health", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status": "healthy"}`))
}))
```

**2. Proper Port Binding**
```go
port := os.Getenv("PORT")
if port == "" {
    port = "8080"
}
addr := ":" + port
log.Printf("Starting server on %s", addr)
log.Fatal(http.ListenAndServe(addr, router))
```

## FAQ

### Q: Why are my routes returning 404?
A: Check method matching, path patterns, and parameter requirements. Use `router.List()` to verify route configuration.

### Q: How do I debug middleware execution order?
A: Add debugging middleware that logs entry/exit at each level to see the execution flow.

### Q: Why isn't my domain-based routing working?
A: Verify the Host header in requests and ensure domain patterns match exactly. Check for port mismatches.

### Q: How can I improve performance with many routes?
A: Use groups for organization, put specific routes first, and profile with benchmarks.

### Q: What's the difference between StdMiddleware and MiddlewareInterface?
A: `StdMiddleware` is the standard function type, while `MiddlewareInterface` provides naming and metadata capabilities.

### Q: How do I handle CORS?
A: Use the built-in CORS middleware or implement custom CORS middleware.

### Q: Can I use RTR with existing Go middleware?
A: Yes, use `rtr.NewAnonymousMiddleware()` to wrap standard Go middleware.

### Q: How do I test my routes?
A: Use `httptest.NewRequest` and `httptest.NewRecorder` for unit testing routes.

## See Also

- [Getting Started Guide](getting_started.md) - Learn basic usage
- [API Reference](api_reference.md) - Complete API documentation
- [Development Guide](development.md) - Development workflow and testing
- [Architecture Documentation](architecture.md) - System design overview
