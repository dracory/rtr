package rtr_test

// This file exercises handler-based APIs and response helpers:
// - Handler types (HTML/JSON/CSS/XML/Text/Static) and their content-types
// - Handler priority rules (standard handler vs HTML/JSON/Static)
// - Response helper functions and ToHandler helpers
// Controllers are covered separately in controller_interfaces_test.go.

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/dracory/rtr"
)

func TestHTMLHandler(t *testing.T) {
	route := rtr.GetHTML("/test", func(w http.ResponseWriter, r *http.Request) string {
		return "<h1>Test HTML</h1>"
	})

	handler := route.GetHandler()
	if handler == nil {
		t.Fatal("Expected handler to be non-nil")
	}

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected Content-Type 'text/html; charset=utf-8', got '%s'", contentType)
	}

	body := w.Body.String()
	if body != "<h1>Test HTML</h1>" {
		t.Errorf("Expected body '<h1>Test HTML</h1>', got '%s'", body)
	}
}

func TestJSONHandler(t *testing.T) {
	route := rtr.GetJSON("/api/test", func(w http.ResponseWriter, r *http.Request) string {
		return `{"message": "test"}`
	})

	handler := route.GetHandler()
	if handler == nil {
		t.Fatal("Expected handler to be non-nil")
	}

	req := httptest.NewRequest("GET", "/api/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	body := w.Body.String()
	if body != `{"message": "test"}` {
		t.Errorf("Expected body '{\"message\": \"test\"}', got '%s'", body)
	}
}

func TestCSSHandler(t *testing.T) {
	route := rtr.GetCSS("/styles.css", func(w http.ResponseWriter, r *http.Request) string {
		return "body { color: red; }"
	})

	handler := route.GetHandler()
	if handler == nil {
		t.Fatal("Expected handler to be non-nil")
	}

	req := httptest.NewRequest("GET", "/styles.css", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/css" {
		t.Errorf("Expected Content-Type 'text/css', got '%s'", contentType)
	}

	body := w.Body.String()
	if body != "body { color: red; }" {
		t.Errorf("Expected body 'body { color: red; }', got '%s'", body)
	}
}

func TestXMLHandler(t *testing.T) {
	route := rtr.GetXML("/data.xml", func(w http.ResponseWriter, r *http.Request) string {
		return `<?xml version="1.0"?><root><item>test</item></root>`
	})

	handler := route.GetHandler()
	if handler == nil {
		t.Fatal("Expected handler to be non-nil")
	}

	req := httptest.NewRequest("GET", "/data.xml", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/xml" {
		t.Errorf("Expected Content-Type 'application/xml', got '%s'", contentType)
	}

	body := w.Body.String()
	expected := `<?xml version="1.0"?><root><item>test</item></root>`
	if body != expected {
		t.Errorf("Expected body '%s', got '%s'", expected, body)
	}
}

func TestTextHandler(t *testing.T) {
	route := rtr.GetText("/robots.txt", func(w http.ResponseWriter, r *http.Request) string {
		return "User-agent: *\nDisallow: /"
	})

	handler := route.GetHandler()
	if handler == nil {
		t.Fatal("Expected handler to be non-nil")
	}

	req := httptest.NewRequest("GET", "/robots.txt", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/plain; charset=utf-8" {
		t.Errorf("Expected Content-Type 'text/plain; charset=utf-8', got '%s'", contentType)
	}

	body := w.Body.String()
	expected := "User-agent: *\nDisallow: /"
	if body != expected {
		t.Errorf("Expected body '%s', got '%s'", expected, body)
	}
}

func TestHandlerPriority(t *testing.T) {
	// Test that Handler takes priority over HTMLHandler
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/test").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Standard Handler"))
		}).
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			return "<h1>HTML Handler</h1>"
		})

	handler := route.GetHandler()
	if handler == nil {
		t.Fatal("Expected handler to be non-nil")
	}

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	body := w.Body.String()
	if body != "Standard Handler" {
		t.Errorf("Expected 'Standard Handler' (Handler should take priority), got '%s'", body)
	}
}

