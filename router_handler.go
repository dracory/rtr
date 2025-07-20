package rtr

import (
	"net/http"
)

// buildHandler constructs the final http.Handler for a given route, applying all relevant middleware
// in the correct order: global -> domain -> group -> route.
func (r *routerImpl) buildHandler(route RouteInterface, groups []GroupInterface, domain DomainInterface) http.Handler {
	// Start with the route's own handler
	finalHandler := http.Handler(http.HandlerFunc(route.GetHandler()))

	// Collect all middlewares
	var allMiddlewares []MiddlewareInterface

	// Global 'before' middlewares
	allMiddlewares = append(allMiddlewares, r.GetBeforeMiddlewares()...)

	// Domain 'before' middlewares
	if domain != nil {
		allMiddlewares = append(allMiddlewares, domain.GetBeforeMiddlewares()...)
	}

	// Group 'before' middlewares (in order from parent to child)
	for _, group := range groups {
		allMiddlewares = append(allMiddlewares, group.GetBeforeMiddlewares()...)
	}

	// Route 'before' middlewares
	allMiddlewares = append(allMiddlewares, route.GetBeforeMiddlewares()...)

	// Route 'after' middlewares
	allMiddlewares = append(allMiddlewares, route.GetAfterMiddlewares()...)

	// Group 'after' middlewares (in reverse order from child to parent)
	for i := len(groups) - 1; i >= 0; i-- {
		allMiddlewares = append(allMiddlewares, groups[i].GetAfterMiddlewares()...)
	}

	// Domain 'after' middlewares
	if domain != nil {
		allMiddlewares = append(allMiddlewares, domain.GetAfterMiddlewares()...)
	}

	// Global 'after' middlewares
	allMiddlewares = append(allMiddlewares, r.GetAfterMiddlewares()...)

	// Apply all middlewares in reverse order to chain them correctly
	for i := len(allMiddlewares) - 1; i >= 0; i-- {
		finalHandler = allMiddlewares[i].Execute(finalHandler)
	}

	return finalHandler
}
