package rtr

import (
	"net/http"
)

// buildHandler constructs the final http.Handler for a given route, applying all relevant middleware
// in the correct order: global before → domain before → group before → route before → handler → route after → group after → domain after → global after.
func (r *routerImpl) buildHandler(route RouteInterface, groups []GroupInterface, domain DomainInterface) http.Handler {
	// Start with the route's own handler
	handler := http.Handler(http.HandlerFunc(route.GetHandler()))

	// Collect all middleware slices in the correct order
	var allMiddlewares [][]MiddlewareInterface

	// 1. Add global before middlewares (outermost)
	allMiddlewares = append(allMiddlewares, r.GetBeforeMiddlewares())

	// 2. Add domain before middlewares
	if domain != nil {
		allMiddlewares = append(allMiddlewares, domain.GetBeforeMiddlewares())
	}

	// 3. Add group before middlewares (from outermost to innermost)
	for _, group := range groups {
		allMiddlewares = append(allMiddlewares, group.GetBeforeMiddlewares())
	}

	// 4. Add route before middlewares (innermost before handler)
	allMiddlewares = append(allMiddlewares, route.GetBeforeMiddlewares())

	// 5. Add route after middlewares (innermost after handler)
	allMiddlewares = append(allMiddlewares, route.GetAfterMiddlewares())

	// 6. Add group after middlewares (from innermost to outermost)
	for i := len(groups) - 1; i >= 0; i-- {
		allMiddlewares = append(allMiddlewares, groups[i].GetAfterMiddlewares())
	}

	// 7. Add domain after middlewares
	if domain != nil {
		allMiddlewares = append(allMiddlewares, domain.GetAfterMiddlewares())
	}

	// 8. Add global after middlewares (outermost)
	allMiddlewares = append(allMiddlewares, r.GetAfterMiddlewares())

	// Apply all middlewares in the correct order using BuildMiddlewareChainFromSlices
	// This will ensure that middlewares are applied in the correct order:
	// 1. Global before (outermost) -> ... -> Route before -> Handler -> Route after -> ... -> Global after (outermost)
	handler = BuildMiddlewareChainFromSlices(handler, allMiddlewares...)

	return handler
}


