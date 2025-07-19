package rtr

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

// MiddlewareInfo represents middleware information for display purposes
type MiddlewareInfo struct {
	Name string
	Func Middleware
}

// List displays the router's routes, groups, domains, and middleware in a formatted table
// This provides an easy way to visualize the router's configuration for debugging and documentation
func (r *routerImpl) List() {
	r.listMiddlewares()
	r.listDomains()
	r.listRoutes()
	r.listGroups()
}

// listMiddlewares displays global middleware in a formatted table
func (r *routerImpl) listMiddlewares() {
	beforeMiddlewares := r.GetBeforeMiddlewares()
	afterMiddlewares := r.GetAfterMiddlewares()
	
	if len(beforeMiddlewares) == 0 && len(afterMiddlewares) == 0 {
		return
	}

	// Before middlewares table
	if len(beforeMiddlewares) > 0 {
		tableMiddleware := table.NewWriter()
		tableMiddleware.AppendHeader(table.Row{"#", "Middleware Name", "Type"})
		
		for index, middleware := range beforeMiddlewares {
			name := GetMiddlewareName(middleware)
			tableMiddleware.AppendRow(table.Row{index + 1, name, "Before"})
		}
		
		tableMiddleware.SetIndexColumn(1)
		tableMiddleware.SetTitle(fmt.Sprintf("GLOBAL BEFORE MIDDLEWARE LIST (TOTAL: %d)", len(beforeMiddlewares)))
		fmt.Println(tableMiddleware.Render())
		fmt.Println()
	}

	// After middlewares table
	if len(afterMiddlewares) > 0 {
		tableMiddleware := table.NewWriter()
		tableMiddleware.AppendHeader(table.Row{"#", "Middleware Name", "Type"})
		
		for index, middleware := range afterMiddlewares {
			name := GetMiddlewareName(middleware)
			tableMiddleware.AppendRow(table.Row{index + 1, name, "After"})
		}
		
		tableMiddleware.SetIndexColumn(1)
		tableMiddleware.SetTitle(fmt.Sprintf("GLOBAL AFTER MIDDLEWARE LIST (TOTAL: %d)", len(afterMiddlewares)))
		fmt.Println(tableMiddleware.Render())
		fmt.Println()
	}
}

// listDomains displays domains and their routes in a formatted table
func (r *routerImpl) listDomains() {
	domains := r.GetDomains()
	if len(domains) == 0 {
		return
	}

	for _, domain := range domains {
		tableDomain := table.NewWriter()
		tableDomain.AppendHeader(table.Row{"#", "Route Path", "Method", "Route Name", "Middleware List"})
		
		routes := domain.GetRoutes()
		for index, route := range routes {
			method := route.GetMethod()
			if method == "" {
				method = "ALL"
			}
			
			// For domains, we just show the route path without domain prefix
			path := route.GetPath()
			name := route.GetName()
			if name == "" {
				name = "unnamed"
			}
			
			middlewareNames := GetRouteMiddlewareNames(route)
			tableDomain.AppendRow(table.Row{index + 1, path, method, name, middlewareNames})
		}
		
		// Add groups within domain
		groups := domain.GetGroups()
		routeIndex := len(routes)
		for _, group := range groups {
			r.addGroupRoutesToTable(tableDomain, group, "", &routeIndex)
		}
		
		tableDomain.SetIndexColumn(1)
		patterns := strings.Join(domain.GetPatterns(), ", ")
		tableDomain.SetTitle(fmt.Sprintf("DOMAIN ROUTES [%s] (TOTAL: %d)", patterns, routeIndex))
		fmt.Println(tableDomain.Render())
		fmt.Println()
	}
}

// listRoutes displays direct routes in a formatted table
func (r *routerImpl) listRoutes() {
	routes := r.GetRoutes()
	if len(routes) == 0 {
		return
	}

	tableRoutes := table.NewWriter()
	tableRoutes.AppendHeader(table.Row{"#", "Route Path", "Method", "Route Name", "Middleware List"})
	
	for index, route := range routes {
		method := route.GetMethod()
		if method == "" {
			method = "ALL"
		}
		
		path := r.GetPrefix() + route.GetPath()
		name := route.GetName()
		if name == "" {
			name = "unnamed"
		}
		
		middlewareNames := GetRouteMiddlewareNames(route)
		tableRoutes.AppendRow(table.Row{index + 1, path, method, name, middlewareNames})
	}
	
	tableRoutes.SetIndexColumn(1)
	tableRoutes.SetTitle(fmt.Sprintf("DIRECT ROUTES LIST (TOTAL: %d)", len(routes)))
	fmt.Println(tableRoutes.Render())
	fmt.Println()
}

