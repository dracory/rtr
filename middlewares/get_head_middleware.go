package middlewares

import (
	"net/http"
	"net/http/httptest"

	"github.com/dracory/rtr"
)

// GetHead creates a middleware that automatically routes undefined HEAD requests to
// GET handlers. This is useful for automatically handling HEAD requests without
// requiring explicit HEAD handlers.
//
// By using this middleware, you are in compliance with the HTTP/1.1 spec (RFC
// 2616), which states that servers MUST support the HEAD method for any URI that
// returns a response body for a GET request.
//
// Additionally, this middleware provides a performance benefit by saving clients
// the overhead of downloading the full response body when they only need
// metadata.
//
// This is a common web practice, and many web frameworks and servers (like
// Express.js, Django, etc.) provide this functionality out of the box. It also
// saves developers from having to implement HEAD handlers separately for every
// route.
//
// Usage:
//
//	router := rtr.NewRouter()
//	router.AddRoute(rtr.NewRoute().
//	  SetMethod("GET").
//	  SetPath("/test").
//	  SetHandler(func(w http.ResponseWriter, r *http.Request) {
//	    w.WriteHeader(http.StatusOK)
//	  }).
//	  AddMiddleware(middlewares.GetHead()))
//
// Parameters:
//   - next: The next handler in the middleware chain.
//
// Returns:
//   - A middleware that automatically routes undefined HEAD requests to GET handlers.
func GetHead() rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("GET-HEAD Middleware").
		SetHandler(getHeadHandler)
}

func getHeadHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only handle HEAD requests
		if r.Method == "HEAD" {
			// Create a response recorder to capture the response
			recorder := httptest.NewRecorder()

			// Change the method to GET
			r.Method = "GET"

			// Process the request with the next handler
			next.ServeHTTP(recorder, r)

			// Copy only the headers (not the body) to the original response
			headers := w.Header()
			for k, v := range recorder.Header() {
				headers[k] = v
			}

			// Set the status code and discard the body
			w.WriteHeader(recorder.Code)
			return
		}

		// For non-HEAD requests, just pass through
		next.ServeHTTP(w, r)
	})
}
