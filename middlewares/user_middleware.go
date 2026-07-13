package middlewares

import (
	"net/http"
	"strings"

	"github.com/dracory/rtr"
)

// UserMiddlewareConfig configures the user middleware
type UserMiddlewareConfig struct {
	// GetUser extracts the authenticated user from the request context.
	// Returns nil if no user is authenticated. Required.
	GetUser func(r *http.Request) UserMiddlewareUser

	// RegistrationEnabled indicates whether registration is enabled
	RegistrationEnabled bool

	// OnNotAuthenticated is called when the user is not authenticated.
	// Typically redirects to login with a flash message.
	OnNotAuthenticated func(w http.ResponseWriter, r *http.Request)

	// OnNotActive is called when the user account is not active.
	// Typically redirects to the home page with an error.
	OnNotActive func(w http.ResponseWriter, r *http.Request)

	// OnRegistrationIncomplete is called when the user hasn't completed registration
	// and registration is enabled. Typically redirects to the registration page.
	OnRegistrationIncomplete func(w http.ResponseWriter, r *http.Request)

	// RegistrationPaths is a list of URL paths that are exempt from the
	// registration completion check (e.g. /profile, /register)
	RegistrationPaths []string

	// RequireRoles is a list of role names. When non-empty, the user must
	// have at least one of the specified roles (via UserWithRole.HasRole).
	RequireRoles []string

	// OnNotAuthorized is called when the user lacks a required role.
	// Typically redirects to the home page with an error.
	OnNotAuthorized func(w http.ResponseWriter, r *http.Request)
}

// UserMiddlewareUser defines the interface for user middleware user operations
type UserMiddlewareUser interface {
	IsActive() bool
	IsRegistrationCompleted() bool
}

// UserWithRole extends UserMiddlewareUser with role-checking capability.
// Implementations provide a HasRole method so the middleware can verify
// arbitrary roles without hardcoding specific role methods.
type UserWithRole interface {
	UserMiddlewareUser
	HasRole(role string) bool
}

// UserMiddleware creates a middleware that checks if the user is authenticated
// and active before allowing access to the protected route.
//
// Required config field: GetUser. If missing, the middleware returns
// HTTP 500 on every request with a descriptive error message.
//
// Business logic:
//  1. user must be authenticated
//  2. user must be active
//  3. user must have completed registration (unless on an exempt path)
//  4. optional role check must pass
func UserMiddleware(config UserMiddlewareConfig) rtr.MiddlewareInterface {
	var configErr string
	if config.GetUser == nil {
		configErr = "user middleware: GetUser is required"
	}

	return rtr.NewMiddleware().
		SetName("User Middleware").
		SetHandler(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if configErr != "" {
					http.Error(w, configErr, http.StatusInternalServerError)
					return
				}

				user := config.GetUser(r)

				if isNilInterface(user) {
					if config.OnNotAuthenticated != nil {
						config.OnNotAuthenticated(w, r)
						return
					}
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				if !user.IsActive() {
					if config.OnNotActive != nil {
						config.OnNotActive(w, r)
						return
					}
					http.Error(w, "Account inactive", http.StatusForbidden)
					return
				}

				if !user.IsRegistrationCompleted() && !isOnRegistrationPath(r.URL.Path, config.RegistrationPaths) {
					if config.RegistrationEnabled {
						if config.OnRegistrationIncomplete != nil {
							config.OnRegistrationIncomplete(w, r)
							return
						}
						http.Error(w, "Registration incomplete", http.StatusForbidden)
						return
					}
				}

				if len(config.RequireRoles) > 0 && !userHasRole(user, config.RequireRoles) {
					if config.OnNotAuthorized != nil {
						config.OnNotAuthorized(w, r)
						return
					}
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}

				next.ServeHTTP(w, r)
			})
		})
}

// userHasRole checks if the user implements UserWithRole and has at least one of the required roles.
func userHasRole(user UserMiddlewareUser, roles []string) bool {
	userWithRole, ok := user.(UserWithRole)
	if !ok {
		return false
	}
	for _, role := range roles {
		if userWithRole.HasRole(role) {
			return true
		}
	}
	return false
}

// isOnRegistrationPath checks if the request path matches any of the registration paths
func isOnRegistrationPath(requestPath string, registrationPaths []string) bool {
	trimmedRequest := strings.Trim(requestPath, "/")

	for _, p := range registrationPaths {
		if strings.Trim(p, "/") == trimmedRequest {
			return true
		}
	}

	return false
}
