---
path: api_reference.md
page-type: reference
summary: Complete API reference for RTR router with all interfaces, types, and methods.
tags: [api, reference, interfaces, methods, types]
created: 2025-12-26
updated: 2025-12-26
version: 1.0.0
---

# RTR Router API Reference

This document provides a comprehensive reference for all RTR router APIs, interfaces, types, and methods.

## Core Types

### Handler Types

#### StdHandler
```go
type StdHandler func(http.ResponseWriter, *http.Request)
```
Standard HTTP handler function signature. Provides full control over HTTP responses.

#### ErrorHandler
```go
type ErrorHandler func(http.ResponseWriter, *http.Request) error
```
Error handler that returns an error. If the error is nil, no content is written.

#### StringHandler
```go
type StringHandler func(http.ResponseWriter, *http.Request) string
```
Returns a string that will be written directly to the response without setting headers.

#### HTMLHandler
```go
type HTMLHandler func(http.ResponseWriter, *http.Request) string
```
Returns HTML content with automatic `Content-Type: text/html; charset=utf-8` header.

#### JSONHandler
```go
type JSONHandler func(http.ResponseWriter, *http.Request) string
```
Returns JSON content with automatic `Content-Type: application/json` header.

#### CSSHandler
```go
type CSSHandler func(http.ResponseWriter, *http.Request) string
```
Returns CSS content with automatic `Content-Type: text/css` header.

#### XMLHandler
```go
type XMLHandler func(http.ResponseWriter, *http.Request) string
```
Returns XML content with automatic `Content-Type: application/xml` header.

#### TextHandler
```go
type TextHandler func(http.ResponseWriter, *http.Request) string
```
Returns plain text content with automatic `Content-Type: text/plain; charset=utf-8` header.

#### JSHandler
```go
type JSHandler func(http.ResponseWriter, *http.Request) string
```
Returns JavaScript content with automatic `Content-Type: application/javascript` header.

### Middleware Types

#### StdMiddleware
```go
type StdMiddleware func(http.Handler) http.Handler
```
Standard middleware function following the Go HTTP middleware pattern.

#### MiddlewareInterface
```go
type MiddlewareInterface interface {
    GetName() string
    SetName(name string) MiddlewareInterface
    GetHandler() StdMiddleware
    SetHandler(handler StdMiddleware) MiddlewareInterface
    Execute(next http.Handler) http.Handler
}
```
Named middleware interface that provides metadata and debugging capabilities.

## Core Interfaces

### RouterInterface

The main router interface for managing routes, groups, and domains.

```go
type RouterInterface interface {
    // Prefix management
    GetPrefix() string
    SetPrefix(prefix string) RouterInterface
    
    // Route management
    AddRoute(route RouteInterface) RouterInterface
    AddRoutes(routes []RouteInterface) RouterInterface
    GetRoutes() []RouteInterface
    
    // Group management
    AddGroup(group GroupInterface) RouterInterface
    AddGroups(groups []GroupInterface) RouterInterface
    GetGroups() []GroupInterface
    
    // Domain management
    AddDomain(domain DomainInterface) RouterInterface
    AddDomains(domains []DomainInterface) RouterInterface
    GetDomains() []DomainInterface
    
    // Middleware management
    AddBeforeMiddlewares(middlewares []MiddlewareInterface) RouterInterface
    AddAfterMiddlewares(middlewares []MiddlewareInterface) RouterInterface
    GetBeforeMiddlewares() []MiddlewareInterface
    GetAfterMiddlewares() []MiddlewareInterface
    
    // HTTP handler
    ServeHTTP(w http.ResponseWriter, r *http.Request)
    
    // Utility methods
    List()
    String() string
}
```

#### Methods

##### NewRouter
```go
func NewRouter() RouterInterface
```
Creates a new router instance with default configuration.

##### NewRouterFromConfig
```go
func NewRouterFromConfig(config RouterConfig) InterfaceInterface
```
Creates a router from a declarative configuration.

### RouteInterface

Interface for managing individual routes.

