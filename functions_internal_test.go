package rtr

import (
	"net/http"
	"reflect"
	"testing"
)

// TestAppendReversed verifies the helper correctly appends the source slice
// to the destination slice in reverse order without mutating inputs.
func TestAppendReversed(t *testing.T) {
	mk := func(_ string) MiddlewareInterface {
		return NewAnonymousMiddleware(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		})
	}

	// Prepare identifiable middlewares (by instance identity)
	a := mk("a")
	b := mk("b")
	c := mk("c")

	t.Run("nil src", func(t *testing.T) {
		dst := []MiddlewareInterface{a}
		var src []MiddlewareInterface

		out := appendReversed(dst, src)

		if !reflect.DeepEqual(out, []MiddlewareInterface{a}) {
			t.Fatalf("expected [a], got %v", out)
		}
	})

	t.Run("empty src", func(t *testing.T) {
		dst := []MiddlewareInterface{a}
		src := []MiddlewareInterface{}

		out := appendReversed(dst, src)

		if !reflect.DeepEqual(out, []MiddlewareInterface{a}) {
			t.Fatalf("expected [a], got %v", out)
		}
	})

	t.Run("single element", func(t *testing.T) {
		dst := []MiddlewareInterface{a}
		src := []MiddlewareInterface{b}

		out := appendReversed(dst, src)

		expected := []MiddlewareInterface{a, b}
		if !reflect.DeepEqual(out, expected) {
			t.Fatalf("expected [a,b], got %v", out)
		}
	})

	t.Run("multiple elements reversed order", func(t *testing.T) {
		dst := []MiddlewareInterface{a}
		src := []MiddlewareInterface{b, c}

		// Expect to append c then b
		out := appendReversed(dst, src)

		expected := []MiddlewareInterface{a, c, b}
		if !reflect.DeepEqual(out, expected) {
			t.Fatalf("expected [a,c,b], got %v", out)
		}

		// Ensure inputs not mutated
		if !reflect.DeepEqual(dst, []MiddlewareInterface{a}) {
			t.Fatalf("dst should not be mutated; got %v", dst)
		}
		if !reflect.DeepEqual(src, []MiddlewareInterface{b, c}) {
			t.Fatalf("src should not be mutated; got %v", src)
		}
	})
}
