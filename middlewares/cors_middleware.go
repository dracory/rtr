package middlewares

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/dracory/rtr"
)

// CORSOptions configures the behaviour of CORSMiddleware.
// It mirrors the fields of go-chi/cors.Options so existing callers can migrate
// by changing only the type name.
type CORSOptions struct {
	// AllowedOrigins is a list of origins a cross-domain request can be
	// executed from. If the special "*" value is present in the list, all
	// origins will be allowed. An origin may contain a wildcard (*) to replace
	// 0 or more characters (i.e.: http://*.domain.com). Usage of wildcards
	// implies a small performance penalty. Only one wildcard can be used per
	// origin. Default value is ["*"].
	AllowedOrigins []string

	// AllowedMethods is a list of methods the client is allowed to use with
	// cross-domain requests. Default value is simple methods (HEAD, GET and POST).
	AllowedMethods []string

	// AllowedHeaders is list of non simple headers the client is allowed to use
	// with cross-domain requests. If the special "*" value is present in the
	// list, all headers will be allowed. Default value is [] but "Origin" is
	// always appended to the list.
	AllowedHeaders []string

	// ExposedHeaders indicates which headers are safe to expose to the API of a
	// CORS API specification.
	ExposedHeaders []string

	// AllowCredentials indicates whether the request can include user
	// credentials like cookies, HTTP authentication or client side SSL
	// certificates.
	AllowCredentials bool

	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached.
	MaxAge int

	// OptionsPassthrough instructs preflight to let other potential next
	// handlers to process the OPTIONS method. Turn this on if your application
	// handles OPTIONS.
	OptionsPassthrough bool
}

// CORSMiddleware returns a middleware that handles CORS requests.
// It is a stdlib-only reimplementation of go-chi/cors.
// By default (via DefaultCORSMiddleware) it allows all origins, common methods,
// and common headers.
func CORSMiddleware(opts CORSOptions) rtr.MiddlewareInterface {
	c := newCORS(opts)
	return rtr.NewMiddleware(
		rtr.WithName("CORS"),
		rtr.WithHandler(func(next http.Handler) http.Handler {
			return c.handler(next)
		}),
	)
}

// DefaultCORSMiddleware returns a CORS middleware with sensible defaults:
// - Allow all origins
// - Allow common HTTP methods
// - Allow common headers
// - Allow credentials
// - Max age: 300 (5 minutes)
func DefaultCORSMiddleware() rtr.MiddlewareInterface {
	return CORSMiddleware(CORSOptions{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutes
	})
}

// wildcard represents an origin pattern with a single "*" placeholder.
type wildcard struct {
	prefix string
	suffix string
}

// match returns true if the given origin matches the wildcard pattern.
func (w wildcard) match(origin string) bool {
	if !strings.HasPrefix(origin, w.prefix) {
		return false
	}
	return strings.HasSuffix(origin[len(w.prefix):], w.suffix)
}

// cors holds the parsed CORS configuration, mirroring go-chi/cors's Cors struct.
type cors struct {
	// Normalized list of plain allowed origins
	allowedOrigins []string

	// List of allowed origins containing wildcards
	allowedWOrigins []wildcard

	// Normalized list of allowed headers (canonical form)
	allowedHeaders []string

	// Normalized list of allowed methods (upper-case)
	allowedMethods []string

	// Normalized list of exposed headers (canonical form)
	exposedHeaders []string

	maxAge int

	// Set to true when allowed origins contains a "*"
	allowedOriginsAll bool

	// Set to true when allowed headers contains a "*"
	allowedHeadersAll bool

	allowCredentials  bool
	optionPassthrough bool
}

