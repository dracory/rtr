package rtr

import (
	"context"
	"net/http"
	"strings"
)

// Handler represents the function that handles a request.
// It is a function type that takes an http.ResponseWriter and an *http.Request as parameters.
// This is the standard Go HTTP handler function signature.
type Handler func(http.ResponseWriter, *http.Request)

// Middleware represents a middleware function.
// It is a function type that takes an http.Handler and returns an http.Handler.
// Middleware functions can be used to process requests before or after they reach the main handler.
type Middleware func(http.Handler) http.Handler

// NewRouter creates and returns a new RouterInterface implementation.
// This is the main entry point for creating a new router.
// By default, it includes recovery middleware to handle panics.
func NewRouter() RouterInterface {
	r := &routerImpl{
		routes:  make([]RouteInterface, 0),
		groups:  make([]GroupInterface, 0),
		domains: make([]DomainInterface, 0),
		prefix:  "",
	}

	// Add recovery middleware by default
	r.AddBeforeMiddlewares(DefaultMiddlewares())
	return r
}

// This is used to create a new route that can be added to a router or group.
func NewRoute() RouteInterface {
	return &routeImpl{}
}

// NewGroup creates and returns a new GroupInterface implementation.
// This is used to create a new route group that can be added to a router.
func NewGroup() GroupInterface {
	return &groupImpl{}
}

// routerImpl implements the RouterInterface
// It represents a router that can handle HTTP requests by matching them to the appropriate route handler.
type routerImpl struct {
	prefix            string
	routes            []RouteInterface
	groups            []GroupInterface
	domains           []DomainInterface
	beforeMiddlewares []Middleware
	// afterMiddlewares are middleware functions that will be executed after any route handler
	afterMiddlewares []Middleware
}

// GetPrefix returns the URL path prefix associated with this router.
// Returns the string representation of the prefix.
func (r *routerImpl) GetPrefix() string {
	return r.prefix
}

// SetPrefix sets the URL path prefix for this router and returns the router for method chaining.
// The prefix will be prepended to all routes in this router.
func (r *routerImpl) SetPrefix(prefix string) RouterInterface {
	r.prefix = prefix
	return r
}

// AddGroup adds a single group to this router and returns the router for method chaining.
// The group's prefix will be combined with the router's prefix for all routes in the group.
func (r *routerImpl) AddGroup(group GroupInterface) RouterInterface {
	r.groups = append(r.groups, group)
	return r
}

// AddGroups adds multiple groups to this router and returns the router for method chaining.
// Each group's prefix will be combined with the router's prefix for all routes in the group.
func (r *routerImpl) AddGroups(groups []GroupInterface) RouterInterface {
	r.groups = append(r.groups, groups...)
	return r
}

// GetGroups returns all groups that belong to this router.
// Returns a slice of GroupInterface implementations.
func (r *routerImpl) GetGroups() []GroupInterface {
	return r.groups
}

// AddRoute adds a single route to this router and returns the router for method chaining.
// The route's path will be prefixed with the router's prefix.
func (r *routerImpl) AddRoute(route RouteInterface) RouterInterface {
	r.routes = append(r.routes, route)
	return r
}

// AddRoutes adds multiple routes to this router and returns the router for method chaining.
// Each route's path will be prefixed with the router's prefix.
func (r *routerImpl) AddRoutes(routes []RouteInterface) RouterInterface {
	r.routes = append(r.routes, routes...)
	return r
}

// GetRoutes returns all routes that belong to this router.
// Returns a slice of RouteInterface implementations.
func (r *routerImpl) GetRoutes() []RouteInterface {
	return r.routes
}

// AddBeforeMiddlewares adds middleware functions to be executed before any route handler.
// The middleware functions will be executed in the order they are added.
// Returns the router for method chaining.
func (r *routerImpl) AddBeforeMiddlewares(middleware []Middleware) RouterInterface {
	r.beforeMiddlewares = append(r.beforeMiddlewares, middleware...)
	return r
}

// GetBeforeMiddlewares returns all middleware functions that will be executed before any route handler.
// Returns a slice of Middleware functions.
func (r *routerImpl) GetBeforeMiddlewares() []Middleware {
	return r.beforeMiddlewares
}

// AddAfterMiddlewares adds middleware functions to be executed after any route handler.
// The middleware functions will be executed in reverse order of how they were added.
// Returns the router for method chaining.
func (r *routerImpl) AddAfterMiddlewares(middleware []Middleware) RouterInterface {
	r.afterMiddlewares = append(r.afterMiddlewares, middleware...)
	return r
}

