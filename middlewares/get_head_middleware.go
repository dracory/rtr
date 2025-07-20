package middlewares

import (
	"net/http"
	"net/http/httptest"

	"github.com/dracory/rtr"
)

// GetHead creates a middleware that automatically routes undefined HEAD requests to GET handlers.
// This is useful for automatically handling HEAD requests without requiring explicit HEAD handlers.
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