```go
type RouteInterface interface {
    // HTTP method and path
    GetMethod() string
    SetMethod(method string) RouteInterface
    GetPath() string
    SetPath(path string) RouteInterface
    
    // Handlers
    GetHandler() StdHandler
    SetHandler(handler StdHandler) RouteInterface
    GetStringHandler() StringHandler
    SetStringHandler(handler StringHandler) RouteInterface
    GetHTMLHandler() HTMLHandler
    SetHTMLHandler(handler HTMLHandler) RouteInterface
    GetJSONHandler() JSONHandler
    SetJSONHandler(handler JSONHandler) RouteInterface
    GetCSSHandler() CSSHandler
    SetCSSHandler(handler CSSHandler) RouteInterface
    GetXMLHandler() XMLHandler
    SetXMLHandler(handler XMLHandler) RouteInterface
    GetTextHandler() TextHandler
    SetTextHandler(handler TextHandler) RouteInterface
    GetJSHandler() JSHandler
    SetJSHandler(handler JSHandler) RouteInterface
    GetErrorHandler() ErrorHandler
    SetErrorHandler(handler ErrorHandler) RouteInterface
    
    // Route metadata
    GetName() string
    SetName(name string) RouteInterface
    GetMetadata() map[string]interface{}
    SetMetadata(metadata map[string]interface{}) RouteInterface
    
    // Middleware
    AddBeforeMiddlewares(middlewares []MiddlewareInterface) RouteInterface
    AddAfterMiddlewares(middlewares []MiddlewareInterface) RouteInterface
    GetBeforeMiddlewares() []MiddlewareInterface
    GetAfterMiddlewares() []MiddlewareInterface
    
    // Utility methods
    String() string
}
```

#### Route Creation Functions

##### NewRoute
```go
func NewRoute() RouteInterface
```
Creates a new empty route instance.

##### HTTP Method Shortcuts
```go
func Get(path string, handler StdHandler) RouteInterface
func Post(path string, handler StdHandler) RouteInterface
func Put(path string, handler StdHandler) RouteInterface
func Delete(path string, handler StdHandler) RouteInterface
func Patch(path string, handler StdHandler) RouteInterface
func Options(path string, handler StdHandler) RouteInterface
func Head(path string, handler StdHandler) RouteInterface
```

##### Specialized Handler Shortcuts
```go
func GetHTML(path string, handler HTMLHandler) RouteInterface
func PostHTML(path string, handler HTMLHandler) RouteInterface
func GetJSON(path string, handler JSONHandler) RouteInterface
func PostJSON(path string, handler JSONHandler) RouteInterface
func GetCSS(path string, handler CSSHandler) RouteInterface
func GetXML(path string, handler XMLHandler) RouteInterface
func GetText(path string, handler TextHandler) RouteInterface
func GetJS(path string, handler JSHandler) RouteInterface
```

### GroupInterface

Interface for managing route groups.

```go
type GroupInterface interface {
    // Prefix management
    GetPrefix() string
    SetPrefix(prefix string) GroupInterface
    
    // Route management
    AddRoute(route RouteInterface) GroupInterface
    AddRoutes(routes []RouteInterface) GroupInterface
    GetRoutes() []RouteInterface
    
    // Nested group management
    AddGroup(group GroupInterface) GroupInterface
    AddGroups(groups []GroupInterface) GroupInterface
    GetGroups() []GroupInterface
    
    // Middleware management
    AddBeforeMiddlewares(middlewares []MiddlewareInterface) GroupInterface
    AddAfterMiddlewares(middlewares []MiddlewareInterface) GroupInterface
    GetBeforeMiddlewares() []MiddlewareInterface
    GetAfterMiddlewares() []MiddlewareInterface
    
    // Group metadata
    GetName() string
    SetName(name string) GroupInterface
    GetMetadata() map[string]interface{}
    SetMetadata(metadata map[string]interface{}) GroupInterface
    
    // Utility methods
    String() string
}
```

#### Group Creation Functions

##### NewGroup
```go
func NewGroup() GroupInterface
```
Creates a new empty group instance.

### DomainInterface

Interface for managing domain-based routing.

```go
type DomainInterface interface {
    // Domain patterns
    GetPatterns() []string
    SetPatterns(patterns []string) DomainInterface
    AddPattern(pattern string) DomainInterface
    
    // Route management
    AddRoute(route RouteInterface) DomainInterface
    AddRoutes(routes []RouteInterface) DomainInterface
    GetRoutes() []RouteInterface
    
    // Group management
    AddGroup(group GroupInterface) DomainInterface
    AddGroups(groups []GroupInterface) DomainInterface
    GetGroups() []GroupInterface
    
    // Middleware management
    AddBeforeMiddlewares(middlewares []MiddlewareInterface) DomainInterface
    AddAfterMiddlewares(middlewares []MiddlewareInterface) DomainInterface
    GetBeforeMiddlewares() []MiddlewareInterface
    GetAfterMiddlewares() []MiddlewareInterface
    
    // Domain metadata
    GetName() string
    SetName(name string) DomainInterface
    GetMetadata() map[string]interface{}
    SetMetadata(metadata map[string]interface{}) DomainInterface
    
    // Utility methods
    Matches(host string) bool
    String() string
}
```

#### Domain Creation Functions

