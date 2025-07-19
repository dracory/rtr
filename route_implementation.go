package rtr

import (
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
	handler Handler

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

	// errorHandler is the error handler function that returns error message and status code
	errorHandler ErrorHandler

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

// GetPath returns the URL path pattern associated with this route.
// Returns the string representation of the path (e.g., "/users", "/api/products").
func (r *routeImpl) GetPath() string {
	return r.path
}

// SetPath sets the URL path pattern for this route and extracts any parameter names.
// This method supports method chaining by returning the RouteInterface.
// The path parameter should be a valid URL path pattern (e.g., "/users/:id").
func (r *routeImpl) SetPath(path string) RouteInterface {
	r.path = path
	r.paramNames = nil
	r.hasOptionalParams = false

	// Extract parameter names from the path
	segments := strings.Split(path, "/")
	for _, segment := range segments {
		if len(segment) > 0 && (segment[0] == ':' || (len(segment) > 1 && segment[0] == ':' && segment[1] == '?')) {
			// Remove the leading ':' and optional '?'
			paramName := strings.TrimLeft(segment, ":")
			if strings.HasSuffix(paramName, "?") {
				paramName = strings.TrimSuffix(paramName, "?")
				r.hasOptionalParams = true
			}
			r.paramNames = append(r.paramNames, paramName)
		}
	}

	return r
}

// GetHandler returns the handler function associated with this route.
// Returns the Handler function that will be called when this route is matched.
// Implements handler prioritization: Handler > StringHandler > HTMLHandler > JSONHandler > CSSHandler > XMLHandler > TextHandler > ErrorHandler
func (r *routeImpl) GetHandler() Handler {
	// Priority 1: Direct Handler
	if r.handler != nil {
		return r.handler
	}

	// Priority 2: StringHandler - convert to standard Handler (no automatic headers)
	if r.stringHandler != nil {
		return ToHandler(r.stringHandler)
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

	// Priority 7: TextHandler - convert to standard Handler with Text headers
	if r.textHandler != nil {
		return func(w http.ResponseWriter, req *http.Request) {
			body := r.textHandler(w, req)
			TextResponse(w, req, body)
		}
	}

	// Priority 8: JSHandler - convert to standard Handler with JavaScript headers
	if r.jsHandler != nil {
		return func(w http.ResponseWriter, req *http.Request) {
			body := r.jsHandler(w, req)
			JSResponse(w, req, body)
		}
	}

	// Priority 9: ErrorHandler - convert to standard Handler
	if r.errorHandler != nil {
		return ErrorHandlerToHandler(r.errorHandler)
	}

	// No handler found
	return nil
}

// SetHandler sets the handler function for this route.
// This method supports method chaining by returning the RouteInterface.
// The handler parameter should be a function that implements the Handler interface.
func (r *routeImpl) SetHandler(handler Handler) RouteInterface {
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
func Get(path string, handler Handler) RouteInterface {
	return NewRoute().SetMethod(http.MethodGet).SetPath(path).SetHandler(handler)
}

// Post creates a new POST route with the given path and handler
// It is a shortcut method that combines setting the method to POST, path, and handler.
func Post(path string, handler Handler) RouteInterface {
	return NewRoute().SetMethod(http.MethodPost).SetPath(path).SetHandler(handler)
}

// Put creates a new PUT route with the given path and handler
// It is a shortcut method that combines setting the method to PUT, path, and handler.
func Put(path string, handler Handler) RouteInterface {
	return NewRoute().SetMethod(http.MethodPut).SetPath(path).SetHandler(handler)
}

// Delete creates a new DELETE route with the given path and handler
// It is a shortcut method that combines setting the method to DELETE, path, and handler.
func Delete(path string, handler Handler) RouteInterface {
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
