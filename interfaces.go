package rtr

import "net/http"

// Handler defines the function signature for HTTP request handlers.
type Handler func(http.ResponseWriter, *http.Request)

// ErrorHandler is a convenience shorthand handler for error responses.
// Returns an error that will be handled appropriately.
// If the returned error is nil, it means no error occurred and the response is successful.
// Content-Type and status codes are left to specific extensions like 404ErrorHandler.
type ErrorHandler func(http.ResponseWriter, *http.Request) error

// StringHandler is a convenience shorthand handler for simple string responses.
// It returns a string that will be written directly to the response without setting any headers.
// The handler is responsible for setting any headers it needs.
type StringHandler func(http.ResponseWriter, *http.Request) string

// HTMLHandler is a convenience shorthand handler that returns HTML content.
// Automatically sets Content-Type: "text/html; charset=utf-8" header.
// Returns an HTML string that will be wrapped with HTMLResponse().
type HTMLHandler StringHandler

// JSONHandler is a convenience shorthand handler that returns JSON content.
// Automatically sets Content-Type: "application/json" header.
// Returns a JSON string that will be wrapped with JSONResponse().
type JSONHandler StringHandler

// CSSHandler is a convenience shorthand handler that returns CSS content.
// Automatically sets Content-Type: "text/css" header.
// Returns a CSS string that will be wrapped with CSSResponse().
type CSSHandler StringHandler

// XMLHandler is a convenience shorthand handler that returns XML content.
// Automatically sets Content-Type: "application/xml" header.
// Returns an XML string that will be wrapped with XMLResponse().
type XMLHandler StringHandler

// TextHandler is a convenience shorthand handler that returns plain text content.
// Automatically sets Content-Type: "text/plain; charset=utf-8" header.
// Returns a plain text string that will be wrapped with TextResponse().
type TextHandler StringHandler

// JSHandler is a convenience shorthand handler that returns JavaScript content.
// Automatically sets Content-Type: "application/javascript" header.
// Returns a JavaScript string that will be wrapped with JSResponse().
type JSHandler StringHandler

// Middleware represents a middleware function.
// It is a function type that takes an http.Handler and returns an http.Handler.
// Middleware functions can be used to process requests before or after they reach the main handler.
type Middleware func(http.Handler) http.Handler

// RouteInterface defines the interface for a single route definition.
// A route represents a mapping between an HTTP method, a URL path pattern, and a handler function.
// Routes can also have associated middleware that will be executed before or after the handler.
type RouteInterface interface {
	// GetMethod returns the HTTP method associated with this route.
	GetMethod() string
	// SetMethod sets the HTTP method for this route and returns the route for method chaining.
	SetMethod(method string) RouteInterface

	// GetPath returns the URL path pattern associated with this route.
	GetPath() string
	// SetPath sets the URL path pattern for this route and returns the route for method chaining.
	SetPath(path string) RouteInterface

	// GetHandler returns the handler function associated with this route.
	GetHandler() Handler
	// SetHandler sets the handler function for this route and returns the route for method chaining.
	SetHandler(handler Handler) RouteInterface

	// GetStringHandler returns the string handler function associated with this route.
	GetStringHandler() StringHandler
	// SetStringHandler sets the string handler function for this route and returns the route for method chaining.
	SetStringHandler(handler StringHandler) RouteInterface

	// GetHTMLHandler returns the HTML handler function associated with this route.
	GetHTMLHandler() HTMLHandler
	// SetHTMLHandler sets the HTML handler function for this route and returns the route for method chaining.
	SetHTMLHandler(handler HTMLHandler) RouteInterface

	// GetJSONHandler returns the JSON handler function associated with this route.
	GetJSONHandler() JSONHandler
	// SetJSONHandler sets the JSON handler function for this route and returns the route for method chaining.
	SetJSONHandler(handler JSONHandler) RouteInterface

	// GetCSSHandler returns the CSS handler function associated with this route.
	GetCSSHandler() CSSHandler
	// SetCSSHandler sets the CSS handler function for this route and returns the route for method chaining.
	SetCSSHandler(handler CSSHandler) RouteInterface

	// GetXMLHandler returns the XML handler function associated with this route.
	GetXMLHandler() XMLHandler
	// SetXMLHandler sets the XML handler function for this route and returns the route for method chaining.
	SetXMLHandler(handler XMLHandler) RouteInterface

	// GetTextHandler returns the text handler function associated with this route.
	GetTextHandler() TextHandler
	// SetTextHandler sets the text handler function for this route and returns the route for method chaining.
	SetTextHandler(handler TextHandler) RouteInterface

	// GetJSHandler returns the JavaScript handler function associated with this route.
	GetJSHandler() JSHandler
	// SetJSHandler sets the JavaScript handler function for this route and returns the route for method chaining.
	SetJSHandler(handler JSHandler) RouteInterface

	// GetErrorHandler returns the error handler function associated with this route.
	GetErrorHandler() ErrorHandler
	// SetErrorHandler sets the error handler function for this route and returns the route for method chaining.
	SetErrorHandler(handler ErrorHandler) RouteInterface

	// GetName returns the name identifier associated with this route.
	GetName() string
	// SetName sets the name identifier for this route and returns the route for method chaining.
	SetName(name string) RouteInterface

	// AddBeforeMiddlewares adds middleware functions to be executed before the route handler.
	// Returns the route for method chaining.
	AddBeforeMiddlewares(middleware []Middleware) RouteInterface
	// GetBeforeMiddlewares returns all middleware functions that will be executed before the route handler.
	GetBeforeMiddlewares() []Middleware

	// AddAfterMiddlewares adds middleware functions to be executed after the route handler.
	// Returns the route for method chaining.
	AddAfterMiddlewares(middleware []Middleware) RouteInterface
	// GetAfterMiddlewares returns all middleware functions that will be executed after the route handler.
	GetAfterMiddlewares() []Middleware
}

