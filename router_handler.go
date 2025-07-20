package rtr

import (
	"net/http"
)

// buildHandler constructs the final http.Handler for a given route, applying all relevant middleware
// in the correct order: global -> domain -> group -> route.
func (r *routerImpl) buildHandler(route RouteInterface, groups []GroupInterface, domain DomainInterface) http.Handler {
	// Start with the route's own handler
	finalHandler := http.Handler(http.HandlerFunc(route.GetHandler()))

	// Chain 'after' middlewares in reverse order (inner to outer)
	finalHandler = r.chainMiddlewares(route.GetAfterMiddlewares(), finalHandler)
	for i := len(groups) - 1; i >= 0; i-- {
		finalHandler = r.chainMiddlewares(groups[i].GetAfterMiddlewares(), finalHandler)
	}
	if domain != nil {
		finalHandler = r.chainMiddlewares(domain.GetAfterMiddlewares(), finalHandler)
	}
	finalHandler = r.chainMiddlewares(r.GetAfterMiddlewares(), finalHandler)

	// Chain 'before' middlewares (outer to inner)
	finalHandler = r.chainMiddlewares(route.GetBeforeMiddlewares(), finalHandler)
	for _, group := range groups {
		finalHandler = r.chainMiddlewares(group.GetBeforeMiddlewares(), finalHandler)
	}
	if domain != nil {
		finalHandler = r.chainMiddlewares(domain.GetBeforeMiddlewares(), finalHandler)
	}
	finalHandler = r.chainMiddlewares(r.GetBeforeMiddlewares(), finalHandler)

	return finalHandler
}

func (r *routerImpl) chainMiddlewares(middlewares []MiddlewareInterface, handler http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i].Execute(handler)
	}
	return handler
}

