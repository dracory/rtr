package rtr

import "net/http"

// ===========================================================================
// = CONSTRUCTOR OPTIONS
// ===========================================================================

// MiddlewareOption defines a function type that can configure a middleware
type MiddlewareOption func(*middlewareImpl)

// WithName returns a MiddlewareOption that sets the name of the middleware
func WithName(name string) MiddlewareOption {
	return func(m *middlewareImpl) {
		m.name = name
	}
}

// WithHandler returns a MiddlewareOption that sets the handler function of the middleware
func WithHandler(handler StdMiddleware) MiddlewareOption {
	return func(m *middlewareImpl) {
		m.handler = handler
	}
}

// ===========================================================================
// = CONSTTRUCTORS
// ===========================================================================

// NewMiddleware creates a new middleware with the provided options.
// Example: NewMiddleware(WithName("auth"), WithHandler(myAuthHandler))
func NewMiddleware(opts ...MiddlewareOption) MiddlewareInterface {
	m := &middlewareImpl{}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// NewAnonymousMiddleware creates a new middleware with the given handler.
// This is a convenience function that's equivalent to NewMiddleware(WithHandler(handler)).
// It's maintained for backward compatibility with existing code.
func NewAnonymousMiddleware(handler StdMiddleware) MiddlewareInterface {
	return NewMiddleware(WithHandler(handler))
}

// ===========================================================================
// = IMPLEMENTATION
// ===========================================================================

// middlewareImpl implements the MiddlewareInterface
// It represents a named middleware that can be applied to routes, groups, or routers.
// This implementation follows the same pattern as routeImpl in the codebase.
type middlewareImpl struct {
	// name is an optional identifier for this middleware, useful for debugging and documentation
	name string

	// handler is the middleware function that will be executed
	handler StdMiddleware
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
// Returns the StdMiddleware function that will be executed.
func (m *middlewareImpl) GetHandler() StdMiddleware {
	return m.handler
}

// SetHandler sets the middleware function and returns the middleware for method chaining.
// The handler parameter should be a valid StdMiddleware function.
func (m *middlewareImpl) SetHandler(handler StdMiddleware) MiddlewareInterface {
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

// MiddlewareFromFunction converts a StdMiddleware function to a MiddlewareInterface.
// This is a convenience function for backward compatibility.
func MiddlewareFromFunction(handler StdMiddleware) MiddlewareInterface {
	return NewAnonymousMiddleware(handler)
}

// MiddlewaresToInterfaces converts a slice of StdMiddleware functions to MiddlewareInterface slice.
// This is useful for migrating existing code that uses []StdMiddleware to []MiddlewareInterface.
func MiddlewaresToInterfaces(middlewares []StdMiddleware) []MiddlewareInterface {
	if middlewares == nil {
		return nil
	}
	result := make([]MiddlewareInterface, 0, len(middlewares))
	for _, mw := range middlewares {
		if mw != nil {
			result = append(result, NewAnonymousMiddleware(mw))
		}
	}
	return result
}

// InterfacesToMiddlewares converts a slice of MiddlewareInterface to StdMiddleware functions.
// This is useful for backward compatibility when you need to work with the underlying functions.
func InterfacesToMiddlewares(interfaces []MiddlewareInterface) []StdMiddleware {
	var middlewares []StdMiddleware
	for _, mw := range interfaces {
		middlewares = append(middlewares, mw.GetHandler())
	}
	return middlewares
}

// AddMiddlewaresToInterfaces converts and adds StdMiddleware functions to a MiddlewareInterface slice
func AddMiddlewaresToInterfaces(interfaces []MiddlewareInterface, middlewares []StdMiddleware) []MiddlewareInterface {
	return append(interfaces, MiddlewaresToInterfaces(middlewares)...)
}

// ExecuteMiddlewareChain executes a chain of MiddlewareInterface in order.
// This is a helper function that applies all middleware in the slice to the final handler.
func ExecuteMiddlewareChain(middlewares []MiddlewareInterface, finalHandler http.Handler) http.Handler {
	// Apply middlewares in reverse order so they execute in the order they're defined
	// when the request comes in
	handler := finalHandler
	for i := len(middlewares) - 1; i >= 0; i-- {
		if middlewares[i] != nil {
			handler = middlewares[i].Execute(handler)
		}
	}
	return handler
}