// GroupInterface defines the interface for a group of routes.
// A group represents a collection of routes that share common properties such as a URL prefix and middleware.
// Groups can also be nested to create hierarchical route structures.
type GroupInterface interface {
	// GetPrefix returns the URL path prefix associated with this group.
	GetPrefix() string
	// SetPrefix sets the URL path prefix for this group and returns the group for method chaining.
	SetPrefix(prefix string) GroupInterface

	// AddRoute adds a single route to this group and returns the group for method chaining.
	AddRoute(route RouteInterface) GroupInterface
	// AddRoutes adds multiple routes to this group and returns the group for method chaining.
	AddRoutes(routes []RouteInterface) GroupInterface
	// GetRoutes returns all routes that belong to this group.
	GetRoutes() []RouteInterface

	// AddGroup adds a single nested group to this group and returns the group for method chaining.
	AddGroup(group GroupInterface) GroupInterface
	// AddGroups adds multiple nested groups to this group and returns the group for method chaining.
	AddGroups(groups []GroupInterface) GroupInterface
	// GetGroups returns all nested groups that belong to this group.
	GetGroups() []GroupInterface

	// AddBeforeMiddlewares adds middleware functions to be executed before any route handler in this group.
	// Returns the group for method chaining.
	AddBeforeMiddlewares(middleware []Middleware) GroupInterface
	// GetBeforeMiddlewares returns all middleware functions that will be executed before any route handler in this group.
	GetBeforeMiddlewares() []Middleware

	// AddAfterMiddlewares adds middleware functions to be executed after any route handler in this group.
	// Returns the group for method chaining.
	AddAfterMiddlewares(middleware []Middleware) GroupInterface
	// GetAfterMiddlewares returns all middleware functions that will be executed after any route handler in this group.
	GetAfterMiddlewares() []Middleware
}

