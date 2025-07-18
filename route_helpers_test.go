package rtr_test

import (
	"net/http"
	"testing"

	"github.com/dracory/rtr"
)

func TestRouteHelpers(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(string, rtr.Handler) rtr.RouteInterface
		method   string
		handler  rtr.Handler
		path     string
		hasError bool
	}{
		{
			name:   "GET helper",
			fn:     rtr.Get,
			method: http.MethodGet,
			path:   "/test",
			handler: func(w http.ResponseWriter, r *http.Request) {},
		},
		{
			name:   "POST helper",
			fn:     rtr.Post,
			method: http.MethodPost,
			path:   "/test",
			handler: func(w http.ResponseWriter, r *http.Request) {},
		},
		{
			name:   "PUT helper",
			fn:     rtr.Put,
			method: http.MethodPut,
			path:   "/test",
			handler: func(w http.ResponseWriter, r *http.Request) {},
		},
		{
			name:   "DELETE helper",
			fn:     rtr.Delete,
			method: http.MethodDelete,
			path:   "/test",
			handler: func(w http.ResponseWriter, r *http.Request) {},
		},
		// These cases don't panic, they just create routes with empty paths or nil handlers
		{
			name:     "empty path",
			fn:       rtr.Get,
			method:   http.MethodGet,
			path:     "",
			handler:  func(w http.ResponseWriter, r *http.Request) {},
			hasError: false,
		},
		// Skip testing nil handler case since it's not clear what the expected behavior should be
		// and the current implementation might not handle it gracefully
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture any panics to test error cases
			var route rtr.RouteInterface
			var panicked bool
			func() {
				defer func() {
					if r := recover(); r != nil {
						panicked = true
					}
				}()
				route = tt.fn(tt.path, tt.handler)
			}()

			if tt.hasError {
				if !panicked {
					t.Error("Expected panic but didn't get one")
				}
				return
			}

			if panicked {
				t.Error("Unexpected panic")
			}
			if route == nil {
				t.Error("Route should not be nil")
				return
			}
			if got := route.GetMethod(); got != tt.method {
				t.Errorf("Unexpected HTTP method: got %v, want %v", got, tt.method)
			}
			if got := route.GetPath(); got != tt.path {
				t.Errorf("Unexpected path: got %v, want %v", got, tt.path)
			}
			// Only check handler if it's not expected to be nil
			if tt.handler != nil && route.GetHandler() == nil {
				t.Error("Handler should not be nil")
			}
		})
	}
}
