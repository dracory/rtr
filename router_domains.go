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
	// First check if the domain matches the request host
	host := req.Host
	if host == "" {
		host = req.URL.Host
	}

	if !domain.Match(host) {
		return nil, nil
	}

	// Check direct routes in domain
	for _, route := range domain.GetRoutes() {
		if match, params := r.routeMatches(route, req); match {
			// Create a copy of the request to avoid mutating the original
			reqCopy := req.Clone(req.Context())
			
			// Add params to request context if any
			if len(params) > 0 {
				ctx := context.WithValue(reqCopy.Context(), ParamsKey, params)
				reqCopy = reqCopy.WithContext(ctx)
			}
			
			// Build the handler with the updated request
			handler := r.buildHandler(route, nil, domain)
			
			// Return a handler that will use the request with the updated context
			return route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.ServeHTTP(w, reqCopy)
			})
		}
	}

	// Check route groups in domain
	for _, group := range domain.GetGroups() {
		if route, handler := r.findMatchingRouteInGroup(group, req, nil); route != nil {
			return route, handler
		}
	}

	return nil, nil
}
