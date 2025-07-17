package rtr_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/rtr"
)

func TestRecoveryMiddleware(t *testing.T) {
	// Create a test handler that panics
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Wrap the panic handler with our recovery middleware
	handler := rtr.RecoveryMiddleware(panicHandler)

	// Call ServeHTTP which should recover from the panic
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	// Check the response body is what we expect
	expected := "Internal Server Error\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestDefaultMiddlewares(t *testing.T) {
	middlewares := rtr.DefaultMiddlewares()

	// Check that we have exactly one middleware (the recovery middleware)
	if len(middlewares) != 1 {
		t.Errorf("expected 1 default middleware, got %d", len(middlewares))
	}
}
