package rtr_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/rtr"
)

func TestBuildMiddlewareChain(t *testing.T) {
	var order []string

	newMiddleware := func(name string) rtr.MiddlewareInterface {
		return rtr.NewAnonymousMiddleware(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, name+" before")
				next.ServeHTTP(w, r)
				order = append(order, name+" after")
			})
		})
	}

	mw1 := newMiddleware("mw1")
	mw2 := newMiddleware("mw2")
	mw3 := newMiddleware("mw3")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
	})

	chain := rtr.BuildMiddlewareChain(handler, []rtr.MiddlewareInterface{mw1, mw2}, []rtr.MiddlewareInterface{mw3})

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	chain.ServeHTTP(w, req)

	expected := []string{"mw1 before", "mw2 before", "mw3 before", "handler", "mw3 after", "mw2 after", "mw1 after"}

	if len(order) != len(expected) {
		t.Fatalf("Expected %d items, got %d", len(expected), len(order))
	}

	for i := range expected {
		if order[i] != expected[i] {
			t.Errorf("Expected %s at index %d, got %s", expected[i], i, order[i])
		}
	}
}