##### NewDomain
```go
func NewDomain(patterns ...string) DomainInterface
```
Creates a new domain with one or more patterns.

## Parameter Functions

### GetParam
```go
func GetParam(r *http.Request, key string) (string, bool)
```
Safely gets a path parameter. Returns the parameter value and a boolean indicating if it exists.

### MustGetParam
```go
func MustGetParam(r *http.Request, key string) string
```
Gets a required path parameter. Panics if the parameter doesn't exist.

### GetParams
```go
func GetParams(r *http.Request) map[string]string
```
Returns all path parameters as a map.

### SetParam
```go
func SetParam(r *http.Request, key, value string)
```
Sets a path parameter (primarily for testing).

## Middleware Functions

### NewMiddleware
```go
func NewMiddleware(name string, handler StdMiddleware) MiddlewareInterface
```
Creates a named middleware instance.

### NewAnonymousMiddleware
```go
func NewAnonymousMiddleware(handler StdMiddleware) MiddlewareInterface
```
Creates an anonymous middleware instance.

### MiddlewaresToInterfaces
```go
func MiddlewaresToInterfaces(middlewares []StdMiddleware) []MiddlewareInterface
```
Converts a slice of standard middleware to middleware interfaces.

### InterfacesToMiddlewares
```go
func InterfacesToMiddlewares(interfaces []MiddlewareInterface) []StdMiddleware
```
Converts a slice of middleware interfaces to standard middleware.

## Configuration Types

### RouterConfig
```go
type RouterConfig struct {
    Name              string
    BeforeMiddleware  []MiddlewareConfig
    AfterMiddleware   []MiddlewareConfig
    Routes            []RouteConfig
    Groups            []GroupConfig
    Domains           []DomainConfig
    Metadata          map[string]interface{}
}
```

### RouteConfig
```go
type RouteConfig struct {
    Method             string
    Path               string
    Name               string
    Handler            StdHandler
    StringHandler      StringHandler
    HTMLHandler        HTMLHandler
    JSONHandler        JSONHandler
    CSSHandler         CSSHandler
    XMLHandler         XMLHandler
    TextHandler        TextHandler
    JSHandler          JSHandler
    ErrorHandler       ErrorHandler
    BeforeMiddleware   []MiddlewareConfig
    AfterMiddleware    []MiddlewareConfig
    Metadata           map[string]interface{}
}
```

### GroupConfig
```go
type GroupConfig struct {
    Prefix             string
    Name               string
    Routes             []RouteConfig
    Groups             []GroupConfig
    BeforeMiddleware   []MiddlewareConfig
    AfterMiddleware    []MiddlewareConfig
    Metadata           map[string]interface{}
}
```

### DomainConfig
```go
type DomainConfig struct {
    Patterns           []string
    Name               string
    Routes             []RouteConfig
    Groups             []GroupConfig
    BeforeMiddleware   []MiddlewareConfig
    AfterMiddleware    []MiddlewareConfig
    Metadata           map[string]interface{}
}
```

### MiddlewareConfig
```go
type MiddlewareConfig struct {
    Name               string
    Handler            StdMiddleware
    Metadata           map[string]interface{}
}
```

## Configuration Functions

### NewMiddlewareConfig
```go
func NewMiddlewareConfig(name string, handler StdMiddleware) MiddlewareConfig
```
Creates a middleware configuration.

### MiddlewareConfigsToInterfaces
```go
func MiddlewareConfigsToInterfaces(configs []MiddlewareConfig) []MiddlewareInterface
```
Converts middleware configurations to middleware interfaces.

### InterfacesToMiddlewareConfigs
```go
func InterfacesToMiddlewareConfigs(interfaces []MiddlewareInterface) []MiddlewareConfig
```
Converts middleware interfaces to configurations.

### StdMiddlewaresToConfigs
```go
func StdMiddlewaresToConfigs(middlewares []StdMiddleware) []MiddlewareConfig
```
Converts standard middleware to configurations.

## Response Helper Functions

### JSONResponse
```go
func JSONResponse(w http.ResponseWriter, r *http.Request, body string)
```
Writes a JSON response with appropriate headers.

### HTMLResponse
```go
func HTMLResponse(w http.ResponseWriter, r *http.Request, body string)
```
Writes an HTML response with appropriate headers.

### CSSResponse
```go
func CSSResponse(w http.ResponseWriter, r *http.Request, body string)
```
Writes a CSS response with appropriate headers.

### XMLResponse
```go
func XMLResponse(w http.ResponseWriter, r *http.Request, body string)
```
Writes an XML response with appropriate headers.

### TextResponse
```go
func TextResponse(w http.ResponseWriter, r *http.Request, body string)
```
Writes a text response with appropriate headers.

