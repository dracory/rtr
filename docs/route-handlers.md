# Route Handlers

RTR Router supports various handler types for different response formats, each automatically handling appropriate HTTP headers.

## Handler Priority

When multiple handlers are set on a route, they are used in this priority order:

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
    return "<h1>About Us</h1><p>Welcome to our site</p>"
})
```

### JSON Handler
Returns JSON with proper content-type header.

```go
router.GetJSON("/api/users", func(w http.ResponseWriter, r *http.Request) string {
    return `{"users": [{"id": 1, "name": "John"}]}`
})
```

### CSS Handler
Returns CSS with proper content-type header.

```go
router.GetCSS("/styles.css", func(w http.ResponseWriter, r *http.Request) string {
    return `body { font-family: Arial; margin: 0; }`
})
```

### XML Handler
Returns XML with proper content-type header.

```go
router.GetXML("/sitemap.xml", func(w http.ResponseWriter, r *http.Request) string {
    return `<?xml version="1.0" encoding="UTF-8"?>
    <urlset>
        <url><loc>https://example.com</loc></url>
    </urlset>`
})
```

### Text Handler
Returns plain text with proper content-type header.

```go
router.GetText("/robots.txt", func(w http.ResponseWriter, r *http.Request) string {
    return "User-agent: *\nDisallow: /private/"
})
```

### JavaScript Handler
Returns JavaScript with proper content-type header.

```go
router.GetJS("/app.js", func(w http.ResponseWriter, r *http.Request) string {
    return "console.log('App initialized');"
})
```

### Error Handler
Handles errors by returning an error value.

```go
router.Get("/secure", func(w http.ResponseWriter, r *http.Request) error {
    if !isAuthenticated(r) {
        return errors.New("unauthorized")
    }
    return nil
})
```

## Response Helpers

Helper functions for standard handlers:

```go
// In a standard handler
r.Get("/api/data", func(w http.ResponseWriter, r *http.Request) {
    // JSON response
    rtr.JSONResponse(w, r, `{"status":"ok"}`)
    
    // HTML response
    rtr.HTMLResponse(w, r, "<h1>Hello</h1>")
    
    // Text response
    rtr.TextResponse(w, r, "Hello World")
    
    // Error response
    rtr.ErrorResponse(w, r, http.StatusBadRequest, "Invalid request")
})
```

## Handler Combinations

You can set multiple handlers on a single route. The router will use the highest priority handler that is set:

```go
// HTML handler takes priority over JSON handler
router.Get("/content")
    .SetHTMLHandler(handleHTML)  // This will be used
    .SetJSONHandler(handleJSON)  // This will be ignored if handleHTML is set
```

## Dynamic Content with Parameters

All handler types work with path parameters:

```go
// HTML with parameters
router.GetHTML("/user/:id", func(w http.ResponseWriter, r *http.Request) string {
    userID := rtr.MustGetParam(r, "id")
    return fmt.Sprintf("<h1>User %s</h1>", userID)
})

// JSON with parameters
router.GetJSON("/api/user/:id", func(w http.ResponseWriter, r *http.Request) string {
    userID := rtr.MustGetParam(r, "id")
    return fmt.Sprintf(`{"id":"%s","name":"User %s"}`, userID, userID)
})
```

## Best Practices

1. Use the most specific handler type for your response format
2. Prefer typed handlers (JSON, HTML) over StringHandler for proper content-type headers
3. Use ErrorHandler for centralized error handling
4. Keep handler functions focused and delegate business logic to separate packages
5. Use response helpers in standard handlers for consistent behavior
