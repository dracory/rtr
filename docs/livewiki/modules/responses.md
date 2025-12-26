---
path: modules/responses.md
page-type: module
summary: Response helper functions and content-type management for HTTP responses.
tags: [module, responses, content-types, helpers, headers]
created: 2025-12-26
updated: 2025-12-26
version: 1.0.0
---

# Responses Module

The responses module provides helper functions for generating HTTP responses with appropriate content-type headers and standardized formatting. These functions are used internally by specialized handlers and can be used directly in standard handlers.

## Overview

Response helpers simplify the process of generating HTTP responses by automatically setting appropriate headers and handling content formatting. They provide a consistent way to generate responses across different content types.

## Response Helper Functions

### JSONResponse

Generates a JSON response with appropriate headers:

```go
func JSONResponse(w http.ResponseWriter, r *http.Request, body string)
```

**Headers Set:**
- `Content-Type: application/json`

**Example:**
```go
router.AddRoute(rtr.Get("/api/data", func(w http.ResponseWriter, r *http.Request) {
    data := `{
        "status": "success",
        "data": {
            "users": 123,
            "active": 45
        },
        "timestamp": "` + time.Now().Format(time.RFC3339) + `"
    }`
    
    rtr.JSONResponse(w, r, data)
}))
```

### HTMLResponse

Generates an HTML response with appropriate headers:

```go
func HTMLResponse(w http.ResponseWriter, r *http.Request, body string)
```

**Headers Set:**
- `Content-Type: text/html; charset=utf-8` (only if not already set)

**Example:**
```go
router.AddRoute(rtr.Get("/page", func(w http.ResponseWriter, r *http.Request) {
    html := `<!DOCTYPE html>
<html>
<head>
    <title>My Page</title>
    <meta charset="utf-8">
</head>
<body>
    <h1>Hello, World!</h1>
    <p>Generated at: ` + time.Now().Format(time.RFC1123) + `</p>
</body>
</html>`
    
    rtr.HTMLResponse(w, r, html)
}))
```

### CSSResponse

Generates a CSS response with appropriate headers:

```go
func CSSResponse(w http.ResponseWriter, r *http.Request, body string)
```

**Headers Set:**
- `Content-Type: text/css`

**Example:**
```go
router.AddRoute(rtr.Get("/styles/theme.css", func(w http.ResponseWriter, r *http.Request) {
    theme := r.URL.Query().Get("theme")
    
    var css string
    switch theme {
    case "dark":
        css = `body {
            background-color: #1a1a1a;
            color: #ffffff;
        }`
    case "light":
        css = `body {
            background-color: #ffffff;
            color: #000000;
        }`
    default:
        css = `body {
            background-color: #f0f0f0;
            color: #333333;
        }`
    }
    
    rtr.CSSResponse(w, r, css)
}))
```

### XMLResponse

Generates an XML response with appropriate headers:

```go
func XMLResponse(w http.ResponseWriter, r *http.Request, body string)
```

**Headers Set:**
- `Content-Type: application/xml`

