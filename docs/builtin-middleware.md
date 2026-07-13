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

## Auth Middleware

The auth middleware reads a session key from a cookie, resolves the session and user from their respective stores (with optional TTL memory cache), and injects both into the request context for downstream handlers.

### Required Config

| Field | Type | Description |
|-------|------|-------------|
| `SessionStore` | `AuthSessionStore` | Session lookup by key |
| `UserStore` | `AuthUserStore` | User lookup by ID |
| `ContextKeyUser` | `any` | Context key for storing the authenticated user |
| `ContextKeySession` | `any` | Context key for storing the session object |
| `CookieName` | `string` | Name of the cookie containing the session key |

Optional fields: `Logger` (error logging), `MemoryCache` (TTL cache for session/user objects).

If any required field is missing, the middleware returns HTTP 500 on every request with a descriptive error message.

### Usage

```go
import (
    "github.com/dracory/rtr"
    "github.com/dracory/rtr/middlewares"
)

authMw := middlewares.AuthMiddleware(middlewares.AuthMiddlewareConfig{
    SessionStore:      mySessionStore,
    UserStore:         myUserStore,
    ContextKeyUser:    contextKeyUser,
    ContextKeySession: contextKeySession,
    CookieName:        "session_key",
    Logger:            myLogger,
    MemoryCache:       ttlcache.New[string, any](),
})

router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{authMw})
```

### Caching

When `MemoryCache` is provided, sessions and users are cached with a 5-minute TTL. On a cache miss, the middleware fetches from the store and caches the result. Expired sessions are never cached — the expiry check runs before the cache write.

### Interfaces

```go
type AuthSessionStore interface {
    SessionFindByKey(ctx context.Context, key string) (AuthSession, error)
}

type AuthSession interface {
    GetUserID() string
    IsExpired() bool
}

type AuthUserStore interface {
    UserFindByID(ctx context.Context, id string) (AuthUser, error)
}

type AuthUser interface {
    IsActive() bool
    IsAdministrator() bool
    IsSuperuser() bool
    IsRegistrationCompleted() bool
}
```

## User Middleware

The user middleware enforces authentication, active status, registration completion, and optional role checks. It expects the user to already be in the request context (typically set by the auth middleware).

### Required Config

| Field | Type | Description |
|-------|------|-------------|
| `GetUser` | `func(r *http.Request) UserMiddlewareUser` | Extracts the user from the request context |

Optional fields: `RegistrationEnabled`, `RegistrationPaths`, `RequireRoles`, and callback functions (`OnNotAuthenticated`, `OnNotActive`, `OnRegistrationIncomplete`, `OnNotAuthorized`).

If `GetUser` is nil, the middleware returns HTTP 500 on every request with a descriptive error message.

### Usage

```go
import (
    "github.com/dracory/rtr"
    "github.com/dracory/rtr/middlewares"
)

userMw := middlewares.UserMiddleware(middlewares.UserMiddlewareConfig{
    GetUser: func(r *http.Request) middlewares.UserMiddlewareUser {
        return middlewares.GetUserFromContext(r.Context())
    },
    RegistrationEnabled: true,
    RegistrationPaths:   []string{"/profile", "/register"},
    OnNotAuthenticated:  redirectToLogin,
    OnNotActive:         redirectToHome,
    OnRegistrationIncomplete: redirectToRegister,
})

router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{userMw})
```

### Role-Based Access Control

Set `RequireRoles` to enforce that the user has at least one of the specified roles. The user type must implement `UserWithRole` (which extends `UserMiddlewareUser` with `HasRole(role string) bool`).

```go
userMw := middlewares.UserMiddleware(middlewares.UserMiddlewareConfig{
    GetUser:      getUserFromContext,
    RequireRoles: []string{"admin", "superuser"},
    OnNotAuthorized: func(w http.ResponseWriter, r *http.Request) {
        http.Error(w, "Forbidden", http.StatusForbidden)
    },
})
```

### Interfaces

```go
type UserMiddlewareUser interface {
    IsActive() bool
    IsRegistrationCompleted() bool
}

type UserWithRole interface {
    UserMiddlewareUser
    HasRole(role string) bool
}
```

### Callbacks

Each check has an optional callback. When a callback is nil, a default HTTP error response is used:

| Check | Callback | Default Response |
|-------|----------|-----------------|
| Not authenticated | `OnNotAuthenticated` | 401 Unauthorized |
| Not active | `OnNotActive` | 403 Forbidden |
| Registration incomplete | `OnRegistrationIncomplete` | 403 Forbidden |
| Missing required role | `OnNotAuthorized` | 403 Forbidden |

## Domain Redirect Middleware

### WWW → Naked Domain

Redirects requests from `www.` subdomains to the naked domain using the incoming request scheme (defaulting to `https` when missing).

```go
import (
    "github.com/dracory/rtr"
    "github.com/dracory/rtr/middlewares"
)

router := rtr.NewRouter()
router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    middlewares.WwwToNakedDomainMiddleware(),
})
```

### Naked → WWW Domain

Redirects naked domains to `www.` subdomains. You can exclude specific hosts (e.g., `localhost`) via `hostExcludes`.

```go
import (
    "github.com/dracory/rtr"
    "github.com/dracory/rtr/middlewares"
)

router := rtr.NewRouter()
router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    middlewares.NakedDomainToWwwMiddleware([]string{"localhost", "127.0.0.1"}),
})
```