// newCORS builds a cors instance from CORSOptions, normalising all values.
// Mirrors the normalisation logic of go-chi/cors.New.
func newCORS(opts CORSOptions) *cors {
	c := &cors{
		exposedHeaders:    convert(opts.ExposedHeaders, http.CanonicalHeaderKey),
		allowCredentials:  opts.AllowCredentials,
		maxAge:            opts.MaxAge,
		optionPassthrough: opts.OptionsPassthrough,
	}

	// Allowed Origins. Spec requires case-sensitive matching, but chi ignores
	// the spec here to be less error-prone; we do the same.
	if len(opts.AllowedOrigins) == 0 {
		// Default is all origins
		c.allowedOriginsAll = true
	} else {
		c.allowedOrigins = []string{}
		c.allowedWOrigins = []wildcard{}
		for _, origin := range opts.AllowedOrigins {
			origin = strings.ToLower(origin)
			if origin == "*" {
				// If "*" is present in the list, turn the whole list into a match all
				c.allowedOriginsAll = true
				c.allowedOrigins = nil
				c.allowedWOrigins = nil
				break
			} else if i := strings.IndexByte(origin, '*'); i >= 0 {
				// Split the origin in two: start and end string without the *
				w := wildcard{origin[0:i], origin[i+1:]}
				c.allowedWOrigins = append(c.allowedWOrigins, w)
			} else {
				c.allowedOrigins = append(c.allowedOrigins, origin)
			}
		}
	}

	// Allowed Headers
	if len(opts.AllowedHeaders) == 0 {
		// Use sensible defaults
		c.allowedHeaders = []string{"Origin", "Accept", "Content-Type"}
	} else {
		// Origin is always appended as some browsers will always request for
		// this header at preflight. Copy into a new slice first to avoid
		// mutating the caller's underlying array when it has spare capacity.
		headers := make([]string, 0, len(opts.AllowedHeaders)+1)
		headers = append(headers, opts.AllowedHeaders...)
		headers = append(headers, "Origin")
		c.allowedHeaders = convert(headers, http.CanonicalHeaderKey)
		for _, h := range opts.AllowedHeaders {
			if h == "*" {
				c.allowedHeadersAll = true
				c.allowedHeaders = nil
				break
			}
		}
	}

	// Allowed Methods
	if len(opts.AllowedMethods) == 0 {
		// Default is spec's "simple" methods
		c.allowedMethods = []string{http.MethodGet, http.MethodPost, http.MethodHead}
	} else {
		c.allowedMethods = convert(opts.AllowedMethods, strings.ToUpper)
	}

	return c
}

// convert applies fn to each element of src and returns the result.
func convert(src []string, fn func(string) string) []string {
	if src == nil {
		return nil
	}
	out := make([]string, len(src))
	for i, v := range src {
		out[i] = fn(v)
	}
	return out
}

// handler creates the middleware handler that applies CORS headers.
// Mirrors go-chi/cors.Cors.Handler.
func (c *cors) handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// null or empty Origin header value is acceptable and it is considered
		// having that header
		_, hasOriginHeader := r.Header["Origin"]

		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" && hasOriginHeader {
			c.handlePreflight(w, r)
			// Preflight requests are standalone and should stop the chain as
			// some other middleware may not handle OPTIONS requests correctly.
			if c.optionPassthrough {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		} else {
			c.handleActualRequest(w, r)
			next.ServeHTTP(w, r)
		}
	})
}

// handlePreflight handles pre-flight CORS requests.
// Mirrors go-chi/cors.Cors.handlePreflight.
func (c *cors) handlePreflight(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	origin := r.Header.Get("Origin")

	if r.Method != http.MethodOptions {
		return
	}

	// Always set Vary headers
	headers.Add("Vary", "Origin")
	headers.Add("Vary", "Access-Control-Request-Method")
	headers.Add("Vary", "Access-Control-Request-Headers")

	if !c.isOriginAllowed(origin) {
		return
	}

	reqMethod := r.Header.Get("Access-Control-Request-Method")
	if !c.isMethodAllowed(reqMethod) {
		return
	}
	reqHeaders := parseHeaderList(r.Header.Get("Access-Control-Request-Headers"))
	if !c.areHeadersAllowed(reqHeaders) {
		return
	}
	headers.Set("Access-Control-Allow-Origin", c.allowOriginValue(origin))
	// Spec says: Since the list of methods can be unbounded, simply returning
	// the method indicated by Access-Control-Request-Method (if supported) can
	// be enough
	headers.Set("Access-Control-Allow-Methods", strings.ToUpper(reqMethod))
	if len(reqHeaders) > 0 {
		// Spec says: Since the list of headers can be unbounded, simply
		// returning supported headers from Access-Control-Request-Headers can
		// be enough
		headers.Set("Access-Control-Allow-Headers", strings.Join(reqHeaders, ", "))
	}
	if c.allowCredentials {
		headers.Set("Access-Control-Allow-Credentials", "true")
	}
	if c.maxAge > 0 {
		headers.Set("Access-Control-Max-Age", strconv.Itoa(c.maxAge))
	}
}

