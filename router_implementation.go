package rtr

import (
	"context"
	"net/http"
	"strings"
)

// NewRouter creates and returns a new RouterInterface implementation.
// This is the main entry point for creating a new router.
// The router starts with no default middlewares - users should add middlewares as needed.
func NewRouter() RouterInterface {
	r := &routerImpl{
		routes:  make([]RouteInterface, 0),
		groups:  make([]GroupInterface, 0),
		domains: make([]DomainInterface, 0),
		prefix:  "",
	}

	// Router starts with no default middlewares
	// Users can add middlewares as needed using AddBeforeMiddlewares() or AddAfterMiddlewares()
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
	beforeMiddlewares []MiddlewareInterface
	// afterMiddlewares are middleware functions that will be executed after any route handler
	afterMiddlewares []MiddlewareInterface
}

var _ RouterInterface = (*routerImpl)(nil)

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

// AddBeforeMiddlewares adds middleware to be executed before any route handler.
// The middleware will be executed in the order they are added.
// Returns the router for method chaining.
func (r *routerImpl) AddBeforeMiddlewares(middleware []MiddlewareInterface) RouterInterface {
	r.beforeMiddlewares = append(r.beforeMiddlewares, middleware...)
	return r
}

// GetBeforeMiddlewares returns all middleware that will be executed before any route handler.
// Returns a slice of MiddlewareInterface.
func (r *routerImpl) GetBeforeMiddlewares() []MiddlewareInterface {
	return r.beforeMiddlewares
}

// AddAfterMiddlewares adds middleware to be executed after any route handler.
// The middleware will be executed in reverse order of how they were added.
// Returns the router for method chaining.
func (r *routerImpl) AddAfterMiddlewares(middleware []MiddlewareInterface) RouterInterface {
	r.afterMiddlewares = append(r.afterMiddlewares, middleware...)
	return r
}

// GetAfterMiddlewares returns all middleware that will be executed after any route handler.
// Returns a slice of MiddlewareInterface.
func (r *routerImpl) GetAfterMiddlewares() []MiddlewareInterface {
	return r.afterMiddlewares
}

func (r *routerImpl) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Find a matching route
	route, handler := r.findMatchingRoute(req)

	// If no route found, check domains
	if route == nil {
		for _, domain := range r.domains {
			route, handler = r.findMatchingRouteInDomain(domain, req)
			if route != nil {
				break
			}
		}
	}

	// If still no route found, return 404
	if route == nil {
		http.NotFound(w, req)
		return
	}

	// Serve the request
	handler.ServeHTTP(w, req)
}

// matchParameterizedRoute checks if a parameterized route matches the request path and extracts parameters
func matchParameterizedRoute(routePath, requestPath string) (bool, map[string]string) {
	routeSegments := strings.Split(routePath, "/")
	requestSegments := strings.Split(requestPath, "/")

	// Detect if the last route segment is a greedy parameter ':name...'
	greedy := false
	if len(routeSegments) > 0 {
		last := routeSegments[len(routeSegments)-1]
		if len(last) > 0 && last[0] == ':' && strings.HasSuffix(last, "...") {
			greedy = true
		}
	}

	hasMoreSegments := len(requestSegments) > len(routeSegments)
	hasLessSegments := len(requestSegments) < len(routeSegments)

	if greedy {
		// Greedy param must capture at least one segment
		if hasLessSegments {
			return false, nil
		}
	} else {
		// If the request has more segments than the route, it can't match
		if hasMoreSegments {
			return false, nil
		}
		// If the request has fewer segments, check if the remaining segments are optional
		if hasLessSegments {
			for i := len(requestSegments); i < len(routeSegments); i++ {
				seg := routeSegments[i]
				if len(seg) == 0 || seg[0] != ':' || !strings.HasSuffix(seg, "?") {
					return false, nil
				}
			}
		}
	}

	params := make(map[string]string)

	if greedy {
		// Match all but the last (greedy) route segment
		for i := 0; i < len(routeSegments)-1; i++ {
			routeSeg := routeSegments[i]
			reqSeg := requestSegments[i]

			if len(routeSeg) > 0 && routeSeg[0] == ':' {
				// Parameter segment
				paramName := strings.TrimLeft(routeSeg, ":")
				paramName = strings.TrimSuffix(paramName, "?")
				params[paramName] = reqSeg
			} else if routeSeg != reqSeg {
				return false, nil
			}
		}

		// Handle the greedy segment
		last := routeSegments[len(routeSegments)-1]
		paramName := strings.TrimLeft(last, ":")
		paramName = strings.TrimSuffix(paramName, "...")
		remainder := strings.Join(requestSegments[len(routeSegments)-1:], "/")
		if remainder == "" {
			return false, nil
		}
		params[paramName] = remainder
		return true, params
	}

	// Non-greedy: Iterate over request segments; any remaining route segments must be optional
	for i, reqSeg := range requestSegments {
		routeSeg := routeSegments[i]

		if len(routeSeg) > 0 && routeSeg[0] == ':' {
			// Parameter segment
			paramName := strings.TrimLeft(routeSeg, ":")
			paramName = strings.TrimSuffix(paramName, "?")
			params[paramName] = reqSeg
		} else if routeSeg != reqSeg {
			return false, nil
		}
	}

	return true, params
}

