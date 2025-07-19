package rtr

// RouteConfig represents a declarative route configuration
type RouteConfig struct {
	Name             string                 `json:"name,omitempty"`
	Method           string                 `json:"method,omitempty"`
	Path             string                 `json:"path"`
	Handler          Handler                `json:"-"`
	BeforeMiddleware []Middleware           `json:"-"`
	AfterMiddleware  []Middleware           `json:"-"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// GroupConfig represents a declarative group configuration
type GroupConfig struct {
	Name             string                 `json:"name,omitempty"`
	Prefix           string                 `json:"prefix"`
	Routes           []RouteConfig          `json:"routes,omitempty"`
	Groups           []GroupConfig          `json:"groups,omitempty"`
	BeforeMiddleware []Middleware           `json:"-"`
	AfterMiddleware  []Middleware           `json:"-"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// DomainConfig represents a declarative domain configuration
type DomainConfig struct {
	Name             string                 `json:"name,omitempty"`
	Patterns         []string               `json:"patterns"`
	Routes           []RouteConfig          `json:"routes,omitempty"`
	Groups           []GroupConfig          `json:"groups,omitempty"`
	BeforeMiddleware []Middleware           `json:"-"`
	AfterMiddleware  []Middleware           `json:"-"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// RouterConfig represents a complete declarative router configuration
type RouterConfig struct {
	Name             string                 `json:"name,omitempty"`
	Prefix           string                 `json:"prefix,omitempty"`
	Routes           []RouteConfig          `json:"routes,omitempty"`
	Groups           []GroupConfig          `json:"groups,omitempty"`
	Domains          []DomainConfig         `json:"domains,omitempty"`
	BeforeMiddleware []Middleware           `json:"-"`
	AfterMiddleware  []Middleware           `json:"-"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// NewRouterFromConfig creates a router from a declarative configuration
func NewRouterFromConfig(config RouterConfig) RouterInterface {
	router := NewRouter()

	// Set prefix if specified
	if config.Prefix != "" {
		router.SetPrefix(config.Prefix)
	}

	// Add global middleware
	if len(config.BeforeMiddleware) > 0 {
		router.AddBeforeMiddlewares(config.BeforeMiddleware)
	}
	if len(config.AfterMiddleware) > 0 {
		router.AddAfterMiddlewares(config.AfterMiddleware)
	}

	// Add direct routes
	for _, routeConfig := range config.Routes {
		route := buildRouteFromConfig(routeConfig)
		router.AddRoute(route)
	}

	// Add groups
	for _, groupConfig := range config.Groups {
		group := buildGroupFromConfig(groupConfig)
		router.AddGroup(group)
	}

	// Add domains
	for _, domainConfig := range config.Domains {
		domain := buildDomainFromConfig(domainConfig)
		router.AddDomain(domain)
	}

	return router
}

// buildRouteFromConfig converts a RouteConfig to a RouteInterface
func buildRouteFromConfig(config RouteConfig) RouteInterface {
	route := NewRoute()

	if config.Name != "" {
		route.SetName(config.Name)
	}
	if config.Method != "" {
		route.SetMethod(config.Method)
	}
	if config.Path != "" {
		route.SetPath(config.Path)
	}
	if config.Handler != nil {
		route.SetHandler(config.Handler)
	}
	if len(config.BeforeMiddleware) > 0 {
		route.AddBeforeMiddlewares(config.BeforeMiddleware)
	}
	if len(config.AfterMiddleware) > 0 {
		route.AddAfterMiddlewares(config.AfterMiddleware)
	}

	return route
}

// buildGroupFromConfig converts a GroupConfig to a GroupInterface
func buildGroupFromConfig(config GroupConfig) GroupInterface {
	group := NewGroup()

	if config.Prefix != "" {
		group.SetPrefix(config.Prefix)
	}
	if len(config.BeforeMiddleware) > 0 {
		group.AddBeforeMiddlewares(config.BeforeMiddleware)
	}
	if len(config.AfterMiddleware) > 0 {
		group.AddAfterMiddlewares(config.AfterMiddleware)
	}

	// Add routes to group
	for _, routeConfig := range config.Routes {
		route := buildRouteFromConfig(routeConfig)
		group.AddRoute(route)
	}

	// Add nested groups
	for _, nestedGroupConfig := range config.Groups {
		nestedGroup := buildGroupFromConfig(nestedGroupConfig)
		group.AddGroup(nestedGroup)
	}

	return group
}

// buildDomainFromConfig converts a DomainConfig to a DomainInterface
func buildDomainFromConfig(config DomainConfig) DomainInterface {
	domain := NewDomain(config.Patterns...)

	if len(config.BeforeMiddleware) > 0 {
		domain.AddBeforeMiddlewares(config.BeforeMiddleware)
	}
	if len(config.AfterMiddleware) > 0 {
		domain.AddAfterMiddlewares(config.AfterMiddleware)
	}

	// Add routes to domain
	for _, routeConfig := range config.Routes {
		route := buildRouteFromConfig(routeConfig)
		domain.AddRoute(route)
	}

	// Add groups to domain
	for _, groupConfig := range config.Groups {
		group := buildGroupFromConfig(groupConfig)
		domain.AddGroup(group)
	}

	return domain
}

// Convenience functions for common route types

// GET creates a GET route configuration
func GET(path string, handler Handler) RouteConfig {
	return RouteConfig{
		Method:  "GET",
		Path:    path,
		Handler: handler,
	}
}

// POST creates a POST route configuration
func POST(path string, handler Handler) RouteConfig {
	return RouteConfig{
		Method:  "POST",
		Path:    path,
		Handler: handler,
	}
}

// PUT creates a PUT route configuration
func PUT(path string, handler Handler) RouteConfig {
	return RouteConfig{
		Method:  "PUT",
		Path:    path,
		Handler: handler,
	}
}

// DELETE creates a DELETE route configuration
func DELETE(path string, handler Handler) RouteConfig {
	return RouteConfig{
		Method:  "DELETE",
		Path:    path,
		Handler: handler,
	}
}

// PATCH creates a PATCH route configuration
func PATCH(path string, handler Handler) RouteConfig {
	return RouteConfig{
		Method:  "PATCH",
		Path:    path,
		Handler: handler,
	}
}

// OPTIONS creates an OPTIONS route configuration
func OPTIONS(path string, handler Handler) RouteConfig {
	return RouteConfig{
		Method:  "OPTIONS",
		Path:    path,
		Handler: handler,
	}
}

// WithName adds a name to a route configuration
func (r RouteConfig) WithName(name string) RouteConfig {
	r.Name = name
	return r
}

// WithBeforeMiddleware adds before middleware to a route configuration
func (r RouteConfig) WithBeforeMiddleware(middleware ...Middleware) RouteConfig {
	r.BeforeMiddleware = append(r.BeforeMiddleware, middleware...)
	return r
}

// WithAfterMiddleware adds after middleware to a route configuration
func (r RouteConfig) WithAfterMiddleware(middleware ...Middleware) RouteConfig {
	r.AfterMiddleware = append(r.AfterMiddleware, middleware...)
	return r
}

// WithMetadata adds metadata to a route configuration
func (r RouteConfig) WithMetadata(key string, value interface{}) RouteConfig {
	if r.Metadata == nil {
		r.Metadata = make(map[string]interface{})
	}
	r.Metadata[key] = value
	return r
}

// Group creates a group configuration
func Group(prefix string, items ...interface{}) GroupConfig {
	group := GroupConfig{
		Prefix: prefix,
		Routes: []RouteConfig{},
		Groups: []GroupConfig{},
	}

	for _, item := range items {
		switch v := item.(type) {
		case RouteConfig:
			group.Routes = append(group.Routes, v)
		case GroupConfig:
			group.Groups = append(group.Groups, v)
		case []RouteConfig:
			group.Routes = append(group.Routes, v...)
		case []GroupConfig:
			group.Groups = append(group.Groups, v...)
		}
	}

	return group
}

// WithName adds a name to a group configuration
func (g GroupConfig) WithName(name string) GroupConfig {
	g.Name = name
	return g
}

// WithBeforeMiddleware adds before middleware to a group configuration
func (g GroupConfig) WithBeforeMiddleware(middleware ...Middleware) GroupConfig {
	g.BeforeMiddleware = append(g.BeforeMiddleware, middleware...)
	return g
}

// WithAfterMiddleware adds after middleware to a group configuration
func (g GroupConfig) WithAfterMiddleware(middleware ...Middleware) GroupConfig {
	g.AfterMiddleware = append(g.AfterMiddleware, middleware...)
	return g
}

// Domain creates a domain configuration
func Domain(patterns []string, items ...interface{}) DomainConfig {
	domain := DomainConfig{
		Patterns: patterns,
		Routes:   []RouteConfig{},
		Groups:   []GroupConfig{},
	}

	for _, item := range items {
		switch v := item.(type) {
		case RouteConfig:
			domain.Routes = append(domain.Routes, v)
		case GroupConfig:
			domain.Groups = append(domain.Groups, v)
		case []RouteConfig:
			domain.Routes = append(domain.Routes, v...)
		case []GroupConfig:
			domain.Groups = append(domain.Groups, v...)
		}
	}

	return domain
}
