# Testing Guide

## Overview

This guide covers testing strategies for your router, routes, and middleware. It includes examples of unit tests, integration tests, and best practices for testing HTTP handlers.

## Table of Contents
- [Testing Routes](#testing-routes)
- [Testing Middleware](#testing-middleware)
- [Testing Error Cases](#testing-error-cases)
- [Mocking Dependencies](#mocking-dependencies)
- [Test Helpers](#test-helpers)
- [Benchmarking](#benchmarking)
- [Best Practices](#best-practices)

## Testing Routes

### Basic Route Test

```go
func TestGetUser(t *testing.T) {
    // Setup router
    r := rtr.NewRouter()
    r.AddRoute(rtr.Get("/users/:id", getUserHandler))
    
    // Create test request
    req := httptest.NewRequest("GET", "/users/123", nil)
    w := httptest.NewRecorder()
    
    // Serve the request
    r.ServeHTTP(w, req)
    
    // Assertions
    if w.Code != http.StatusOK {
        t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
    }
    
    var response map[string]interface{}
    if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
        t.Fatalf("failed to parse response: %v", err)
    }
    
    if response["id"] != "123" {
        t.Errorf("expected user ID 123, got %s", response["id"])
    }
}
```

### Testing Route Parameters

```go
func TestGetUserWithParams(t *testing.T) {
    r := rtr.NewRouter()
    r.AddRoute(rtr.Get("/users/:id", getUserHandler))
    
    req := httptest.NewRequest("GET", "/users/123?name=test", nil)
    w := httptest.NewRecorder()
    
    r.ServeHTTP(w, req)
    
    // Test URL parameters
    if w.Code != http.StatusOK {
        t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
    }
}
```

## Testing Middleware

### Testing Authentication Middleware

```go
func TestAuthMiddleware(t *testing.T) {
    tests := []struct {
        name           string
        authHeader     string
        expectedStatus int
    }{
        {"Valid token", "Bearer valid-token", http.StatusOK},
        {"Missing token", "", http.StatusUnauthorized},
        {"Invalid token", "Bearer invalid-token", http.StatusForbidden},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create a test handler that requires authentication
            handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(http.StatusOK)
                w.Write([]byte("OK"))
            })
            
            // Apply the auth middleware
            req := httptest.NewRequest("GET", "/protected", nil)
            if tt.authHeader != "" {
                req.Header.Set("Authorization", tt.authHeader)
            }
            
            w := httptest.NewRecorder()
            authMiddleware(handler).ServeHTTP(w, req)
            
            if w.Code != tt.expectedStatus {
                t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
            }
        })
    }
}
```

## Testing Error Cases

### Testing 404 Not Found

```go
func TestNotFound(t *testing.T) {
    r := rtr.NewRouter()
    r.AddRoute(rtr.Get("/users", listUsersHandler))
    
    req := httptest.NewRequest("GET", "/nonexistent", nil)
    w := httptest.NewRecorder()
    
    r.ServeHTTP(w, req)
    
    if w.Code != http.StatusNotFound {
        t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
    }
}
```

## Mocking Dependencies

### Using Testify Mock

```go
// UserService is an interface that our handler depends on
type UserService interface {
    GetUser(id string) (*User, error)
}

// MockUserService is a mock implementation of UserService
type MockUserService struct {
    mock.Mock
}

func (m *MockUserService) GetUser(id string) (*User, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

func TestGetUserHandler(t *testing.T) {
    // Setup mock
    mockService := new(MockUserService)
    mockService.On("GetUser", "123").Return(&User{ID: "123", Name: "Test User"}, nil)
    
    // Create handler with mock dependency
    handler := &UserHandler{service: mockService}
    
    // Test the handler
    req := httptest.NewRequest("GET", "/users/123", nil)
    w := httptest.NewRecorder()
    
    handler.ServeHTTP(w, req)
    
    // Assertions
    if w.Code != http.StatusOK {
        t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
    }
    
    // Verify mock expectations
    mockService.AssertExpectations(t)
}
```

## Test Helpers

### Create Test Router

```go
func newTestRouter() *rtr.Router {
    r := rtr.NewRouter()
    // Add common test routes and middleware here
    return r
}

func TestWithTestRouter(t *testing.T) {
    r := newTestRouter()
    // Test with the common router setup
}
```

### Assert JSON Response

```go
func assertJSON(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedJSON string) {
    t.Helper()
    
    if w.Code != expectedStatus {
        t.Errorf("expected status %d, got %d", expectedStatus, w.Code)
    }
    
    var got, want interface{}
    if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
        t.Fatalf("failed to parse response: %v", err)
    }
    
    if err := json.Unmarshal([]byte(expectedJSON), &want); err != nil {
        t.Fatalf("invalid expected JSON: %v", err)
    }
    
    if !reflect.DeepEqual(got, want) {
        t.Errorf("response does not match expected\ngot:  %+v\nwant: %+v", got, want)
    }
}
```

## Benchmarking

### Benchmark Route Matching

```go
func BenchmarkRouteMatching(b *testing.B) {
    r := rtr.NewRouter()
    
    // Add many routes to test matching performance
    for i := 0; i < 1000; i++ {
        path := fmt.Sprintf("/users/%d", i)
        r.AddRoute(rtr.Get(path, func(w http.ResponseWriter, r *http.Request) {}))
    }
    
    req := httptest.NewRequest("GET", "/users/999", nil)
    w := httptest.NewRecorder()
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        r.ServeHTTP(w, req)
    }
}
```

## Best Practices

1. **Table-Driven Tests**: Use table-driven tests for testing multiple scenarios.
2. **Test Coverage**: Aim for high test coverage, especially for critical paths.
3. **Parallel Testing**: Use `t.Parallel()` for independent tests to speed up test execution.
4. **Cleanup**: Clean up any test resources using `t.Cleanup()`.
5. **Test Helpers**: Create helper functions for common test assertions and setups.
6. **Benchmarking**: Include benchmarks for performance-critical paths.
7. **Golden Files**: Use golden files for testing complex responses.
8. **Mock External Services**: Mock external dependencies to make tests reliable and fast.

For more advanced testing patterns, see [Advanced Testing](./advanced-testing.md).