// listGroups displays route groups and their routes in a formatted table
func (r *routerImpl) listGroups() {
	groups := r.GetGroups()
	if len(groups) == 0 {
		return
	}

	for _, group := range groups {
		tableGroup := table.NewWriter()
		tableGroup.AppendHeader(table.Row{"#", "Route Path", "Method", "Route Name", "Middleware List"})
		
		routeIndex := 0
		r.addGroupRoutesToTable(tableGroup, group, r.GetPrefix(), &routeIndex)
		
		tableGroup.SetIndexColumn(1)
		tableGroup.SetTitle(fmt.Sprintf("GROUP ROUTES [%s] (TOTAL: %d)", group.GetPrefix(), routeIndex))
		fmt.Println(tableGroup.Render())
		fmt.Println()
	}
}

// addGroupRoutesToTable recursively adds group routes to a table
func (r *routerImpl) addGroupRoutesToTable(t table.Writer, group GroupInterface, parentPath string, routeIndex *int) {
	groupPath := parentPath + group.GetPrefix()
	
	// Add direct routes in this group
	routes := group.GetRoutes()
	for _, route := range routes {
		*routeIndex++
		method := route.GetMethod()
		if method == "" {
			method = "ALL"
		}
		
		path := groupPath + route.GetPath()
		name := route.GetName()
		if name == "" {
			name = "unnamed"
		}
		
		// Combine group and route middleware
		middlewareNames := getCombinedMiddlewareNames(group, route)
		t.AppendRow(table.Row{*routeIndex, path, method, name, middlewareNames})
	}
	
	// Recursively add nested groups
	nestedGroups := group.GetGroups()
	for _, nestedGroup := range nestedGroups {
		r.addGroupRoutesToTable(t, nestedGroup, groupPath, routeIndex)
	}
}

// GetMiddlewareName attempts to get a readable name for a middleware function
func GetMiddlewareName(middleware Middleware) string {
	if middleware == nil {
		return "nil"
	}
	
	// Get the function name using reflection
	funcValue := reflect.ValueOf(middleware)
	if funcValue.Kind() != reflect.Func {
		return "invalid"
	}
	
	// Get the function pointer and try to get its name
	funcPtr := runtime.FuncForPC(funcValue.Pointer())
	if funcPtr == nil {
		return "anonymous"
	}
	
	fullName := funcPtr.Name()
	
	// Clean up the function name for better readability
	parts := strings.Split(fullName, ".")
	if len(parts) > 0 {
		name := parts[len(parts)-1]
		// Remove common suffixes
		name = strings.TrimSuffix(name, "-fm")
		name = strings.TrimSuffix(name, ".func1")
		return name
	}
	
	return "unnamed"
}

// GetRouteMiddlewareNames gets middleware names for a route
func GetRouteMiddlewareNames(route RouteInterface) []string {
	var names []string
	
	// Before middlewares
	for _, middleware := range route.GetBeforeMiddlewares() {
		name := GetMiddlewareName(middleware)
		names = append(names, name+" (before)")
	}
	
	// After middlewares
	for _, middleware := range route.GetAfterMiddlewares() {
		name := GetMiddlewareName(middleware)
		names = append(names, name+" (after)")
	}
	
	if len(names) == 0 {
		return []string{"none"}
	}
	
	return names
}

// getCombinedMiddlewareNames gets combined middleware names for group and route
func getCombinedMiddlewareNames(group GroupInterface, route RouteInterface) []string {
	var names []string
	
	// Group before middlewares
	for _, middleware := range group.GetBeforeMiddlewares() {
		name := GetMiddlewareName(middleware)
		names = append(names, name+" (group-before)")
	}
	
	// Route before middlewares
	for _, middleware := range route.GetBeforeMiddlewares() {
		name := GetMiddlewareName(middleware)
		names = append(names, name+" (route-before)")
	}
	
	// Route after middlewares
	for _, middleware := range route.GetAfterMiddlewares() {
		name := GetMiddlewareName(middleware)
		names = append(names, name+" (route-after)")
	}
	
	// Group after middlewares
	for _, middleware := range group.GetAfterMiddlewares() {
		name := GetMiddlewareName(middleware)
		names = append(names, name+" (group-after)")
	}
	
	if len(names) == 0 {
		return []string{"none"}
	}
	
	return names
}
