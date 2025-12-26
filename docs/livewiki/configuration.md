---
path: configuration.md
page-type: reference
summary: Complete guide to configuring RTR router including declarative and imperative approaches.
tags: [configuration, declarative, imperative, setup, options]
created: 2025-12-26
updated: 2025-12-26
version: 1.0.0
---

# RTR Router Configuration

This document covers all aspects of configuring the RTR router, including both declarative and imperative configuration approaches.

## Configuration Approaches

RTR supports two main configuration approaches:

### 1. Imperative Configuration
Programmatic configuration using method calls and fluent APIs.

### 2. Declarative Configuration
Configuration-as-code using data structures that can be serialized.

## Imperative Configuration

### Basic Router Setup

```go
package main

import (
    "net/http"
    "github.com/dracory/rtr"
)

func main() {
    // Create a new router
    router := rtr.NewRouter()
    
    // Configure routes, groups, and middleware
    setupRoutes(router)
    setupMiddleware(router)
    
    // Start server
    http.ListenAndServe(":8080", router)
}

func setupRoutes(router rtr.RouterInterface) {
    // Add individual routes
    router.AddRoute(rtr.Get("/", homeHandler))
    router.AddRoute(rtr.Get("/health", healthHandler))
    
    // Add multiple routes at once
    router.AddRoutes([]rtr.RouteInterface{
        rtr.Get("/users", usersHandler),
        rtr.Post("/users", createUserHandler),
        rtr.Get("/products", productsHandler),
    })
}

func setupMiddleware(router rtr.RouterInterface) {
    // Add global middleware
    router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
        rtr.NewAnonymousMiddleware(loggingMiddleware),
        rtr.NewAnonymousMiddleware(recoveryMiddleware),
    })
}
```

### Route Configuration

#### Individual Routes

```go
// Using shortcut methods
router.AddRoute(rtr.Get("/api/users", handleUsers))

// Using method chaining
router.AddRoute(rtr.NewRoute().
    SetMethod("GET").
    SetPath("/api/users").
    SetHandler(handleUsers).
    SetName("List Users").
    AddBeforeMiddlewares([]rtr.MiddlewareInterface{
        rtr.NewAnonymousMiddleware(authMiddleware),
    }))
```

#### Specialized Handler Routes

```go
// JSON handler
router.AddRoute(rtr.GetJSON("/api/status", func(w http.ResponseWriter, r *http.Request) string {
    return `{"status": "ok", "version": "1.0.0"}`
}))

// HTML handler
router.AddRoute(rtr.GetHTML("/page", func(w http.ResponseWriter, r *http.Request) string {
    return "<h1>Hello, World!</h1>"
}))

// CSS handler
router.AddRoute(rtr.GetCSS("/style.css", func(w http.ResponseWriter, r *http.Request) string {
    return "body { font-family: Arial; }"
}))
```

### Group Configuration

#### Basic Groups

```go
// Create a group
apiGroup := rtr.NewGroup().SetPrefix("/api/v1")

// Add routes to group
apiGroup.AddRoute(rtr.Get("/users", usersHandler))
apiGroup.AddRoute(rtr.Post("/users", createUserHandler))

// Add group to router
router.AddGroup(apiGroup)
```

#### Groups with Middleware

```go
// Group with middleware
apiGroup := rtr.NewGroup().
    SetPrefix("/api/v1").
    SetName("API v1").
    AddBeforeMiddlewares([]rtr.MiddlewareInterface{
        rtr.NewAnonymousMiddleware(apiKeyMiddleware),
        rtr.NewAnonymousMiddleware(rateLimitMiddleware),
    })

apiGroup.AddRoute(rtr.Get("/users", usersHandler))
router.AddGroup(apiGroup)
```

#### Nested Groups

```go
// Parent group
apiGroup := rtr.NewGroup().SetPrefix("/api")

// Child group
v1Group := rtr.NewGroup().
    SetPrefix("/v1").
    AddBeforeMiddlewares([]rtr.MiddlewareInterface{
        rtr.NewAnonymousMiddleware(versionMiddleware),
    })

// Add nested structure
apiGroup.AddGroup(v1Group)
router.AddGroup(apiGroup)

// Routes will be accessible at /api/v1/...
```

### Domain Configuration

#### Basic Domains

```go
// Create domain for API
apiDomain := rtr.NewDomain("api.example.com")
apiDomain.AddRoute(rtr.Get("/users", apiUsersHandler))

// Create domain for web
webDomain := rtr.NewDomain("www.example.com")
webDomain.AddRoute(rtr.Get("/", webHomeHandler))

// Add domains to router
router.AddDomain(apiDomain)
router.AddDomain(webDomain)
```

