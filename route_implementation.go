package rtr

import (
	"io/fs"
	"net/http"
	"strings"
)

// RouteImpl implements the RouteInterface
// It represents a single route definition with its associated properties and middleware.
// A route defines how a specific HTTP request should be handled, including the HTTP method,
// URL path, handler function, and any middleware that should be applied before or after the handler.
type routeImpl struct {
	// method specifies the HTTP method for this route (e.g., "GET", "POST", "PUT", "DELETE")
	method string

	// path specifies the URL path pattern for this route (e.g., "/users", "/api/products")
	path string

	// paramNames stores the names of path parameters in order of appearance
	paramNames []string

	// hasOptionalParams indicates if the route contains any optional parameters
	hasOptionalParams bool

	// handler is the function that will be called when this route is matched
	handler StdHandler

	// stringHandler is the simple string handler function that returns a string without setting headers
	stringHandler StringHandler

	// htmlHandler is the HTML handler function that returns HTML string
	htmlHandler HTMLHandler

	// jsonHandler is the JSON handler function that returns JSON string
	jsonHandler JSONHandler

	// cssHandler is the CSS handler function that returns CSS string
	cssHandler CSSHandler

	// xmlHandler is the XML handler function that returns XML string
	xmlHandler XMLHandler

	// textHandler is the text handler function that returns plain text string
	textHandler TextHandler

	// jsHandler is the JavaScript handler function that returns JavaScript string
	jsHandler JSHandler

	// staticHandler is the static handler function that serves static files
	staticHandler StaticHandler

	// errorHandler is the error handler function that returns error message and status code
	errorHandler ErrorHandler

	// controller is the standard controller that implements ControllerInterface
	controller ControllerInterface

	// htmlController is the HTML controller that implements HTMLControllerInterface
	htmlController HTMLControllerInterface

	// jsonController is the JSON controller that implements JSONControllerInterface
	jsonController JSONControllerInterface

	// textController is the text controller that implements TextControllerInterface
	textController TextControllerInterface

	// name is an optional identifier for this route, useful for route generation and debugging
	name string

	// beforeMiddlewares are middleware that will be executed before the route handler
	beforeMiddlewares []MiddlewareInterface

	// afterMiddlewares are middleware that will be executed after the route handler
	afterMiddlewares []MiddlewareInterface
}

var _ RouteInterface = (*routeImpl)(nil)

// GetMethod returns the HTTP method associated with this route.
// Returns the string representation of the HTTP method (e.g., "GET", "POST").
func (r *routeImpl) GetMethod() string {
	return r.method
}

// SetMethod sets the HTTP method for this route.
// This method supports method chaining by returning the RouteInterface.
// The method parameter should be a valid HTTP method string (e.g., "GET", "POST").
func (r *routeImpl) SetMethod(method string) RouteInterface {
	r.method = method
	return r
}

// normalizeBracesToColon converts brace-style parameters to colon syntax so the
// router can parse one consistent format.
// Examples:
//
//	/users/{id}      -> /users/:id
//	/users/{id?}     -> /users/:id?
//	/files/{path...} -> /files/:path...
func normalizeBracesToColon(path string) string {
	if path == "" {
		return path
	}
	segments := strings.Split(path, "/")
	for i, seg := range segments {
		if normalized, ok := normalizeBraceSegment(seg); ok {
			segments[i] = normalized
		}
	}
	return strings.Join(segments, "/")
}

// normalizeBraceSegment converts a single brace-wrapped segment to colon syntax.
// Returns the normalized segment and true if normalization occurred; otherwise the
// original segment and false.
func normalizeBraceSegment(seg string) (string, bool) {
	if !(strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}")) {
		return seg, false
	}
	inner := seg[1 : len(seg)-1]
	if inner == "" {
		// Invalid empty name; leave segment unchanged
		return seg, false
	}
	suffix := ""
	if strings.HasSuffix(inner, "...") {
		suffix = "..."
		inner = strings.TrimSuffix(inner, "...")
	} else if strings.HasSuffix(inner, "?") {
		suffix = "?"
		inner = strings.TrimSuffix(inner, "?")
	}

	// Preserve any optional ('?') or greedy ('...') suffixes by returning as-is
	if idx := strings.Index(inner, ":"); idx != -1 {
		inner = inner[:idx]
	}

	return ":" + inner + suffix, true
}

// GetPath returns the URL path pattern associated with this route.
// Returns the string representation of the path (e.g., "/users", "/api/products").
func (r *routeImpl) GetPath() string {
	return r.path
}

