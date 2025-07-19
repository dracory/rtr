package rtr

import "net/http"

// middlewareImpl implements the MiddlewareInterface
// It represents a named middleware that can be applied to routes, groups, or routers.
// This implementation follows the same pattern as routeImpl in the codebase.
type middlewareImpl struct {
	// name is an optional identifier for this middleware, useful for debugging and documentation
	name string

	// handler is the middleware function that will be executed
	handler Middleware
}

var _ MiddlewareInterface = (*middlewareImpl)(nil)

// GetName returns the name identifier associated with this middleware.
// Returns the string name of the middleware, which may be empty for anonymous middleware.
func (m *middlewareImpl) GetName() string {
	return m.name
}

// SetName sets the name identifier for this middleware and returns the middleware for method chaining.
// The name parameter can be used for middleware identification and debugging.
func (m *middlewareImpl) SetName(name string) MiddlewareInterface {
	m.name = name
	return m
}

// GetHandler returns the underlying middleware function.
// Returns the Middleware function that will be executed.
func (m *middlewareImpl) GetHandler() Middleware {
	return m.handler
}

// SetHandler sets the middleware function and returns the middleware for method chaining.
// The handler parameter should be a valid Middleware function.
func (m *middlewareImpl) SetHandler(handler Middleware) MiddlewareInterface {
	m.handler = handler
	return m
}

// Execute applies the middleware to the given handler.
// This is equivalent to calling GetHandler()(next).
func (m *middlewareImpl) Execute(next http.Handler) http.Handler {
	if m.handler == nil {
		return next
	}
	return m.handler(next)
}

// NewMiddleware creates a new named middleware with the given name and handler.
// This is the main factory function for creating named middleware.
func NewMiddleware(name string, handler Middleware) MiddlewareInterface {
	return &middlewareImpl{
		name:    name,
		handler: handler,
	}
}

// NewAnonymousMiddleware creates a new middleware without a name.
// This is useful for backward compatibility with existing code that uses anonymous middleware.
func NewAnonymousMiddleware(handler Middleware) MiddlewareInterface {
	return &middlewareImpl{
		name:    "",
		handler: handler,
	}
}

// MiddlewareFromFunction converts a Middleware function to a MiddlewareInterface.
// This is a convenience function for backward compatibility.
func MiddlewareFromFunction(handler Middleware) MiddlewareInterface {
	return NewAnonymousMiddleware(handler)
}

// MiddlewaresToInterfaces converts a slice of Middleware functions to MiddlewareInterface slice.
// This is useful for migrating existing code that uses []Middleware to []MiddlewareInterface.
func MiddlewaresToInterfaces(middlewares []Middleware) []MiddlewareInterface {
	var interfaces []MiddlewareInterface
	for _, mw := range middlewares {
		interfaces = append(interfaces, NewAnonymousMiddleware(mw))
	}
	return interfaces
}

// InterfacesToMiddlewares converts a slice of MiddlewareInterface to Middleware functions.
// This is useful for backward compatibility when you need to work with the underlying functions.
func InterfacesToMiddlewares(interfaces []MiddlewareInterface) []Middleware {
	var middlewares []Middleware
	for _, mw := range interfaces {
		middlewares = append(middlewares, mw.GetHandler())
	}
	return middlewares
}

// AddMiddlewaresToInterfaces converts and adds Middleware functions to a MiddlewareInterface slice
func AddMiddlewaresToInterfaces(interfaces []MiddlewareInterface, middlewares []Middleware) []MiddlewareInterface {
	return append(interfaces, MiddlewaresToInterfaces(middlewares)...)
}

// ExecuteMiddlewareChain executes a chain of MiddlewareInterface in order.
// This is a helper function that applies all middleware in the slice to the final handler.
func ExecuteMiddlewareChain(middlewares []MiddlewareInterface, finalHandler http.Handler) http.Handler {
	// Start with the final handler
	handler := finalHandler

	// Apply middleware in reverse order (last middleware wraps first)
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i].Execute(handler)
	}

	return handler
}
