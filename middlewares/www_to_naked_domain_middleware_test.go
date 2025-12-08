package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWwwToNakedDomainMiddleware_RedirectsWww(t *testing.T) {
	m := WwwToNakedDomainMiddleware()

	nextCalled := false
	handler := m.Execute(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	}))

	req := httptest.NewRequest(http.MethodGet, "/some/path?x=1", nil)
	req.Host = "www.example.com"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if nextCalled {
		t.Fatalf("next handler should not be called on www host")
	}

	if rr.Code != http.StatusTemporaryRedirect {
		t.Fatalf("expected status %d, got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	wantLocation := "https://example.com/some/path?x=1"
	if rr.Header().Get("Location") != wantLocation {
		t.Fatalf("expected redirect to %q, got %q", wantLocation, rr.Header().Get("Location"))
	}
}

func TestWwwToNakedDomainMiddleware_PassesThroughNonWww(t *testing.T) {
	m := WwwToNakedDomainMiddleware()

	nextCalled := false
	handler := m.Execute(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/ok", nil)
	req.Host = "example.com"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if !nextCalled {
		t.Fatal("expected next handler to be called for non-www host")
	}

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}
