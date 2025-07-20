package rtr

import (
	"encoding/json"
	"fmt"
)

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
	Status      string        `json:"status,omitempty"`
	Hosts       []string      `json:"hosts,omitempty"`
	Items       []ItemInterface `json:"items,omitempty"` // Can contain both Groups and Routes in sequence
	Middlewares []string      `json:"middlewares,omitempty"`
	Name        string        `json:"name,omitempty"`
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

// MarshalJSON implements the json.Marshaler interface for Domain
func (d Domain) MarshalJSON() ([]byte, error) {
	// Create an alias to avoid infinite recursion
	type Alias Domain
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Alias
	}{
		Type: TypeDomain,
		Alias: (*Alias)(&d),
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface for Domain
func (d *Domain) UnmarshalJSON(data []byte) error {
	type Alias Domain
	aux := &struct {
		*Alias
		Items []json.RawMessage `json:"items,omitempty"`
	}{
		Alias: (*Alias)(d),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Process Items
	for _, itemData := range aux.Items {
		// First unmarshal just the type field to determine the concrete type
		var typeOnly struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(itemData, &typeOnly); err != nil {
			return err
		}

		switch typeOnly.Type {
		case TypeRoute:
			var route Route
			if err := json.Unmarshal(itemData, &route); err != nil {
				return err
			}
			d.Items = append(d.Items, route)
		case TypeGroup:
			var group Group
			if err := json.Unmarshal(itemData, &group); err != nil {
				return err
			}
			d.Items = append(d.Items, group)
		default:
			return fmt.Errorf("unknown item type: %s", typeOnly.Type)
		}
	}

	return nil
}

// Group represents a group of routes
type Group struct {
	Status      string   `json:"status,omitempty"`
	Prefix      string   `json:"prefix,omitempty"`
	Routes      []Route  `json:"routes,omitempty"`
	Middlewares []string `json:"middlewares,omitempty"`
	Name        string   `json:"name,omitempty"`
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

// MarshalJSON implements the json.Marshaler interface for Group
func (g Group) MarshalJSON() ([]byte, error) {
	// Create an alias to avoid infinite recursion
	type Alias Group
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  TypeGroup,
		Alias: (*Alias)(&g),
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface for Group
func (g *Group) UnmarshalJSON(data []byte) error {
	type Alias Group
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(g),
	}

	return json.Unmarshal(data, &aux)
}

// Route represents a single route configuration
type Route struct {
	Status       string   `json:"status,omitempty"`
	Path         string   `json:"path,omitempty"`
	Method       string   `json:"method,omitempty"`
	Handler      string   `json:"handler,omitempty"`      // Standard HTTP handler reference
	HTMLHandler  string   `json:"htmlHandler,omitempty"`  // HTML handler reference
	JSONHandler  string   `json:"jsonHandler,omitempty"`  // JSON handler reference
	CSSHandler   string   `json:"cssHandler,omitempty"`   // CSS handler reference
	XMLHandler   string   `json:"xmlHandler,omitempty"`   // XML handler reference
	TextHandler  string   `json:"textHandler,omitempty"`  // Text handler reference
	JSHandler    string   `json:"jsHandler,omitempty"`    // JavaScript handler reference
	ErrorHandler string   `json:"errorHandler,omitempty"` // Error handler reference
	Name         string   `json:"name,omitempty"`
	Middlewares  []string `json:"middlewares,omitempty"`
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

// MarshalJSON implements the json.Marshaler interface for Route
func (r Route) MarshalJSON() ([]byte, error) {
	// Create an alias to avoid infinite recursion
	type Alias Route
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  TypeRoute,
		Alias: (*Alias)(&r),
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface for Route
func (r *Route) UnmarshalJSON(data []byte) error {
	type Alias Route
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	return json.Unmarshal(data, &aux)
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
