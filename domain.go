package rtr

import (
	"fmt"
	"strings"
)

// domainImpl implements the DomainInterface
type domainImpl struct {
	patterns          []string
	routes            []RouteInterface
	groups            []GroupInterface
	beforeMiddlewares []Middleware
	afterMiddlewares  []Middleware
}

// NewDomain creates a new domain with the given patterns
func NewDomain(patterns ...string) DomainInterface {
	return &domainImpl{
		patterns: patterns,
		routes:   make([]RouteInterface, 0),
		groups:   make([]GroupInterface, 0),
	}
}

// GetPatterns returns the domain patterns that this domain matches against
func (d *domainImpl) GetPatterns() []string {
	return d.patterns
}

// SetPatterns sets the domain patterns for this domain and returns the domain for method chaining
func (d *domainImpl) SetPatterns(patterns ...string) DomainInterface {
	d.patterns = patterns
	return d
}

// AddRoute adds a route to this domain and returns the domain for method chaining
func (d *domainImpl) AddRoute(route RouteInterface) DomainInterface {
	d.routes = append(d.routes, route)
	return d
}

// AddRoutes adds multiple routes to this domain and returns the domain for method chaining
func (d *domainImpl) AddRoutes(routes []RouteInterface) DomainInterface {
	d.routes = append(d.routes, routes...)
	return d
}

// GetRoutes returns all routes that belong to this domain
func (d *domainImpl) GetRoutes() []RouteInterface {
	return d.routes
}

// AddGroup adds a group to this domain and returns the domain for method chaining
func (d *domainImpl) AddGroup(group GroupInterface) DomainInterface {
	d.groups = append(d.groups, group)
	return d
}

// AddGroups adds multiple groups to this domain and returns the domain for method chaining
func (d *domainImpl) AddGroups(groups []GroupInterface) DomainInterface {
	d.groups = append(d.groups, groups...)
	return d
}

// GetGroups returns all groups that belong to this domain
func (d *domainImpl) GetGroups() []GroupInterface {
	return d.groups
}

// AddBeforeMiddlewares adds middleware functions to be executed before any route handler in this domain
// Returns the domain for method chaining
func (d *domainImpl) AddBeforeMiddlewares(middleware []Middleware) DomainInterface {
	d.beforeMiddlewares = append(d.beforeMiddlewares, middleware...)
	return d
}

// GetBeforeMiddlewares returns all middleware functions that will be executed before any route handler in this domain
func (d *domainImpl) GetBeforeMiddlewares() []Middleware {
	return d.beforeMiddlewares
}

// AddAfterMiddlewares adds middleware functions to be executed after any route handler in this domain
// Returns the domain for method chaining
func (d *domainImpl) AddAfterMiddlewares(middleware []Middleware) DomainInterface {
	d.afterMiddlewares = append(d.afterMiddlewares, middleware...)
	return d
}

// GetAfterMiddlewares returns all middleware functions that will be executed after any route handler in this domain
func (d *domainImpl) GetAfterMiddlewares() []Middleware {
	return d.afterMiddlewares
}

// Match checks if the given host matches any of this domain's patterns
// The host can include a port (e.g., "example.com:8080"), and patterns can optionally specify ports.
// Port matching rules:
// - If pattern includes a port (e.g., "example.com:8080"), it must match exactly
// - Use "*" as port in pattern (e.g., "example.com:*") to match any port
// - If pattern has no port, it matches any port for that host
func (d *domainImpl) Match(host string) bool {
	if host == "" {
		return false
	}

	// Handle IPv6 addresses which use square brackets for addresses with ports (e.g., [::1]:8080)
	if strings.HasPrefix(host, "[") {
		// Extract IPv6 address and port if present
		if closeBracket := strings.Index(host, "]"); closeBracket != -1 {
			if len(host) > closeBracket+1 && host[closeBracket+1] == ':' {
				// Has port: [::1]:8080
				host = host[1:closeBracket] + ":" + host[closeBracket+2:]
			} else {
				// No port: [::1]
				host = host[1:closeBracket]
			}
		}
	}

	for _, pattern := range d.patterns {
		if d.matchesPattern(host, pattern) {
			return true
		}
	}
	return false
}

// matchesPattern checks if the host matches the given pattern
func (d *domainImpl) matchesPattern(host, pattern string) bool {
	// Split pattern into host and port parts
	patternHost, patternPort, _ := strings.Cut(pattern, ":")
	hostName, hostPort, _ := strings.Cut(host, ":")

	// If pattern specifies a port, it must match exactly (or be a wildcard)
	if patternPort != "" {
		if patternPort != "*" && patternPort != hostPort {
			return false
		}
	}

	// Exact host match
	if patternHost == hostName {
		return true
	}

	// Wildcard subdomain (e.g., "*.example.com")
	if strings.HasPrefix(patternHost, "*.") {
		patternHost = strings.TrimPrefix(patternHost, "*.") 
		// Only match if host ends with .pattern (e.g., "sub.example.com" matches "*.example.com")
		// But don't match if host is exactly the pattern (e.g., "example.com" should not match "*.example.com")
		return strings.HasSuffix(hostName, "."+patternHost)
	}

	return false
}

// String returns a string representation of the domain
func (d *domainImpl) String() string {
	return fmt.Sprintf("Domain(patterns=%v)", d.patterns)
}
