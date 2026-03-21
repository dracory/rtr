package middlewares

import (
	"net/http"

	"github.com/dracory/rtr"
)

// HTTPSRedirectConfig provides configuration for HTTPS redirect middleware
type HTTPSRedirectConfig struct {
	// SkipLocalhost skips HTTPS redirect for localhost and local development
	SkipLocalhost bool
	// TrustedProxies contains list of trusted proxy IPs for X-Forwarded-Proto checking
	TrustedProxies []string
	// CustomSkipFunc allows custom logic to skip HTTPS redirect
	CustomSkipFunc func(r *http.Request) bool
}

// DefaultHTTPSRedirectConfig returns a default configuration
func DefaultHTTPSRedirectConfig() *HTTPSRedirectConfig {
	return &HTTPSRedirectConfig{
		SkipLocalhost:  true,
		TrustedProxies: []string{},
		CustomSkipFunc: nil,
	}
}

// NewHTTPSRedirectMiddleware creates middleware that redirects HTTP requests to HTTPS
func NewHTTPSRedirectMiddleware(config *HTTPSRedirectConfig) rtr.MiddlewareInterface {
	if config == nil {
		config = DefaultHTTPSRedirectConfig()
	}

	return rtr.NewMiddleware().
		SetName("HTTPS Redirect Middleware").
		SetHandler(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check custom skip function first
				if config.CustomSkipFunc != nil && config.CustomSkipFunc(r) {
					next.ServeHTTP(w, r)
					return
				}

				// Skip redirection for localhost if configured
				if config.SkipLocalhost && isLocalhost(r.Host) {
					next.ServeHTTP(w, r)
					return
				}

				// Check if already HTTPS
				if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
					next.ServeHTTP(w, r)
					return
				}

				// Redirect to HTTPS version of same URL
				httpsURL := "https://" + r.Host + r.URL.Path
				if r.URL.RawQuery != "" {
					httpsURL += "?" + r.URL.RawQuery
				}
				http.Redirect(w, r, httpsURL, http.StatusMovedPermanently)
			})
		})
}

// isLocalhost checks if the host is a local development address
func isLocalhost(host string) bool {
	return host == "localhost" ||
		host == "127.0.0.1" ||
		host == "0.0.0.0" ||
		len(host) > 6 && host[len(host)-6:] == ".local" || // Ends with .local
		(len(host) > 7 && host[:7] == "127.0.0.") || // Starts with 127.0.0.
		(len(host) > 4 && host[:4] == "127.") || // Starts with 127. (localhost range)
		(len(host) > 8 && host[:8] == "192.168.") || // Starts with 192.168. (private network)
		(len(host) > 7 && host[:7] == "10.0.0.") || // Starts with 10.0.0. (private network)
		(len(host) > 3 && host[:3] == "10.") // Starts with 10. (private network 10.0.0.0/8)
}