// SetPath sets the URL path pattern for this route and extracts any parameter names.
// This method supports method chaining by returning the RouteInterface.
// The path parameter should be a valid URL path pattern (e.g., "/users/:id").
func (r *routeImpl) SetPath(path string) RouteInterface {
	// Normalize brace-style parameters to colon syntax to keep a single matching engine
	// e.g., {id} -> :id, {id?} -> :id?, {path...} -> :path...
	normalizedPath := normalizeBracesToColon(path)
	r.path = normalizedPath
	r.paramNames = nil
	r.hasOptionalParams = false

	// Extract parameter names from the path
	segments := strings.SplitSeq(normalizedPath, "/")
	for segment := range segments {
		if segment == "" {
			continue
		}
		if segment[0] != ':' {
			continue
		}
		if strings.HasSuffix(segment, "?") {
			r.hasOptionalParams = true
		}

		// Remove the leading ':' and optional/greedy suffixes
		paramName := strings.TrimPrefix(segment, ":")
		paramName = strings.TrimSuffix(paramName, "?")
		paramName = strings.TrimSuffix(paramName, "...")
		r.paramNames = append(r.paramNames, paramName)
	}

	return r
}

// GetHandler returns the handler function associated with this route.
// Returns the Handler function that will be called when this route is matched.
// Implements handler prioritization: Handler > StringHandler > HTMLHandler > JSONHandler > CSSHandler > XMLHandler > TextHandler > ErrorHandler
func (r *routeImpl) GetHandler() StdHandler {
	// Priority 1: Direct Handler
	if r.handler != nil {
		return r.handler
	}

	// Priority 2: StringHandler - convert to standard Handler (no automatic headers)
	if r.stringHandler != nil {
		return ToStdHandler(r.stringHandler)
	}

	// Priority 3: HTMLHandler - convert to standard Handler with HTML headers
	if r.htmlHandler != nil {
		return func(w http.ResponseWriter, req *http.Request) {
			body := r.htmlHandler(w, req)
			HTMLResponse(w, req, body)
		}
	}

	// Priority 4: JSONHandler - convert to standard Handler with JSON headers
	if r.jsonHandler != nil {
		return func(w http.ResponseWriter, req *http.Request) {
			body := r.jsonHandler(w, req)
			JSONResponse(w, req, body)
		}
	}

	// Priority 5: CSSHandler - convert to standard Handler with CSS headers
	if r.cssHandler != nil {
		return func(w http.ResponseWriter, req *http.Request) {
			body := r.cssHandler(w, req)
			CSSResponse(w, req, body)
		}
	}

	// Priority 6: XMLHandler - convert to standard Handler with XML headers
	if r.xmlHandler != nil {
		return func(w http.ResponseWriter, req *http.Request) {
			body := r.xmlHandler(w, req)
			XMLResponse(w, req, body)
		}
	}

	// Priority 7: StaticHandler - convert to standard Handler for static file serving
	if r.staticHandler != nil {
		return func(w http.ResponseWriter, req *http.Request) {
			staticDir := r.staticHandler(w, req)
			urlPrefix := r.path
			if strings.HasSuffix(urlPrefix, "/*") {
				urlPrefix = strings.TrimSuffix(urlPrefix, "/*")
			}
			StaticFileServer(staticDir, urlPrefix)(w, req)
		}
	}

	// Priority 8: TextHandler - convert to standard Handler with Text headers
	if r.textHandler != nil {
		return func(w http.ResponseWriter, req *http.Request) {
			body := r.textHandler(w, req)
			TextResponse(w, req, body)
		}
	}

	// Priority 9: JSHandler - convert to standard Handler with JavaScript headers
	if r.jsHandler != nil {
		return func(w http.ResponseWriter, req *http.Request) {
			body := r.jsHandler(w, req)
			JSResponse(w, req, body)
		}
	}

	// Priority 10: ErrorHandler - convert to standard Handler
	if r.errorHandler != nil {
		return ErrorHandlerToHandler(r.errorHandler)
	}

	// Priority 10: ControllerInterface - convert to standard Handler (no automatic headers)
	if r.controller != nil {
		return r.controller.Handler
	}

	// Priority 11: HTMLControllerInterface - convert to standard Handler with HTML headers
	if r.htmlController != nil {
		return func(w http.ResponseWriter, req *http.Request) {
			body := r.htmlController.Handler(w, req)
			HTMLResponse(w, req, body)
		}
	}

	// Priority 12: JSONControllerInterface - convert to standard Handler with JSON headers
	if r.jsonController != nil {
		return func(w http.ResponseWriter, req *http.Request) {
			body := r.jsonController.Handler(w, req)
			JSONResponse(w, req, body)
		}
	}

	// Priority 13: TextControllerInterface - convert to standard Handler with text headers
	if r.textController != nil {
		return func(w http.ResponseWriter, req *http.Request) {
			body := r.textController.Handler(w, req)
			TextResponse(w, req, body)
		}
	}

	// No handler found
	return nil
}

