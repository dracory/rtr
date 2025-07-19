package rtr

// groupImpl implements the GroupInterface
// It represents a group of routes that share common properties such as a URL prefix and middleware.
// Route groups allow for organizing related routes together and applying common middleware to all routes in the group.
// Groups can also be nested to create hierarchical route structures.
type groupImpl struct {
	// prefix is the URL path prefix that will be prepended to all routes in this group
	prefix string

	// routes contains all the routes that belong to this group
	routes []RouteInterface
	// groups contains all the nested groups that belong to this group
	groups []GroupInterface

	// beforeMiddlewares are middleware that will be executed before any route handler in this group
	beforeMiddlewares []MiddlewareInterface
	// afterMiddlewares are middleware that will be executed after any route handler in this group
	afterMiddlewares []MiddlewareInterface
}

var _ GroupInterface = (*groupImpl)(nil)

// GetPrefix returns the URL path prefix associated with this group.
// Returns the string representation of the prefix (e.g., "/api", "/admin").
func (g *groupImpl) GetPrefix() string {
	return g.prefix
}

// SetPrefix sets the URL path prefix for this group.
// This method supports method chaining by returning the GroupInterface.
// The prefix parameter should be a valid URL path prefix.
func (g *groupImpl) SetPrefix(prefix string) GroupInterface {
	g.prefix = prefix
	return g
}

// AddRoute adds a single route to this group.
// This method supports method chaining by returning the GroupInterface.
// The route parameter should be a valid RouteInterface implementation.
func (g *groupImpl) AddRoute(route RouteInterface) GroupInterface {
	g.routes = append(g.routes, route)
	return g
}

// AddRoutes adds multiple routes to this group.
// This method supports method chaining by returning the GroupInterface.
// The routes parameter should be a slice of valid RouteInterface implementations.
func (g *groupImpl) AddRoutes(routes []RouteInterface) GroupInterface {
	g.routes = append(g.routes, routes...)
	return g
}

// GetRoutes returns all routes that belong to this group.
// Returns a slice of RouteInterface implementations.
func (g *groupImpl) GetRoutes() []RouteInterface {
	return g.routes
}

// AddGroup adds a single nested group to this group.
// This method supports method chaining by returning the GroupInterface.
// The group parameter should be a valid GroupInterface implementation.
func (g *groupImpl) AddGroup(group GroupInterface) GroupInterface {
	g.groups = append(g.groups, group)
	return g
}

// AddGroups adds multiple nested groups to this group.
// This method supports method chaining by returning the GroupInterface.
// The groups parameter should be a slice of valid GroupInterface implementations.
func (g *groupImpl) AddGroups(groups []GroupInterface) GroupInterface {
	g.groups = append(g.groups, groups...)
	return g
}

// GetGroups returns all nested groups that belong to this group.
// Returns a slice of GroupInterface implementations.
func (g *groupImpl) GetGroups() []GroupInterface {
	return g.groups
}

// AddBeforeMiddlewares adds middleware to be executed before any route handler in this group.
// This method supports method chaining by returning the GroupInterface.
// The middleware parameter should be a slice of MiddlewareInterface.
// These middleware will be executed in the order they are added.
func (g *groupImpl) AddBeforeMiddlewares(middleware []MiddlewareInterface) GroupInterface {
	g.beforeMiddlewares = append(g.beforeMiddlewares, middleware...)
	return g
}

// GetBeforeMiddlewares returns all middleware that will be executed before any route handler in this group.
// Returns a slice of MiddlewareInterface in the order they will be executed.
func (g *groupImpl) GetBeforeMiddlewares() []MiddlewareInterface {
	return g.beforeMiddlewares
}

// AddAfterMiddlewares adds middleware to be executed after any route handler in this group.
// This method supports method chaining by returning the GroupInterface.
// The middleware parameter should be a slice of MiddlewareInterface.
// These middleware will be executed in the order they are added.
func (g *groupImpl) AddAfterMiddlewares(middleware []MiddlewareInterface) GroupInterface {
	g.afterMiddlewares = append(g.afterMiddlewares, middleware...)
	return g
}

// GetAfterMiddlewares returns all middleware that will be executed after any route handler in this group.
// Returns a slice of MiddlewareInterface in the order they will be executed.
func (g *groupImpl) GetAfterMiddlewares() []MiddlewareInterface {
	return g.afterMiddlewares
}

// GetHandler returns nil as groups do not have a single handler.
// This method is implemented to satisfy the RouteInterface but is not used for groups.
func (g *groupImpl) GetHandler() Handler {
	return nil // Groups do not have a single handler
}