### JSResponse
```go
func JSResponse(w http.ResponseWriter, r *http.Request, body string)
```
Writes a JavaScript response with appropriate headers.

## Declarative Configuration Helpers

### HTTP Method Helpers
```go
func GET(path string, handler StdHandler) RouteConfig
func POST(path string, handler StdHandler) RouteConfig
func PUT(path string, handler StdHandler) RouteConfig
func DELETE(path string, handler StdHandler) RouteConfig
func PATCH(path string, handler StdHandler) RouteConfig
func OPTIONS(path string, handler StdHandler) RouteConfig
func HEAD(path string, handler StdHandler) RouteConfig
```

### Specialized Handler Helpers
```go
func GET_HTML(path string, handler HTMLHandler) RouteConfig
func POST_HTML(path string, handler HTMLHandler) RouteConfig
func GET_JSON(path string, handler JSONHandler) RouteConfig
func POST_JSON(path string, handler JSONHandler) RouteConfig
func GET_CSS(path string, handler CSSHandler) RouteConfig
func GET_XML(path string, handler XMLHandler) RouteConfig
func GET_TEXT(path string, handler TextHandler) RouteConfig
func GET_JS(path string, handler JSHandler) RouteConfig
```

### Group Helper
```go
func Group(prefix string, routes ...RouteConfig) GroupConfig
```
Creates a group configuration with the given prefix and routes.

### Route Configuration Chaining

Route configurations support method chaining for fluent API:

```go
route := rtr.GET("/users", handler).
    WithName("List Users").
    WithBeforeMiddleware(rtr.NewAnonymousMiddleware(authMiddleware)).
    WithMetadata("version", "1.0").
    WithMetadata("protected", "true")
```

#### Available Chain Methods

```go
func (rc RouteConfig) WithName(name string) RouteConfig
func (rc RouteConfig) WithBeforeMiddleware(middleware ...MiddlewareInterface) RouteConfig
func (rc RouteConfig) WithAfterMiddleware(middleware ...MiddlewareInterface) RouteConfig
func (rc RouteConfig) WithMetadata(key string, value interface{}) RouteConfig
```

## Utility Functions

### List
```go
func List(router RouterInterface)
```
Displays a formatted list of all routes, groups, domains, and middleware.

### String
```go
func String(router RouterInterface) string
```
Returns a string representation of the router configuration.

## Constants

### HTTP Methods
```go
const MethodGet     = "GET"
const MethodPost    = "POST"
const MethodPut     = "PUT"
const MethodDelete  = "DELETE"
const MethodPatch   = "PATCH"
const MethodOptions = "OPTIONS"
const MethodHead    = "HEAD"
```

### Content Types
```go
const ContentTypeJSON = "application/json"
const ContentTypeHTML = "text/html; charset=utf-8"
const ContentTypeCSS  = "text/css"
const ContentTypeXML  = "application/xml"
const ContentTypeText = "text/plain; charset=utf-8"
const ContentTypeJS   = "application/javascript"
```

## Error Types

RTR uses standard Go error types. Common errors include:

- Route not found (404)
- Method not allowed (405)
- Parameter extraction errors
- Middleware execution errors

## Examples

### Basic Router Setup
```go
router := rtr.NewRouter()
router.AddRoute(rtr.Get("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, World!"))
}))
```

### Advanced Configuration
```go
config := rtr.RouterConfig{
    Name: "My API",
    BeforeMiddleware: []rtr.MiddlewareConfig{
        rtr.NewMiddlewareConfig("Recovery", middlewares.RecoveryMiddleware),
    },
    Routes: []rtr.RouteConfig{
        rtr.GET("/", homeHandler).WithName("Home"),
        rtr.GET_JSON("/api/status", statusHandler).WithName("Status"),
    },
    Groups: []rtr.GroupConfig{
        rtr.Group("/api/v1",
            rtr.GET("/users", usersHandler).WithName("List Users"),
            rtr.POST("/users", createUserHandler).WithName("Create User"),
        ).WithName("API Group"),
    },
}

router := rtr.NewRouterFromConfig(config)
```

### Middleware Usage
```go
// Standard middleware
loggingMiddleware := func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("Request: %s %s\n", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

// Named middleware
namedMiddleware := rtr.NewMiddleware("Logger", loggingMiddleware)

// Add to router
router.AddBeforeMiddlewares([]rtr.MiddlewareInterface{namedMiddleware})
```

## See Also

- [Getting Started Guide](getting_started.md) - Learn how to use RTR
- [Architecture Documentation](architecture.md) - Understand the system design
- [Middleware Guide](modules/middleware.md) - Middleware system details
- [Configuration Guide](configuration.md) - Configuration options and patterns
