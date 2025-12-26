---
path: modules/handlers.md
page-type: module
summary: Handler types, execution, and response generation for different content formats.
tags: [module, handlers, responses, content-types, execution]
created: 2025-12-26
updated: 2025-12-26
version: 1.0.0
---

# Handlers Module

The handlers module provides a comprehensive system for handling HTTP requests with different response types. It supports standard HTTP handlers, specialized content handlers, and automatic content-type management.

## Overview

Handlers are the core components that process HTTP requests and generate responses. RTR supports multiple handler types to simplify common response patterns while maintaining full flexibility for complex scenarios.

## Handler Types

### Standard Handler (StdHandler)

The standard Go HTTP handler with full control over the response:

```go
type StdHandler func(http.ResponseWriter, *http.Request)
```

**Example:**
```go
func(w http.ResponseWriter, r *http.Request) {
    // Full control over headers
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("X-Custom-Header", "value")
    
    // Control over status code
    w.WriteHeader(http.StatusOK)
    
    // Control over body
    w.Write([]byte(`{"message": "success"}`))
}
```

**Use Cases:**
- Complex response logic
- Custom headers
- Conditional responses
- Error handling
- Streaming responses

### String Handler (StringHandler)

Returns a string that is written directly to the response without setting headers:

```go
type StringHandler func(http.ResponseWriter, *http.Request) string
```

**Example:**
```go
func(w http.ResponseWriter, r *http.Request) string {
    // Set headers manually if needed
    w.Header().Set("Content-Type", "text/plain")
    
    // Return content string
    return "Hello, World!"
}
```

**Use Cases:**
- Simple text responses
- Custom header management
- Template rendering with manual headers

### HTML Handler (HTMLHandler)

Returns HTML content with automatic `Content-Type: text/html; charset=utf-8` header:

```go
type HTMLHandler StringHandler
```

**Example:**
```go
func(w http.ResponseWriter, r *http.Request) string {
    return `<!DOCTYPE html>
<html>
<head>
    <title>My Page</title>
    <meta charset="utf-8">
</head>
<body>
    <h1>Hello, World!</h1>
</body>
</html>`
}
```

**Use Cases:**
- Web pages
- HTML templates
- Server-side rendering
- Static HTML content

### JSON Handler (JSONHandler)

Returns JSON content with automatic `Content-Type: application/json` header:

```go
type JSONHandler StringHandler
```

**Example:**
```go
func(w http.ResponseWriter, r *http.Request) string {
    return `{
    "status": "ok",
    "version": "1.0.0",
    "timestamp": "` + time.Now().Format(time.RFC3339) + `"
}`
}
```

**Use Cases:**
- REST API responses
- AJAX responses
- Configuration data
- Status endpoints

### CSS Handler (CSSHandler)

Returns CSS content with automatic `Content-Type: text/css` header:

```go
type CSSHandler StringHandler
```

**Example:**
```go
func(w http.ResponseWriter, r *http.Request) string {
    return `body {
    font-family: Arial, sans-serif;
    background-color: #f0f0f0;
    margin: 0;
    padding: 20px;
}

h1 {
    color: #333;
    border-bottom: 2px solid #007acc;
}`
}
```

**Use Cases:**
- Dynamic stylesheets
- Theme switching
- CSS generation
- Asset serving

### XML Handler (XMLHandler)

Returns XML content with automatic `Content-Type: application/xml` header:

```go
type XMLHandler StringHandler
```

**Example:**
```go
func(w http.ResponseWriter, r *http.Request) string {
    return `<?xml version="1.0" encoding="UTF-8"?>
<users>
    <user id="1">
        <name>John Doe</name>
        <email>john@example.com</email>
    </user>
    <user id="2">
        <name>Jane Smith</name>
        <email>jane@example.com</email>
    </user>
</users>`
}
```

**Use Cases:**
- API responses (XML)
- Sitemaps
- RSS feeds
- Configuration files

### Text Handler (TextHandler)

Returns plain text content with automatic `Content-Type: text/plain; charset=utf-8` header:

```go
type TextHandler StringHandler
```

**Example:**
```go
func(w http.ResponseWriter, r *http.Request) string {
    return `User-agent: *