// GetAfterMiddlewares returns all middleware functions that will be executed after any route handler.
// Returns a slice of Middleware functions.
func (r *routerImpl) GetAfterMiddlewares() []Middleware {
	return r.afterMiddlewares
}

func (r *routerImpl) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Create a handler chain by wrapping the final handler with middlewares
	var matchedHandler http.Handler

	// First, check if the request matches any domain
	if domain := r.findMatchingDomain(req.Host); domain != nil {
		// Try to find a matching route within the domain
		if _, handler := r.findMatchingRouteInDomain(domain, req); handler != nil {
			matchedHandler = handler
		}
	}

	// If no domain matched or no route in domain matched, check global routes
	if matchedHandler == nil {
		if _, handler := r.findMatchingRoute(req); handler != nil {
			matchedHandler = handler
		}
	}

	// If still no handler matched, return 404
	if matchedHandler == nil {
		http.NotFound(w, req)
		return
	}

	// Execute the handler chain
	matchedHandler.ServeHTTP(w, req)
}

// matchParameterizedRoute checks if a parameterized route matches the request path and extracts parameters
func matchParameterizedRoute(routePath, requestPath string, paramNames []string) (bool, map[string]string) {
	routeSegments := strings.Split(routePath, "/")
	requestSegments := strings.Split(requestPath, "/")
	hasMoreSegments := len(requestSegments) > len(routeSegments)
	hasLessSegments := len(requestSegments) < len(routeSegments)

	// If the request has more segments than the route, it can't match
	if hasMoreSegments {
		return false, nil
	}

	// If the request has fewer segments, check if the remaining segments are optional
	if hasLessSegments {
		// Check if all remaining segments are optional parameters
		for i := len(requestSegments); i < len(routeSegments); i++ {
			seg := routeSegments[i]
			if len(seg) == 0 || seg[0] != ':' || !strings.HasSuffix(seg, "?") {
				return false, nil
			}
		}
	}

	params := make(map[string]string)
	paramIndex := 0

	// Iterate over request segments using range
	// Any remaining route segments must be optional (checked above)
	for i, reqSeg := range requestSegments {
		routeSeg := routeSegments[i]

		// Handle parameter segments (starting with ':')
		if len(routeSeg) > 0 && routeSeg[0] == ':' {
			// Get the parameter name and check if it's optional
			paramName := strings.TrimLeft(routeSeg, ":")
			isOptional := strings.HasSuffix(paramName, "?")

			// If the segment is empty and the parameter is optional, skip it
			if reqSeg == "" && isOptional {
				continue
			}

			// Clean up the parameter name if it was optional
			if isOptional {
				paramName = strings.TrimSuffix(paramName, "?")
			}

			// Store the parameter value with its name
			params[paramName] = reqSeg
			paramIndex++
		} else if routeSeg != reqSeg {
			// If it's not a parameter and segments don't match, the route doesn't match
			return false, nil
		}
	}

	return true, params
}

// findMatchingRoute attempts to find a route that matches the request
// It returns the matched route and an http.Handler that includes all middlewares
func (r *routerImpl) findMatchingRoute(req *http.Request) (RouteInterface, http.Handler) {
	// Check direct routes on the router
	for _, route := range r.routes {
		if match, params := r.routeMatches(route, req); match {
			// Add params to request context if any
			if len(params) > 0 {
				ctx := context.WithValue(req.Context(), ParamsKey, params)
				req = req.WithContext(ctx)
			}
			return route, r.wrapWithMiddlewares(route, req)
		}
	}

	// Check routes in groups
	for _, group := range r.groups {
		if route, handler := r.findMatchingRouteInGroup(group, req, ""); route != nil {
			return route, handler
		}
	}

	return nil, nil
}

