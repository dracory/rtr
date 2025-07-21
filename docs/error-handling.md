# Error Handling

## Overview

This guide covers error handling patterns in the router, including built-in error handling, custom error handlers, and best practices for error management in your application.

## Table of Contents
- [Built-in Error Handling](#built-in-error-handling)
- [Custom Error Handlers](#custom-error-handlers)
- [Error Handler Interface](#error-handler-interface)
- [Best Practices](#best-practices)
- [Common Patterns](#common-patterns)

## Built-in Error Handling

The router includes several built-in error handling features:

### Panic Recovery

The router includes a recovery middleware that catches panics and returns a 500 Internal Server Error response:

```go
router := rtr.NewRouter() // Recovery middleware is added by default
```

### Not Found Handler

Handle 404 Not Found errors:

```go
router.SetNotFoundHandler(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)
    w.Write([]byte("404 - Not Found"))
})
```

### Method Not Allowed

Handle 405 Method Not Allowed errors:

```go
router.SetMethodNotAllowedHandler(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusMethodNotAllowed)
    w.Write([]byte("405 - Method Not Allowed"))
})
```

## Custom Error Handlers

### Route-Level Error Handlers

Each route can have its own error handler:

```go
rtr.NewRoute().
    SetMethod("GET").
    SetPath("/users/:id").
    SetHandler(getUserHandler).
    SetErrorHandler(func(w http.ResponseWriter, r *http.Request) error {
        // Handle errors from the route handler
        if err := getUserHandler(w, r); err != nil {
            if errors.Is(err, ErrUserNotFound) {
                w.WriteHeader(http.StatusNotFound)
                w.Write([]byte("User not found"))
                return nil
            }
            return err // Will be handled by the global error handler
        }
        return nil
    })
```

### Global Error Handler

Set a global error handler to catch all unhandled errors:

```go
router.SetErrorHandler(func(w http.ResponseWriter, r *http.Request, err error) {
    log.Printf("Error handling %s %s: %v", r.Method, r.URL.Path, err)
    
    var status int
    switch {
    case errors.Is(err, ErrNotFound):
        status = http.StatusNotFound
    case errors.Is(err, ErrUnauthorized):
        status = http.StatusUnauthorized
    default:
        status = http.StatusInternalServerError
    }
    
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "error":   http.StatusText(status),
        "message": err.Error(),
        "status":  status,
    })
})
```

## Error Handler Interface

The router uses the `ErrorHandler` type for error handling:

```go
type ErrorHandler func(http.ResponseWriter, *http.Request) error
```

## Best Practices

1. **Use Custom Error Types**: Create custom error types for different error conditions.
2. **Handle Errors at the Right Level**: Handle errors as close to their source as possible.
3. **Log Errors**: Always log errors with enough context to debug issues.
4. **Return Appropriate Status Codes**: Use the correct HTTP status codes for different error conditions.
5. **Sanitize Error Messages**: Don't expose sensitive information in error messages.

## Common Patterns

### Error Wrapping

Use `fmt.Errorf` with `%w` to wrap errors with additional context:

```go
func getUser(id string) (*User, error) {
    user, err := db.GetUser(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user %s: %w", id, err)
    }
    return user, nil
}
```

### Error Responses

Use a consistent error response format:

```go
type ErrorResponse struct {
    Status  int    `json:"status"`
    Message string `json:"message,omitempty"`
    Details string `json:"details,omitempty"`
}

func writeError(w http.ResponseWriter, status int, message string, details string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(ErrorResponse{
        Status:  status,
        Message: message,
        Details: details,
    })
}
```

For more advanced error handling patterns, see [Advanced Error Handling](./advanced-error-handling.md).
