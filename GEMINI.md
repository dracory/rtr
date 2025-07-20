# GEMINI.md for `dracory/rtr` Project

This file provides comprehensive documentation and guidelines for the `dracory/rtr` Go HTTP router project.

## Project Overview

`dracory/rtr` is a high-performance, feature-rich HTTP router for Go applications. It provides a flexible and intuitive API for building web applications with robust routing capabilities.

## Core Features

### 1. Route Management
- Support for all standard HTTP methods (GET, POST, PUT, DELETE, etc.)
- Path parameter extraction
- Flexible route matching
- Handler prioritization system

### 2. Middleware System
- **MiddlewareInterface**: Named middleware with metadata support
- **StdMiddleware**: Standard Go middleware pattern
- Middleware chaining and composition
- Built-in recovery middleware
- Support for both global and route-specific middleware

### 3. Route Organization
- Route grouping with shared prefixes
- Nested groups for hierarchical routing
- Domain-based routing
- Declarative configuration

### 4. Handler Types
- **StdHandler**: Standard `http.HandlerFunc`
- **StringHandler**: Returns string responses
- **HTMLHandler**: HTML responses with proper headers
- **JSONHandler**: JSON responses with proper headers
- **CSSHandler**: CSS content delivery
- **XMLHandler**: XML responses
- **TextHandler**: Plain text responses
- **ErrorHandler**: Standardized error handling

## Architecture

### Core Interfaces

1. **RouterInterface**
   - Main router implementation
   - Manages routes, groups, and domains
   - Implements `http.Handler`

2. **RouteInterface**
   - Represents individual routes
   - Handles HTTP method, path, and handler associations

3. **GroupInterface**
   - Groups related routes
   - Supports shared middleware and prefixes
   - Enables nested grouping

4. **DomainInterface**
   - Handles domain-based routing
   - Supports pattern matching for hostnames

5. **MiddlewareInterface**
   - Defines named middleware
   - Supports metadata and configuration
   - Provides chainable API

## Middleware System

### Middleware Execution Order

The middleware execution in `dracory/rtr` follows a specific sequence to ensure predictable request processing. The complete order is:

1. **Global Before Middleware**
   - Added via `router.AddBeforeMiddlewares()`
   - Executes first, in the order they were added

2. **Domain Before Middleware**
   - Added via `domain.AddBeforeMiddlewares()`
   - Only executes if the request matches the domain's host pattern
   - Multiple matching domains will execute in the order they were registered

3. **Group Before Middleware**
   - Added via `group.AddBeforeMiddlewares()`
   - Executes from outermost to innermost group
   - For nested groups, parent group middleware runs before child group middleware

4. **Route Before Middleware**
   - Added via `route.AddBeforeMiddlewares()`
   - Executes in the order they were added to the specific route

5. **Route Handler**
   - The actual route handler function processes the request

6. **Route After Middleware**
   - Added via `route.AddAfterMiddlewares()`
   - Executes in reverse order (last added, first executed)

7. **Group After Middleware**
   - Added via `group.AddAfterMiddlewares()`
   - Executes from innermost to outermost group
   - For nested groups, child group middleware runs before parent group middleware

8. **Domain After Middleware**
   - Added via `domain.AddAfterMiddlewares()`
   - Only executes if the request matches the domain's host pattern
   - Multiple matching domains will execute in reverse order of registration

9. **Global After Middleware**
   - Added via `router.AddAfterMiddlewares()`
   - Executes last, in reverse order they were added

### Visual Representation