**Example:**
```go
router.AddRoute(rtr.Get("/sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
    sitemap := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
    <url>
        <loc>https://example.com/</loc>
        <lastmod>` + time.Now().Format("2006-01-02") + `</lastmod>
        <changefreq>daily</changefreq>
        <priority>1.0</priority>
    </url>
    <url>
        <loc>https://example.com/about</loc>
        <lastmod>` + time.Now().Format("2006-01-02") + `</lastmod>
        <changefreq>monthly</changefreq>
        <priority>0.8</priority>
    </url>
</urlset>`
    
    rtr.XMLResponse(w, r, sitemap)
}))
```

### TextResponse

Generates a plain text response with appropriate headers:

```go
func TextResponse(w http.ResponseWriter, r *http.Request, body string)
```

**Headers Set:**
- `Content-Type: text/plain; charset=utf-8`

**Example:**
```go
router.AddRoute(rtr.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
    robots := `User-agent: *
Disallow: /admin/
Disallow: /private/
Allow: /
Sitemap: https://example.com/sitemap.xml`
    
    rtr.TextResponse(w, r, robots)
}))
```

### JSResponse

Generates a JavaScript response with appropriate headers:

```go
func JSResponse(w http.ResponseWriter, r *http.Request, body string)
```

**Headers Set:**
- `Content-Type: application/javascript`

**Example:**
```go
router.AddRoute(rtr.Get("/config.js", func(w http.ResponseWriter, r *http.Request) {
    config := `// Generated configuration
window.APP_CONFIG = {
    apiUrl: "` + os.Getenv("API_URL") + `",
    version: "` + os.Getenv("APP_VERSION") + `",
    environment: "` + os.Getenv("ENVIRONMENT") + `",
    debug: ` + os.Getenv("DEBUG") + `
};

console.log('Configuration loaded:', window.APP_CONFIG);`
    
    rtr.JSResponse(w, r, config)
}))
```

## Advanced Response Patterns

### Conditional Responses

```go
router.AddRoute(rtr.Get("/api/data", func(w http.ResponseWriter, r *http.Request) {
    accept := r.Header.Get("Accept")
    data := map[string]interface{}{
        "status": "success",
        "data":   "example data",
    }
    
    switch {
    case strings.Contains(accept, "application/json"):
        jsonData, _ := json.Marshal(data)
        rtr.JSONResponse(w, r, string(jsonData))
        
    case strings.Contains(accept, "application/xml"):
        xmlData := `<response><status>success</status><data>example data</data></response>`
        rtr.XMLResponse(w, r, xmlData)
        
    case strings.Contains(accept, "text/html"):
        htmlData := `<html><body><h1>Success</h1><p>example data</p></body></html>`
        rtr.HTMLResponse(w, r, htmlData)
        
    default:
        textData := "Status: success\nData: example data"
        rtr.TextResponse(w, r, textData)
    }
}))
```

### Template Responses

```go
func renderTemplate(w http.ResponseWriter, r *http.Request, templateName string, data interface{}) {
    tmpl, err := template.ParseFiles("templates/" + templateName)
    if err != nil {
        http.Error(w, "Template error", http.StatusInternalServerError)
        return
    }
    
    var buf bytes.Buffer
    err = tmpl.Execute(&buf, data)
    if err != nil {
        http.Error(w, "Template execution error", http.StatusInternalServerError)
        return
    }
    
    // Determine content type based on template extension
    switch filepath.Ext(templateName) {
    case ".html":
        rtr.HTMLResponse(w, r, buf.String())
    case ".xml":
        rtr.XMLResponse(w, r, buf.String())
    default:
        rtr.TextResponse(w, r, buf.String())
    }
}

router.AddRoute(rtr.Get("/page/:name", func(w http.ResponseWriter, r *http.Request) {
    pageName := rtr.MustGetParam(r, "name")
    data := map[string]interface{}{
        "Title":   strings.Title(pageName),
        "User":    getCurrentUser(r),
        "Time":    time.Now(),
    }
    
    renderTemplate(w, r, pageName+".html", data)
}))
```

### Streaming Responses

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
        chunk := fmt.Sprintf("Chunk %d\n", i)
        w.Write([]byte(chunk))
        flusher.Flush()
        time.Sleep(500 * time.Millisecond)
    }
    
    w.Write([]byte("Stream complete\n"))
    flusher.Flush()
}))
```

### File Downloads

```go
router.AddRoute(rtr.Get("/download/:filename", func(w http.ResponseWriter, r *http.Request) {
    filename := rtr.MustGetParam(r, "filename")
    
    // Validate filename
    if !isValidFilename(filename) {
        http.Error(w, "Invalid filename", http.StatusBadRequest)
        return
    }
    
    // Set download headers
    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
    w.Header().Set("Content-Type", "application/octet-stream")
    
    // Serve file content
    content := fmt.Sprintf("Content of file: %s\nGenerated at: %s", 
        filename, time.Now().Format(time.RFC3339))
    
    w.Write([]byte(content))
}))
```

## Error Responses

### Standard Error Responses

```go
func errorResponse(w http.ResponseWriter, r *http.Request, message string, code int) {
    w.WriteHeader(code)
    
    // Return error in requested format if possible
    accept := r.Header.Get("Accept")
    
    switch {
    case strings.Contains(accept, "application/json"):
        errorJSON := fmt.Sprintf(`{"error": "%s", "code": %d}`, message, code)
        rtr.JSONResponse(w, r, errorJSON)
        
    case strings.Contains(accept, "application/xml"):
        errorXML := fmt.Sprintf(`<error><message>%s</message><code>%d</code></error>`, message, code)
        rtr.XMLResponse(w, r, errorXML)
        
    default:
        rtr.TextResponse(w, r, fmt.Sprintf("Error %d: %s", code, message))
    }
}

router.AddRoute(rtr.Get("/error", func(w http.ResponseWriter, r *http.Request) {
    errorResponse(w, r, "This is a test error", http.StatusBadRequest)
}))
```

