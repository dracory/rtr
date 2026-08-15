package middlewares

import (
	"net"
	"net/http"
	"strings"

	"github.com/dracory/rtr"
)

// realIPHeaders lists the headers consulted in priority order when determining
// the client's real IP address. The first header that yields a valid IP wins.
var realIPHeaders = []string{
	"True-Client-IP",
	"X-Real-IP",
	"X-Forwarded-For",
}

// RealIPMiddleware returns a middleware that sets the client's real IP address
// in r.RemoteAddr based on proxy headers. It consults True-Client-IP,
// X-Real-IP, and X-Forwarded-For (in that priority order) and takes the first
// valid IP address. If none of the headers contain a valid IP, RemoteAddr is
// left unchanged.
func RealIPMiddleware() rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Real IP").
		SetHandler(realIPHandler)
}

// realIPHandler is the stdlib-only implementation of the real-IP extraction
// behaviour previously provided by chi's middleware.RealIP.
func realIPHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ip := extractRealIP(r); ip != "" {
			r.RemoteAddr = ip
		}
		next.ServeHTTP(w, r)
	})
}

// extractRealIP inspects the proxy headers on the request and returns the first
// valid IP address found, or an empty string if none are valid.
func extractRealIP(r *http.Request) string {
	for _, header := range realIPHeaders {
		value := r.Header.Get(header)
		if value == "" {
			continue
		}
		// X-Forwarded-For may contain a comma-separated list of IPs; take the
		// first (leftmost) entry which represents the original client.
		if comma := strings.IndexByte(value, ','); comma >= 0 {
			value = value[:comma]
		}
		value = strings.TrimSpace(value)
		if isValidIP(value) {
			return value
		}
	}
	return ""
}

// isValidIP returns true if the given string parses as a valid IPv4 or IPv6
// address.
func isValidIP(s string) bool {
	return net.ParseIP(s) != nil
}
