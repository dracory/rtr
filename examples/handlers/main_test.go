package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	rtr "github.com/dracory/rtr"
)

func TestHandlerExamples(t *testing.T) {
	// Create a new router with all the routes
	r := rtr.NewRouter()

	// Add all the routes (same as in main.go)
	setupRoutes(r)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedType   string
		expectedBody   string
	}{
		{
			name:           "Root page",
			method:         "GET",
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedType:   "text/html; charset=utf-8",
			expectedBody:   "Handler Types Example",
		},
		{
			name:           "Traditional handler",
			method:         "GET",
			path:           "/traditional",
			expectedStatus: http.StatusOK,
			expectedType:   "text/html; charset=utf-8",
			expectedBody:   "Welcome to RTR Router!",
		},
		{
			name:           "HTML handler",
			method:         "GET",
			path:           "/html",
			expectedStatus: http.StatusOK,
			expectedType:   "text/html; charset=utf-8",
			expectedBody:   "HTML Handler",
		},
		{
			name:           "JSON handler - users",
			method:         "GET",
			path:           "/api/users",
			expectedStatus: http.StatusOK,
			expectedType:   "application/json",
			expectedBody:   "Alice",
		},
		{
			name:           "JSON handler - status",
			method:         "GET",
			path:           "/api/status",
			expectedStatus: http.StatusOK,
			expectedType:   "application/json",
			expectedBody:   "rtr-router",
		},
		{
			name:           "CSS handler",
			method:         "GET",
			path:           "/styles.css",
			expectedStatus: http.StatusOK,
			expectedType:   "text/css",
			expectedBody:   "font-family: Arial",
		},
		{
			name:           "XML handler",
			method:         "GET",
			path:           "/api/data.xml",
			expectedStatus: http.StatusOK,
			expectedType:   "application/xml",
			expectedBody:   "<?xml version",
		},
		{
			name:           "Text handler",
			method:         "GET",
			path:           "/robots.txt",
			expectedStatus: http.StatusOK,
			expectedType:   "text/plain; charset=utf-8",
			expectedBody:   "User-agent:",
		},
		{
			name:           "Priority demo (HTML wins)",
			method:         "GET",
			path:           "/priority-demo",
			expectedStatus: http.StatusOK,
			expectedType:   "text/html; charset=utf-8",
			expectedBody:   "HTML Handler",
		},
		{
			name:           "HTML with parameters",
			method:         "GET",
			path:           "/user/123",
			expectedStatus: http.StatusOK,
			expectedType:   "text/html; charset=utf-8",
			expectedBody:   "User ID: 123",
		},
		{
			name:           "JSON with parameters",
			method:         "GET",
			path:           "/api/user/456",
			expectedStatus: http.StatusOK,
			expectedType:   "application/json",
			expectedBody:   "\"id\": \"456\"",
		},
		{
			name:           "JS handler",
			method:         "GET",
			path:           "/script.js",
			expectedStatus: http.StatusOK,
			expectedType:   "application/javascript",
			expectedBody:   "console.log",
		},
		{
			name:           "String handler",
			method:         "GET",
			path:           "/raw",
			expectedStatus: http.StatusOK,
			expectedType:   "", // StringHandler doesn't set Content-Type automatically
			expectedBody:   "raw string response",
		},
		{
			name:           "Error handler - success case",
			method:         "GET",
			path:           "/error-demo",
			expectedStatus: http.StatusOK,
			expectedType:   "application/json",
			expectedBody:   "Success! No error occurred",
		},
		{
			name:           "Error handler - error case",
			method:         "GET",
			path:           "/error-demo?fail=true",
			expectedStatus: http.StatusInternalServerError,
			expectedType:   "application/json", // Set by the error handler
			expectedBody:   "simulated internal server error",
		},
		{
			name:           "Error handler - 404 case",
			method:         "GET",
			path:           "/not-found-demo",
			expectedStatus: http.StatusNotFound,
			expectedType:   "application/json",
			expectedBody:   "Resource not found",
		},
		{
			name:           "ToHandler function demo",
			method:         "GET",
			path:           "/to-handler-demo",
			expectedStatus: http.StatusOK,
			expectedType:   "text/html; charset=utf-8",
			expectedBody:   "ToHandler Function Demo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check content type (skip for StringHandler which doesn't set headers automatically)
			if tt.expectedType != "" {
				contentType := w.Header().Get("Content-Type")
				if contentType != tt.expectedType {
					t.Errorf("Expected Content-Type %s, got %s", tt.expectedType, contentType)
				}
			}

			// Check body contains expected content
			body := w.Body.String()
			if !strings.Contains(body, tt.expectedBody) {
				t.Errorf("Expected body to contain %q, got %q", tt.expectedBody, body)
			}
		})
	}
}

