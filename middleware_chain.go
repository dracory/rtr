package rtr

import "net/http"

// BuildMiddlewareChain builds a middleware chain from a series of middleware slices.
func BuildMiddlewareChain(handler http.Handler, middlewares ...[]MiddlewareInterface) http.Handler {
	for _, mws := range middlewares {
		for i := len(mws) - 1; i >= 0; i-- {
			handler = mws[i].Execute(handler)
		}
	}
	return handler
}
