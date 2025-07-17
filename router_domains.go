package rtr

import (
	"net/http"
)

// AddDomain adds a domain to this router and returns the router for method chaining
func (r *routerImpl) AddDomain(domain DomainInterface) RouterInterface {
	r.domains = append(r.domains, domain)
	return r
}

// AddDomains adds multiple domains to this router and returns the router for method chaining
func (r *routerImpl) AddDomains(domains []DomainInterface) RouterInterface {
	r.domains = append(r.domains, domains...)
	return r
}

// GetDomains returns all domains that belong to this router
func (r *routerImpl) GetDomains() []DomainInterface {
	return r.domains
}

// findMatchingDomain finds the first domain that matches the given host
func (r *routerImpl) findMatchingDomain(host string) DomainInterface {
	for _, domain := range r.domains {
		if domain.Match(host) {
			return domain
		}
	}
	return nil
}

// findMatchingRouteInDomain finds a route that matches the request within a domain
func (r *routerImpl) findMatchingRouteInDomain(domain DomainInterface, req *http.Request) (RouteInterface, http.Handler) {
	// Check direct routes in domain
	for _, route := range domain.GetRoutes() {
		if r.routeMatches(route, req) {
			return route, r.wrapWithDomainMiddlewares(route, domain, req)
		}
	}

	// Check route groups in domain
	for _, group := range domain.GetGroups() {
		if route, _ := r.findMatchingRouteInGroup(group, req, ""); route != nil {
			return route, r.wrapWithDomainGroupMiddlewares(route, group, domain, req)
		}
	}

	return nil, nil
}

// wrapWithDomainMiddlewares wraps a route's handler with domain and route middlewares
func (r *routerImpl) wrapWithDomainMiddlewares(route RouteInterface, domain DomainInterface, req *http.Request) http.Handler {
	// Start with the route's handler
	var handler http.Handler = http.HandlerFunc(route.GetHandler())

	// Helper function to apply middlewares with proper type conversion
	applyMiddlewares := func(handler http.Handler, middlewares []Middleware) http.Handler {
		for i := range middlewares {
			mw := middlewares[i]
			handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mw(handler).ServeHTTP(w, r)
			})
		}
		return handler
	}

	// Apply route's before middlewares in order
	handler = applyMiddlewares(handler, route.GetBeforeMiddlewares())

	// Apply domain's before middlewares in order
	handler = applyMiddlewares(handler, domain.GetBeforeMiddlewares())

	// Apply router's before middlewares in order
	handler = applyMiddlewares(handler, r.beforeMiddlewares)

	// Create a final handler that wraps the chain with after middlewares
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Call the handler chain
		handler.ServeHTTP(w, req)

		// After handler execution, run after middlewares in reverse order
		var afterHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

		// Helper function to apply middlewares in reverse order
		applyAfterMiddlewares := func(handler http.Handler, middlewares []Middleware) http.Handler {
			for i := len(middlewares) - 1; i >= 0; i-- {
				mw := middlewares[i]
				handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					mw(handler).ServeHTTP(w, r)
				})
			}
			return handler
		}

		// Apply route's after middlewares in reverse order
		afterHandler = applyAfterMiddlewares(afterHandler, route.GetAfterMiddlewares())

		// Apply domain's after middlewares in reverse order
		afterHandler = applyAfterMiddlewares(afterHandler, domain.GetAfterMiddlewares())

		// Apply router's after middlewares in reverse order
		afterHandler = applyAfterMiddlewares(afterHandler, r.afterMiddlewares)

		// Execute the after middlewares
		afterHandler.ServeHTTP(w, req)
	})
}

// wrapWithDomainGroupMiddlewares wraps a route's handler with domain, group, and route middlewares
func (r *routerImpl) wrapWithDomainGroupMiddlewares(route RouteInterface, group GroupInterface, domain DomainInterface, req *http.Request) http.Handler {
	// Start with a handler that calls the route's handler function
	var handler http.Handler = http.HandlerFunc(route.GetHandler())

	// Helper function to apply middlewares with proper type conversion
	applyMiddlewares := func(handler http.Handler, middlewares []Middleware) http.Handler {
		for i := range middlewares {
			mw := middlewares[i]
			handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mw(handler).ServeHTTP(w, r)
			})
		}
		return handler
	}

	// Apply route's before middlewares in order
	handler = applyMiddlewares(handler, route.GetBeforeMiddlewares())

	// Apply group's before middlewares in order
	handler = applyMiddlewares(handler, group.GetBeforeMiddlewares())

	// Apply domain's before middlewares in order
	handler = applyMiddlewares(handler, domain.GetBeforeMiddlewares())

	// Apply router's before middlewares in order
	handler = applyMiddlewares(handler, r.beforeMiddlewares)

	// Create a final handler that wraps the chain with after middlewares
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Call the handler chain
		handler.ServeHTTP(w, req)

		// After handler execution, run after middlewares in reverse order
		var afterHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

		// Helper function to apply middlewares in reverse order
		applyAfterMiddlewares := func(handler http.Handler, middlewares []Middleware) http.Handler {
			for i := len(middlewares) - 1; i >= 0; i-- {
				mw := middlewares[i]
				handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					mw(handler).ServeHTTP(w, r)
				})
			}
			return handler
		}

		// Apply route's after middlewares in reverse order
		afterHandler = applyAfterMiddlewares(afterHandler, route.GetAfterMiddlewares())

		// Apply group's after middlewares in reverse order
		afterHandler = applyAfterMiddlewares(afterHandler, group.GetAfterMiddlewares())

		// Apply domain's after middlewares in reverse order
		afterHandler = applyAfterMiddlewares(afterHandler, domain.GetAfterMiddlewares())

		// Apply router's after middlewares in reverse order
		afterHandler = applyAfterMiddlewares(afterHandler, r.afterMiddlewares)

		// Execute the after middlewares
		afterHandler.ServeHTTP(w, req)
	})
}
