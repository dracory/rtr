# Handler Types Example

This example demonstrates all the different handler types available in the rtr router, showcasing how each handler type simplifies response generation by automatically setting appropriate Content-Type headers and handling response writing.

## Features

- **Complete Handler Coverage** - All 9 handler types: Handler, StringHandler, HTMLHandler, JSONHandler, CSSHandler, XMLHandler, TextHandler, JSHandler, ErrorHandler
- **ToHandler Function Demo** - Shows how to convert StringHandler to standard Handler
- **Handler Priority System** - Demonstrates how handlers are prioritized when multiple are set
- **URL Parameters** - Shows parameter extraction with different handler types
- **Interactive Web Interface** - Beautiful HTML overview of all examples
- **Error Handling Examples** - Proper error responses with ErrorHandler
- **Comprehensive Testing** - Full test coverage for all handler types

## Quick Start

1. **Run the example:**
   ```bash
   go run main.go
   ```

2. **Open your browser:**
   ```
   http://localhost:8080
   ```

3. **Explore the examples:**
   - Click through the interactive interface to see each handler type in action
   - Or use the direct endpoints listed below

4. **Run tests:**
   ```bash
   go test -v
   ```

## Complete Handler Types

### 1. Traditional Handler
Standard `func(w http.ResponseWriter, r *http.Request)` handler with full control over the response.

```go
r.AddRoute(rtr.NewRoute().
    SetPath("/traditional").
    SetHandler(func(w http.ResponseWriter, req *http.Request) {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        w.Write([]byte("<h1>Full Control</h1>"))
    }))
```

**Endpoint:** `GET /traditional`

### 2. StringHandler
Returns string without automatically setting any headers. Gives you full control over headers.

```go
r.AddRoute(rtr.NewRoute().
    SetPath("/raw").
    SetStringHandler(func(w http.ResponseWriter, req *http.Request) string {
        w.Header().Set("X-Custom-Header", "Raw Response")
        return "Raw string response without automatic Content-Type headers."
    }))
```

**Endpoint:** `GET /raw`

### 3. HTMLHandler
Returns HTML string, automatically sets `Content-Type: text/html; charset=utf-8`.

```go
r.AddRoute(rtr.NewRoute().
    SetPath("/html").
    SetHTMLHandler(func(w http.ResponseWriter, req *http.Request) string {
        return `<h1>HTML Handler</h1><p>Just return HTML string!</p>`
    }))
```

**Endpoints:**
- `GET /html` - Static HTML example
- `GET /user/:id` - Dynamic HTML with URL parameters

### 4. JSONHandler
Returns JSON string, automatically sets `Content-Type: application/json`.

```go
r.AddRoute(rtr.NewRoute().
    SetPath("/api/users").
    SetJSONHandler(func(w http.ResponseWriter, req *http.Request) string {
        return `{"users": ["alice", "bob", "charlie"]}`
    }))
```

**Endpoints:**
- `GET /api/users` - Static JSON example
- `GET /api/status` - Dynamic JSON with timestamp
- `GET /api/user/:id` - JSON with URL parameters

### 5. CSSHandler
Returns CSS string, automatically sets `Content-Type: text/css`.

```go
r.AddRoute(rtr.NewRoute().
    SetPath("/styles.css").
    SetCSSHandler(func(w http.ResponseWriter, req *http.Request) string {
        return `body { font-family: Arial, sans-serif; }`
    }))
```

**Endpoint:** `GET /styles.css`

### 6. XMLHandler
Returns XML string, automatically sets `Content-Type: application/xml`.

```go
r.AddRoute(rtr.NewRoute().
    SetPath("/api/data.xml").
    SetXMLHandler(func(w http.ResponseWriter, req *http.Request) string {
        return `<?xml version="1.0"?><data><users>...</users></data>`
    }))
```

**Endpoint:** `GET /api/data.xml`

### 7. TextHandler
Returns plain text, automatically sets `Content-Type: text/plain; charset=utf-8`.

```go
r.AddRoute(rtr.NewRoute().
    SetPath("/robots.txt").
    SetTextHandler(func(w http.ResponseWriter, req *http.Request) string {
        return `User-agent: *\nDisallow: /admin/`
    }))
```

**Endpoint:** `GET /robots.txt`

### 8. JSHandler
Returns JavaScript string, automatically sets `Content-Type: application/javascript`.

```go
r.AddRoute(rtr.NewRoute().
    SetPath("/script.js").
    SetJSHandler(func(w http.ResponseWriter, req *http.Request) string {
        return `console.log('Hello from RTR Router!');`
    }))
```

**Endpoint:** `GET /script.js`

### 9. ErrorHandler
Returns error for proper error handling. Allows you to return errors that can be handled appropriately.