// findMatchingRoute attempts to find a route that matches the request
// It returns the matched route and an http.Handler that includes all middlewares
func (r *routerImpl) findMatchingRoute(req *http.Request) (RouteInterface, http.Handler) {
	// Create a copy of the request to avoid mutating the original
	reqCopy := req.Clone(req.Context())

	// Check direct routes on the router
	for _, route := range r.routes {
		// Create a fresh copy of the request for each route check
		rc := reqCopy.Clone(reqCopy.Context())

		if match, params := r.routeMatches(route, rc); match {
			// Add params to request context if any
			if len(params) > 0 {
				ctx := context.WithValue(rc.Context(), ParamsKey, params)
				rc = rc.WithContext(ctx)
			}

			// Create a handler that will use the correct request with parameters
			handler := r.buildHandler(route, nil, nil)

			// Return a handler that will use the request with the updated context
			return route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.ServeHTTP(w, rc)
			})
		}
	}

	// Check routes in groups
	for _, group := range r.groups {
		if route, handler := r.findMatchingRouteInGroup(group, reqCopy, nil); route != nil {
			return route, handler
		}
	}

	return nil, nil
}

// findMatchingRouteInGroup recursively searches for a matching route in a group and its subgroups
func (r *routerImpl) findMatchingRouteInGroup(group GroupInterface, req *http.Request, parentGroups []GroupInterface) (RouteInterface, http.Handler) {
	// Combine parent path with group prefix
	currentGroups := append(parentGroups, group)

	// Check routes in the current group
	for _, route := range group.GetRoutes() {
		// Create a full path for matching
		fullPath := ""
		for _, g := range currentGroups {
			fullPath += g.GetPrefix()
		}
		fullPath += route.GetPath()

		// Create a temporary route for matching
		tempRoute := &routeImpl{
			method:     route.GetMethod(),
			path:       fullPath,
			handler:    route.GetHandler(),
			paramNames: route.(*routeImpl).paramNames,
		}

		// Create a copy of the request for this route check
		reqCopy := req.Clone(req.Context())

		if match, params := r.routeMatches(tempRoute, reqCopy); match {
			// Create a new request with the updated context containing the parameters
			if len(params) > 0 {
				ctx := context.WithValue(reqCopy.Context(), ParamsKey, params)
				reqCopy = reqCopy.WithContext(ctx)
			}

			// Find the domain that matches this request (if any)
			var domain DomainInterface
			host := req.Host
			if host != "" {
				for _, d := range r.domains {
					if d.Match(host) {
						domain = d
						break
					}
				}
			}

			// Create a handler that will use the correct request with parameters
			handler := r.buildHandler(route, currentGroups, domain)

			// Return a handler that will use the request with the updated context
			return route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.ServeHTTP(w, reqCopy)
			})
		}
	}

	// Check subgroups
	for _, subgroup := range group.GetGroups() {
		if route, handler := r.findMatchingRouteInGroup(subgroup, req, currentGroups); route != nil {
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
	return matchParameterizedRoute(routePath, requestPath)
}
