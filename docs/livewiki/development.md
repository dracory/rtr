---
path: development.md
page-type: tutorial
summary: Development workflow, testing, and contributing to the RTR router project.
tags: [development, testing, contributing, workflow, setup]
created: 2025-12-26
updated: 2025-12-26
version: 1.0.0
---

# RTR Router Development Guide

This guide covers development workflow, testing practices, and contributing to the RTR router project.

## Development Setup

### Prerequisites

- Go 1.25 or higher
- Git
- Make (optional, for build scripts)

### Repository Setup

```bash
# Clone the repository
git clone https://github.com/dracory/rtr.git
cd rtr

# Install dependencies
go mod download

# Verify installation
go test ./...
```

### Development Environment

#### IDE Configuration

**VS Code Settings (`.vscode/settings.json`):**
```json
{
    "go.toolsManagement.autoUpdate": true,
    "go.useLanguageServer": true,
    "go.lintTool": "golangci-lint",
    "go.lintFlags": ["--fast"],
    "go.testFlags": ["-v"],
    "go.coverOnSave": true,
    "go.coverageDecorator": {
        "type": "gutter",
        "coveredHighlightColor": "rgba(64,128,64,0.5)",
        "uncoveredHighlightColor": "rgba(128,64,64,0.25)"
    }
}
```

**Recommended Extensions:**
- Go (golang.go)
- GitLens
- Better Comments
- Error Lens

#### Makefile

```makefile
.PHONY: test test-verbose test-race benchmark lint build clean

# Run all tests
test:
	go test ./...

# Run tests with verbose output
test-verbose:
	go test -v ./...

# Run tests with race detection
test-race:
	go test -race ./...

# Run benchmarks
benchmark:
	go test -bench=. -benchmem ./...

# Run linter
lint:
	golangci-lint run

# Build the project
build:
	go build ./...

# Clean build artifacts
clean:
	go clean -cache
	go clean -testcache

# Run integration tests
integration:
	go test -tags=integration ./...

# Generate coverage report
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run example tests
examples:
	go test ./examples/...
```

## Testing Strategy

### Test Organization

```
router_test.go              # Core router functionality
route_test.go               # Route-specific tests
group_test.go               # Group functionality tests
domain_test.go              # Domain routing tests
middleware_test.go          # Middleware system tests
handlers_test.go            # Handler type tests
params_test.go              # Parameter extraction tests
integration_test.go         # End-to-end tests
examples/                   # Example applications with tests
```

### Test Categories

#### Unit Tests

Test individual components in isolation:

```go
func TestRouteCreation(t *testing.T) {
    route := rtr.NewRoute().
        SetMethod("GET").
        SetPath("/users").
        SetHandler(func(w http.ResponseWriter, r *http.Request) {})
    
    assert.Equal(t, "GET", route.GetMethod())
    assert.Equal(t, "/users", route.GetPath())
    assert.NotNil(t, route.GetHandler())
}
```

#### Integration Tests

Test component interactions:

```go
func TestRouterIntegration(t *testing.T) {
    router := rtr.NewRouter()
    router.AddRoute(rtr.Get("/test", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    }))
    
    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()
    
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
}
```

#### Benchmark Tests

Performance testing:

```go
func BenchmarkRouteMatching(b *testing.B) {
    router := setupBenchmarkRouter()
    req := httptest.NewRequest("GET", "/users/123", nil)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)
    }
}
```

### Test Utilities

#### Helper Functions

```go
// test_helpers.go
package rtr_test

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/stretchr/testify/assert"
)

func createTestRequest(method, path string) *http.Request {
    return httptest.NewRequest(method, path, nil)
}

func executeRequest(router rtr.RouterInterface, method, path string) *httptest.ResponseRecorder {
    req := createTestRequest(method, path)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    return w
}

func assertResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedBody string) {
    assert.Equal(t, expectedStatus, w.Code)
    assert.Equal(t, expectedBody, w.Body.String())
}
```

#### Test Fixtures

```go
// fixtures.go
package rtr_test

func setupTestRouter() rtr.RouterInterface {
    router := rtr.NewRouter()
    
    // Add test routes
    router.AddRoute(rtr.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("home"))
    }))
    
    router.AddRoute(rtr.GetJSON("/api/status", func(w http.ResponseWriter, r *http.Request) string {
        return `{"status": "ok"}`
    }))
    
    // Add test group
    apiGroup := rtr.NewGroup().SetPrefix("/api/v1")
    apiGroup.AddRoute(rtr.Get("/users", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("users"))
    }))
    router.AddGroup(apiGroup)
    
    return router
}

func setupMiddlewareRouter() rtr.RouterInterface {
    router := rtr.NewRouter()
    
    // Add test middleware
    router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
        rtr.NewAnonymousMiddleware(func(next http.Handler) http.Handler {
            return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.Header().Set("X-Test-Middleware", "before")
                next.ServeHTTP(w, r)
            })
        }),
    })
    
    router.AddRoute(rtr.Get("/test", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("test"))
    }))
    
    return router
}
```

