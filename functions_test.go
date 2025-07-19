package rtr

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestToHandler(t *testing.T) {
	tests := []struct {
		name         string
		handlerFunc  func(http.ResponseWriter, *http.Request) string
		expectedBody string
		expectedCT   string
		setupHeaders func(http.ResponseWriter)
	}{
		{
			name: "simple string handler",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) string {
				return "Hello, World!"
			},
			expectedBody: "Hello, World!",
			expectedCT:   "text/plain; charset=utf-8", // Go automatically detects plain text
		},
		{
			name: "HTML string handler",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) string {
				return "<h1>Hello HTML</h1>"
			},
			expectedBody: "<h1>Hello HTML</h1>",
			expectedCT:   "text/html; charset=utf-8", // Go automatically detects HTML
		},
		{
			name: "JSON string handler",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) string {
				return `{"message": "Hello JSON"}`
			},
			expectedBody: `{"message": "Hello JSON"}`,
			expectedCT:   "text/plain; charset=utf-8", // Go treats JSON as plain text unless Content-Type is set
		},
		{
			name: "handler that sets its own headers",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) string {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-Custom-Header", "test-value")
				return `{"message": "With headers"}`
			},
			expectedBody: `{"message": "With headers"}`,
			expectedCT:   "application/json",
		},
		{
			name: "empty string handler",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) string {
				return ""
			},
			expectedBody: "",
			expectedCT:   "text/plain; charset=utf-8", // Go sets this even for empty content
		},
		{
			name: "multiline string handler",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) string {
				return "Line 1\nLine 2\nLine 3"
			},
			expectedBody: "Line 1\nLine 2\nLine 3",
			expectedCT:   "text/plain; charset=utf-8", // Go detects as plain text
		},
		{
			name: "handler with URL parameters",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) string {
				// Simulate getting a parameter (in real usage this would come from context)
				return "User ID: 123"
			},
			expectedBody: "User ID: 123",
			expectedCT:   "text/plain; charset=utf-8", // Go detects as plain text
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the handler using ToHandler
			handler := ToHandler(tt.handlerFunc)

			// Create test request and response recorder
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			// Execute the handler
			handler(w, req)

			// Check the response body
			if w.Body.String() != tt.expectedBody {
				t.Errorf("Expected body '%s', got '%s'", tt.expectedBody, w.Body.String())
			}

			// Check Content-Type header if expected
			contentType := w.Header().Get("Content-Type")
			if contentType != tt.expectedCT {
				t.Errorf("Expected Content-Type '%s', got '%s'", tt.expectedCT, contentType)
			}

			// Check status code (should be 200 by default)
			if w.Code != http.StatusOK {
				t.Errorf("Expected status code 200, got %d", w.Code)
			}
		})
	}
}

func TestToHandlerWithCustomHeaders(t *testing.T) {
	// Test that ToHandler preserves custom headers set by the string handler
	handler := ToHandler(func(w http.ResponseWriter, r *http.Request) string {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("X-Custom", "test-value")
		return "Content with custom headers"
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	// Check that all custom headers are preserved
	expectedHeaders := map[string]string{
		"Content-Type":  "text/plain; charset=utf-8",
		"Cache-Control": "no-cache",
		"X-Custom":      "test-value",
	}

	for headerName, expectedValue := range expectedHeaders {
		actualValue := w.Header().Get(headerName)
		if actualValue != expectedValue {
			t.Errorf("Expected header %s to be '%s', got '%s'", headerName, expectedValue, actualValue)
		}
	}

	// Check body
	expectedBody := "Content with custom headers"
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, w.Body.String())
	}
}

func TestToHandlerWithStatusCode(t *testing.T) {
	// Test that ToHandler allows the string handler to set custom status codes
	handler := ToHandler(func(w http.ResponseWriter, r *http.Request) string {
		w.WriteHeader(http.StatusCreated)
		return "Resource created"
	})

	req := httptest.NewRequest("POST", "/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	// Check status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Check body
	expectedBody := "Resource created"
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, w.Body.String())
	}
}

func TestToHandlerNilHandler(t *testing.T) {
	// Test behavior with nil handler (should panic when called, which is expected)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected ToHandler with nil handler to panic when called, but it didn't")
		}
	}()

	// Create handler with nil function
	handler := ToHandler(nil)
	
	// This should panic when we try to call the nil function
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	handler(w, req) // This should panic
}