```go
r.AddRoute(rtr.NewRoute().
    SetPath("/error-demo").
    SetErrorHandler(func(w http.ResponseWriter, req *http.Request) error {
        if someCondition {
            w.WriteHeader(http.StatusInternalServerError)
            return fmt.Errorf("something went wrong")
        }
        
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"status": "ok"}`))
        return nil // nil means success
    }))
```

**Endpoints:**
- `GET /error-demo` - Success case (returns nil error)
- `GET /error-demo?fail=true` - Error case (returns error)
- `GET /not-found-demo` - 404 error example

### 10. ToHandler Function
Utility function that converts any `StringHandler` to a standard `Handler`.

```go
r.AddRoute(rtr.NewRoute().
    SetPath("/to-handler-demo").
    SetHandler(rtr.ToHandler(func(w http.ResponseWriter, req *http.Request) string {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        return `<h1>Converted Handler</h1>`
    })))
```

**Endpoint:** `GET /to-handler-demo`

## Handler Priority System

When multiple handlers are set on the same route, the router uses this priority order:

1. **Handler** (traditional HTTP handler) - Highest priority
2. **StringHandler** 
3. **HTMLHandler**
4. **JSONHandler**
5. **CSSHandler**
6. **XMLHandler**
7. **TextHandler**
8. **JSHandler**
9. **ErrorHandler** - Lowest priority

### Priority Demo

The `/priority-demo` endpoint demonstrates this by setting both HTMLHandler and JSONHandler:

```go
r.AddRoute(rtr.NewRoute().
    SetPath("/priority-demo").
    SetHTMLHandler(func(w http.ResponseWriter, req *http.Request) string {
        return "<h1>HTML Handler</h1><p>HTMLHandler has higher priority</p>"
    }).
    SetJSONHandler(func(w http.ResponseWriter, req *http.Request) string {
        return `{"message": "This won't be returned due to priority"}`
    }))
```

**Result:** HTMLHandler executes because it has higher priority than JSONHandler.

## All Available Endpoints

| Endpoint | Handler Type | Description |
|----------|--------------|-------------|
| `GET /` | Traditional | Interactive overview of all examples |
| `GET /traditional` | Handler | Standard HTTP handler example |
| `GET /raw` | StringHandler | Raw string without automatic headers |
| `GET /html` | HTMLHandler | Static HTML generation |
| `GET /user/:id` | HTMLHandler | Dynamic HTML with parameters |
| `GET /api/users` | JSONHandler | Static JSON response |
| `GET /api/status` | JSONHandler | Dynamic JSON with timestamp |
| `GET /api/user/:id` | JSONHandler | JSON response with parameters |
| `GET /styles.css` | CSSHandler | CSS stylesheet generation |
| `GET /api/data.xml` | XMLHandler | XML data response |
| `GET /robots.txt` | TextHandler | Plain text file |
| `GET /script.js` | JSHandler | JavaScript file generation |
| `GET /error-demo` | ErrorHandler | Success case (nil error) |
| `GET /error-demo?fail=true` | ErrorHandler | Error case (returns error) |
| `GET /not-found-demo` | ErrorHandler | 404 error example |
| `GET /to-handler-demo` | ToHandler | StringHandler converted to Handler |
| `GET /priority-demo` | Multiple | Handler priority demonstration |

## URL Parameters

All handler types support URL parameters through the standard router parameter extraction:

```go
// HTMLHandler with parameters
r.AddRoute(rtr.NewRoute().
    SetPath("/user/:id").
    SetHTMLHandler(func(w http.ResponseWriter, req *http.Request) string {
        userID := rtr.MustGetParam(req, "id")
        return fmt.Sprintf("<h1>User %s</h1>", userID)
    }))

// JSONHandler with parameters
r.AddRoute(rtr.NewRoute().
    SetPath("/api/user/:id").
    SetJSONHandler(func(w http.ResponseWriter, req *http.Request) string {
        userID := rtr.MustGetParam(req, "id")
        return fmt.Sprintf(`{"user_id": "%s"}`, userID)
    }))
```

## Testing

The example includes comprehensive tests covering:

- **All Handler Types** - Verify each handler type works correctly
- **Content-Type Headers** - Ensure proper headers are set automatically
- **URL Parameters** - Test parameter extraction with different handlers
- **Handler Priority** - Verify priority system works as expected
- **Error Handling** - Test ErrorHandler success and error cases
- **ToHandler Function** - Verify conversion works properly
- **Response Content** - Validate actual response content

Run tests with:
```bash
go test -v                    # Run all tests
go test -run TestHTML         # Run specific handler tests
go test -cover               # Run with coverage
```

## Benefits of Specialized Handlers

### Simplified Code
Instead of manually setting headers and writing responses:

```go
// Traditional approach
func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}

// With JSONHandler
func(w http.ResponseWriter, r *http.Request) string {
    return `{"message": "Hello World"}`
}
```

### Automatic Content-Type Handling
- No need to remember correct MIME types
- Consistent header setting across your application
- Reduced boilerplate code

### Better Error Handling
ErrorHandler provides a clean way to handle errors:

```go
// Traditional error handling
func(w http.ResponseWriter, r *http.Request) {
    if err := doSomething(); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Error occurred"))
        return
    }
    w.Write([]byte("Success"))
}

