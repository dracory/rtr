package rtr

import "net/http"

// ToHandler converts any string-returning handler to a standard Handler.
// It simply writes the returned string to the response without setting any headers.
// The string handler is responsible for setting any headers it needs.
func ToHandler(handler func(http.ResponseWriter, *http.Request) string) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(handler(w, r)))
	}
}
