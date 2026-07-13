package middlewares

import (
	"testing"
)

func TestIsNilInterface_NilValue(t *testing.T) {
	if !isNilInterface(nil) {
		t.Fatal("expected true for untyped nil")
	}
}

func TestIsNilInterface_TypedNilPointer(t *testing.T) {
	var p *int = nil
	if !isNilInterface(p) {
		t.Fatal("expected true for typed nil pointer")
	}
}

func TestIsNilInterface_NonNilPointer(t *testing.T) {
	n := 42
	if isNilInterface(&n) {
		t.Fatal("expected false for non-nil pointer")
	}
}

func TestIsNilInterface_NonNilStruct(t *testing.T) {
	s := struct{}{}
	if isNilInterface(s) {
		t.Fatal("expected false for non-nil struct")
	}
}

func TestIsNilInterface_NilInterface(t *testing.T) {
	var s AuthSession = nil
	if !isNilInterface(s) {
		t.Fatal("expected true for nil interface value")
	}
}

func TestIsNilInterface_TypedNilInterface(t *testing.T) {
	var session *mockAuthSession = nil
	var s AuthSession = session
	if !isNilInterface(s) {
		t.Fatal("expected true for interface holding typed nil")
	}
}

func TestIsNilInterface_NonNilInterface(t *testing.T) {
	var s AuthSession = &mockAuthSession{userID: "abc"}
	if isNilInterface(s) {
		t.Fatal("expected false for non-nil interface")
	}
}