func TestToHandlerReturnType(t *testing.T) {
	// Test that ToHandler returns the correct Handler type
	stringHandler := func(w http.ResponseWriter, r *http.Request) string {
		return "test"
	}

	handler := ToHandler(stringHandler)

	// Verify it implements the Handler interface
	var _ Handler = handler

	// Verify it's actually a function with the right signature
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// This should not panic
	handler(w, req)

	if w.Body.String() != "test" {
		t.Errorf("Expected body 'test', got '%s'", w.Body.String())
	}
}

func TestErrorHandlerToHandler(t *testing.T) {
	tests := []struct {
		name         string
		errorHandler ErrorHandler
		expectedBody string
		expectedCode int
	}{
		{
			name: "error handler returns error",
			errorHandler: func(w http.ResponseWriter, r *http.Request) error {
				return errors.New("something went wrong")
			},
			expectedBody: "something went wrong",
			expectedCode: http.StatusOK, // ErrorHandler doesn't set status codes
		},
		{
			name: "error handler returns nil",
			errorHandler: func(w http.ResponseWriter, r *http.Request) error {
				return nil
			},
			expectedBody: "", // No output when error is nil
			expectedCode: http.StatusOK,
		},
		{
			name: "error handler with custom error message",
			errorHandler: func(w http.ResponseWriter, r *http.Request) error {
				return errors.New("user not found")
			},
			expectedBody: "user not found",
			expectedCode: http.StatusOK,
		},
		{
			name: "error handler that sets headers before returning error",
			errorHandler: func(w http.ResponseWriter, r *http.Request) error {
				w.Header().Set("X-Custom-Header", "error-occurred")
				return errors.New("internal error")
			},
			expectedBody: "internal error",
			expectedCode: http.StatusOK,
		},
		{
			name: "error handler that sets headers and returns nil",
			errorHandler: func(w http.ResponseWriter, r *http.Request) error {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-Success", "true")
				return nil
			},
			expectedBody: "",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the handler using ErrorHandlerToHandler
			handler := ErrorHandlerToHandler(tt.errorHandler)

			// Create test request and response recorder
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			// Execute the handler
			handler(w, req)

			// Check the response body
			if w.Body.String() != tt.expectedBody {
				t.Errorf("Expected body '%s', got '%s'", tt.expectedBody, w.Body.String())
			}

			// Check status code
			if w.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, w.Code)
			}
		})
	}
}

func TestErrorHandlerToHandlerWithCustomHeaders(t *testing.T) {
	// Test that ErrorHandlerToHandler preserves custom headers set by the error handler
	handler := ErrorHandlerToHandler(func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Error-Code", "USER_NOT_FOUND")
		return errors.New(`{"error": "user not found"}`)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	// Check that custom headers are preserved
	expectedHeaders := map[string]string{
		"Content-Type":  "application/json",
		"X-Error-Code": "USER_NOT_FOUND",
	}

	for headerName, expectedValue := range expectedHeaders {
		actualValue := w.Header().Get(headerName)
		if actualValue != expectedValue {
			t.Errorf("Expected header %s to be '%s', got '%s'", headerName, expectedValue, actualValue)
		}
	}

	// Check body
	expectedBody := `{"error": "user not found"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, w.Body.String())
	}
}

func TestErrorHandlerToHandlerNilHandler(t *testing.T) {
	// Test behavior with nil handler (should panic when called, which is expected)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected ErrorHandlerToHandler with nil handler to panic when called, but it didn't")
		}
	}()

	// Create handler with nil function
	handler := ErrorHandlerToHandler(nil)
	
	// This should panic when we try to call the nil function
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	handler(w, req) // This should panic
}

func TestErrorHandlerToHandlerReturnType(t *testing.T) {
	// Test that ErrorHandlerToHandler returns the correct Handler type
	errorHandler := func(w http.ResponseWriter, r *http.Request) error {
		return errors.New("test error")
	}

	handler := ErrorHandlerToHandler(errorHandler)

	// Verify it implements the Handler interface
	var _ Handler = handler

	// Verify it's actually a function with the right signature
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// This should not panic
	handler(w, req)

	if w.Body.String() != "test error" {
		t.Errorf("Expected body 'test error', got '%s'", w.Body.String())
	}
}