### Validation Errors

```go
func validationErrorResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
    w.WriteHeader(http.StatusBadRequest)
    
    errorsJSON, _ := json.Marshal(map[string]interface{}{
        "error":  "Validation failed",
        "fields": errors,
    })
    
    rtr.JSONResponse(w, r, string(errorsJSON))
}

router.AddRoute(rtr.Post("/users", func(w http.ResponseWriter, r *http.Request) {
    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        validationErrorResponse(w, r, map[string]string{
            "body": "Invalid JSON format",
        })
        return
    }
    
    validationErrors := validateUser(user)
    if len(validationErrors) > 0 {
        validationErrorResponse(w, r, validationErrors)
        return
    }
    
    // Process valid user
    rtr.JSONResponse(w, r, `{"status": "user created"}`)
}))
```

## Response Caching

### Cache Control Headers

```go
func cachedResponse(w http.ResponseWriter, r *http.Request, body string, contentType string, maxAge int) {
    w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
    w.Header().Set("ETag", generateETag(body))
    
    // Check if client has cached version
    if ifNoneMatch := r.Header.Get("If-None-Match"); ifNoneMatch == generateETag(body) {
        w.WriteHeader(http.StatusNotModified)
        return
    }
    
    switch contentType {
    case "json":
        rtr.JSONResponse(w, r, body)
    case "html":
        rtr.HTMLResponse(w, r, body)
    case "css":
        rtr.CSSResponse(w, r, body)
    case "xml":
        rtr.XMLResponse(w, r, body)
    case "js":
        rtr.JSResponse(w, r, body)
    default:
        rtr.TextResponse(w, r, body)
    }
}

router.AddRoute(rtr.Get("/static/data.json", func(w http.ResponseWriter, r *http.Request) {
    data := `{"data": "cached content", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`
    cachedResponse(w, r, data, "json", 300) // 5 minutes cache
}))
```

## Response Compression

### Gzip Compression

```go
func gzipMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
            next.ServeHTTP(w, r)
            return
        }
        
        w.Header().Set("Content-Encoding", "gzip")
        gzipWriter := gzip.NewWriter(w)
        defer gzipWriter.Close()
        
        wrappedWriter := &gzipResponseWriter{
            ResponseWriter: w,
            gzipWriter:    gzipWriter,
        }
        
        next.ServeHTTP(wrappedWriter, r)
    })
}

type gzipResponseWriter struct {
    http.ResponseWriter
    gzipWriter *gzip.Writer
}

func (gw *gzipResponseWriter) Write(data []byte) (int, error) {
    return gw.gzipWriter.Write(data)
}

// Usage
router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
    rtr.NewAnonymousMiddleware(gzipMiddleware),
})
```

## Testing Responses

### Unit Testing Response Helpers

```go
func TestJSONResponse(t *testing.T) {
    w := httptest.NewRecorder()
    r := httptest.NewRequest("GET", "/test", nil)
    
    body := `{"message": "test"}`
    rtr.JSONResponse(w, r, body)
    
    assert.Equal(t, http.StatusOK, w.Code)
    assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
    assert.Equal(t, body, w.Body.String())
}

func TestHTMLResponse(t *testing.T) {
    w := httptest.NewRecorder()
    r := httptest.NewRequest("GET", "/test", nil)
    
    body := "<h1>Test</h1>"
    rtr.HTMLResponse(w, r, body)
    
    assert.Equal(t, http.StatusOK, w.Code)
    assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
    assert.Equal(t, body, w.Body.String())
}
```

### Integration Testing Responses

```go
func TestResponseIntegration(t *testing.T) {
    router := rtr.NewRouter()
    
    router.AddRoute(rtr.Get("/json", func(w http.ResponseWriter, r *http.Request) {
        rtr.JSONResponse(w, r, `{"status": "ok"}`)
    }))
    
    router.AddRoute(rtr.Get("/html", func(w http.ResponseWriter, r *http.Request) {
        rtr.HTMLResponse(w, r, "<h1>OK</h1>")
    }))
    
    tests := []struct {
        path           string
        expectedBody   string
        expectedType  string
    }{
        {"/json", `{"status": "ok"}`, "application/json"},
        {"/html", "<h1>OK</h1>", "text/html; charset=utf-8"},
    }
    
    for _, tt := range tests {
        t.Run(tt.path, func(t *testing.T) {
            req := httptest.NewRequest("GET", tt.path, nil)
            w := httptest.NewRecorder()
            
            router.ServeHTTP(w, req)
            
            assert.Equal(t, http.StatusOK, w.Code)
            assert.Equal(t, tt.expectedType, w.Header().Get("Content-Type"))
            assert.Equal(t, tt.expectedBody, w.Body.String())
        })
    }
}
```

## Performance Considerations

### Response Generation Efficiency

```go
// Good: Use strings.Builder for large responses
func efficientJSONResponse(w http.ResponseWriter, r *http.Request, data []Item) {
    var buf strings.Builder
    buf.WriteString(`{"items": [`)
    
    for i, item := range data {
        if i > 0 {
            buf.WriteString(",")
        }
        buf.WriteString(fmt.Sprintf(`{"id": %d, "name": "%s"}`, item.ID, item.Name))
    }
    
    buf.WriteString("]}")
    rtr.JSONResponse(w, r, buf.String())
}

// Avoid: String concatenation in loops
func inefficientJSONResponse(w http.ResponseWriter, r *http.Request, data []Item) {
    result := `{"items": [`
    
    for i, item := range data {
        if i > 0 {
            result += ","
        }
        result += fmt.Sprintf(`{"id": %d, "name": "%s"}`, item.ID, item.Name)
    }
    
    result += "]}"
    rtr.JSONResponse(w, r, result)
}
```

### Memory Usage

- **Large Responses**: Consider streaming for very large responses
- **Template Caching**: Cache parsed templates for better performance
- **Response Buffers**: Use buffers for response construction

## Best Practices

### 1. Use Appropriate Response Helpers

```go
// Good: Use JSON helper for JSON responses
router.AddRoute(rtr.Get("/api/data", func(w http.ResponseWriter, r *http.Request) {
    rtr.JSONResponse(w, r, `{"data": "value"}`)
}))

// Good: Use HTML helper for HTML responses
router.AddRoute(rtr.Get("/page", func(w http.ResponseWriter, r *http.Request) {
    rtr.HTMLResponse(w, r, "<h1>Page</h1>")
}))
```

### 2. Handle Content Negotiation

```go
func contentNegotiatedResponse(w http.ResponseWriter, r *http.Request, data interface{}) {
    accept := r.Header.Get("Accept")
    
    switch {
    case strings.Contains(accept, "application/json"):
        jsonData, _ := json.Marshal(data)
        rtr.JSONResponse(w, r, string(jsonData))
    case strings.Contains(accept, "text/html"):
        // Render HTML template
        rtr.HTMLResponse(w, r, renderHTML(data))
    default:
        rtr.TextResponse(w, r, fmt.Sprintf("%+v", data))
    }
}
```

### 3. Set Appropriate Headers

```go
func apiResponse(w http.ResponseWriter, r *http.Request, data string) {
    // Set cache headers
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("X-API-Version", "1.0")
    
    // Use response helper
    rtr.JSONResponse(w, r, data)
}
```

### 4. Handle Errors Consistently

```go
func handleError(w http.ResponseWriter, r *http.Request, err error, code int) {
    w.WriteHeader(code)
    
    errorResponse := map[string]interface{}{
        "error":     err.Error(),
        "code":      code,
        "timestamp": time.Now().UTC().Format(time.RFC3339),
    }
    
    errorJSON, _ := json.Marshal(errorResponse)
    rtr.JSONResponse(w, r, string(errorJSON))
}
```

### 5. Optimize for Performance

```go
// Good: Pre-allocate buffers for known sizes
func optimizedResponse(w http.ResponseWriter, r *http.Request, items []Item) {
    buf := make([]byte, 0, len(items)*100) // Estimate size
    
    // Build response efficiently
    // ...
    
    rtr.JSONResponse(w, r, string(buf))
}
```

## See Also

- [Handlers Module](handlers.md) - Handler types and execution
- [Middleware Module](middleware.md) - Response modification middleware
- [Parameters Module](parameters.md) - Parameter extraction and validation
- [API Reference](../api_reference.md) - Complete API documentation