// setupRoutes adds all the routes to the router (extracted from main function for testing)
func setupRoutes(r rtr.RouterInterface) {
	// Root route
	r.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/").
		SetHandler(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte("Handler Types Example"))
		}))

	// Traditional Handler
	r.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/traditional").
		SetHandler(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte("<h1>Welcome to RTR Router!</h1><p>This is a traditional handler.</p>"))
		}))

	// HTMLHandler
	r.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/html").
		SetHTMLHandler(func(w http.ResponseWriter, req *http.Request) string {
			return "<h1>HTML Handler</h1><p>Generated by HTMLHandler</p>"
		}))

	// JSONHandler - users
	r.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/api/users").
		SetJSONHandler(func(w http.ResponseWriter, req *http.Request) string {
			return `{"users": [{"name": "Alice"}]}`
		}))

	// JSONHandler - status
	r.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/api/status").
		SetJSONHandler(func(w http.ResponseWriter, req *http.Request) string {
			return `{"server": "rtr-router"}`
		}))

	// CSSHandler
	r.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/styles.css").
		SetCSSHandler(func(w http.ResponseWriter, req *http.Request) string {
			return "body { font-family: Arial; }"
		}))

	// XMLHandler
	r.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/api/data.xml").
		SetXMLHandler(func(w http.ResponseWriter, req *http.Request) string {
			return `<?xml version="1.0"?><data></data>`
		}))

	// TextHandler
	r.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/robots.txt").
		SetTextHandler(func(w http.ResponseWriter, req *http.Request) string {
			return "User-agent: *"
		}))

	// Priority demo
	r.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/priority-demo").
		SetHTMLHandler(func(w http.ResponseWriter, req *http.Request) string {
			return "<h1>HTML Handler</h1>"
		}).
		SetJSONHandler(func(w http.ResponseWriter, req *http.Request) string {
			return `{"message": "JSON"}`
		}))

	// HTML with parameters
	r.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/user/:id").
		SetHTMLHandler(func(w http.ResponseWriter, req *http.Request) string {
			userID := rtr.MustGetParam(req, "id")
			return "<p>User ID: " + userID + "</p>"
		}))

	// JSON with parameters
	r.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/api/user/:id").
		SetJSONHandler(func(w http.ResponseWriter, req *http.Request) string {
			userID := rtr.MustGetParam(req, "id")
			return `{"id": "` + userID + `"}`
		}))

	// JSHandler
	r.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/script.js").
		SetJSHandler(func(w http.ResponseWriter, req *http.Request) string {
			return "console.log('Hello from RTR Router!');"
		}))

	// StringHandler
	r.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/raw").
		SetStringHandler(func(w http.ResponseWriter, req *http.Request) string {
			w.Header().Set("X-Custom-Header", "Raw Response")
			return "This is a raw string response without automatic Content-Type headers."
		}))

	// ErrorHandler - success case
	r.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/error-demo").
		SetErrorHandler(func(w http.ResponseWriter, req *http.Request) error {
			if req.URL.Query().Get("fail") == "true" {
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"error": "simulated internal server error"}`))
				return fmt.Errorf("simulated internal server error")
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"message": "Success! No error occurred.", "status": "ok"}`))
			return nil
		}))

	// ErrorHandler - 404 case
	r.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/not-found-demo").
		SetErrorHandler(func(w http.ResponseWriter, req *http.Request) error {
			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"error": "Resource not found", "code": 404}`))
			return fmt.Errorf("resource not found")
		}))

	// ToHandler function demo
	r.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/to-handler-demo").
		SetHandler(rtr.ToStdHandler(func(w http.ResponseWriter, req *http.Request) string {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("X-Converted-Handler", "true")
			return `<h1>ToHandler Function Demo</h1><p>This was converted using rtr.ToHandler()!</p>`
		})))
}