#### Wildcard Domains

```go
// Wildcard subdomain
tenantDomain := rtr.NewDomain("*.example.com")
tenantDomain.AddRoute(rtr.Get("/data", tenantDataHandler))

// Multiple patterns
multiDomain := rtr.NewDomain("example.com", "api.example.com", "*.example.org")
multiDomain.AddRoute(rtr.Get("/health", healthHandler))
```

#### Domain with Middleware

```go
adminDomain := rtr.NewDomain("admin.example.com").
    SetName("Admin Domain").
    AddBeforeMiddlewares([]rtr.MiddlewareInterface{
        rtr.NewAnonymousMiddleware(adminAuthMiddleware),
        rtr.NewAnonymousMiddleware(loggingMiddleware),
    })

adminDomain.AddRoute(rtr.Get("/dashboard", adminDashboardHandler))
router.AddDomain(adminDomain)
```

### Middleware Configuration

#### Global Middleware

```go
// Before middleware (executed before handler)
router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    rtr.NewMiddleware("Recovery", recoveryMiddleware),
    rtr.NewMiddleware("Logger", loggingMiddleware),
    rtr.NewMiddleware("CORS", corsMiddleware),
})

// After middleware (executed after handler)
router.AddAfterMiddlewares([]rtr.MiddlewareInterface{
    rtr.NewMiddleware("ResponseTime", responseTimeMiddleware),
})
```

#### Named Middleware

```go
// Create named middleware for reuse
authMiddleware := rtr.NewMiddleware("Auth", authenticationMiddleware)
loggingMiddleware := rtr.NewMiddleware("Logger", loggingMiddleware)

// Use in multiple places
router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{authMiddleware})
apiGroup.AddBeforeMiddlewares([]rtr.MiddlewareInterface{loggingMiddleware})
```

## Declarative Configuration

### Basic Declarative Setup

```go
config := rtr.RouterConfig{
    Name: "My Application",
    BeforeMiddleware: []rtr.MiddlewareConfig{
        rtr.NewMiddlewareConfig("Recovery", middlewares.RecoveryMiddleware),
        rtr.NewMiddlewareConfig("Logger", middlewares.LoggerMiddleware),
    },
    AfterMiddleware: []rtr.MiddlewareConfig{
        rtr.NewMiddlewareConfig("ResponseTime", responseTimeMiddleware),
    },
    Routes: []rtr.RouteConfig{
        rtr.GET("/", homeHandler).WithName("Home"),
        rtr.GET("/health", healthHandler).WithName("Health Check"),
        rtr.GET_JSON("/api/status", statusHandler).WithName("API Status"),
    },
}

router := rtr.NewRouterFromConfig(config)
```

### Advanced Declarative Configuration

```go
config := rtr.RouterConfig{
    Name: "E-commerce API",
    Metadata: map[string]interface{}{
        "version": "2.0.0",
        "environment": "production",
    },
    
    // Global middleware
    BeforeMiddleware: []rtr.MiddlewareConfig{
        rtr.NewMiddlewareConfig("Recovery", middlewares.RecoveryMiddleware),
        rtr.NewMiddlewareConfig("CORS", middlewares.CorsMiddleware),
        rtr.NewMiddlewareConfig("RateLimit", middlewares.RateLimitByIPMiddleware),
    },
    
    // Direct routes
    Routes: []rtr.RouteConfig{
        rtr.GET("/", homeHandler).
            WithName("Home").
            WithMetadata("public", true),
            
        rtr.GET("/health", healthHandler).
            WithName("Health Check").
            WithMetadata("public", true),
            
        rtr.POST("/auth/login", loginHandler).
            WithName("Login").
            WithBeforeMiddleware(rtr.NewAnonymousMiddleware(rateLimitLoginMiddleware)),
    },
    
    // Route groups
    Groups: []rtr.GroupConfig{
        rtr.Group("/api/v1",
            rtr.GET("/users", usersHandler).
                WithName("List Users").
                WithBeforeMiddleware(rtr.NewAnonymousMiddleware(authMiddleware)),
                
            rtr.POST("/users", createUserHandler).
                WithName("Create User").
                WithBeforeMiddleware(rtr.NewAnonymousMiddleware(authMiddleware)),
                
            rtr.GET("/products", productsHandler).
                WithName("List Products").
                WithMetadata("cache", "5m"),
        ).
        WithName("API v1").
        WithBeforeMiddleware(
            rtr.NewAnonymousMiddleware(apiKeyMiddleware),
            rtr.NewAnonymousMiddleware(requestIDMiddleware),
        ).
        WithAfterMiddleware(
            rtr.NewAnonymousMiddleware(responseLoggerMiddleware),
        ),
        
        rtr.Group("/admin",
            rtr.GET("/dashboard", adminDashboardHandler).
                WithName("Admin Dashboard"),
                
            rtr.GET("/users", adminUsersHandler).
                WithName("Admin Users"),
        ).
        WithName("Admin Panel").
        WithBeforeMiddleware(
            rtr.NewAnonymousMiddleware(adminAuthMiddleware),
            rtr.NewAnonymousMiddleware(adminLoggingMiddleware),
        ),
    },
    
    // Domain-specific routing
    Domains: []rtr.DomainConfig{
        rtr.DomainConfig{
            Patterns: []string{"api.example.com"},
            Name: "API Domain",
            BeforeMiddleware: []rtr.MiddlewareConfig{
                rtr.NewMiddlewareConfig("APIKey", apiKeyMiddleware),
            },
            Routes: []rtr.RouteConfig{
                rtr.GET("/v1/status", apiStatusHandler).WithName("API Status"),
            },
        },
        rtr.DomainConfig{
            Patterns: []string{"www.example.com", "example.com"},
            Name: "Web Domain",
            Routes: []rtr.RouteConfig{
                rtr.GET_HTML("/", webHomeHandler).WithName("Web Home"),
                rtr.GET_HTML("/about", webAboutHandler).WithName("About"),
            },
        },
    },
}

router := rtr.NewRouterFromConfig(config)
```

