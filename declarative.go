package rtr

// Status constants for declarative configuration
const (
	StatusEnabled  = "enabled"
	StatusDisabled = "disabled"
)

// Type constants for declarative configuration
const (
	TypeRoute      = "route"
	TypeGroup      = "group"
	TypeDomain     = "domain"
	TypeMiddleware = "middleware"
)

// HTTP method constants for declarative configuration
const (
	MethodGET     = "GET"
	MethodPOST    = "POST"
	MethodPUT     = "PUT"
	MethodPATCH   = "PATCH"
	MethodDELETE  = "DELETE"
	MethodHEAD    = "HEAD"
	MethodOPTIONS = "OPTIONS"
)

type NameInterface interface {
	GetName() string
}

type StatusInterface interface {
	GetStatus() string
}

// ItemInterface represents a common interface for groups and routes
type ItemInterface interface {
	NameInterface
	StatusInterface
	GetMiddlewares() []string
}

var _ ItemInterface = (*Domain)(nil)
var _ ItemInterface = (*Group)(nil)
var _ ItemInterface = (*Route)(nil)

// Domain represents a domain-specific routing configuration
type Domain struct {
	Status      string
	Hosts       []string
	Items       []ItemInterface // Can contain both Groups and Routes in sequence
	Middlewares []string
	Name        string
}

// GetName implements ItemInterface
func (d Domain) GetName() string {
	return d.Name
}

// GetStatus implements ItemInterface
func (d Domain) GetStatus() string {
	return d.Status
}

// GetMiddlewares implements ItemInterface
func (d Domain) GetMiddlewares() []string {
	return d.Middlewares
}

// Group represents a group of routes
type Group struct {
	Status      string
	Prefix      string
	Routes      []Route
	Middlewares []string
	Name        string
}

// GetName implements ItemInterface
func (g Group) GetName() string {
	return g.Name
}

// GetStatus implements ItemInterface
func (g Group) GetStatus() string {
	return g.Status
}

// GetMiddlewares implements ItemInterface
func (g Group) GetMiddlewares() []string {
	return g.Middlewares
}

// Route represents a single route configuration
type Route struct {
	Status      string
	Path        string
	Method      string
	Handler      string // Standard HTTP handler reference
	HTMLHandler  string // HTML handler reference
	JSONHandler  string // JSON handler reference
	CSSHandler   string // CSS handler reference
	XMLHandler   string // XML handler reference
	TextHandler  string // Text handler reference
	JSHandler    string // JavaScript handler reference
	ErrorHandler string // Error handler reference
	Name        string
	Middlewares []string
}

// GetName implements ItemInterface
func (r Route) GetName() string {
	return r.Name
}

// GetStatus implements ItemInterface
func (r Route) GetStatus() string {
	return r.Status
}

// GetMiddlewares implements ItemInterface
func (r Route) GetMiddlewares() []string {
	return r.Middlewares
}

// Middleware represents middleware configuration
type Middleware struct {
	Name    string
	Handler string
	Config  map[string]any
}

// HandlerRegistry maps string names to actual functions
type HandlerRegistry struct {
	routes     map[string]RouteInterface
	middleware map[string]MiddlewareInterface
}

// NewHandlerRegistry creates a new handler registry
func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		routes:     make(map[string]RouteInterface),
		middleware: make(map[string]MiddlewareInterface),
	}
}

// AddRoute adds a route to the registry
func (r *HandlerRegistry) AddRoute(route RouteInterface) {
	r.routes[route.GetName()] = route
}

// AddMiddleware adds middleware to the registry
func (r *HandlerRegistry) AddMiddleware(middleware MiddlewareInterface) {
	r.middleware[middleware.GetName()] = middleware
}

// FindRoute finds a route by name
func (r *HandlerRegistry) FindRoute(name string) RouteInterface {
	route, exists := r.routes[name]
	if !exists {
		return nil
	}
	return route
}

// FindMiddleware finds middleware by name
func (r *HandlerRegistry) FindMiddleware(name string) MiddlewareInterface {
	middleware, exists := r.middleware[name]
	if !exists {
		return nil
	}
	return middleware
}

// RemoveRoute removes a route handler from the registry
func (r *HandlerRegistry) RemoveRoute(name string) {
	delete(r.routes, name)
}

// RemoveMiddleware removes a middleware factory from the registry
func (r *HandlerRegistry) RemoveMiddleware(name string) {
	delete(r.middleware, name)
}