// findMatchingRouteInGroup recursively searches for a matching route in a group and its subgroups
func (r *routerImpl) findMatchingRouteInGroup(group GroupInterface, req *http.Request, parentPath string) (RouteInterface, http.Handler) {
	// Combine parent path with group prefix
	groupPath := parentPath + group.GetPrefix()

	// Check routes in the current group
	for _, route := range group.GetRoutes() {
		// Create a copy of the route with adjusted path
		adjustedRoute := &routeImpl{
			method:            route.GetMethod(),
			path:              groupPath + route.GetPath(),
			handler:           route.GetHandler(),
			name:              route.GetName(),
			beforeMiddlewares: route.GetBeforeMiddlewares(),
			afterMiddlewares:  route.GetAfterMiddlewares(),
			paramNames:        route.(*routeImpl).paramNames,
		}

		if match, params := r.routeMatches(adjustedRoute, req); match {
			// Add params to request context if any
			if len(params) > 0 {
				ctx := context.WithValue(req.Context(), ParamsKey, params)
				req = req.WithContext(ctx)
			}
			// Create a handler chain with group middlewares and route middlewares
			return route, r.wrapWithGroupMiddlewares(route, group, req, parentPath)
		}
	}

	// Check subgroups
	for _, subgroup := range group.GetGroups() {
		if route, handler := r.findMatchingRouteInGroup(subgroup, req, groupPath); route != nil {
			return route, handler
		}
	}

	return nil, nil
}

// routeMatches checks if a route matches the request method and path
func (r *routerImpl) routeMatches(route RouteInterface, req *http.Request) (bool, map[string]string) {
	// Check if method matches
	if route.GetMethod() != req.Method && route.GetMethod() != "" {
		return false, nil
	}

	routePath := r.prefix + route.GetPath()
	requestPath := req.URL.Path

	// Handle catch-all routes
	if routePath == "/*" || routePath == "/**" {
		return true, nil
	}

	// Handle wildcard patterns at the end of the path
	if len(routePath) > 2 && routePath[len(routePath)-2:] == "/*" {
		// Check if the base path matches
		basePath := routePath[:len(routePath)-2]
		if len(requestPath) >= len(basePath) && requestPath[:len(basePath)] == basePath {
			return true, nil
		}
		return false, nil
	}

	// If no parameters, do exact matching
	if len(route.(*routeImpl).paramNames) == 0 {
		return routePath == requestPath, nil
	}

	// Handle parameterized routes
	return matchParameterizedRoute(routePath, requestPath, route.(*routeImpl).paramNames)
}

// wrapWithMiddlewares wraps a route's handler with its middlewares and the router's middlewares
func (r *routerImpl) wrapWithMiddlewares(route RouteInterface, req *http.Request) http.Handler {
	// Start with the route's handler
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Use the request from the parameter to preserve the context with parameters
		route.GetHandler()(w, req)
	})

	// Apply route's after middlewares (in reverse order)
	for i := len(route.GetAfterMiddlewares()) - 1; i >= 0; i-- {
		handler = route.GetAfterMiddlewares()[i](handler)
	}

	// Apply router's after middlewares (in reverse order)
	for i := len(r.afterMiddlewares) - 1; i >= 0; i-- {
		handler = r.afterMiddlewares[i](handler)
	}

	// Apply route's before middlewares
	for _, middleware := range route.GetBeforeMiddlewares() {
		handler = middleware(handler)
	}

	// Apply router's before middlewares
	for _, middleware := range r.beforeMiddlewares {
		handler = middleware(handler)
	}

	return handler
}

// wrapWithGroupMiddlewares wraps a route's handler with its middlewares, the group's middlewares, and the router's middlewares
func (r *routerImpl) wrapWithGroupMiddlewares(route RouteInterface, group GroupInterface, req *http.Request, parentPath string) http.Handler {
	// Start with the route's handler
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route.GetHandler()(w, r)
	})

	// Apply route's after middlewares (in reverse order)
	for i := len(route.GetAfterMiddlewares()) - 1; i >= 0; i-- {
		handler = route.GetAfterMiddlewares()[i](handler)
	}

	// Apply group's after middlewares (in reverse order)
	for i := len(group.GetAfterMiddlewares()) - 1; i >= 0; i-- {
		handler = group.GetAfterMiddlewares()[i](handler)
	}

	// Apply router's after middlewares (in reverse order)
	for i := len(r.afterMiddlewares) - 1; i >= 0; i-- {
		handler = r.afterMiddlewares[i](handler)
	}

	// Apply route's before middlewares
	for _, middleware := range route.GetBeforeMiddlewares() {
		handler = middleware(handler)
	}

	// Apply group's before middlewares
	for _, middleware := range group.GetBeforeMiddlewares() {
		handler = middleware(handler)
	}

	// Apply router's before middlewares
	for _, middleware := range r.beforeMiddlewares {
		handler = middleware(handler)
	}

	return handler
}