// SetHandler sets the handler function for this route.
// This method supports method chaining by returning the RouteInterface.
// The handler parameter should be a function that implements the Handler interface.
func (r *routeImpl) SetHandler(handler StdHandler) RouteInterface {
	r.handler = handler
	return r
}

// GetStringHandler returns the string handler function associated with this route.
// Returns the StringHandler function that will be called when this route is matched.
func (r *routeImpl) GetStringHandler() StringHandler {
	return r.stringHandler
}

// SetStringHandler sets the string handler function for this route.
// This method supports method chaining by returning the RouteInterface.
// The handler parameter should be a function that returns a string without setting headers.
func (r *routeImpl) SetStringHandler(handler StringHandler) RouteInterface {
	r.stringHandler = handler
	return r
}

// GetHTMLHandler returns the HTML handler function associated with this route.
// Returns the HTMLHandler function that will be called when this route is matched.
func (r *routeImpl) GetHTMLHandler() HTMLHandler {
	return r.htmlHandler
}

// SetHTMLHandler sets the HTML handler function for this route.
// This method supports method chaining by returning the RouteInterface.
// The handler parameter should be a function that returns HTML string.
func (r *routeImpl) SetHTMLHandler(handler HTMLHandler) RouteInterface {
	r.htmlHandler = handler
	return r
}

// GetJSONHandler returns the JSON handler function associated with this route.
// Returns the JSONHandler function that will be called when this route is matched.
func (r *routeImpl) GetJSONHandler() JSONHandler {
	return r.jsonHandler
}

// SetJSONHandler sets the JSON handler function for this route.
// This method supports method chaining by returning the RouteInterface.
// The handler parameter should be a function that returns JSON string.
func (r *routeImpl) SetJSONHandler(handler JSONHandler) RouteInterface {
	r.jsonHandler = handler
	return r
}

// GetCSSHandler returns the CSS handler function associated with this route.
// Returns the CSSHandler function that will be called when this route is matched.
func (r *routeImpl) GetCSSHandler() CSSHandler {
	return r.cssHandler
}

// SetCSSHandler sets the CSS handler function for this route.
// This method supports method chaining by returning the RouteInterface.
// The handler parameter should be a function that returns CSS string.
func (r *routeImpl) SetCSSHandler(handler CSSHandler) RouteInterface {
	r.cssHandler = handler
	return r
}

// GetXMLHandler returns the XML handler function associated with this route.
// Returns the XMLHandler function that will be called when this route is matched.
func (r *routeImpl) GetXMLHandler() XMLHandler {
	return r.xmlHandler
}

// SetXMLHandler sets the XML handler function for this route.
// This method supports method chaining by returning the RouteInterface.
// The handler parameter should be a function that returns XML string.
func (r *routeImpl) SetXMLHandler(handler XMLHandler) RouteInterface {
	r.xmlHandler = handler
	return r
}

// GetTextHandler returns the text handler function associated with this route.
// Returns the TextHandler function that will be called when this route is matched.
func (r *routeImpl) GetTextHandler() TextHandler {
	return r.textHandler
}

// SetTextHandler sets the text handler function for this route.
// This method supports method chaining by returning the RouteInterface.
// The handler parameter should be a function that returns plain text string.
func (r *routeImpl) SetTextHandler(handler TextHandler) RouteInterface {
	r.textHandler = handler
	return r
}

// GetJSHandler returns the JavaScript handler function associated with this route.
// Returns the JSHandler function that will be called when this route is matched.
func (r *routeImpl) GetJSHandler() JSHandler {
	return r.jsHandler
}

// SetJSHandler sets the JavaScript handler function for this route.
// This method supports method chaining by returning the RouteInterface.
// The handler parameter should be a function that returns JavaScript string.
func (r *routeImpl) SetJSHandler(handler JSHandler) RouteInterface {
	r.jsHandler = handler
	return r
}