func TestHTMLHandlerPriorityOverJSON(t *testing.T) {
	// Test that HTMLHandler takes priority over JSONHandler
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/test").
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			return "<h1>HTML Handler</h1>"
		}).
		SetJSONHandler(func(w http.ResponseWriter, r *http.Request) string {
			return `{"message": "JSON Handler"}`
		})

	handler := route.GetHandler()
	if handler == nil {
		t.Fatal("Expected handler to be non-nil")
	}

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	body := w.Body.String()
	if body != "<h1>HTML Handler</h1>" {
		t.Errorf("Expected '<h1>HTML Handler</h1>' (HTMLHandler should take priority), got '%s'", body)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected Content-Type 'text/html; charset=utf-8', got '%s'", contentType)
	}
}

func TestNoHandlerReturnsNil(t *testing.T) {
	route := rtr.NewRoute().SetMethod("GET").SetPath("/test")

	handler := route.GetHandler()
	if handler != nil {
		t.Error("Expected handler to be nil when no handlers are set")
	}
}

func TestResponseHelpers(t *testing.T) {
	tests := []struct {
		name         string
		responseFunc func(http.ResponseWriter, *http.Request, string)
		body         string
		expectedCT   string
	}{
		{
			name:         "HTMLResponse",
			responseFunc: rtr.HTMLResponse,
			body:         "<h1>Test</h1>",
			expectedCT:   "text/html; charset=utf-8",
		},
		{
			name:         "JSONResponse",
			responseFunc: rtr.JSONResponse,
			body:         `{"test": true}`,
			expectedCT:   "application/json",
		},
		{
			name:         "CSSResponse",
			responseFunc: rtr.CSSResponse,
			body:         "body { color: red; }",
			expectedCT:   "text/css",
		},
		{
			name:         "XMLResponse",
			responseFunc: rtr.XMLResponse,
			body:         `<?xml version="1.0"?><root/>`,
			expectedCT:   "application/xml",
		},
		{
			name:         "TextResponse",
			responseFunc: rtr.TextResponse,
			body:         "Plain text",
			expectedCT:   "text/plain; charset=utf-8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			tt.responseFunc(w, req, tt.body)

			if w.Body.String() != tt.body {
				t.Errorf("Expected body '%s', got '%s'", tt.body, w.Body.String())
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != tt.expectedCT {
				t.Errorf("Expected Content-Type '%s', got '%s'", tt.expectedCT, contentType)
			}
		})
	}
}

func TestHTMLResponsePreservesExistingContentType(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Set a custom Content-Type before calling HTMLResponse
	w.Header().Set("Content-Type", "text/html; charset=iso-8859-1")

	rtr.HTMLResponse(w, req, "<h1>Test</h1>")

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html; charset=iso-8859-1" {
		t.Errorf("Expected Content-Type to be preserved as 'text/html; charset=iso-8859-1', got '%s'", contentType)
	}
}

func TestToHandlerHelpers(t *testing.T) {
	tests := []struct {
		name         string
		handlerFunc  func() rtr.StdHandler
		expectedBody string
		expectedCT   string
	}{
		{
			name: "ToHandler simple string",
			handlerFunc: func() rtr.StdHandler {
				return rtr.ToStdHandler(func(w http.ResponseWriter, r *http.Request) string {
					return "<h1>HTML Test</h1>"
				})
			},
			expectedBody: "<h1>HTML Test</h1>",
			expectedCT:   "text/html; charset=utf-8", // Go automatically detects HTML content
		},
		{
			name: "ToHandler with manual headers",
			handlerFunc: func() rtr.StdHandler {
				return rtr.ToStdHandler(func(w http.ResponseWriter, r *http.Request) string {
					w.Header().Set("Content-Type", "application/json")
					return `{"test": "json"}`
				})
			},
			expectedBody: `{"test": "json"}`,
			expectedCT:   "application/json",
		},
		{
			name: "ToHandler plain text",
			handlerFunc: func() rtr.StdHandler {
				return rtr.ToStdHandler(func(w http.ResponseWriter, r *http.Request) string {
					return "Plain text content"
				})
			},
			expectedBody: "Plain text content",
			expectedCT:   "text/plain; charset=utf-8", // Go automatically detects plain text content
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			handler := tt.handlerFunc()
			handler(w, req)

			if w.Body.String() != tt.expectedBody {
				t.Errorf("Expected body '%s', got '%s'", tt.expectedBody, w.Body.String())
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != tt.expectedCT {
				t.Errorf("Expected Content-Type '%s', got '%s'", tt.expectedCT, contentType)
			}
		})
	}
}

