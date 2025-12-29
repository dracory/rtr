package rtr_test

// This file exercises controller-based APIs:
// - Standard, HTML, JSON, and Text controllers
// - Priority rules between controllers and direct handlers
// - Controller integration with the router
// Handler-based tests live in handlers_test.go.

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/rtr"
)

// Test controller implementations

type standardController struct{}

func (c *standardController) Handler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Standard Controller Response"))
}

type htmlController struct{}

func (c *htmlController) Handler(w http.ResponseWriter, r *http.Request) string {
	return "<h1>HTML Controller Response</h1>"
}

type jsonController struct{}

func (c *jsonController) Handler(w http.ResponseWriter, r *http.Request) string {
	return `{"message": "JSON Controller Response"}`
}

type textController struct{}

func (c *textController) Handler(w http.ResponseWriter, r *http.Request) string {
	return "Text Controller Response"
}

// Tests

func TestControllerInterface(t *testing.T) {
	controller := &standardController{}
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/test").
		SetController(controller)

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

	body := w.Body.String()
	if body != "Standard Controller Response" {
		t.Errorf("Expected body 'Standard Controller Response', got '%s'", body)
	}

	// Note: Go's http.ResponseWriter automatically detects Content-Type when Write() is called
	// For ControllerInterface, we don't explicitly set headers, so Go detects it as text/plain
	contentType := w.Header().Get("Content-Type")
	if contentType != "text/plain; charset=utf-8" {
		t.Errorf("Expected Content-Type 'text/plain; charset=utf-8' (auto-detected by Go), got '%s'", contentType)
	}
}

func TestHTMLControllerInterface(t *testing.T) {
	controller := &htmlController{}
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/test").
		SetHTMLController(controller)

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
	if body != "<h1>HTML Controller Response</h1>" {
		t.Errorf("Expected body '<h1>HTML Controller Response</h1>', got '%s'", body)
	}
}

func TestJSONControllerInterface(t *testing.T) {
	controller := &jsonController{}
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/api/test").
		SetJSONController(controller)

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
	expected := `{"message": "JSON Controller Response"}`
	if body != expected {
		t.Errorf("Expected body '%s', got '%s'", expected, body)
	}
}

func TestTextControllerInterface(t *testing.T) {
	controller := &textController{}
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/robots.txt").
		SetTextController(controller)

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
	if body != "Text Controller Response" {
		t.Errorf("Expected body 'Text Controller Response', got '%s'", body)
	}
}

func TestControllerPriority(t *testing.T) {
	// Test that direct Handler takes priority over ControllerInterface
	controller := &standardController{}
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/test").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("Direct Handler"))
		}).
		SetController(controller)

	handler := route.GetHandler()
	if handler == nil {
		t.Fatal("Expected handler to be non-nil")
	}

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	body := w.Body.String()
	if body != "Direct Handler" {
		t.Errorf("Expected 'Direct Handler' (Handler should take priority over Controller), got '%s'", body)
	}
}

func TestHTMLControllerPriorityOverJSON(t *testing.T) {
	// Test that HTMLControllerInterface takes priority over JSONControllerInterface
	htmlCtrl := &htmlController{}
	jsonCtrl := &jsonController{}

	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/test").
		SetHTMLController(htmlCtrl).
		SetJSONController(jsonCtrl)

	handler := route.GetHandler()
	if handler == nil {
		t.Fatal("Expected handler to be non-nil")
	}

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	body := w.Body.String()
	if body != "<h1>HTML Controller Response</h1>" {
		t.Errorf("Expected '<h1>HTML Controller Response</h1>' (HTMLController should take priority), got '%s'", body)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected Content-Type 'text/html; charset=utf-8', got '%s'", contentType)
	}
}

func TestControllerIntegrationWithRouter(t *testing.T) {
	// Test that controllers work correctly when integrated with the router
	router := rtr.NewRouter()

	htmlCtrl := &htmlController{}
	jsonCtrl := &jsonController{}
	textCtrl := &textController{}

	router.AddRoute(rtr.NewRoute().
		SetMethod("GET").
		SetPath("/html").
		SetHTMLController(htmlCtrl))

	router.AddRoute(rtr.NewRoute().
		SetMethod("GET").
		SetPath("/json").
		SetJSONController(jsonCtrl))

	router.AddRoute(rtr.NewRoute().
		SetMethod("GET").
		SetPath("/text").
		SetTextController(textCtrl))

	tests := []struct {
		name         string
		path         string
		expectedBody string
		expectedCT   string
	}{
		{
			name:         "HTML Controller",
			path:         "/html",
			expectedBody: "<h1>HTML Controller Response</h1>",
			expectedCT:   "text/html; charset=utf-8",
		},
		{
			name:         "JSON Controller",
			path:         "/json",
			expectedBody: `{"message": "JSON Controller Response"}`,
			expectedCT:   "application/json",
		},
		{
			name:         "Text Controller",
			path:         "/text",
			expectedBody: "Text Controller Response",
			expectedCT:   "text/plain; charset=utf-8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d", w.Code)
			}

			body := w.Body.String()
			if body != tt.expectedBody {
				t.Errorf("Expected body '%s', got '%s'", tt.expectedBody, body)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != tt.expectedCT {
				t.Errorf("Expected Content-Type '%s', got '%s'", tt.expectedCT, contentType)
			}
		})
	}
}

func TestControllerGetterMethods(t *testing.T) {
	// Test that getter methods return the correct controllers
	stdCtrl := &standardController{}
	htmlCtrl := &htmlController{}
	jsonCtrl := &jsonController{}
	textCtrl := &textController{}

	route := rtr.NewRoute().
		SetController(stdCtrl).
		SetHTMLController(htmlCtrl).
		SetJSONController(jsonCtrl).
		SetTextController(textCtrl)

	if route.GetController() != stdCtrl {
		t.Error("GetController() did not return the correct controller")
	}

	if route.GetHTMLController() != htmlCtrl {
		t.Error("GetHTMLController() did not return the correct controller")
	}

	if route.GetJSONController() != jsonCtrl {
		t.Error("GetJSONController() did not return the correct controller")
	}

	if route.GetTextController() != textCtrl {
		t.Error("GetTextController() did not return the correct controller")
	}
}
