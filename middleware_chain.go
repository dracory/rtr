package rtr

import "net/http"

// BuildMiddlewareChain builds a middleware chain from a single slice of middlewares.
// The middlewares are applied in order, with the first middleware in the slice being the outermost.
func BuildMiddlewareChain(handler http.Handler, middlewares []MiddlewareInterface) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i].Execute(handler)
	}
	return handler
}

// BuildMiddlewareChainFromSlices builds a middleware chain from multiple slices of middlewares.
// The slices are processed in order, with the first slice's middlewares being the outermost.
// Within each slice, middlewares are applied in order, with the first middleware in the slice being the outermost.
func BuildMiddlewareChainFromSlices(handler http.Handler, middlewareSlices ...[]MiddlewareInterface) http.Handler {
	// Process slices in reverse order (first slice is outermost)
	for i := len(middlewareSlices) - 1; i >= 0; i-- {
		slice := middlewareSlices[i]
		// Process middlewares in reverse order (first middleware in slice is outermost)
		for j := len(slice) - 1; j >= 0; j-- {
			handler = slice[j].Execute(handler)
		}
	}
	return handler
}