Disallow: /admin/
Disallow: /private/
Allow: /
Sitemap: https://example.com/sitemap.xml`
}
```

**Use Cases:**
- Plain text files
- Configuration files
- Documentation
- Log files

### JavaScript Handler (JSHandler)

Returns JavaScript content with automatic `Content-Type: application/javascript` header:

```go
type JSHandler StringHandler
```

**Example:**
```go
func(w http.ResponseWriter, r *http.Request) string {
    return `// Generated JavaScript
(function() {
    'use strict';
    
    console.log('Hello from RTR Router!');
    
    function initApp() {
        document.addEventListener('DOMContentLoaded', function() {
            console.log('App initialized');
        });
    }
    
    initApp();
})();`
}
```

**Use Cases:**
- Dynamic JavaScript
- Configuration scripts
- Client-side logic
- Asset serving

### Error Handler (ErrorHandler)

Returns an error that is handled automatically:

```go
type ErrorHandler func(http.ResponseWriter, *http.Request) error
```

**Example:**
```go
func(w http.ResponseWriter, r *http.Request) error {
    userID := rtr.MustGetParam(r, "id")
    
    user, err := getUserByID(userID)
    if err != nil {
        return fmt.Errorf("user not found: %s", userID)
    }
    
    // Success case - no error, no output
    return nil
}
```

**Use Cases:**
- Error-prone operations
- Validation logic
- Database operations
- File operations

## Handler Priority

When multiple handlers are set on a route, the router uses this priority order:

1. **StdHandler** (highest priority)
2. **HTMLHandler**
3. **JSONHandler**
4. **CSSHandler**
5. **XMLHandler**
6. **TextHandler**
7. **JSHandler**
8. **ErrorHandler** (lowest priority)

### Priority Example

```go
route := rtr.NewRoute().
    SetHandler(standardHandler).           // Priority 1
    SetJSONHandler(jsonHandler).           // Priority 3
    SetHTMLHandler(htmlHandler).           // Priority 2
    SetErrorHandler(errorHandler)           // Priority 8

// Only standardHandler will execute
```

## Handler Creation

### Route Handler Assignment

```go
// Standard handler
router.AddRoute(rtr.Get("/standard", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Standard response"))
}))

// JSON handler
router.AddRoute(rtr.GetJSON("/api/status", func(w http.ResponseWriter, r *http.Request) string {
    return `{"status": "ok"}`
}))

// HTML handler
router.AddRoute(rtr.GetHTML("/page", func(w http.ResponseWriter, r *http.Request) string {
    return "<h1>Hello, World!</h1>"
}))
```

### Method Chaining

```go
route := rtr.NewRoute().
    SetMethod("GET").
    SetPath("/api/data").
    SetJSONHandler(func(w http.ResponseWriter, r *http.Request) string {
        return `{"data": "value"}`
    }).
    SetName("Data Endpoint")
```

### Configuration-Based

```go
config := rtr.RouteConfig{
    Method: "GET",
    Path:   "/api/users",
    JSONHandler: func(w http.ResponseWriter, r *http.Request) string {
        return `{"users": []}`
    },
    Name: "List Users",
}
```

## Response Helper Functions

RTR provides helper functions that can be used in standard handlers:

### JSONResponse

```go
func JSONResponse(w http.ResponseWriter, r *http.Request, body string)
```

```go
router.AddRoute(rtr.Get("/manual-json", func(w http.ResponseWriter, r *http.Request) {
    rtr.JSONResponse(w, r, `{"message": "Manual JSON"}`)
    // Equivalent to JSONHandler but with more control
}))
```

### HTMLResponse

```go
func HTMLResponse(w http.ResponseWriter, r *http.Request, body string)
```

```go
router.AddRoute(rtr.Get("/manual-html", func(w http.ResponseWriter, r *http.Request) {
    // Custom logic
    if isMobile(r) {
        rtr.HTMLResponse(w, r, "<h1>Mobile View</h1>")
    } else {
        rtr.HTMLResponse(w, r, "<h1>Desktop View</h1>")
    }
}))
```

### Other Response Helpers

```go
rtr.CSSResponse(w, r, "body { color: red; }")
rtr.XMLResponse(w, r, "<?xml version='1.0'?><root></root>")
rtr.TextResponse(w, r, "Plain text")
rtr.JSResponse(w, r, "console.log('Hello');")
```

## Advanced Handler Patterns

### Conditional Responses

```go
router.AddRoute(rtr.Get("/conditional", func(w http.ResponseWriter, r *http.Request) {
    accept := r.Header.Get("Accept")
    
    switch {
    case strings.Contains(accept, "application/json"):
        rtr.JSONResponse(w, r, `{"format": "json"}`)
    case strings.Contains(accept, "text/html"):
        rtr.HTMLResponse(w, r, "<h1>HTML Format</h1>")
    default:
        rtr.TextResponse(w, r, "Plain text format")
    }
}))
```

### Template-Based Handlers

```go
func templateHandler(templateName string) rtr.HTMLHandler {
    return func(w http.ResponseWriter, r *http.Request) string {
        // Load template
        tmpl, err := template.ParseFiles("templates/" + templateName)
        if err != nil {
            http.Error(w, "Template error", http.StatusInternalServerError)
            return ""
        }
        
        // Execute template
        var buf bytes.Buffer
        err = tmpl.Execute(&buf, map[string]interface{}{
            "Title": "My Page",
            "User":  getCurrentUser(r),
        })
        if err != nil {
            http.Error(w, "Template execution error", http.StatusInternalServerError)
            return ""
        }
        
        return buf.String()
    }
}

