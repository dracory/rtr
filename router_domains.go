package rtr

import (
	"context"
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
		if match, params := r.routeMatches(route, req); match {
			// Add params to request context if any
			if len(params) > 0 {
				ctx := context.WithValue(req.Context(), "params", params)
				req = req.WithContext(ctx)
			}
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
	handler := http.Handler(http.HandlerFunc(route.GetHandler()))

	// Apply route's before middlewares in order
	for i := len(route.GetBeforeMiddlewares()) - 1; i >= 0; i-- {
		mw := route.GetBeforeMiddlewares()[i]
		handler = mw.Execute(handler)
	}

	// Apply domain's before middlewares in order
	for i := len(domain.GetBeforeMiddlewares()) - 1; i >= 0; i-- {
		mw := domain.GetBeforeMiddlewares()[i]
		handler = mw.Execute(handler)
	}

	// Apply router's before middlewares in order
	for i := len(r.beforeMiddlewares) - 1; i >= 0; i-- {
		mw := r.beforeMiddlewares[i]
		handler = mw.Execute(handler)
	}

	// Create a final handler that wraps the chain with after middlewares
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Call the handler chain
		handler.ServeHTTP(w, req)

		// After handler execution, run after middlewares in reverse order
		afterHandler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

		// Apply route's after middlewares in reverse order
		for _, mw := range route.GetAfterMiddlewares() {
			afterHandler = mw.Execute(afterHandler)
		}

		// Apply domain's after middlewares in reverse order
		for _, mw := range domain.GetAfterMiddlewares() {
			afterHandler = mw.Execute(afterHandler)
		}

		// Apply router's after middlewares in reverse order
		for _, mw := range r.afterMiddlewares {
			afterHandler = mw.Execute(afterHandler)
		}

		// Execute the after middlewares if any exist
		hasAfterMiddlewares := len(route.GetAfterMiddlewares()) > 0 || 
			len(domain.GetAfterMiddlewares()) > 0 || 
			len(r.afterMiddlewares) > 0

		if hasAfterMiddlewares {
			afterHandler.ServeHTTP(w, req)
		}
	})
}

// wrapWithDomainGroupMiddlewares wraps a route's handler with domain, group, and route middlewares
func (r *routerImpl) wrapWithDomainGroupMiddlewares(route RouteInterface, group GroupInterface, domain DomainInterface, req *http.Request) http.Handler {
	// Start with the route's handler
	handler := http.Handler(http.HandlerFunc(route.GetHandler()))

	// Apply route's before middlewares in order
	for i := len(route.GetBeforeMiddlewares()) - 1; i >= 0; i-- {
		mw := route.GetBeforeMiddlewares()[i]
		handler = mw.Execute(handler)
	}

	// Apply group's before middlewares in order
	for i := len(group.GetBeforeMiddlewares()) - 1; i >= 0; i-- {
		mw := group.GetBeforeMiddlewares()[i]
		handler = mw.Execute(handler)
	}

	// Apply domain's before middlewares in order
	for i := len(domain.GetBeforeMiddlewares()) - 1; i >= 0; i-- {
		mw := domain.GetBeforeMiddlewares()[i]
		handler = mw.Execute(handler)
	}

	// Apply router's before middlewares in order
	for i := len(r.beforeMiddlewares) - 1; i >= 0; i-- {
		mw := r.beforeMiddlewares[i]
		handler = mw.Execute(handler)
	}

	// Create a final handler that wraps the chain with after middlewares
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Call the handler chain
		handler.ServeHTTP(w, req)

		// After handler execution, run after middlewares in reverse order
		afterHandler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

		// Apply route's after middlewares in reverse order
		for _, mw := range route.GetAfterMiddlewares() {
			afterHandler = mw.Execute(afterHandler)
		}

		// Apply group's after middlewares in reverse order
		for _, mw := range group.GetAfterMiddlewares() {
			afterHandler = mw.Execute(afterHandler)
		}

		// Apply domain's after middlewares in reverse order
		for _, mw := range domain.GetAfterMiddlewares() {
			afterHandler = mw.Execute(afterHandler)
		}

		// Apply router's after middlewares in reverse order
		for _, mw := range r.afterMiddlewares {
			afterHandler = mw.Execute(afterHandler)
		}

		// Execute the after middlewares if any exist
		hasAfterMiddlewares := len(route.GetAfterMiddlewares()) > 0 ||
			len(group.GetAfterMiddlewares()) > 0 ||
			len(domain.GetAfterMiddlewares()) > 0 ||
			len(r.afterMiddlewares) > 0

		if hasAfterMiddlewares {
			afterHandler.ServeHTTP(w, req)
		}
	})
}
