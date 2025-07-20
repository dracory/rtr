package middlewares

import (
	"net/http"
	"path"
	"strings"

	"github.com/dracory/rtr"
)

// CleanPathMiddleware creates a new middleware that cleans up double slashes in URL paths.
// For example, it converts "/users//1" or "//users////1" to "/users/1".
// This should typically be added early in the middleware chain.
func CleanPathMiddleware() rtr.MiddlewareInterface {
	// return &cleanPathMiddleware{
	// 	name:    "Clean Path Middleware",
	// 	handler: cleanPathHandler(),
	// }
	return rtr.NewMiddleware().
		SetName("Clean Path Middleware").
		SetHandler(cleanPathHandler())
}

// cleanPathHandler returns the actual middleware function that cleans the URL path.
func cleanPathHandler() rtr.StdMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check for nil request or URL
			if r == nil || r.URL == nil {
				if next != nil {
					next.ServeHTTP(w, r)
				}
				return
			}

			// Get the path to clean, prefer RawPath if available (for URL-encoded paths)
			pathToClean := r.URL.Path
			if r.URL.RawPath != "" {
				pathToClean = r.URL.RawPath
			}

			// Check if the path needs cleaning (contains //)
			if !needsCleaning(pathToClean) {
				next.ServeHTTP(w, r)
				return
			}

			// Clean the path using path.Clean but preserve trailing slash
			wasTrailingSlash := len(pathToClean) > 1 && strings.HasSuffix(pathToClean, "/")
			cleanPath := path.Clean("/" + pathToClean)

			// Restore trailing slash if it was present
			if wasTrailingSlash && cleanPath != "/" {
				cleanPath += "/"
			}

			// Handle root path edge case
			if cleanPath == "/" && pathToClean != "/" {
				cleanPath = "/"
			}

			// Preserve the query string if present
			if r.URL.RawQuery != "" {
				cleanPath = cleanPath + "?" + r.URL.RawQuery
			}

			// Only redirect if the path was actually changed
			if cleanPath != pathToClean || (r.URL.RawQuery != "" && r.URL.RawQuery != r.URL.Query().Encode()) {
				http.Redirect(w, r, cleanPath, http.StatusMovedPermanently)
				return
			}

			// Call next handler if it exists
			if next != nil {
				next.ServeHTTP(w, r)
			}
		})
	}
}

// needsCleaning checks if the path contains any double slashes that need cleaning
func needsCleaning(p string) bool {
	if p == "" {
		return false
	}
	return strings.Contains(p, "//")
}
