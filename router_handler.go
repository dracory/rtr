package rtr

import (
	"net/http"
)

// buildHandler constructs the final http.Handler for a given route,
// applying all relevant middleware in the correct order:
// global before → domain before → group before → route before →
// handler →
// route after → group after → domain after → global after.
func (r *routerImpl) buildHandler(route RouteInterface, groups []GroupInterface, domain DomainInterface) http.Handler {
	// Count total middlewares to pre-allocate slices
	totalMiddlewares := 0
	
	// Count route middlewares
	totalMiddlewares += len(route.GetAfterMiddlewares()) + len(route.GetBeforeMiddlewares())
	
	// Count group middlewares
	for _, group := range groups {
		totalMiddlewares += len(group.GetAfterMiddlewares()) + len(group.GetBeforeMiddlewares())
	}
	
	// Count domain middlewares
	if domain != nil {
		totalMiddlewares += len(domain.GetAfterMiddlewares()) + len(domain.GetBeforeMiddlewares())
	}
	
	// Count global middlewares
	globalAfter := r.GetAfterMiddlewares()
	globalBefore := r.GetBeforeMiddlewares()
	totalMiddlewares += len(globalAfter) + len(globalBefore)

	// Pre-allocate a single slice for all middlewares in execution order
	allMiddlewares := make([]MiddlewareInterface, 0, totalMiddlewares)
	
	// 1. Add global before middlewares (outermost first)
	allMiddlewares = append(allMiddlewares, globalBefore...)
	
	// 2. Add domain before middlewares
	if domain != nil {
		allMiddlewares = append(allMiddlewares, domain.GetBeforeMiddlewares()...)
	}
	
	// 3. Add group before middlewares (from outermost to innermost)
	for _, group := range groups {
		allMiddlewares = append(allMiddlewares, group.GetBeforeMiddlewares()...)
	}
	
	// 4. Add route before middlewares (innermost)
	allMiddlewares = append(allMiddlewares, route.GetBeforeMiddlewares()...)
	
	// 5. Add route after middlewares (innermost)
	allMiddlewares = append(allMiddlewares, route.GetAfterMiddlewares()...)
	
	// 6. Add group after middlewares (from innermost to outermost)
	for i := len(groups) - 1; i >= 0; i-- {
		allMiddlewares = append(allMiddlewares, groups[i].GetAfterMiddlewares()...)
	}
	
	// 7. Add domain after middlewares
	if domain != nil {
		allMiddlewares = append(allMiddlewares, domain.GetAfterMiddlewares()...)
	}
	
	// 8. Add global after middlewares (outermost)
	allMiddlewares = append(allMiddlewares, globalAfter...)
	
	// Apply all middlewares in reverse order to build the handler chain
	handler := http.Handler(http.HandlerFunc(route.GetHandler()))
	for i := len(allMiddlewares) - 1; i >= 0; i-- {
		if allMiddlewares[i] != nil {
			handler = allMiddlewares[i].Execute(handler)
		}
	}
	
	return handler
}
