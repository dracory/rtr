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

	// Handle multiple patterns (comma-separated)
	for _, pattern := range d.patterns {
		// Check each individual pattern in the comma-separated list
		for _, singlePattern := range strings.Split(pattern, ",") {
			singlePattern = strings.TrimSpace(singlePattern)
			if d.matchesPattern(host, singlePattern) {
				return true
			}
		}
	}
	return false
}

// matchesPattern checks if the host matches the given pattern
func (d *domainImpl) matchesPattern(host, pattern string) bool {
	// Handle IPv6 addresses with ports (e.g., [::1]:8080)
	var hostName, hostPort string
	if strings.HasPrefix(host, "[") {
		// IPv6 address with port: [::1]:8080
		if closeBracket := strings.Index(host, "]"); closeBracket != -1 {
			hostName = host[:closeBracket+1] // Keep the brackets for IPv6
			if len(host) > closeBracket+1 && host[closeBracket+1] == ':' {
				hostPort = host[closeBracket+2:]
			}
		}
	} else {
		// Regular domain or IPv4 with port
		hostName, hostPort, _ = strings.Cut(host, ":")
	}

	// Handle pattern (which might be an IPv6 address with port)
	var patternHost, patternPort string
	if strings.HasPrefix(pattern, "[") {
		// IPv6 pattern: [::1]:8080 or [::1]:*
		if closeBracket := strings.Index(pattern, "]"); closeBracket != -1 {
			patternHost = pattern[:closeBracket+1] // Keep the brackets for IPv6
			if len(pattern) > closeBracket+1 && pattern[closeBracket+1] == ':' {
				patternPort = pattern[closeBracket+2:]
			}
		}
	} else {
		// Regular pattern: example.com:8080 or example.com:*
		patternHost, patternPort, _ = strings.Cut(pattern, ":")
	}

	// If pattern specifies a port, it must match exactly (or be a wildcard)
	if patternPort != "" {
		// If host doesn't have a port but pattern requires one, no match
		if hostPort == "" {
			return false
		}
		// If pattern port is not a wildcard and doesn't match host port, no match
		if patternPort != "*" && patternPort != hostPort {
			return false
		}
	}

	// Handle IPv6 pattern matching
	if strings.HasPrefix(patternHost, "[") && strings.HasSuffix(patternHost, "]") {
		// For IPv6, the entire address must match exactly
		return patternHost == hostName
	}

	// Handle exact host match (for both IPv4 and regular domains)
	if patternHost == hostName {
		return true
	}

	// Handle wildcard subdomains (e.g., "*.example.com")
	if strings.HasPrefix(patternHost, "*.") {
		patternHost = strings.TrimPrefix(patternHost, "*.") // Remove the wildcard prefix
		// Only match if host ends with .pattern and has at least one dot
		// This ensures "example.com" doesn't match "*.example.com"
		return strings.Contains(hostName, ".") && strings.HasSuffix(hostName, "."+patternHost)
	}

	return false
}

// String returns a string representation of the domain
func (d *domainImpl) String() string {
	return fmt.Sprintf("Domain(patterns=%v)", d.patterns)
}
