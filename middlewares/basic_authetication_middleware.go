package middlewares

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"

	"github.com/dracory/rtr"
)

// BasicAuthenticationMiddleware creates a new middleware that enforces HTTP Basic Authentication
// using the provided username and password.
func BasicAuthenticationMiddleware(username string, password string) rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Basic Authentication").
		SetHandler(basicAuthHandler(username, password))
}

// basicAuthHandler returns a handler function that performs HTTP Basic Authentication
func basicAuthHandler(expectedUsername, expectedPassword string) func(http.Handler) http.Handler {
	// Validate input parameters
	if expectedUsername == "" || expectedPassword == "" {
		// Log the error and return a handler that always returns 500
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "User name or password is empty", http.StatusBadRequest)
			})
		}
	}

	// sendUnauthorized sends a 401 Unauthorized response with the proper WWW-Authenticate header
	sendUnauthorized := func(w http.ResponseWriter) {
		if w == nil {
			return
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	return func(next http.Handler) http.Handler {
		// If next is nil, use http.NotFoundHandler as a safe default
		if next == nil {
			next = http.NotFoundHandler()
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r == nil {
				http.Error(w, "Bad Request: nil request", http.StatusBadRequest)
				return
			}

			if w == nil {
				// If we can't write the response, there's nothing we can do
				return
			}

			submittedUsername, submittedPassword, ok := r.BasicAuth()

			if !ok || submittedUsername == "" || submittedPassword == "" {
				sendUnauthorized(w)
				return
			}

			// Hash the credentials for constant-time comparison
			expectedUsernameHash := sha256.Sum256([]byte(expectedUsername))
			expectedPasswordHash := sha256.Sum256([]byte(expectedPassword))
			submittedUsernameHash := sha256.Sum256([]byte(submittedUsername))
			submittedPasswordHash := sha256.Sum256([]byte(submittedPassword))

			// Use constant-time comparison to avoid timing attacks
			usernameMatch := (subtle.ConstantTimeCompare(submittedUsernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(submittedPasswordHash[:], expectedPasswordHash[:]) == 1)

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}

			sendUnauthorized(w)
		})
	}
}