// With ErrorHandler
func(w http.ResponseWriter, r *http.Request) error {
    if err := doSomething(); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return err
    }
    w.Write([]byte("Success"))
    return nil
}
```

### Better Code Organization
- Clear separation of concerns
- Handler type indicates response format
- Easier to understand and maintain

## Advanced Usage

### Combining with Middleware

All handler types work seamlessly with middleware:

```go
r.AddRoute(rtr.NewRoute().
    SetPath("/api/data").
    SetJSONHandler(jsonHandler).
    SetBeforeMiddleware(authMiddleware, loggingMiddleware))
```

### Dynamic Content Generation

Handlers can generate dynamic content using any Go functionality:

```go
SetJSONHandler(func(w http.ResponseWriter, req *http.Request) string {
    data := fetchFromDatabase()
    return generateJSON(data)
})
```

### Custom Headers with StringHandler

StringHandler gives you full control over headers:

```go
SetStringHandler(func(w http.ResponseWriter, req *http.Request) string {
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("X-Custom-Header", "value")
    return "Response with custom headers"
})
```

### Error Handling Patterns

ErrorHandler supports various error handling patterns:

```go
SetErrorHandler(func(w http.ResponseWriter, req *http.Request) error {
    // Validation
    if err := validateRequest(req); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return err
    }
    
    // Business logic
    result, err := processRequest(req)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return err
    }
    
    // Success
    w.Header().Set("Content-Type", "application/json")
    w.Write(result)
    return nil
})
```

### ToHandler Use Cases

ToHandler is useful when you need to convert a StringHandler to a standard Handler:

```go
// When you have a StringHandler but need a Handler
stringHandler := func(w http.ResponseWriter, r *http.Request) string {
    return "Hello World"
}

// Convert to Handler
handler := rtr.ToHandler(stringHandler)

// Use with middleware that expects Handler
middlewareFunc(handler)
```

## Best Practices

1. **Choose the Right Handler Type** - Use the handler that matches your response format
2. **Use ErrorHandler for Error Cases** - Provides cleaner error handling patterns
3. **StringHandler for Custom Headers** - When you need full control over headers
4. **Keep Handlers Simple** - Complex logic should be in separate functions
5. **Handle Parameters Safely** - Validate and sanitize URL parameters
6. **Test All Handler Types** - Ensure proper behavior and headers
7. **Document Handler Behavior** - Make it clear what each endpoint returns
8. **Use ToHandler When Needed** - For converting StringHandlers to Handlers

## Handler Type Selection Guide

| Use Case | Recommended Handler | Reason |
|----------|-------------------|---------|
| HTML pages | HTMLHandler | Automatic Content-Type, clean syntax |
| JSON APIs | JSONHandler | Automatic Content-Type, clean syntax |
| CSS files | CSSHandler | Proper MIME type, browser compatibility |
| JavaScript files | JSHandler | Proper MIME type, browser compatibility |
| XML responses | XMLHandler | Proper MIME type for XML consumers |
| Plain text | TextHandler | Proper encoding, simple syntax |
| Custom headers needed | StringHandler | Full control over response headers |
| Error handling | ErrorHandler | Clean error patterns, proper status codes |
| Complex responses | Handler | Full control when needed |
| Converting handlers | ToHandler | Bridge between StringHandler and Handler |

## Troubleshooting

### Common Issues

1. **Wrong Content-Type**
   - Ensure you're using the correct handler type for your content
   - Traditional handlers require manual header setting

2. **Handler Not Executing**
   - Check handler priority if multiple handlers are set
   - Verify route path and method match

3. **Parameters Not Found**
   - Ensure parameter names in path match extraction calls
   - Use `rtr.GetParam()` for optional parameters

4. **Headers Not Set**
   - StringHandler doesn't set headers automatically
   - Use other handler types for automatic header setting

5. **Error Handling Issues**
   - ErrorHandler requires returning error for error cases
   - Return nil for success cases

### Debug Tips

1. **Check Handler Priority** - Use `/priority-demo` to understand precedence
2. **Inspect Headers** - Use browser dev tools to verify Content-Type
3. **Test Incrementally** - Start with simple handlers and add complexity
4. **Use Logging** - Add logging to see which handlers execute
5. **Test Error Cases** - Verify ErrorHandler behavior with different scenarios

## Related Examples

- [Declarative Example](../declarative/) - Declarative router configuration
- [Domain Example](../domain/) - Domain-based routing
- [Basic Example](../basic/) - Simple router usage

## Contributing

When contributing to this example:

1. Add tests for new handler types
2. Update the interactive HTML interface
3. Document new features in this README
4. Follow Go conventions and best practices
5. Ensure all handler types are covered