// handleActualRequest handles simple cross-origin requests, actual request or
// redirects. Mirrors go-chi/cors.Cors.handleActualRequest.
func (c *cors) handleActualRequest(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	// null Origin header value is acceptable and it is considered having that
	// header
	_, hasOriginHeader := r.Header["Origin"]

	// Always set Vary
	headers.Add("Vary", "Origin")

	if !hasOriginHeader {
		return
	}
	origin := r.Header.Get("Origin")
	if !c.isOriginAllowed(origin) {
		return
	}

	// Spec doesn't instruct to check the allowed methods for simple
	// cross-origin requests, but we think it's a nice feature to have control.
	if !c.isMethodAllowed(r.Method) {
		return
	}
	headers.Set("Access-Control-Allow-Origin", c.allowOriginValue(origin))
	if len(c.exposedHeaders) > 0 {
		headers.Set("Access-Control-Expose-Headers", strings.Join(c.exposedHeaders, ", "))
	}
	if c.allowCredentials {
		headers.Set("Access-Control-Allow-Credentials", "true")
	}
}

// isOriginAllowed checks if a given origin is allowed to perform cross-domain
// requests on the endpoint.
func (c *cors) isOriginAllowed(origin string) bool {
	if c.allowedOriginsAll {
		return true
	}
	origin = strings.ToLower(origin)
	for _, o := range c.allowedOrigins {
		if o == origin {
			return true
		}
	}
	for _, w := range c.allowedWOrigins {
		if w.match(origin) {
			return true
		}
	}
	return false
}

// allowOriginValue returns the value to set on the Access-Control-Allow-Origin
// response header. When all origins are allowed AND credentials are enabled,
// the CORS spec forbids the wildcard "*" — instead the request's specific
// Origin must be echoed back. See https://www.w3.org/TR/cors/#resource-requests
func (c *cors) allowOriginValue(origin string) string {
	if c.allowedOriginsAll {
		if c.allowCredentials {
			return origin
		}
		return "*"
	}
	return origin
}

// isMethodAllowed checks if a given method can be used as part of a
// cross-domain request on the endpoint.
func (c *cors) isMethodAllowed(method string) bool {
	if len(c.allowedMethods) == 0 {
		// If no method allowed, always return false, even for preflight request
		return false
	}
	method = strings.ToUpper(method)
	if method == http.MethodOptions {
		// Always allow preflight requests
		return true
	}
	for _, m := range c.allowedMethods {
		if m == method {
			return true
		}
	}
	return false
}

// areHeadersAllowed checks if a given list of headers are allowed to be used
// within a cross-domain request.
func (c *cors) areHeadersAllowed(requestedHeaders []string) bool {
	if c.allowedHeadersAll || len(requestedHeaders) == 0 {
		return true
	}
	for _, header := range requestedHeaders {
		header = http.CanonicalHeaderKey(header)
		found := false
		for _, h := range c.allowedHeaders {
			if h == header {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// parseHeaderList splits a comma-separated header list and trims whitespace
// from each entry. Mirrors go-chi/cors.parseHeaderList.
func parseHeaderList(headerList string) []string {
	if headerList == "" {
		return nil
	}
	parts := strings.Split(headerList, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
