package rtr

import "net/http"

// RouteImpl implements the RouteInterface
// It represents a single route definition with its associated properties and middleware.
// A route defines how a specific HTTP request should be handled, including the HTTP method,
// URL path, handler function, and any middleware that should be applied before or after the handler.
type routeImpl struct {
	// method specifies the HTTP method for this route (e.g., "GET", "POST", "PUT", "DELETE")
	method string

	// path specifies the URL path pattern for this route (e.g., "/users", "/api/products")
	path string

	// handler is the function that will be called when this route is matched
	handler Handler

	// name is an optional identifier for this route, useful for route generation and debugging
	name string

	// beforeMiddlewares are middleware functions that will be executed before the route handler
	beforeMiddlewares []Middleware

	// afterMiddlewares are middleware functions that will be executed after the route handler
	afterMiddlewares []Middleware
}

// GetMethod returns the HTTP method associated with this route.
// Returns the string representation of the HTTP method (e.g., "GET", "POST").
func (r *routeImpl) GetMethod() string {
	return r.method
}

// SetMethod sets the HTTP method for this route.
// This method supports method chaining by returning the RouteInterface.
// The method parameter should be a valid HTTP method string (e.g., "GET", "POST").
func (r *routeImpl) SetMethod(method string) RouteInterface {
	r.method = method
	return r
}

// GetPath returns the URL path pattern associated with this route.
// Returns the string representation of the path (e.g., "/users", "/api/products").
func (r *routeImpl) GetPath() string {
	return r.path
}

// SetPath sets the URL path pattern for this route.
// This method supports method chaining by returning the RouteInterface.
// The path parameter should be a valid URL path pattern.
func (r *routeImpl) SetPath(path string) RouteInterface {
	r.path = path
	return r
}

// GetHandler returns the handler function associated with this route.
// Returns the Handler function that will be called when this route is matched.
func (r *routeImpl) GetHandler() Handler {
	return r.handler
}

// SetHandler sets the handler function for this route.
// This method supports method chaining by returning the RouteInterface.
// The handler parameter should be a function that implements the Handler interface.
func (r *routeImpl) SetHandler(handler Handler) RouteInterface {
	r.handler = handler
	return r
}

// GetName returns the name identifier associated with this route.
// Returns the string name of the route, which may be empty if not set.
func (r *routeImpl) GetName() string {
	return r.name
}

// SetName sets the name identifier for this route.
// This method supports method chaining by returning the RouteInterface.
// The name parameter can be used for route identification and debugging.
func (r *routeImpl) SetName(name string) RouteInterface {
	r.name = name
	return r
}

// AddBeforeMiddlewares adds middleware functions to be executed before the route handler.
// This method supports method chaining by returning the RouteInterface.
// The middleware parameter should be a slice of Middleware functions.
// These middleware functions will be executed in the order they are added.
func (r *routeImpl) AddBeforeMiddlewares(middleware []Middleware) RouteInterface {
	r.beforeMiddlewares = append(r.beforeMiddlewares, middleware...)
	return r
}

// GetBeforeMiddlewares returns all middleware functions that will be executed before the route handler.
// Returns a slice of Middleware functions in the order they will be executed.
func (r *routeImpl) GetBeforeMiddlewares() []Middleware {
	return r.beforeMiddlewares
}

// AddAfterMiddlewares adds middleware functions to be executed after the route handler.
// This method supports method chaining by returning the RouteInterface.
// The middleware parameter should be a slice of Middleware functions.
// These middleware functions will be executed in the order they are added.
func (r *routeImpl) AddAfterMiddlewares(middleware []Middleware) RouteInterface {
	r.afterMiddlewares = append(r.afterMiddlewares, middleware...)
	return r
}

// GetAfterMiddlewares returns all middleware functions that will be executed after the route handler.
// Returns a slice of Middleware functions in the order they will be executed.
func (r *routeImpl) GetAfterMiddlewares() []Middleware {
	return r.afterMiddlewares
}

// Get creates a new GET route with the given path and handler
// It is a shortcut method that combines setting the method to GET, path, and handler.
func Get(path string, handler Handler) RouteInterface {
	return NewRoute().SetMethod(http.MethodGet).SetPath(path).SetHandler(handler)
}

// Post creates a new POST route with the given path and handler
// It is a shortcut method that combines setting the method to POST, path, and handler.
func Post(path string, handler Handler) RouteInterface {
	return NewRoute().SetMethod(http.MethodPost).SetPath(path).SetHandler(handler)
}

// Put creates a new PUT route with the given path and handler
// It is a shortcut method that combines setting the method to PUT, path, and handler.
func Put(path string, handler Handler) RouteInterface {
	return NewRoute().SetMethod(http.MethodPut).SetPath(path).SetHandler(handler)
}

// Delete creates a new DELETE route with the given path and handler
// It is a shortcut method that combines setting the method to DELETE, path, and handler.
func Delete(path string, handler Handler) RouteInterface {
	return NewRoute().SetMethod(http.MethodDelete).SetPath(path).SetHandler(handler)
}