### Configuration from JSON/YAML

#### JSON Configuration

```json
{
  "name": "My Application",
  "before_middleware": [
    {
      "name": "Recovery",
      "handler": "RecoveryMiddleware"
    },
    {
      "name": "Logger", 
      "handler": "LoggerMiddleware"
    }
  ],
  "routes": [
    {
      "method": "GET",
      "path": "/",
      "name": "Home",
      "handler": "homeHandler"
    }
  ],
  "groups": [
    {
      "prefix": "/api/v1",
      "name": "API v1",
      "routes": [
        {
          "method": "GET",
          "path": "/users",
          "name": "List Users",
          "handler": "usersHandler"
        }
      ]
    }
  ]
}
```

#### Loading Configuration

```go
import "encoding/json"

func loadConfigFromFile(filename string) (rtr.RouterConfig, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return rtr.RouterConfig{}, err
    }
    
    var config rtr.RouterConfig
    err = json.Unmarshal(data, &config)
    return config, err
}

func main() {
    config, err := loadConfigFromFile("router-config.json")
    if err != nil {
        log.Fatal(err)
    }
    
    router := rtr.NewRouterFromConfig(config)
    http.ListenAndServe(":8080", router)
}
```

## Configuration Patterns

### Environment-Based Configuration

```go
func getRouterConfig() rtr.RouterConfig {
    env := os.Getenv("APP_ENV")
    
    baseConfig := rtr.RouterConfig{
        Name: "My Application",
        BeforeMiddleware: []rtr.MiddlewareConfig{
            rtr.NewMiddlewareConfig("Recovery", middlewares.RecoveryMiddleware),
        },
    }
    
    switch env {
    case "development":
        baseConfig.BeforeMiddleware = append(baseConfig.BeforeMiddleware,
            rtr.NewMiddlewareConfig("DevLogger", devLoggerMiddleware),
        )
        baseConfig.Routes = append(baseConfig.Routes,
            rtr.GET("/debug", debugHandler).WithName("Debug"),
        )
        
    case "production":
        baseConfig.BeforeMiddleware = append(baseConfig.BeforeMiddleware,
            rtr.NewMiddlewareConfig("ProdLogger", prodLoggerMiddleware),
            rtr.NewMiddlewareConfig("Metrics", metricsMiddleware),
        )
        
    default:
        baseConfig.BeforeMiddleware = append(baseConfig.BeforeMiddleware,
            rtr.NewMiddlewareConfig("BasicLogger", basicLoggerMiddleware),
        )
    }
    
    return baseConfig
}
```

### Modular Configuration

