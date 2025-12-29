package rtr

import "net/http"

// HTMLResponse responds with HTML content and sets the appropriate Content-Type header.
// It sets the Content-Type to "text/html; charset=utf-8" if not already set.
func HTMLResponse(w http.ResponseWriter, r *http.Request, body string) {
	contentType := w.Header().Get("Content-Type")
	if contentType == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	}
	_, _ = w.Write([]byte(body))
}

// JSONResponse responds with JSON content and sets the appropriate Content-Type header.
// It sets the Content-Type to "application/json".
func JSONResponse(w http.ResponseWriter, r *http.Request, body string) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(body))
}

// CSSResponse responds with CSS content and sets the appropriate Content-Type header.
// It sets the Content-Type to "text/css".
func CSSResponse(w http.ResponseWriter, r *http.Request, body string) {
	w.Header().Set("Content-Type", "text/css")
	_, _ = w.Write([]byte(body))
}

// XMLResponse responds with XML content and sets the appropriate Content-Type header.
// It sets the Content-Type to "application/xml".
func XMLResponse(w http.ResponseWriter, r *http.Request, body string) {
	w.Header().Set("Content-Type", "application/xml")
	_, _ = w.Write([]byte(body))
}

// TextResponse responds with plain text content and sets the appropriate Content-Type header.
// It sets the Content-Type to "text/plain; charset=utf-8".
func TextResponse(w http.ResponseWriter, r *http.Request, body string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(body))
}

// JSResponse responds with JavaScript content and sets the appropriate Content-Type header.
// It sets the Content-Type to "application/javascript".
func JSResponse(w http.ResponseWriter, r *http.Request, body string) {
	w.Header().Set("Content-Type", "application/javascript")
	_, _ = w.Write([]byte(body))
}