// GetStaticHandler returns the static handler function associated with this route.
// Returns the StaticHandler function that will be called when this route is matched.
func (r *routeImpl) GetStaticHandler() StaticHandler {
	return r.staticHandler
}

// SetStaticHandler sets the static handler function for this route.
// This method supports method chaining by returning the RouteInterface.
// The handler parameter should be a function that returns the file path relative to static directory.
func (r *routeImpl) SetStaticHandler(handler StaticHandler) RouteInterface {
	r.staticHandler = handler
	return r
}

// GetErrorHandler returns the error handler function associated with this route.
// Returns the ErrorHandler function that will be called when this route is matched.
func (r *routeImpl) GetErrorHandler() ErrorHandler {
	return r.errorHandler
}

// SetErrorHandler sets the error handler function for this route.
// This method supports method chaining by returning the RouteInterface.
// The handler parameter should be a function that returns an error (nil means no error).
func (r *routeImpl) SetErrorHandler(handler ErrorHandler) RouteInterface {
	r.errorHandler = handler
	return r
}

// GetController returns the controller associated with this route.
// Returns the ControllerInterface implementation that will be called when this route is matched.
func (r *routeImpl) GetController() ControllerInterface {
	return r.controller
}

// SetController sets the controller for this route.
// This method supports method chaining by returning the RouteInterface.
// The controller parameter should implement the ControllerInterface.
func (r *routeImpl) SetController(controller ControllerInterface) RouteInterface {
	r.controller = controller
	return r
}

// GetHTMLController returns the HTML controller associated with this route.
// Returns the HTMLControllerInterface implementation that will be called when this route is matched.
func (r *routeImpl) GetHTMLController() HTMLControllerInterface {
	return r.htmlController
}

// SetHTMLController sets the HTML controller for this route.
// This method supports method chaining by returning the RouteInterface.
// The controller parameter should implement the HTMLControllerInterface.
func (r *routeImpl) SetHTMLController(controller HTMLControllerInterface) RouteInterface {
	r.htmlController = controller
	return r
}

// GetJSONController returns the JSON controller associated with this route.
// Returns the JSONControllerInterface implementation that will be called when this route is matched.
func (r *routeImpl) GetJSONController() JSONControllerInterface {
	return r.jsonController
}

// SetJSONController sets the JSON controller for this route.
// This method supports method chaining by returning the RouteInterface.
// The controller parameter should implement the JSONControllerInterface.
func (r *routeImpl) SetJSONController(controller JSONControllerInterface) RouteInterface {
	r.jsonController = controller
	return r
}

// GetTextController returns the text controller associated with this route.
// Returns the TextControllerInterface implementation that will be called when this route is matched.
func (r *routeImpl) GetTextController() TextControllerInterface {
	return r.textController
}

// SetTextController sets the text controller for this route.
// This method supports method chaining by returning the RouteInterface.
// The controller parameter should implement the TextControllerInterface.
func (r *routeImpl) SetTextController(controller TextControllerInterface) RouteInterface {
	r.textController = controller
	return r
}

// GetName returns the name identifier associated with this route.
// Returns the string name of the route, which may be empty if not set.
func (r *routeImpl) GetName() string {
	return r.name
}

// SetName sets the name identifier for this route.
// This method supports method chaining by returning the RouteInterface.
// The name parameter can be used for route identification and debugging.
func (r *routeImpl) SetName(name string) RouteInterface {
	r.name = name
	return r
}

// AddBeforeMiddlewares adds middleware to be executed before the route handler.
// This method supports method chaining by returning the RouteInterface.
// The middleware parameter should be a slice of MiddlewareInterface implementations.
// These middleware will be executed in the order they are added.
func (r *routeImpl) AddBeforeMiddlewares(middleware []MiddlewareInterface) RouteInterface {
	r.beforeMiddlewares = append(r.beforeMiddlewares, middleware...)
	return r
}

// GetBeforeMiddlewares returns all middleware that will be executed before the route handler.
// Returns a slice of MiddlewareInterface implementations in the order they will be executed.
func (r *routeImpl) GetBeforeMiddlewares() []MiddlewareInterface {
	return r.beforeMiddlewares
}

// AddAfterMiddlewares adds middleware to be executed after the route handler.
// This method supports method chaining by returning the RouteInterface.
// The middleware parameter should be a slice of MiddlewareInterface implementations.
// These middleware will be executed in the order they are added.
func (r *routeImpl) AddAfterMiddlewares(middleware []MiddlewareInterface) RouteInterface {
	r.afterMiddlewares = append(r.afterMiddlewares, middleware...)
	return r
}

