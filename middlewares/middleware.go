package middlewares

// This file provides backward compatibility exports for the main rtr package.
// These are deprecated and users should use the direct functions from this package instead.

import "net/http"

// Middleware represents a middleware function that wraps an http.Handler.
// This type alias is provided for backward compatibility.
type Middleware func(http.Handler) http.Handler