```go
// Base configuration
func baseConfig() rtr.RouterConfig {
    return rtr.RouterConfig{
        BeforeMiddleware: []rtr.MiddlewareConfig{
            rtr.NewMiddlewareConfig("Recovery", middlewares.RecoveryMiddleware),
        },
        Routes: []rtr.RouteConfig{
            rtr.GET("/health", healthHandler).WithName("Health"),
        },
    }
}

// API configuration
func apiConfig() rtr.GroupConfig {
    return rtr.Group("/api/v1",
        rtr.GET("/users", usersHandler).WithName("List Users"),
        rtr.POST("/users", createUserHandler).WithName("Create User"),
    ).
    WithName("API").
    WithBeforeMiddleware(
        rtr.NewAnonymousMiddleware(authMiddleware),
    )
}

// Admin configuration
func adminConfig() rtr.GroupConfig {
    return rtr.Group("/admin",
        rtr.GET("/dashboard", adminDashboardHandler).WithName("Dashboard"),
    ).
    WithName("Admin").
    WithBeforeMiddleware(
        rtr.NewAnonymousMiddleware(adminAuthMiddleware),
    )
}

func main() {
    config := baseConfig()
    config.Groups = append(config.Groups, apiConfig(), adminConfig())
    
    router := rtr.NewRouterFromConfig(config)
    http.ListenAndServe(":8080", router)
}
```

## Configuration Validation

### Custom Validation

```go
func validateConfig(config rtr.RouterConfig) error {
    // Check for duplicate route names
    routeNames := make(map[string]bool)
    for _, route := range config.Routes {
        if route.Name != "" {
            if routeNames[route.Name] {
                return fmt.Errorf("duplicate route name: %s", route.Name)
            }
            routeNames[route.Name] = true
        }
    }
    
    // Validate group prefixes
    for _, group := range config.Groups {
        if !strings.HasPrefix(group.Prefix, "/") {
            return fmt.Errorf("group prefix must start with '/': %s", group.Prefix)
        }
    }
    
    return nil
}

func main() {
    config := getRouterConfig()
    
    if err := validateConfig(config); err != nil {
        log.Fatal("Invalid configuration:", err)
    }
    
    router := rtr.NewRouterFromConfig(config)
    http.ListenAndServe(":8080", router)
}
```

## Configuration Best Practices

### 1. Use Named Middleware

```go
// Good: Named middleware for debugging
rtr.NewMiddlewareConfig("Auth", authMiddleware)

// Avoid: Anonymous middleware in config
rtr.NewAnonymousMiddleware(authMiddleware)
```

### 2. Organize by Feature

```go
// Group related routes
apiGroup := rtr.Group("/api/v1",
    rtr.GET("/users", usersHandler),
    rtr.POST("/users", createUserHandler),
    rtr.GET("/products", productsHandler),
).WithName("API")
```

### 3. Use Metadata for Documentation

```go
rtr.GET("/users", usersHandler).
    WithName("List Users").
    WithMetadata("description", "Returns a list of all users").
    WithMetadata("deprecated", false).
    WithMetadata("version", "1.0")
```

### 4. Environment-Specific Configuration

```go
if os.Getenv("DEBUG") == "true" {
    config.Routes = append(config.Routes,
        rtr.GET("/debug/routes", debugRoutesHandler).WithName("Debug Routes"),
    )
}
```

### 5. Configuration Security

```go
// Don't store sensitive data in configuration
// Use environment variables instead

apiKey := os.Getenv("API_KEY")
if apiKey == "" {
    log.Fatal("API_KEY environment variable required")
}

authMiddleware := createAuthMiddleware(apiKey)
```

## Configuration Migration

### From Imperative to Declarative

```go
// Imperative approach
router := rtr.NewRouter()
router.AddRoute(rtr.Get("/users", usersHandler))
router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    rtr.NewAnonymousMiddleware(loggingMiddleware),
})

// Convert to declarative
config := rtr.RouterConfig{
    Routes: []rtr.RouteConfig{
        rtr.GET("/users", usersHandler).WithName("Users"),
    },
    BeforeMiddleware: []rtr.MiddlewareConfig{
        rtr.NewMiddlewareConfig("Logger", loggingMiddleware),
    },
}

router := rtr.NewRouterFromConfig(config)
```

### Hybrid Approach

```go
// Start with declarative base
config := rtr.RouterConfig{
    Routes: []rtr.RouteConfig{
        rtr.GET("/", homeHandler).WithName("Home"),
    },
}

router := rtr.NewRouterFromConfig(config)

// Add dynamic routes imperatively
if featureEnabled {
    router.AddRoute(rtr.Get("/beta", betaHandler))
}
```

## See Also

- [Getting Started Guide](getting_started.md) - Learn basic configuration
- [Architecture Documentation](architecture.md) - Understand configuration system design
- [API Reference](api_reference.md) - Complete configuration API documentation
- [Development Guide](development.md) - Configuration for development and testing
