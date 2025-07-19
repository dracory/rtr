package rtr

// RouteConfig represents a declarative route configuration
type RouteConfig struct {
	Name                  string                 `json:"name,omitempty"`
	Method                string                 `json:"method,omitempty"`
	Path                  string                 `json:"path"`
	Handler               StdHandler             `json:"-"`
	ErrorHandler          ErrorHandler           `json:"-"`
	HTMLHandler           HTMLHandler            `json:"-"`
	JSONHandler           JSONHandler            `json:"-"`
	CSSHandler            CSSHandler             `json:"-"`
	XMLHandler            XMLHandler             `json:"-"`
	TextHandler           TextHandler            `json:"-"`
	BeforeMiddleware      []MiddlewareInterface  `json:"-"`
	AfterMiddleware       []MiddlewareInterface  `json:"-"`
	NamedBeforeMiddleware []MiddlewareInterface  `json:"-"`
	NamedAfterMiddleware  []MiddlewareInterface  `json:"-"`
	Metadata              map[string]interface{} `json:"metadata,omitempty"`
}

var _ RouteInterface = (*RouteConfig)(nil)

// GroupConfig represents a declarative group configuration
type GroupConfig struct {
	Name             string                 `json:"name,omitempty"`
	Prefix           string                 `json:"prefix"`
	Routes           []RouteConfig          `json:"routes,omitempty"`
	Groups           []GroupConfig          `json:"groups,omitempty"`
	BeforeMiddleware []MiddlewareInterface  `json:"-"`
	AfterMiddleware  []MiddlewareInterface  `json:"-"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

var _ GroupInterface = (*GroupConfig)(nil)

// DomainConfig represents a declarative domain configuration
type DomainConfig struct {
	Name             string                 `json:"name,omitempty"`
	Patterns         []string               `json:"patterns"`
	Routes           []RouteConfig          `json:"routes,omitempty"`
	Groups           []GroupConfig          `json:"groups,omitempty"`
	BeforeMiddleware []MiddlewareInterface  `json:"-"`
	AfterMiddleware  []MiddlewareInterface  `json:"-"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// RouterConfig represents a complete declarative router configuration
type RouterConfig struct {
	Name             string                 `json:"name,omitempty"`
	Prefix           string                 `json:"prefix,omitempty"`
	Routes           []RouteConfig          `json:"routes,omitempty"`
	Groups           []GroupConfig          `json:"groups,omitempty"`
	Domains          []DomainConfig         `json:"domains,omitempty"`
	BeforeMiddleware []StdMiddleware        `json:"-"`
	AfterMiddleware  []StdMiddleware        `json:"-"`
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
		router.AddBeforeMiddlewares(MiddlewaresToInterfaces(config.BeforeMiddleware))
	}
	if len(config.AfterMiddleware) > 0 {
		router.AddAfterMiddlewares(MiddlewaresToInterfaces(config.AfterMiddleware))
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
	if config.HTMLHandler != nil {
		route.SetHTMLHandler(config.HTMLHandler)
	}
	if config.JSONHandler != nil {
		route.SetJSONHandler(config.JSONHandler)
	}
	if config.CSSHandler != nil {
		route.SetCSSHandler(config.CSSHandler)
	}
	if config.XMLHandler != nil {
		route.SetXMLHandler(config.XMLHandler)
	}
	if config.TextHandler != nil {
		route.SetTextHandler(config.TextHandler)
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
func GET(path string, handler StdHandler) RouteConfig {
	return RouteConfig{
		Method:  "GET",
		Path:    path,
		Handler: handler,
	}
}

// POST creates a POST route configuration
func POST(path string, handler StdHandler) RouteConfig {
	return RouteConfig{
		Method:  "POST",
		Path:    path,
		Handler: handler,
	}
}

// PUT creates a PUT route configuration
func PUT(path string, handler StdHandler) RouteConfig {
	return RouteConfig{
		Method:  "PUT",
		Path:    path,
		Handler: handler,
	}
}

// DELETE creates a DELETE route configuration
func DELETE(path string, handler StdHandler) RouteConfig {
	return RouteConfig{
		Method:  "DELETE",
		Path:    path,
		Handler: handler,
	}
}

// PATCH creates a PATCH route configuration
func PATCH(path string, handler StdHandler) RouteConfig {
	return RouteConfig{
		Method:  "PATCH",
		Path:    path,
		Handler: handler,
	}
}

// OPTIONS creates an OPTIONS route configuration
func OPTIONS(path string, handler StdHandler) RouteConfig {
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
func (r RouteConfig) WithBeforeMiddleware(middleware ...StdMiddleware) RouteConfig {
	r.BeforeMiddleware = append(r.BeforeMiddleware, MiddlewaresToInterfaces(middleware)...)
	return r
}

// WithAfterMiddleware adds after middleware to a route configuration
func (r RouteConfig) WithAfterMiddleware(middleware ...StdMiddleware) RouteConfig {
	r.AfterMiddleware = append(r.AfterMiddleware, MiddlewaresToInterfaces(middleware)...)
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
func (g GroupConfig) WithBeforeMiddleware(middleware ...MiddlewareInterface) GroupConfig {
	g.BeforeMiddleware = append(g.BeforeMiddleware, middleware...)
	return g
}

// WithAfterMiddleware adds after middleware to a group configuration
func (g GroupConfig) WithAfterMiddleware(middleware ...MiddlewareInterface) GroupConfig {
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

var _ DomainInterface = (*DomainConfig)(nil)

// RouteInterface implementation for RouteConfig

// GetMethod returns the HTTP method associated with this route.
func (r *RouteConfig) GetMethod() string {
	return r.Method
}

// SetMethod sets the HTTP method for this route and returns the route for method chaining.
func (r *RouteConfig) SetMethod(method string) RouteInterface {
	r.Method = method
	return r
}

// GetPath returns the URL path pattern associated with this route.
func (r *RouteConfig) GetPath() string {
	return r.Path
}

// SetPath sets the URL path pattern for this route and returns the route for method chaining.
func (r *RouteConfig) SetPath(path string) RouteInterface {
	r.Path = path
	return r
}

// GetHandler returns the handler function associated with this route.
func (r *RouteConfig) GetHandler() StdHandler {
	return r.Handler
}

// SetHandler sets the handler function for this route and returns the route for method chaining.
func (r *RouteConfig) SetHandler(handler StdHandler) RouteInterface {
	r.Handler = handler
	return r
}

// GetErrorHandler returns the error handler function associated with this route.
func (r *RouteConfig) GetErrorHandler() ErrorHandler {
	return r.ErrorHandler
}

// SetErrorHandler sets the error handler function for this route and returns the route for method chaining.
func (r *RouteConfig) SetErrorHandler(handler ErrorHandler) RouteInterface {
	r.ErrorHandler = handler
	return r
}

// GetStringHandler returns the string handler function associated with this route.
func (r *RouteConfig) GetStringHandler() StringHandler {
	// RouteConfig doesn't have a StringHandler field, return nil
	return nil
}

// SetStringHandler sets the string handler function for this route and returns the route for method chaining.
func (r *RouteConfig) SetStringHandler(handler StringHandler) RouteInterface {
	// RouteConfig doesn't have a StringHandler field, this is a no-op
	return r
}

// GetHTMLHandler returns the HTML handler function associated with this route.
func (r *RouteConfig) GetHTMLHandler() HTMLHandler {
	return r.HTMLHandler
}

// SetHTMLHandler sets the HTML handler function for this route and returns the route for method chaining.
func (r *RouteConfig) SetHTMLHandler(handler HTMLHandler) RouteInterface {
	r.HTMLHandler = handler
	return r
}

// GetJSONHandler returns the JSON handler function associated with this route.
func (r *RouteConfig) GetJSONHandler() JSONHandler {
	return r.JSONHandler
}

// SetJSONHandler sets the JSON handler function for this route and returns the route for method chaining.
func (r *RouteConfig) SetJSONHandler(handler JSONHandler) RouteInterface {
	r.JSONHandler = handler
	return r
}

// GetCSSHandler returns the CSS handler function associated with this route.
func (r *RouteConfig) GetCSSHandler() CSSHandler {
	return r.CSSHandler
}

// SetCSSHandler sets the CSS handler function for this route and returns the route for method chaining.
func (r *RouteConfig) SetCSSHandler(handler CSSHandler) RouteInterface {
	r.CSSHandler = handler
	return r
}

// GetXMLHandler returns the XML handler function associated with this route.
func (r *RouteConfig) GetXMLHandler() XMLHandler {
	return r.XMLHandler
}

// SetXMLHandler sets the XML handler function for this route and returns the route for method chaining.
func (r *RouteConfig) SetXMLHandler(handler XMLHandler) RouteInterface {
	r.XMLHandler = handler
	return r
}

// GetTextHandler returns the text handler function associated with this route.
func (r *RouteConfig) GetTextHandler() TextHandler {
	return r.TextHandler
}

// SetTextHandler sets the text handler function for this route and returns the route for method chaining.
func (r *RouteConfig) SetTextHandler(handler TextHandler) RouteInterface {
	r.TextHandler = handler
	return r
}

// GetJSHandler returns the JavaScript handler function associated with this route.
func (r *RouteConfig) GetJSHandler() JSHandler {
	// RouteConfig doesn't have a JSHandler field, return nil
	return nil
}

// SetJSHandler sets the JavaScript handler function for this route and returns the route for method chaining.
func (r *RouteConfig) SetJSHandler(handler JSHandler) RouteInterface {
	// RouteConfig doesn't have a JSHandler field, this is a no-op
	return r
}

// GetName returns the name associated with this route.
func (r *RouteConfig) GetName() string {
	return r.Name
}

// SetName sets the name for this route and returns the route for method chaining.
func (r *RouteConfig) SetName(name string) RouteInterface {
	r.Name = name
	return r
}

// GetMetadata returns the metadata associated with this route.
func (r *RouteConfig) GetMetadata() map[string]interface{} {
	return r.Metadata
}

// SetMetadata sets the metadata for this route and returns the route for method chaining.
func (r *RouteConfig) SetMetadata(metadata map[string]interface{}) RouteInterface {
	r.Metadata = metadata
	return r
}

// AddBeforeMiddlewares adds middleware functions to be executed before the route handler.
func (r *RouteConfig) AddBeforeMiddlewares(middleware []MiddlewareInterface) RouteInterface {
	r.BeforeMiddleware = append(r.BeforeMiddleware, middleware...)
	return r
}

// GetBeforeMiddlewares returns all middleware functions that will be executed before the route handler.
func (r *RouteConfig) GetBeforeMiddlewares() []MiddlewareInterface {
	return r.BeforeMiddleware
}

// AddAfterMiddlewares adds middleware functions to be executed after the route handler.
func (r *RouteConfig) AddAfterMiddlewares(middleware []MiddlewareInterface) RouteInterface {
	r.AfterMiddleware = append(r.AfterMiddleware, middleware...)
	return r
}

// GetAfterMiddlewares returns all middleware functions that will be executed after the route handler.
func (r *RouteConfig) GetAfterMiddlewares() []MiddlewareInterface {
	return r.AfterMiddleware
}

// GroupInterface implementation for GroupConfig

// GetPrefix returns the URL path prefix associated with this group.
func (g *GroupConfig) GetPrefix() string {
	return g.Prefix
}

// SetPrefix sets the URL path prefix for this group and returns the group for method chaining.
func (g *GroupConfig) SetPrefix(prefix string) GroupInterface {
	g.Prefix = prefix
	return g
}

// AddRoute adds a single route to this group and returns the group for method chaining.
func (g *GroupConfig) AddRoute(route RouteInterface) GroupInterface {
	// Convert RouteInterface to RouteConfig if possible
	if routeConfig, ok := route.(*RouteConfig); ok {
		g.Routes = append(g.Routes, *routeConfig)
	} else {
		// Create a new RouteConfig from the RouteInterface
		routeConfig := RouteConfig{
			Name:             route.GetName(),
			Method:           route.GetMethod(),
			Path:             route.GetPath(),
			Handler:          route.GetHandler(),
			ErrorHandler:     route.GetErrorHandler(),
			HTMLHandler:      route.GetHTMLHandler(),
			JSONHandler:      route.GetJSONHandler(),
			CSSHandler:       route.GetCSSHandler(),
			XMLHandler:       route.GetXMLHandler(),
			TextHandler:      route.GetTextHandler(),
			BeforeMiddleware: route.GetBeforeMiddlewares(),
			AfterMiddleware:  route.GetAfterMiddlewares(),
			Metadata:         make(map[string]interface{}),
		}
		g.Routes = append(g.Routes, routeConfig)
	}
	return g
}

// AddRoutes adds multiple routes to this group and returns the group for method chaining.
func (g *GroupConfig) AddRoutes(routes []RouteInterface) GroupInterface {
	for _, route := range routes {
		g.AddRoute(route)
	}
	return g
}

// GetRoutes returns all routes that belong to this group.
func (g *GroupConfig) GetRoutes() []RouteInterface {
	routes := make([]RouteInterface, len(g.Routes))
	for i, route := range g.Routes {
		routes[i] = &route
	}
	return routes
}

// AddGroup adds a single nested group to this group and returns the group for method chaining.
func (g *GroupConfig) AddGroup(group GroupInterface) GroupInterface {
	// Convert GroupInterface to GroupConfig if possible
	if groupConfig, ok := group.(*GroupConfig); ok {
		g.Groups = append(g.Groups, *groupConfig)
	} else {
		// Create a new GroupConfig from the GroupInterface
		groupConfig := GroupConfig{
			Prefix:           group.GetPrefix(),
			Routes:           []RouteConfig{},
			Groups:           []GroupConfig{},
			BeforeMiddleware: group.GetBeforeMiddlewares(),
			AfterMiddleware:  group.GetAfterMiddlewares(),
		}
		// Convert routes
		for _, route := range group.GetRoutes() {
			groupConfig.AddRoute(route)
		}
		// Convert nested groups
		for _, nestedGroup := range group.GetGroups() {
			groupConfig.AddGroup(nestedGroup)
		}
		g.Groups = append(g.Groups, groupConfig)
	}
	return g
}

// AddGroups adds multiple nested groups to this group and returns the group for method chaining.
func (g *GroupConfig) AddGroups(groups []GroupInterface) GroupInterface {
	for _, group := range groups {
		g.AddGroup(group)
	}
	return g
}

// GetGroups returns all nested groups that belong to this group.
func (g *GroupConfig) GetGroups() []GroupInterface {
	groups := make([]GroupInterface, len(g.Groups))
	for i, group := range g.Groups {
		groups[i] = &group
	}
	return groups
}

// AddBeforeMiddlewares adds middleware functions to be executed before any route handler in this group.
func (g *GroupConfig) AddBeforeMiddlewares(middleware []MiddlewareInterface) GroupInterface {
	g.BeforeMiddleware = append(g.BeforeMiddleware, middleware...)
	return g
}

// GetBeforeMiddlewares returns all middleware functions that will be executed before any route handler in this group.
func (g *GroupConfig) GetBeforeMiddlewares() []MiddlewareInterface {
	return g.BeforeMiddleware
}

// AddAfterMiddlewares adds middleware functions to be executed after any route handler in this group.
func (g *GroupConfig) AddAfterMiddlewares(middleware []MiddlewareInterface) GroupInterface {
	g.AfterMiddleware = append(g.AfterMiddleware, middleware...)
	return g
}

// GetAfterMiddlewares returns all middleware functions that will be executed after any route handler in this group.
func (g *GroupConfig) GetAfterMiddlewares() []MiddlewareInterface {
	return g.AfterMiddleware
}

// DomainInterface implementation for DomainConfig

// GetPatterns returns the domain patterns that this domain matches against
func (d *DomainConfig) GetPatterns() []string {
	return d.Patterns
}

// SetPatterns sets the domain patterns for this domain and returns the domain for method chaining
func (d *DomainConfig) SetPatterns(patterns ...string) DomainInterface {
	d.Patterns = patterns
	return d
}

// AddRoute adds a route to this domain and returns the domain for method chaining
func (d *DomainConfig) AddRoute(route RouteInterface) DomainInterface {
	// Convert RouteInterface to RouteConfig if possible
	if routeConfig, ok := route.(*RouteConfig); ok {
		d.Routes = append(d.Routes, *routeConfig)
	} else {
		// Create a new RouteConfig from the RouteInterface
		routeConfig := RouteConfig{
			Name:             route.GetName(),
			Method:           route.GetMethod(),
			Path:             route.GetPath(),
			Handler:          route.GetHandler(),
			ErrorHandler:     route.GetErrorHandler(),
			HTMLHandler:      route.GetHTMLHandler(),
			JSONHandler:      route.GetJSONHandler(),
			CSSHandler:       route.GetCSSHandler(),
			XMLHandler:       route.GetXMLHandler(),
			TextHandler:      route.GetTextHandler(),
			BeforeMiddleware: route.GetBeforeMiddlewares(),
			AfterMiddleware:  route.GetAfterMiddlewares(),
			Metadata:         make(map[string]interface{}),
		}
		d.Routes = append(d.Routes, routeConfig)
	}
	return d
}

// AddRoutes adds multiple routes to this domain and returns the domain for method chaining
func (d *DomainConfig) AddRoutes(routes []RouteInterface) DomainInterface {
	for _, route := range routes {
		d.AddRoute(route)
	}
	return d
}

// GetRoutes returns all routes that belong to this domain
func (d *DomainConfig) GetRoutes() []RouteInterface {
	routes := make([]RouteInterface, len(d.Routes))
	for i, route := range d.Routes {
		routes[i] = &route
	}
	return routes
}

// AddGroup adds a group to this domain and returns the domain for method chaining
func (d *DomainConfig) AddGroup(group GroupInterface) DomainInterface {
	// Convert GroupInterface to GroupConfig if possible
	if groupConfig, ok := group.(*GroupConfig); ok {
		d.Groups = append(d.Groups, *groupConfig)
	} else {
		// Create a new GroupConfig from the GroupInterface
		groupConfig := GroupConfig{
			Prefix:           group.GetPrefix(),
			Routes:           []RouteConfig{},
			Groups:           []GroupConfig{},
			BeforeMiddleware: group.GetBeforeMiddlewares(),
			AfterMiddleware:  group.GetAfterMiddlewares(),
		}
		// Convert routes
		for _, route := range group.GetRoutes() {
			groupConfig.AddRoute(route)
		}
		// Convert nested groups
		for _, nestedGroup := range group.GetGroups() {
			groupConfig.AddGroup(nestedGroup)
		}
		d.Groups = append(d.Groups, groupConfig)
	}
	return d
}

// AddGroups adds multiple groups to this domain and returns the domain for method chaining
func (d *DomainConfig) AddGroups(groups []GroupInterface) DomainInterface {
	for _, group := range groups {
		d.AddGroup(group)
	}
	return d
}

// GetGroups returns all groups that belong to this domain
func (d *DomainConfig) GetGroups() []GroupInterface {
	groups := make([]GroupInterface, len(d.Groups))
	for i, group := range d.Groups {
		groups[i] = &group
	}
	return groups
}

// AddBeforeMiddlewares adds middleware functions to be executed before any route handler in this domain
func (d *DomainConfig) AddBeforeMiddlewares(middleware []MiddlewareInterface) DomainInterface {
	d.BeforeMiddleware = append(d.BeforeMiddleware, middleware...)
	return d
}

// GetBeforeMiddlewares returns all middleware functions that will be executed before any route handler in this domain
func (d *DomainConfig) GetBeforeMiddlewares() []MiddlewareInterface {
	return d.BeforeMiddleware
}

// AddAfterMiddlewares adds middleware functions to be executed after any route handler in this domain
func (d *DomainConfig) AddAfterMiddlewares(middleware []MiddlewareInterface) DomainInterface {
	d.AfterMiddleware = append(d.AfterMiddleware, middleware...)
	return d
}

// GetAfterMiddlewares returns all middleware functions that will be executed after any route handler in this domain
func (d *DomainConfig) GetAfterMiddlewares() []MiddlewareInterface {
	return d.AfterMiddleware
}

// Match checks if the given host matches any of this domain's patterns
func (d *DomainConfig) Match(host string) bool {
	for _, pattern := range d.Patterns {
		if pattern == host {
			return true
		}
		// Simple wildcard matching - could be enhanced with more sophisticated pattern matching
		if pattern == "*" {
			return true
		}
	}
	return false
}
