package rtr_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
				return rtr.ToHandler(func(w http.ResponseWriter, r *http.Request) string {
					return "<h1>HTML Test</h1>"
				})
			},
			expectedBody: "<h1>HTML Test</h1>",
			expectedCT:   "text/html; charset=utf-8", // Go automatically detects HTML content
		},
		{
			name: "ToHandler with manual headers",
			handlerFunc: func() rtr.StdHandler {
				return rtr.ToHandler(func(w http.ResponseWriter, r *http.Request) string {
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
				return rtr.ToHandler(func(w http.ResponseWriter, r *http.Request) string {
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