// DomainInterface defines the interface for a domain that can have routes and groups.
// A domain represents a collection of routes and groups that are only accessible
// when the request's Host header matches the domain's patterns.
type DomainInterface interface {
	// GetPatterns returns the domain patterns that this domain matches against
	GetPatterns() []string

	// SetPatterns sets the domain patterns for this domain and returns the domain for method chaining
	SetPatterns(patterns ...string) DomainInterface

	// AddRoute adds a route to this domain and returns the domain for method chaining
	AddRoute(route RouteInterface) DomainInterface

	// AddRoutes adds multiple routes to this domain and returns the domain for method chaining
	AddRoutes(routes []RouteInterface) DomainInterface

	// GetRoutes returns all routes that belong to this domain
	GetRoutes() []RouteInterface

	// AddGroup adds a group to this domain and returns the domain for method chaining
	AddGroup(group GroupInterface) DomainInterface

	// AddGroups adds multiple groups to this domain and returns the domain for method chaining
	AddGroups(groups []GroupInterface) DomainInterface

	// GetGroups returns all groups that belong to this domain
	GetGroups() []GroupInterface

	// AddBeforeMiddlewares adds middleware functions to be executed before any route handler in this domain
	// Returns the domain for method chaining
	AddBeforeMiddlewares(middleware []Middleware) DomainInterface

	// GetBeforeMiddlewares returns all middleware functions that will be executed before any route handler in this domain
	GetBeforeMiddlewares() []Middleware

	// AddAfterMiddlewares adds middleware functions to be executed after any route handler in this domain
	// Returns the domain for method chaining
	AddAfterMiddlewares(middleware []Middleware) DomainInterface

	// GetAfterMiddlewares returns all middleware functions that will be executed after any route handler in this domain
	GetAfterMiddlewares() []Middleware

	// Match checks if the given host matches any of this domain's patterns
	Match(host string) bool
}

// RouterInterface defines the interface for a router that can handle HTTP requests.
// A router is responsible for matching incoming HTTP requests to the appropriate route handler
// and executing any associated middleware.
type RouterInterface interface {
	// GetPrefix returns the URL path prefix associated with this router.
	GetPrefix() string
	// SetPrefix sets the URL path prefix for this router and returns the router for method chaining.
	// The prefix will be prepended to all routes in this router.
	SetPrefix(prefix string) RouterInterface

	// AddGroup adds a single group to this router and returns the router for method chaining.
	// The group's prefix will be combined with the router's prefix for all routes in the group.
	AddGroup(group GroupInterface) RouterInterface
	// AddGroups adds multiple groups to this router and returns the router for method chaining.
	// Each group's prefix will be combined with the router's prefix for all routes in the group.
	AddGroups(groups []GroupInterface) RouterInterface
	// GetGroups returns all groups that belong to this router.
	// Returns a slice of GroupInterface implementations.
	GetGroups() []GroupInterface

	// AddRoute adds a single route to this router and returns the router for method chaining.
	// The route's path will be prefixed with the router's prefix.
	AddRoute(route RouteInterface) RouterInterface
	// AddRoutes adds multiple routes to this router and returns the router for method chaining.
	// Each route's path will be prefixed with the router's prefix.
	AddRoutes(routes []RouteInterface) RouterInterface
	// GetRoutes returns all routes that belong to this router.
	// Returns a slice of RouteInterface implementations.
	GetRoutes() []RouteInterface

	// AddBeforeMiddlewares adds middleware functions to be executed before any route handler.
	// The middleware functions will be executed in the order they are added.
	// Returns the router for method chaining.
	AddBeforeMiddlewares(middleware []Middleware) RouterInterface
	// GetBeforeMiddlewares returns all middleware functions that will be executed before any route handler.
	// Returns a slice of Middleware functions.
	GetBeforeMiddlewares() []Middleware

	// AddAfterMiddlewares adds middleware functions to be executed after any route handler.
	// The middleware functions will be executed in reverse order of how they were added.
	// Returns the router for method chaining.
	AddAfterMiddlewares(middleware []Middleware) RouterInterface
	// GetAfterMiddlewares returns all middleware functions that will be executed after any route handler.
	// Returns a slice of Middleware functions.
	GetAfterMiddlewares() []Middleware

	// AddDomain adds a domain to this router and returns the router for method chaining
	AddDomain(domain DomainInterface) RouterInterface

	// AddDomains adds multiple domains to this router and returns the router for method chaining
	AddDomains(domains []DomainInterface) RouterInterface

	// GetDomains returns all domains that belong to this router
	GetDomains() []DomainInterface

	// List displays the router's configuration in formatted tables for debugging and documentation
	// Shows global middleware, domains, direct routes, and route groups
	List()

	// ServeHTTP implements the http.Handler interface.
	// It matches the incoming request to the appropriate route and executes the handler.
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
