package rtr

import (
	"net/http"
)

// buildHandler constructs the final http.Handler for a given route, applying all relevant middleware
// in the correct order: global before → domain before → group before → route before → handler → route after → group after → domain after → global after.
func (r *routerImpl) buildHandler(route RouteInterface, groups []GroupInterface, domain DomainInterface) http.Handler {
	// Start with the route's own handler
	handler := http.Handler(http.HandlerFunc(route.GetHandler()))

	// 1. First, collect all after middlewares in the order they should be applied (innermost first)
	var afterMiddlewares []MiddlewareInterface

	// 1.1 Add route after middlewares (innermost)
	afterMiddlewares = append(afterMiddlewares, route.GetAfterMiddlewares()...)

	// 1.2 Add group after middlewares (from innermost to outermost)
	for i := len(groups) - 1; i >= 0; i-- {
		afterMiddlewares = append(afterMiddlewares, groups[i].GetAfterMiddlewares()...)
	}

	// 1.3 Add domain after middlewares
	if domain != nil {
		afterMiddlewares = append(afterMiddlewares, domain.GetAfterMiddlewares()...)
	}

	// 1.4 Add global after middlewares (outermost)
	afterMiddlewares = append(afterMiddlewares, r.GetAfterMiddlewares()...)

	// 2. Now collect all before middlewares in the order they should be applied (outermost first)
	var beforeMiddlewares []MiddlewareInterface

	// 2.1 Add global before middlewares (outermost)
	beforeMiddlewares = append(beforeMiddlewares, r.GetBeforeMiddlewares()...)

	// 2.2 Add domain before middlewares
	if domain != nil {
		beforeMiddlewares = append(beforeMiddlewares, domain.GetBeforeMiddlewares()...)
	}

	// 2.3 Add group before middlewares (from outermost to innermost)
	for _, group := range groups {
		beforeMiddlewares = append(beforeMiddlewares, group.GetBeforeMiddlewares()...)
	}

	// 2.4 Add route before middlewares (innermost)
	beforeMiddlewares = append(beforeMiddlewares, route.GetBeforeMiddlewares()...)

	// 3. First, wrap the handler with after middlewares (in reverse order)
	for i := len(afterMiddlewares) - 1; i >= 0; i-- {
		handler = afterMiddlewares[i].Execute(handler)
	}

	// 4. Then wrap with before middlewares (in reverse order)
	for i := len(beforeMiddlewares) - 1; i >= 0; i-- {
		handler = beforeMiddlewares[i].Execute(handler)
	}

	return handler
}