### Testing Patterns

#### Route Testing

```go
func TestRouteHandlers(t *testing.T) {
    tests := []struct {
        name           string
        method         string
        path           string
        expectedStatus int
        expectedBody   string
        setupRouter    func() rtr.RouterInterface
    }{
        {
            name:           "standard handler",
            method:         "GET",
            path:           "/",
            expectedStatus: http.StatusOK,
            expectedBody:   "home",
            setupRouter:    setupTestRouter,
        },
        {
            name:           "JSON handler",
            method:         "GET",
            path:           "/api/status",
            expectedStatus: http.StatusOK,
            expectedBody:   `{"status": "ok"}`,
            setupRouter:    setupTestRouter,
        },
        {
            name:           "group route",
            method:         "GET",
            path:           "/api/v1/users",
            expectedStatus: http.StatusOK,
            expectedBody:   "users",
            setupRouter:    setupTestRouter,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            router := tt.setupRouter()
            w := executeRequest(router, tt.method, tt.path)
            
            assertResponse(t, w, tt.expectedStatus, tt.expectedBody)
        })
    }
}
```

#### Middleware Testing

```go
func TestMiddlewareExecution(t *testing.T) {
    router := setupMiddlewareRouter()
    w := executeRequest(router, "GET", "/test")
    
    assert.Equal(t, http.StatusOK, w.Code)
    assert.Equal(t, "test", w.Body.String())
    assert.Equal(t, "before", w.Header().Get("X-Test-Middleware"))
}

func TestMiddlewareOrder(t *testing.T) {
    var executionOrder []string
    
    router := rtr.NewRouter()
    router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
        rtr.NewAnonymousMiddleware(func(next http.Handler) http.Handler {
            return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                executionOrder = append(executionOrder, "global-before")
                next.ServeHTTP(w, r)
            })
        }),
    })
    
    group := rtr.NewGroup().SetPrefix("/api")
    group.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
        rtr.NewAnonymousMiddleware(func(next http.Handler) http.Handler {
            return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                executionOrder = append(executionOrder, "group-before")
                next.ServeHTTP(w, r)
            })
        }),
    })
    
    group.AddRoute(rtr.Get("/test", func(w http.ResponseWriter, r *http.Request) {
        executionOrder = append(executionOrder, "handler")
        w.WriteHeader(http.StatusOK)
    }))
    
    router.AddGroup(group)
    
    executeRequest(router, "GET", "/api/test")
    
    expected := []string{"global-before", "group-before", "handler"}
    assert.Equal(t, expected, executionOrder)
}
```

#### Parameter Testing

```go
func TestParameterExtraction(t *testing.T) {
    router := rtr.NewRouter()
    
    router.AddRoute(rtr.Get("/users/:id", func(w http.ResponseWriter, r *http.Request) {
        id := rtr.MustGetParam(r, "id")
        w.Write([]byte(id))
    }))
    
    tests := []struct {
        path         string
        expectedBody string
    }{
        {"/users/123", "123"},
        {"/users/abc", "abc"},
        {"/users/user-123", "user-123"},
    }
    
    for _, tt := range tests {
        t.Run(tt.path, func(t *testing.T) {
            w := executeRequest(router, "GET", tt.path)
            assertResponse(t, w, http.StatusOK, tt.expectedBody)
        })
    }
}
```

## Code Quality

### Linting

Install and configure golangci-lint:

```bash
# Install
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run
golangci-lint run
```

**`.golangci.yml` configuration:**
```yaml
run:
  timeout: 5m
  tests: true

linters:
  enable:
    - gofmt
    - goimports
    - govet
    - errcheck
    - staticcheck
    - unused
    - gosimple
    - structcheck
    - varcheck
    - ineffassign
    - deadcode
    - typecheck
    - gosec
    - misspell
    - unconvert
    - dupl
    - goconst
    - gocyclo

linters-settings:
  goconst:
    min-len: 3
    min-occurrences: 3
  gocyclo:
    min-complexity: 10
  dupl:
    threshold: 100
```

