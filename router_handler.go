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
	
	// 5. AFTER middlewares (non-onion execution): desired execution order is
	//    route → groups (outer→inner) → domain → global.
	//    Because we build by reverse-wrapping below, append OUTER→INNER and
	//    reverse each slice so execution preserves definition order.

	// 5.1 Global after (outermost) — append REVERSED
	for i := len(globalAfter) - 1; i >= 0; i-- {
		allMiddlewares = append(allMiddlewares, globalAfter[i])
	}

	// 5.2 Domain after — append REVERSED
	if domain != nil {
		da := domain.GetAfterMiddlewares()
		for i := len(da) - 1; i >= 0; i-- {
			allMiddlewares = append(allMiddlewares, da[i])
		}
	}

	// 5.3 Groups after from OUTERMOST to INNERMOST — each slice REVERSED
	for i := 0; i < len(groups); i++ {
		ga := groups[i].GetAfterMiddlewares()
		for j := len(ga) - 1; j >= 0; j-- {
			allMiddlewares = append(allMiddlewares, ga[j])
		}
	}

	// 5.4 Route after (innermost) — append REVERSED
	ra := route.GetAfterMiddlewares()
	for i := len(ra) - 1; i >= 0; i-- {
		allMiddlewares = append(allMiddlewares, ra[i])
	}
	
	// Apply all middlewares in reverse order to build the handler chain
	handler := http.Handler(http.HandlerFunc(route.GetHandler()))
	for i := len(allMiddlewares) - 1; i >= 0; i-- {
		if allMiddlewares[i] != nil {
			handler = allMiddlewares[i].Execute(handler)
		}
	}
	
	return handler
}
