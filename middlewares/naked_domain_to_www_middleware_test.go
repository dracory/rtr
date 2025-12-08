package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/rtr/middlewares"
)

func TestNakedDomainToWwwMiddleware_RedirectsNaked(t *testing.T) {
	m := middlewares.NakedDomainToWwwMiddleware(nil)

	nextCalled := false
	handler := m.Execute(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	}))

	req := httptest.NewRequest(http.MethodGet, "/some/path?x=1", nil)
	req.Host = "example.com"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if nextCalled {
		t.Fatalf("next handler should not be called on naked host")
	}

	if rr.Code != http.StatusPermanentRedirect {
		t.Fatalf("expected status %d, got %d", http.StatusPermanentRedirect, rr.Code)
	}

	wantLocation := "https://www.example.com/some/path?x=1"
	if rr.Header().Get("Location") != wantLocation {
		t.Fatalf("expected redirect to %q, got %q", wantLocation, rr.Header().Get("Location"))
	}
}

func TestNakedDomainToWwwMiddleware_PassesThroughWww(t *testing.T) {
	m := middlewares.NakedDomainToWwwMiddleware(nil)

	nextCalled := false
	handler := m.Execute(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/ok", nil)
	req.Host = "www.example.com"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if !nextCalled {
		t.Fatal("expected next handler to be called for www host")
	}

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestNakedDomainToWwwMiddleware_RespectsExcludes(t *testing.T) {
	m := middlewares.NakedDomainToWwwMiddleware([]string{"localhost", "127.0.0.1"})

	nextCalled := false
	handler := m.Execute(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusAccepted)
	}))

	req := httptest.NewRequest(http.MethodGet, "/ok", nil)
	req.Host = "localhost"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if !nextCalled {
		t.Fatal("expected next handler to be called for excluded host")
	}

	if rr.Code != http.StatusAccepted {
		t.Fatalf("expected status %d, got %d", http.StatusAccepted, rr.Code)
	}
}