// GetAfterMiddlewares returns all middleware that will be executed after the route handler.
// Returns a slice of MiddlewareInterface implementations in the order they will be executed.
func (r *routeImpl) GetAfterMiddlewares() []MiddlewareInterface {
	return r.afterMiddlewares
}

// Get creates a new GET route with the given path and handler
// It is a shortcut method that combines setting the method to GET, path, and handler.
func Get(path string, handler StdHandler) RouteInterface {
	return NewRoute().SetMethod(http.MethodGet).SetPath(path).SetHandler(handler)
}

// Post creates a new POST route with the given path and handler
// It is a shortcut method that combines setting the method to POST, path, and handler.
func Post(path string, handler StdHandler) RouteInterface {
	return NewRoute().SetMethod(http.MethodPost).SetPath(path).SetHandler(handler)
}

// Put creates a new PUT route with the given path and handler
// It is a shortcut method that combines setting the method to PUT, path, and handler.
func Put(path string, handler StdHandler) RouteInterface {
	return NewRoute().SetMethod(http.MethodPut).SetPath(path).SetHandler(handler)
}

// Delete creates a new DELETE route with the given path and handler
// It is a shortcut method that combines setting the method to DELETE, path, and handler.
func Delete(path string, handler StdHandler) RouteInterface {
	return NewRoute().SetMethod(http.MethodDelete).SetPath(path).SetHandler(handler)
}

// GetHTML creates a new GET route with the given path and HTML handler
// It is a shortcut method that combines setting the method to GET, path, and HTML handler.
func GetHTML(path string, handler HTMLHandler) RouteInterface {
	return NewRoute().SetMethod(http.MethodGet).SetPath(path).SetHTMLHandler(handler)
}

// PostHTML creates a new POST route with the given path and HTML handler
// It is a shortcut method that combines setting the method to POST, path, and HTML handler.
func PostHTML(path string, handler HTMLHandler) RouteInterface {
	return NewRoute().SetMethod(http.MethodPost).SetPath(path).SetHTMLHandler(handler)
}

// GetJSON creates a new GET route with the given path and JSON handler
// It is a shortcut method that combines setting the method to GET, path, and JSON handler.
func GetJSON(path string, handler JSONHandler) RouteInterface {
	return NewRoute().SetMethod(http.MethodGet).SetPath(path).SetJSONHandler(handler)
}

// PostJSON creates a new POST route with the given path and JSON handler
// It is a shortcut method that combines setting the method to POST, path, and JSON handler.
func PostJSON(path string, handler JSONHandler) RouteInterface {
	return NewRoute().SetMethod(http.MethodPost).SetPath(path).SetJSONHandler(handler)
}

// GetCSS creates a new GET route with the given path and CSS handler
// It is a shortcut method that combines setting the method to GET, path, and CSS handler.
func GetCSS(path string, handler CSSHandler) RouteInterface {
	return NewRoute().SetMethod(http.MethodGet).SetPath(path).SetCSSHandler(handler)
}

// GetXML creates a new GET route with the given path and XML handler
// It is a shortcut method that combines setting the method to GET, path, and XML handler.
func GetXML(path string, handler XMLHandler) RouteInterface {
	return NewRoute().SetMethod(http.MethodGet).SetPath(path).SetXMLHandler(handler)
}

// GetText creates a new GET route with the given path and text handler
// It is a shortcut method that combines setting the method to GET, path, and text handler.
func GetText(path string, handler TextHandler) RouteInterface {
	return NewRoute().SetMethod(http.MethodGet).SetPath(path).SetTextHandler(handler)
}

// GetStatic creates a new GET route with the given path and static handler
// It is a shortcut method that combines setting the method to GET, path, and static handler.
func GetStatic(path string, handler StaticHandler) RouteInterface {
	return NewRoute().SetMethod(http.MethodGet).SetPath(path).SetStaticHandler(handler)
}

func GetStaticFS(path string, fsys fs.FS) RouteInterface {
	return NewRoute().SetMethod(http.MethodGet).SetPath(path).SetHandler(func(w http.ResponseWriter, r *http.Request) {
		urlPrefix := path
		if strings.HasSuffix(urlPrefix, "/*") {
			urlPrefix = strings.TrimSuffix(urlPrefix, "/*")
		}
		StaticFileServerFS(fsys, urlPrefix)(w, r)
	})
}