router.AddRoute(rtr.GetHTML("/home", templateHandler("home.html")))
```

### Streaming Handlers

```go
router.AddRoute(rtr.Get("/stream", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    w.Header().Set("Transfer-Encoding", "chunked")
    
    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "Streaming not supported", http.StatusInternalServerError)
        return
    }
    
    for i := 0; i < 10; i++ {
        fmt.Fprintf(w, "Chunk %d\n", i)
        flusher.Flush()
        time.Sleep(500 * time.Millisecond)
    }
}))
```

### Error Handling Handlers

```go
router.AddRoute(rtr.Get("/risky", func(w http.ResponseWriter, r *http.Request) error {
    operation := r.URL.Query().Get("operation")
    
    switch operation {
    case "fail":
        return errors.New("operation failed")
    case "panic":
        panic("simulated panic")
    case "success":
        fmt.Fprintln(w, "Operation succeeded")
        return nil
    default:
        return errors.New("invalid operation")
    }
}))
```

## Handler Testing

### Unit Testing Handlers

```go
func无用 TestJSONHandlerler(t *testing.T) {
    handler := func(w http.ResponseWriter, r *http.Request) string {
        return `{"status": "ok"}`
    }
    
    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()
    
    handler(w, req)
    
    assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
    assert.Equal(t, `{"status": "ok"}`, w.Body.String())
}

func TestErrorHandler(t *testing.T) {
    handler := func(w http.ResponseWriter, r *http.Request) error {
        return errors.New("test error")
    }
    
    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()
    
    err := handler(w, req)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "test error")
}
```

### Integration Testing Handlers

```go
func TestHandlerIntegration(t *testing.T) {
    router := rtr.NewRouter()
    
    router.AddRoute(rtr.GetJSON("/api/test", func(w http.ResponseWriter, r *http.Request) string {
        return `{"message": "test"}`
    }))
    
    req := httptest.NewRequest("GET", "/api/test", nil)
    w := httptest.NewRecorder()
    
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
    assert.Equal(t, `{"message": "test"}`, w.Body.String())
}
```

## Performance Considerations

### Handler Efficiency

```go
// Good: Efficient string concatenation
func efficientHandler(w http.ResponseWriter, r *http.Request) string {
    var buf strings.Builder
    buf.WriteString(`{"users": [`)
    for i, user := range users {
        if i > 0 {
            buf.WriteString(",")
        }
        buf.WriteString(fmt.Sprintf(`{"id": "%d", "name": "%s"}`, user.ID, user.Name))
    }
    buf.WriteString("]}")
    return buf.String()
}

// Avoid: Inefficient concatenation
func inefficientHandler(w http.ResponseWriter, r *http.Request) string {
    result := `{"users": [`
    for i, user := range users {
        if i > 0 {
            result += ","
        }
        result += fmt.Sprintf(`{"id": "%d", "name": "%s"}`, user.ID, user.Name)
    }
    result += "]}"
    return result
}
```

### Memory Usage

- **String Handlers**: Return strings, consider large responses
- **Template Handlers**: Use buffers for template execution
- **Streaming Handlers**: Use flushers for real-time data

## Best Practices

### 1. Choose Appropriate Handler Type

```go
// Good: Use JSON handler for JSON responses
router.AddRoute(rtr.GetJSON("/api/data", jsonHandler))

// Good: Use HTML handler for HTML responses
router.AddRoute(rtr.GetHTML("/page", htmlHandler))

// Good: Use standard handler for complex logic
router.AddRoute(rtr.Get("/complex", complexHandler))
```

### 2. Handle Errors Gracefully

```go
func robustHandler(w http.ResponseWriter, r *http.Request) string {
    data, err := fetchData()
    if err != nil {
        http.Error(w, "Data fetch failed", http.StatusInternalServerError)
        return ""
    }
    return formatJSON(data)
}
```

### 3. Validate Input

```go
func validatedHandler(w http.ResponseWriter, r *http.Request) string {
    id := rtr.MustGetParam(r, "id")
    
    if !isValidID(id) {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return ""
    }
    
    return processData(id)
}
```

### 4. Use Response Helpers

```go
// Good: Use response helpers
func helperHandler(w http.ResponseWriter, r *http.Request) {
    rtr.JSONResponse(w, r, `{"status": "ok"}`)
}

// Alternative: Use specialized handlers
func specializedHandler(w http.ResponseWriter, r *http.Request) string {
    return `{"status": "ok"}`
}
```

### 5. Optimize for Content Type

```go
// Good: Match handler type to content
router.AddRoute(rtr.GetJSON("/api/data", func(w http.ResponseWriter, r *http.Request) string {
    return `{"data": "value"}`
}))

// Avoid: Manual content-type setting in wrong handler
router.AddRoute(rtr.Get("/api/data", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json") // Unnecessary
    w.Write([]byte(`{"data": "value"}`))
}))
```

## See Also

- [Routes Module](routes.md) - Route management and handler assignment
- [Responses Module](responses.md) - Response helper functions
- [Parameters Module](parameters.md) - Parameter extraction and validation
- [API Reference](../api_reference.md) - Complete API documentation
