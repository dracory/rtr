# Performance Considerations

## Overview

This guide provides insights into the performance characteristics of the router and offers best practices for optimizing your application.

## Table of Contents
- [Routing Performance](#routing-performance)
- [Middleware Overhead](#middleware-overhead)
- [Concurrency](#concurrency)
- [Memory Usage](#memory-usage)
- [Caching Strategies](#caching-strategies)
- [Benchmarking](#benchmarking)
- [Production Recommendations](#production-recommendations)

## Routing Performance

The router uses a linear search for route matching, which has O(n) complexity where n is the number of routes. For most applications with a reasonable number of routes, this is perfectly adequate. However, for applications with thousands of routes, consider the following:

- **Route Order**: Place more frequently accessed routes first
- **Use Groups**: Group related routes to reduce search space
- **Avoid Overlapping Patterns**: Be mindful of route patterns that could cause excessive matching attempts

## Middleware Overhead

Each middleware adds a small overhead to request processing. To minimize impact:

```go
// Instead of adding middleware to every route:
router.AddBeforeMiddlewares([]rtr.Middleware{loggingMiddleware, authMiddleware})

// Consider scoping middleware to specific routes/groups:
adminGroup := rtr.NewGroup("/admin")
adminGroup.AddBeforeMiddlewares([]rtr.Middleware{authMiddleware, adminOnlyMiddleware})
```

## Concurrency

The router is safe for concurrent use. However, be aware of:

1. **Global State**: Avoid modifying global state in handlers without proper synchronization
2. **Middleware Initialization**: Initialize middleware dependencies before starting the server
3. **Connection Pooling**: Use connection pools for database and other external services

## Memory Usage

To minimize memory usage:

1. **Reuse Objects**: Use `sync.Pool` for frequently allocated objects
2. **Stream Large Responses**: For large responses, use streaming instead of buffering in memory
3. **Limit Request Size**: Enforce reasonable limits on request body sizes

## Caching Strategies

### Route Caching

For high-traffic applications, consider implementing a route cache:

```go
var routeCache = make(map[string]http.Handler)
var mu sync.RWMutex

func cachedHandler(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cacheKey := r.Method + " " + r.URL.Path
        
        // Check cache
        mu.RLock()
        if h, ok := routeCache[cacheKey]; ok {
            mu.RUnlock()
            h.ServeHTTP(w, r)
            return
        }
        mu.RUnlock()
        
        // Cache miss - process and cache
        recorder := httptest.NewRecorder()
        handler.ServeHTTP(recorder, r)
        
        // Only cache successful responses
        if recorder.Code >= 200 && recorder.Code < 300 {
            mu.Lock()
            routeCache[cacheKey] = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                for k, v := range recorder.Header() {
                    w.Header()[k] = v
                }
                w.WriteHeader(recorder.Code)
                w.Write(recorder.Body.Bytes())
            })
            mu.Unlock()
        }
        
        // Write the response
        for k, v := range recorder.Header() {
            w.Header()[k] = v
        }
        w.WriteHeader(recorder.Code)
        w.Write(recorder.Body.Bytes())
    })
}
```

## Benchmarking

### Basic Benchmark

```go
func BenchmarkRouter(b *testing.B) {
    router := setupTestRouter()
    req := httptest.NewRequest("GET", "/users/123", nil)
    w := httptest.NewRecorder()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        router.ServeHTTP(w, req)
    }
}
```

### Benchmark with Concurrency

```go
func BenchmarkRouterParallel(b *testing.B) {
    router := setupTestRouter()
    
    b.RunParallel(func(pb *testing.PB) {
        req := httptest.NewRequest("GET", "/users/123", nil)
        w := httptest.NewRecorder()
        
        for pb.Next() {
            router.ServeHTTP(w, req)
        }
    })
}
```

## Production Recommendations

1. **Enable HTTP/2**: For better performance with HTTPS
2. **Use a Reverse Proxy**: Like Nginx or Caddy in front of your application
3. **Monitor Performance**: Use tools like Prometheus and Grafana
4. **Profile Regularly**: Use Go's built-in profiling tools
5. **Set Timeouts**: Always set read/write timeouts
6. **Connection Limits**: Limit the number of concurrent connections
7. **Compression**: Enable response compression

## Monitoring and Metrics

### Basic Metrics with Prometheus

```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    requestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
    
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: []float64{0.1, 0.5, 1, 2.5, 5, 10},
        },
        []string{"method", "path"},
    )
)

func init() {
    prometheus.MustRegister(requestsTotal)
    prometheus.MustRegister(requestDuration)
}

func metricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Create a response writer that captures the status code
        rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
        
        // Process the request
        next.ServeHTTP(rw, r)
        
        // Record metrics
        duration := time.Since(start).Seconds()
        status := strconv.Itoa(rw.status)
        
        requestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
        requestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
    })
}

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
    http.ResponseWriter
    status int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.status = code
    rw.ResponseWriter.WriteHeader(code)
}

// In your main function:
func main() {
    router := rtr.NewRouter()
    
    // Add metrics endpoint
    router.AddRoute(rtr.Get("/metrics", promhttp.Handler().ServeHTTP))
    
    // Add metrics middleware to all routes
    router.AddBeforeMiddlewares([]rtr.Middleware{metricsMiddleware})
    
    // Add your other routes...
    
    log.Fatal(http.ListenAndServe(":8080", router))
}
```

For more advanced performance optimization techniques, see [Advanced Performance Tuning](./advanced-performance.md).
