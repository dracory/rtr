# Request Handlers

The router supports multiple handler types for different response formats, each automatically handling appropriate HTTP headers.

## Handler Priority

When multiple handlers are set on a route, they're used in this priority order:

| Priority | Handler Type   | Description                          | Auto Headers  |
|----------|----------------|--------------------------------------|--------------:|
| 1        | `Handler`      | Standard HTTP handler                | None          |
| 2        | `StringHandler`| Returns string without setting headers| None          |
| 3        | `HTMLHandler`  | Returns HTML content                 | text/html     |
| 4        | `JSONHandler`  | Returns JSON response                | application/json |
| 5        | `CSSHandler`   | Returns CSS styles                   | text/css      |
| 6        | `XMLHandler`   | Returns XML content                  | application/xml |
| 7        | `TextHandler`  | Returns plain text                   | text/plain    |
| 8        | `JSHandler`    | Returns JavaScript                   | application/javascript |
| 9        | `ErrorHandler` | Returns error responses              | None          |

## Handler Types

### Standard Handler
Full control over the HTTP response.

```go
router.Get("/users", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("X-Custom-Header", "value")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"ok"}`))
})
```

### String Handler
Returns a string without setting content-type headers.

```go
router.Get("/greet", func(w http.ResponseWriter, r *http.Request) string {
    return "Hello, World!"
})
```

### HTML Handler
Returns HTML with proper content-type header.

```go
router.GetHTML("/about", func(w http.ResponseWriter, r *http.Request) string {
    return "<h1>About</h1><p>Welcome to our site</p>"
})
```

### JSON Handler
Returns JSON with proper content-type header.

```go
router.GetJSON("/api/users", func(w http.ResponseWriter, r *http.Request) string {
    return `{"users": [{"id": 1, "name": "John"}]}`
})
```

### Error Handler
Handles error returns from handlers.

```go
router.Get("/secure", func(w http.ResponseWriter, r *http.Request) (string, error) {
    if !isAuthorized(r) {
        return "", errors.New("unauthorized")
    }
    return "Welcome admin!", nil
}).SetErrorHandler(func(err error, w http.ResponseWriter, r *http.Request) {
    if err.Error() == "unauthorized" {
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte("Access denied"))
    } else {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Something went wrong"))
    }
})
```

## Best Practices

1. Use the most specific handler type for your response format
2. Prefer typed handlers (JSON, HTML) over StringHandler for proper content-type headers
3. Use ErrorHandler for centralized error handling
4. Keep handler functions focused and delegate business logic to separate packages