### Code Formatting

Use standard Go formatting:

```bash
# Format code
go fmt ./...

# Organize imports
goimports -w .
```

### Documentation

#### Public API Documentation

All public APIs must have documentation:

```go
// NewRouter creates a new router instance with default configuration.
// The router implements http.Handler interface and can be used directly
// with http.ListenAndServe.
//
// Example:
//     router := rtr.NewRouter()
//     router.AddRoute(rtr.Get("/", handler))
//     http.ListenAndServe(":8080", router)
func NewRouter() RouterInterface {
    return &routerImpl{}
}
```

#### Package Documentation

Each package should have package documentation:

```go
// Package rtr provides a high-performance HTTP router for Go applications.
// It supports middleware chaining, route grouping, domain-based routing,
// and multiple handler types for different response formats.
//
// Basic usage:
//
//     router := rtr.NewRouter()
//     router.AddRoute(rtr.Get("/", func(w http.ResponseWriter, r *http.Request) {
//         w.Write([]byte("Hello, World!"))
//     }))
//     http.ListenAndServe(":8080", router)
package rtr
```

## Contributing

### Workflow

1. **Fork the repository**
2. **Create a feature branch**
3. **Make your changes**
4. **Add tests**
5. **Run the test suite**
6. **Submit a pull request**

```bash
# Create feature branch
git checkout -b feature/new-feature

# Make changes and commit
git add .
git commit -m "feat: add new feature"

# Push to fork
git push origin feature/new-feature
```

### Commit Message Convention

Use conventional commits:

```
feat: add new feature
fix: resolve bug in middleware execution
docs: update API documentation
test: add tests for parameter extraction
refactor: improve route matching algorithm
chore: update dependencies
```

### Pull Request Guidelines

#### PR Template

```markdown
## Description
Brief description of the change

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Added new tests for new functionality

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No breaking changes (or documented)
```

#### Review Criteria

- **Functionality**: Does the code work as intended?
- **Tests**: Are there adequate tests?
- **Documentation**: Is the code documented?
- **Style**: Does it follow project conventions?
- **Performance**: Does it maintain or improve performance?

## Performance Testing

### Benchmark Suite

Run comprehensive benchmarks:

```bash
# Run all benchmarks
go test -bench=. -benchmem ./...

# Run specific benchmarks
go test -bench=BenchmarkRouteMatching -benchmem ./...

# Compare benchmarks
go test -bench=. -benchmem -cpuprofile=cpu.prof -memprofile=mem.prof ./...
```

### Performance Guidelines

#### Route Matching

- Use efficient data structures
- Minimize string allocations
- Optimize parameter extraction
- Consider cache-friendly layouts

#### Middleware

- Keep middleware lightweight
- Avoid unnecessary allocations
- Use object pools for frequently allocated objects
- Profile critical paths

#### Memory Usage

- Reuse buffers where possible
- Avoid memory leaks
- Monitor garbage collection impact
- Use sync.Pool for object reuse

## Debugging

### Debug Routes

Use the built-in listing functionality:

```go
router := setupRouter()
router.List() // Prints detailed router configuration
```

### Debug Middleware

Add debugging middleware:

```go
router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    rtr.NewAnonymousMiddleware(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            fmt.Printf("Before: %s %s\n", r.Method, r.URL.Path)
            start := time.Now()
            next.ServeHTTP(w, r)
            duration := time.Since(start)
            fmt.Printf("After: %v\n", duration)
        })
    }),
})
```

### Debug Parameters

Add parameter logging:

```go
router.AddRoute(rtr.Get("/users/:id", func(w http.ResponseWriter, r *http.Request) {
    params := rtr.GetParams(r)
    fmt.Printf("Parameters: %+v\n", params)
    
    id := rtr.MustGetParam(r, "id")
    w.Write([]byte("User ID: " + id))
}))
```

## Release Process

### Version Management

Use semantic versioning:

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Checklist

1. **Update version number**
2. **Update CHANGELOG.md**
3. **Run full test suite**
4. **Update documentation**
5. **Create git tag**
6. **Create GitHub release**
7. **Update Go module version**

```bash
# Create release tag
git tag -a v1.2.3 -m "Release version 1.2.3"
git push origin v1.2.3

# Update Go module
go mod tidy
```

## See Also

- [Getting Started Guide](getting_started.md) - Learn basic usage
- [API Reference](api_reference.md) - Complete API documentation
- [Architecture Documentation](architecture.md) - System design overview
- [Troubleshooting Guide](troubleshooting.md) - Common issues and solutions
