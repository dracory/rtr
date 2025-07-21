# Built-in Middleware

## Recovery Middleware

The router includes a built-in recovery middleware that catches panics in your handlers and returns a 500 Internal Server Error response instead of crashing the server. This middleware is added by default when you create a new router with `NewRouter()`.

### Usage

```go
// This is automatically added when you create a new router
router := router.NewRouter()

// But you can also add it manually if needed
router.AddBeforeMiddlewares([]router.Middleware{router.RecoveryMiddleware})
```

### Custom Recovery Handler

You can provide a custom recovery handler to handle panics in a specific way:

```go
customRecovery := router.NewRecoveryMiddleware(func(w http.ResponseWriter, r *http.Request, err interface{}) {
    log.Printf("Recovered from panic: %v", err)
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
})

router.AddBeforeMiddlewares([]router.Middleware{customRecovery})
```

## CORS Middleware

The router includes a CORS (Cross-Origin Resource Sharing) middleware to handle cross-origin requests.

### Basic Usage

```go
import "github.com/dracory/router/middlewares"

// Enable CORS for all routes
cors := middlewares.CORS()
router.AddBeforeMiddlewares([]router.Middleware{cors})
```

### Custom CORS Configuration

```go
cors := middlewares.CORS(
    middlewares.AllowOrigins([]string{"https://example.com"}),
    middlewares.AllowMethods([]string{"GET", "POST"}),
    middlewares.AllowHeaders([]string{"Content-Type"}),
    middlewares.AllowCredentials(true),
    middlewares.MaxAge(3600),
)
```

## Logging Middleware

### Basic Usage

```go
import "github.com/dracory/router/middlewares"

// Add request logging
logger := middlewares.Logger()
router.AddBeforeMiddlewares([]router.Middleware{logger})
```

### Custom Logger

```go
customLogger := middlewares.LoggerWithConfig(middlewares.LoggerConfig{
    Format: "${method} ${uri} - ${status} - ${latency}\n",
    Output: os.Stdout,
})
```

## Rate Limiting Middleware

### Basic Usage

```go
import "github.com/dracory/router/middlewares"

// Allow 100 requests per minute per IP
limiter := middlewares.RateLimit(100, time.Minute)
router.AddBeforeMiddlewares([]router.Middleware{limiter})
```

### Custom Rate Limiter

```go
import "golang.org/x/time/rate"

customLimiter := middlewares.RateLimitWithConfig(middlewares.RateLimitConfig{
    Limiter: rate.NewLimiter(rate.Every(time.Second), 10), // 10 requests per second
    KeyFunc: func(r *http.Request) string {
        return r.RemoteAddr // Rate limit by IP
    },
})
```

## Request ID Middleware

### Basic Usage

```go
import "github.com/dracory/router/middlewares"

// Add request ID to each request
requestID := middlewares.RequestID()
router.AddBeforeMiddlewares([]router.Middleware{requestID})
```

### Custom Request ID Generator

```go
customRequestID := middlewares.RequestIDWithConfig(middlewares.RequestIDConfig{
    Generator: func() string {
        return "req_" + someUUIDGenerator()
    },
    Header: "X-Request-ID",
})
```

## Secure Middleware

### Basic Usage

```go
import "github.com/dracory/router/middlewares"

// Add security headers
secure := middlewares.Secure()
router.AddBeforeMiddlewares([]router.Middleware{secure})
```

### Custom Security Headers

```go
customSecure := middlewares.SecureWithConfig(middlewares.SecureConfig{
    XSSProtection:         "1; mode=block",
    ContentTypeNosniff:    "nosniff",
    XFrameOptions:         "DENY",
    HSTSMaxAge:            31536000, // 1 year
    ContentSecurityPolicy: "default-src 'self'",
})
```

## Timeout Middleware

### Basic Usage

```go
import "github.com/dracory/router/middlewares"

// Add 10-second timeout to all requests
timeout := middlewares.Timeout(10 * time.Second)
router.AddBeforeMiddlewares([]router.Middleware{timeout})
```

### Custom Timeout Handler

```go
customTimeout := middlewares.TimeoutWithConfig(middlewares.TimeoutConfig{
    Timeout: 5 * time.Second,
    Handler: func(w http.ResponseWriter, r *http.Request) {
        http.Error(w, "Request timeout", http.StatusRequestTimeout)
    },
})
```