```
Incoming Request
        ↓
┌───────────────────────┐
│  Global Before (1..N) │
└──────────┬────────────┘
           ↓
┌───────────────────────┐
│  Domain Before (1..N) │
└──────────┬────────────┘
           ↓
┌───────────────────────┐
│  Group Before (Outer) │
└──────────┬────────────┘
           ↓
┌───────────────────────┐
│  Group Before (Inner) │
└──────────┬────────────┘
           ↓
┌───────────────────────┐
│  Route Before (1..N)  │
└──────────┬────────────┘
           ↓
┌───────────────────────┐
│     Route Handler     │
└──────────┬────────────┘
           ↓
┌───────────────────────┐
│  Route After (N..1)   │
└──────────┬────────────┘
           ↓
┌───────────────────────┐
│  Group After (Inner)  │
└──────────┬────────────┘
           ↓
┌───────────────────────┐
│  Group After (Outer)  │
└──────────┬────────────┘
           ↓
┌───────────────────────┐
│  Domain After (N..1)  │
└──────────┬────────────┘
           ↓
┌───────────────────────┐
│  Global After (N..1)  │
└──────────┬────────────┘
           ↓
     Response Sent
```

### Important Notes

- Middleware at the same level executes in the order it was added
- The execution order is always from global → domain → group → route → handler → route → group → domain → global
- After middleware executes in reverse order (last added, first executed) within their respective scopes
- If any middleware in the "before" chain writes a response and doesn't call `next.ServeHTTP()`, the remaining middleware in that chain and the route handler will be skipped

### Middleware Types
1. **Before Middleware**: Executed before the route handler
2. **After Middleware**: Executed after the route handler
3. **Recovery Middleware**: Handles panics in handlers

## Usage Examples

### Basic Route Definition
```go
router := rtr.NewRouter()
router.AddRoute(rtr.NewRoute().
    SetMethod("GET").
    SetPath("/hello").
    SetStringHandler(func(w http.ResponseWriter, r *http.Request) string {
        return "Hello, World!"
    }))
```

### Using Middleware
```go
auth := rtr.NewMiddleware("auth", func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Authentication logic here
        next.ServeHTTP(w, r)
    })
})

router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{auth})
```

### Route Groups
```go
api := rtr.NewGroup().
    SetPrefix("/api").
    AddBeforeMiddlewares([]rtr.MiddlewareInterface{auth})

api.AddRoute(rtr.NewRoute().
    SetMethod("GET").
    SetPath("/users").
    SetJSONHandler(func(w http.ResponseWriter, r *http.Request) string {
        return `{"status": "success"}`
    }))

router.AddGroup(api)
```

### Domain-based Routing
```go
domain := rtr.NewDomain("api.example.com").
    AddBeforeMiddlewares([]rtr.MiddlewareInterface{rateLimiter})

domain.AddRoute(rtr.NewRoute().
    SetMethod("GET").
    SetPath("/status").
    SetJSONHandler(func(w http.ResponseWriter, r *http.Request) string {
        return `{"status": "ok"}`
    }))

router.AddDomain(domain)
```

## Development Guidelines

### Code Style
- Follow standard Go formatting (`go fmt`)
- Use `goimports` for import organization
- Write clear, descriptive variable and function names
- Document all exported functions and types
- Keep functions small and focused

### Error Handling
- Always check and handle errors
- Provide meaningful error messages
- Use custom error types when appropriate

### Testing
- Write tests for all public APIs
- Use table-driven tests where applicable
- Test edge cases and error conditions
- Maintain good test coverage

## Best Practices

1. **Middleware**
   - Keep middleware focused on a single responsibility
   - Use named middleware for better debugging
   - Document middleware dependencies

2. **Routing**
   - Group related routes together
   - Use meaningful route prefixes
   - Document route requirements and behaviors

3. **Performance**
   - Minimize middleware overhead
   - Use sync.Pool for frequently allocated objects
   - Avoid expensive operations in hot paths

4. **Security**
   - Validate all inputs
   - Use context for request-scoped values
   - Implement proper CORS and CSRF protection

## Project Structure

- `interfaces.go`: Core interface definitions
- `router_implementation.go`: Main router implementation
- `route.go`: Route handling logic
- `group_implementation.go`: Route group implementation
- `middleware_implementation.go`: Middleware system
- `domain.go`: Domain-based routing
- `declarative.go`: Declarative configuration
- `list.go`: Route listing and debugging

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run tests and linters
6. Submit a pull request

## License

MIT License - See LICENSE file for details
