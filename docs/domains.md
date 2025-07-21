# Domain-Based Routing

## Overview

Domain-based routing allows you to handle requests differently based on the `Host` header. This is useful for:

- Multi-tenant applications
- API versioning via subdomains
- Different behavior for different domains/subdomains
- Custom domain support in SaaS applications

## Basic Usage

### Creating a Domain

```go
router := rtr.NewRouter()

// Create a domain for api.example.com
apiDomain := rtr.NewDomain("api.example.com")

// Add routes to the domain
apiDomain.AddRoute(rtr.NewRoute().
    SetMethod("GET").
    SetPath("/users").
    SetHandler(apiUsersHandler))

// Add the domain to the router
router.AddDomain(apiDomain)
```

### Wildcard Domains

Use `*` as a wildcard for subdomains:

```go
// Matches any subdomain of example.com
wildcardDomain := rtr.NewDomain("*.example.com")
```

### Multiple Domain Patterns

A domain can match multiple patterns:

```go
// Matches both www.example.com and example.com
mainDomain := rtr.NewDomain("example.com", "www.example.com")
```

## Domain-Specific Middleware

Add middleware that only runs for specific domains:

```go
// Create a domain for the API
apiDomain := rtr.NewDomain("api.example.com")

// Add API-specific middleware
apiDomain.AddBeforeMiddlewares([]rtr.Middleware{
    apiAuthMiddleware,
    rateLimiterMiddleware,
})

// Add the domain to the router
router.AddDomain(apiDomain)
```

## Example: Multi-Tenant Application

```go
router := rtr.NewRouter()

// Main domain
mainDomain := rtr.NewDomain("example.com", "www.example.com")
mainDomain.AddRoute(rtr.Get("/", homeHandler))
router.AddDomain(mainDomain)

// API subdomain
apiDomain := rtr.NewDomain("api.example.com")
apiDomain.AddBeforeMiddlewares([]rtr.Middleware{apiAuthMiddleware})
apiDomain.AddRoute(rtr.Get("/users", apiUsersHandler))
router.AddDomain(apiDomain)

// Admin subdomain
adminDomain := rtr.NewDomain("admin.example.com")
adminDomain.AddBeforeMiddlewares([]rtr.Middleware{authMiddleware, adminOnlyMiddleware})
adminDomain.AddRoute(rtr.Get("/dashboard", adminDashboardHandler))
router.AddDomain(adminDomain)

// Wildcard for tenant subdomains
tenantDomain := rtr.NewDomain("*.example.com")
tenantDomain.AddBeforeMiddlewares([]rtr.Middleware{tenantResolverMiddleware})
tenantDomain.AddRoute(rtr.Get("/", tenantHomeHandler))
router.AddDomain(tenantDomain)
```

## Best Practices

1. **Order Matters**: More specific domains should be added before more general ones.
2. **Wildcards**: Use wildcards carefully as they match any subdomain.
3. **Middleware**: Use domain-specific middleware for cross-cutting concerns.
4. **Performance**: Be mindful of the number of domains and their middleware.
5. **Testing**: Test domain matching thoroughly, especially with wildcards.

## Common Patterns

### API Versioning

```go
// v1 API
v1Domain := rtr.NewDomain("api-v1.example.com")
v1Domain.AddRoute(rtr.Get("/users", v1UsersHandler))

// v2 API
v2Domain := rtr.NewDomain("api-v2.example.com")
v2Domain.AddRoute(rtr.Get("/users", v2UsersHandler))
```

### Multi-Tenant Applications

```go
// Main application
mainDomain := rtr.NewDomain("example.com", "www.example.com")
mainDomain.AddRoute(rtr.Get("/", homeHandler))

// Tenant-specific subdomains
tenantDomain := rtr.NewDomain("*.example.com")
tenantDomain.AddBeforeMiddlewares([]rtr.Middleware{tenantResolverMiddleware})
tenantDomain.AddRoute(rtr.Get("/", tenantHomeHandler))
```

For more advanced domain routing patterns, see [Advanced Domain Routing](./advanced-domain-routing.md).
