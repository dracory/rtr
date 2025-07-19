package rtr

import "net/http"

// ToHandler converts any string-returning handler to a standard Handler.
// It simply writes the returned string to the response without setting any headers.
// The string handler is responsible for setting any headers it needs.
//
// Parameters:
//   - handler: The string handler function to convert.
//
// Returns:
//   - A standard Handler function that writes the returned string to the response.
func ToHandler(handler StringHandler) StdHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(handler(w, r)))
	}
}

// ErrorHandlerToHandler converts an ErrorHandler to a standard Handler.
// If the error handler returns an error, it writes the error message to the response.
// If the error handler returns nil, it does nothing.
//
// Parameters:
//   - handler: The error handler function to convert.
//
// Returns:
//   - A standard Handler function that writes the error message to the response if an error is returned.
func ErrorHandlerToHandler(handler ErrorHandler) StdHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
	}
}
