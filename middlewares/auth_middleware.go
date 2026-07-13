package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/dracory/rtr"
	"github.com/jellydator/ttlcache/v3"
)

const (
	authSessionCacheTTL    = 5 * time.Minute
	authUserCacheTTL       = 5 * time.Minute
	authSessionCachePrefix = "auth:session:"
	authUserCachePrefix    = "auth:user:"
)

// AuthMiddlewareConfig configures the auth middleware
type AuthMiddlewareConfig struct {
	// SessionStore provides session lookup by key. Required.
	SessionStore AuthSessionStore
	// UserStore provides user lookup by ID. Required.
	UserStore AuthUserStore
	// Logger for error logging. Optional.
	Logger AuthLogger
	// MemoryCache is an optional TTL cache for session/user objects.
	MemoryCache *ttlcache.Cache[string, any]
	// ContextKeyUser is the context key used to store the authenticated user. Required.
	ContextKeyUser any
	// ContextKeySession is the context key used to store the session object. Required.
	ContextKeySession any
	// CookieName is the name of the cookie containing the session key. Required.
	CookieName string
}

// AuthSessionStore defines the interface for session store operations
type AuthSessionStore interface {
	SessionFindByKey(ctx context.Context, key string) (AuthSession, error)
}

// AuthSession defines the interface for session operations
type AuthSession interface {
	GetUserID() string
	IsExpired() bool
}

// AuthUserStore defines the interface for user store operations
type AuthUserStore interface {
	UserFindByID(ctx context.Context, id string) (AuthUser, error)
}

// AuthUser defines the interface for user operations
type AuthUser interface {
	IsActive() bool
	IsAdministrator() bool
	IsSuperuser() bool
	IsRegistrationCompleted() bool
}

// AuthLogger defines the interface for logging
type AuthLogger interface {
	Error(msg string, args ...any)
}

// AuthMiddleware creates a middleware that adds the authenticated user and session
// to the request context.
//
// Business logic:
//  1. Checks if the user session key exists in the incoming request cookie
//  2. Retrieves the session using the session key (with optional memory cache)
//  3. Checks the session is not expired
//  4. Retrieves the user using the user ID from the session (with optional memory cache)
//  5. Stores the user and session object in the request context
func AuthMiddleware(config AuthMiddlewareConfig) rtr.MiddlewareInterface {
	var configErr string
	switch {
	case config.SessionStore == nil:
		configErr = "auth middleware: SessionStore is required"
	case config.UserStore == nil:
		configErr = "auth middleware: UserStore is required"
	case config.ContextKeyUser == nil:
		configErr = "auth middleware: ContextKeyUser is required"
	case config.ContextKeySession == nil:
		configErr = "auth middleware: ContextKeySession is required"
	case config.CookieName == "":
		configErr = "auth middleware: CookieName is required"
	}

	return rtr.NewMiddleware().
		SetName("Auth Middleware").
		SetHandler(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if next == nil || r == nil {
					return
				}

				if configErr != "" {
					http.Error(w, configErr, http.StatusInternalServerError)
					return
				}

				sessionKey := authMiddlewareSessionKey(r, config.CookieName, config.Logger)
				if sessionKey == "" {
					next.ServeHTTP(w, r)
					return
				}

				session, ok := authResolveSession(w, r, next, config, sessionKey)
				if !ok {
					return
				}

				userID := session.GetUserID()
				if userID == "" {
					next.ServeHTTP(w, r)
					return
				}

				user, ok := authResolveUser(w, r, next, config, userID)
				if !ok {
					return
				}

				ctx := context.WithValue(r.Context(), config.ContextKeyUser, user)
				ctx = context.WithValue(ctx, config.ContextKeySession, session)

				next.ServeHTTP(w, r.WithContext(ctx))
			})
		})
}

// authResolveSession retrieves a valid (non-expired) session from cache or store.
// Returns (session, true) on success, or (nil, false) after calling next on failure.
func authResolveSession(w http.ResponseWriter, r *http.Request, next http.Handler, config AuthMiddlewareConfig, sessionKey string) (AuthSession, bool) {
	session := authCacheGetSession(config.MemoryCache, sessionKey)

	if isNilInterface(session) {
		var err error
		session, err = config.SessionStore.SessionFindByKey(r.Context(), sessionKey)

		if err != nil {
			if config.Logger != nil {
				config.Logger.Error("auth_middleware", "error", err.Error())
			}
			next.ServeHTTP(w, r)
			return nil, false
		}

		if isNilInterface(session) {
			next.ServeHTTP(w, r)
			return nil, false
		}

		if !session.IsExpired() {
			authCacheSetSession(config.MemoryCache, sessionKey, session)
		}
	}

	if session.IsExpired() {
		next.ServeHTTP(w, r)
		return nil, false
	}

	return session, true
}

// authResolveUser retrieves a valid user from cache or store.
// Returns (user, true) on success, or (nil, false) after calling next on failure.
func authResolveUser(w http.ResponseWriter, r *http.Request, next http.Handler, config AuthMiddlewareConfig, userID string) (AuthUser, bool) {
	user := authCacheGetUser(config.MemoryCache, userID)

	if isNilInterface(user) {
		fetchedUser, err := config.UserStore.UserFindByID(r.Context(), userID)

		if err != nil {
			if config.Logger != nil {
				config.Logger.Error("auth_middleware", "error", err.Error())
			}
			next.ServeHTTP(w, r)
			return nil, false
		}

		if isNilInterface(fetchedUser) {
			next.ServeHTTP(w, r)
			return nil, false
		}

		authCacheSetUser(config.MemoryCache, userID, fetchedUser)
		user = fetchedUser
	}

	return user, true
}

// authMiddlewareSessionKey extracts the session key from the request cookie
func authMiddlewareSessionKey(r *http.Request, cookieName string, logger AuthLogger) string {
	if cookieName == "" {
		return ""
	}

	cookie, err := r.Cookie(cookieName)
	if err != nil {
		if err != http.ErrNoCookie && logger != nil {
			logger.Error("auth_middleware", "error", err.Error())
		}
	}

	if cookie == nil {
		return ""
	}

	return cookie.Value
}

func authCacheGetSession(cache *ttlcache.Cache[string, any], sessionKey string) AuthSession {
	if cache == nil {
		return nil
	}

	item := cache.Get(authSessionCachePrefix + sessionKey)

	if item == nil {
		return nil
	}

	session, ok := item.Value().(AuthSession)

	if !ok || isNilInterface(session) {
		return nil
	}

	if session.IsExpired() {
		return nil
	}

	return session
}

func authCacheSetSession(cache *ttlcache.Cache[string, any], sessionKey string, session AuthSession) {
	if cache == nil || isNilInterface(session) {
		return
	}

	cache.Set(authSessionCachePrefix+sessionKey, session, authSessionCacheTTL)
}

func authCacheGetUser(cache *ttlcache.Cache[string, any], userID string) AuthUser {
	if cache == nil {
		return nil
	}

	item := cache.Get(authUserCachePrefix + userID)

	if item == nil {
		return nil
	}

	user, ok := item.Value().(AuthUser)

	if !ok || isNilInterface(user) {
		return nil
	}

	return user
}

func authCacheSetUser(cache *ttlcache.Cache[string, any], userID string, user AuthUser) {
	if cache == nil || isNilInterface(user) {
		return
	}

	cache.Set(authUserCachePrefix+userID, user, authUserCacheTTL)
}