func TestStaticHandler(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	// Create a test file
	testContent := "body { color: blue; }"
	testFile := tempDir + "/test.css"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/static/*").
		SetStaticHandler(func(w http.ResponseWriter, r *http.Request) string {
			return tempDir // Return the static directory path
		})

	handler := route.GetHandler()
	if handler == nil {
		t.Fatal("Expected handler to be non-nil")
	}

	// Test serving the static file
	req := httptest.NewRequest("GET", "/static/test.css", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if body != testContent {
		t.Errorf("Expected body '%s', got '%s'", testContent, body)
	}

	// Check that Content-Type is set correctly for CSS
	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/css") {
		t.Errorf("Expected Content-Type to contain 'text/css', got '%s'", contentType)
	}
}

func TestStaticHandlerPreventsDirectoryTraversal(t *testing.T) {
	tempDir := t.TempDir()

	// Create a test file
	testContent := "safe content"
	testFile := tempDir + "/safe.txt"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/static/*").
		SetStaticHandler(func(w http.ResponseWriter, r *http.Request) string {
			return tempDir
		})

	handler := route.GetHandler()

	// Test directory traversal attempt
	req := httptest.NewRequest("GET", "/static/../../../etc/passwd", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for directory traversal attempt, got %d", w.Code)
	}
}

func TestStaticHandlerFileNotFound(t *testing.T) {
	tempDir := t.TempDir()

	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/static/*").
		SetStaticHandler(func(w http.ResponseWriter, r *http.Request) string {
			return tempDir
		})

	handler := route.GetHandler()

	// Test requesting non-existent file
	req := httptest.NewRequest("GET", "/static/nonexistent.txt", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for non-existent file, got %d", w.Code)
	}
}

func TestStaticHandlerPriority(t *testing.T) {
	// Test that StaticHandler takes priority over TextHandler
	tempDir := t.TempDir()
	testContent := "static file content"
	testFile := tempDir + "/test.txt"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/static/*").
		SetStaticHandler(func(w http.ResponseWriter, r *http.Request) string {
			return tempDir
		}).
		SetTextHandler(func(w http.ResponseWriter, r *http.Request) string {
			return "text handler content"
		})

	handler := route.GetHandler()
	if handler == nil {
		t.Fatal("Expected handler to be non-nil")
	}

	req := httptest.NewRequest("GET", "/static/test.txt", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	body := w.Body.String()
	if body != testContent {
		t.Errorf("Expected '%s' (StaticHandler should take priority), got '%s'", testContent, body)
	}
}

func TestStaticHandlerGetterSetter(t *testing.T) {
	route := rtr.NewRoute()

	// Test initial state
	if route.GetStaticHandler() != nil {
		t.Error("Expected initial StaticHandler to be nil")
	}

	// Test setter
	staticHandler := func(w http.ResponseWriter, r *http.Request) string {
		return "/static"
	}

	returnedRoute := route.SetStaticHandler(staticHandler)
	if returnedRoute != route {
		t.Error("SetStaticHandler should return the same route instance")
	}

	// Test getter
	if route.GetStaticHandler() == nil {
		t.Error("Expected StaticHandler to be non-nil after setting")
	}

	// Test that the handler works
	handler := route.GetHandler()
	if handler == nil {
		t.Fatal("Expected handler to be non-nil")
	}
}

func TestStaticFSHandler(t *testing.T) {
	f := fstest.MapFS{
		"style.css":  {Data: []byte("body{color:red}")},
		"index.html": {Data: []byte("<h1>Home</h1>")},
	}

	route := rtr.GetStaticFS("/static/*", f)
	handler := route.GetHandler()
	if handler == nil {
		t.Fatal("Expected handler to be non-nil")
	}

	req := httptest.NewRequest("GET", "/static/style.css", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if body := w.Body.String(); body != "body{color:red}" {
		t.Errorf("Expected body 'body{color:red}', got '%s'", body)
	}
}

func TestStaticFSHandlerNotFound(t *testing.T) {
	f := fstest.MapFS{}

	route := rtr.GetStaticFS("/static/*", f)
	handler := route.GetHandler()
	if handler == nil {
		t.Fatal("Expected handler to be non-nil")
	}

	req := httptest.NewRequest("GET", "/static/missing.txt", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}
