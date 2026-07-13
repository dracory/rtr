package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jellydator/ttlcache/v3"
)

// --- Mock types ---

type mockAuthSession struct {
	userID  string
	expired bool
}

func (s *mockAuthSession) GetUserID() string { return s.userID }
func (s *mockAuthSession) IsExpired() bool   { return s.expired }

type mockAuthUser struct {
	active     bool
	admin      bool
	superuser  bool
	registered bool
}

func (u *mockAuthUser) IsActive() bool                { return u.active }
func (u *mockAuthUser) IsAdministrator() bool         { return u.admin }
func (u *mockAuthUser) IsSuperuser() bool             { return u.superuser }
func (u *mockAuthUser) IsRegistrationCompleted() bool { return u.registered }
func (u *mockAuthUser) HasRole(role string) bool {
	switch role {
	case "admin":
		return u.admin
	case "superuser":
		return u.superuser
	default:
		return false
	}
}

// mockAuthUserNoRoles implements UserMiddlewareUser but NOT UserWithRole.
type mockAuthUserNoRoles struct {
	active     bool
	registered bool
}

func (u *mockAuthUserNoRoles) IsActive() bool                { return u.active }
func (u *mockAuthUserNoRoles) IsRegistrationCompleted() bool { return u.registered }

type mockAuthSessionStore struct {
	session AuthSession
	err     error
}

func (s *mockAuthSessionStore) SessionFindByKey(ctx context.Context, key string) (AuthSession, error) {
	return s.session, s.err
}

type mockAuthUserStore struct {
	user AuthUser
	err  error
}

func (s *mockAuthUserStore) UserFindByID(ctx context.Context, id string) (AuthUser, error) {
	return s.user, s.err
}

type mockAuthLogger struct {
	lastMsg string
}

func (l *mockAuthLogger) Error(msg string, args ...any) {
	l.lastMsg = msg
}

type contextKey string

const testUserKey contextKey = "user"
const testSessionKey contextKey = "session"

// typedNilAuthSessionStore returns a typed nil AuthSession from the store.
type typedNilAuthSessionStore struct{}

func (s *typedNilAuthSessionStore) SessionFindByKey(ctx context.Context, key string) (AuthSession, error) {
	var session *mockAuthSession = nil
	return session, nil
}

// typedNilAuthUserStore returns a typed nil AuthUser from the store.
type typedNilAuthUserStore struct{}

func (s *typedNilAuthUserStore) UserFindByID(ctx context.Context, id string) (AuthUser, error) {
	var user *mockAuthUser = nil
	return user, nil
}

// --- Auth Middleware Tests ---

func TestAuthMiddleware_NoSessionKey(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	mw := AuthMiddleware(AuthMiddlewareConfig{
		SessionStore:      &mockAuthSessionStore{},
		UserStore:         &mockAuthUserStore{},
		Logger:            &mockAuthLogger{},
		ContextKeyUser:    testUserKey,
		ContextKeySession: testSessionKey,
		CookieName:        "session",
	})

	handler := mw.GetHandler()(next)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(w, r)

	if !called {
		t.Fatal("next handler should have been called when no session key")
	}
}

func TestAuthMiddleware_WithValidSession(t *testing.T) {
	user := &mockAuthUser{active: true, registered: true}
	session := &mockAuthSession{userID: "user-123", expired: false}

	called := false
	var ctxUser any
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		ctxUser = r.Context().Value(testUserKey)
	})

	mw := AuthMiddleware(AuthMiddlewareConfig{
		SessionStore:      &mockAuthSessionStore{session: session},
		UserStore:         &mockAuthUserStore{user: user},
		Logger:            &mockAuthLogger{},
		ContextKeyUser:    testUserKey,
		ContextKeySession: testSessionKey,
		CookieName:        "session",
	})

	handler := mw.GetHandler()(next)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{Name: "session", Value: "valid-key"})
	handler.ServeHTTP(w, r)

	if !called {
		t.Fatal("next handler should have been called")
	}
	if ctxUser == nil {
		t.Fatal("user should be in context")
	}
}

func TestAuthMiddleware_ExpiredSession(t *testing.T) {
	session := &mockAuthSession{userID: "user-123", expired: true}

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	mw := AuthMiddleware(AuthMiddlewareConfig{
		SessionStore:      &mockAuthSessionStore{session: session},
		UserStore:         &mockAuthUserStore{},
		Logger:            &mockAuthLogger{},
		ContextKeyUser:    testUserKey,
		ContextKeySession: testSessionKey,
		CookieName:        "session",
	})

	handler := mw.GetHandler()(next)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{Name: "session", Value: "expired-key"})
	handler.ServeHTTP(w, r)

	if !called {
		t.Fatal("next handler should have been called for expired session")
	}
}

func TestAuthMiddleware_SessionStoreNotEnabled(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true })

	mw := AuthMiddleware(AuthMiddlewareConfig{
		SessionStore:      nil,
		UserStore:         &mockAuthUserStore{},
		ContextKeyUser:    testUserKey,
		ContextKeySession: testSessionKey,
		CookieName:        "session",
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	mw.GetHandler()(next).ServeHTTP(w, r)

	if called {
		t.Fatal("next should NOT have been called with missing SessionStore")
	}
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestAuthMiddleware_UserStoreNotInitialized(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true })

	mw := AuthMiddleware(AuthMiddlewareConfig{
		SessionStore:      &mockAuthSessionStore{},
		UserStore:         nil,
		ContextKeyUser:    testUserKey,
		ContextKeySession: testSessionKey,
		CookieName:        "session",
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	mw.GetHandler()(next).ServeHTTP(w, r)

	if called {
		t.Fatal("next should NOT have been called with missing UserStore")
	}
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestAuthMiddleware_EmptyCookieName(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true })

	mw := AuthMiddleware(AuthMiddlewareConfig{
		SessionStore:      &mockAuthSessionStore{},
		UserStore:         &mockAuthUserStore{},
		ContextKeyUser:    testUserKey,
		ContextKeySession: testSessionKey,
		CookieName:        "",
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	mw.GetHandler()(next).ServeHTTP(w, r)

	if called {
		t.Fatal("next should NOT have been called with empty CookieName")
	}
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestAuthMiddleware_MemoryCache(t *testing.T) {
	user := &mockAuthUser{active: true, registered: true}
	session := &mockAuthSession{userID: "user-123", expired: false}

	cache := ttlcache.New[string, any]()
	go cache.Start()
	defer cache.Stop()

	authCacheSetSession(cache, "cached-key", session)
	authCacheSetUser(cache, "user-123", user)
	time.Sleep(10 * time.Millisecond)

	called := false
	var ctxUser any
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		ctxUser = r.Context().Value(testUserKey)
	})

	mw := AuthMiddleware(AuthMiddlewareConfig{
		SessionStore:      &mockAuthSessionStore{}, // won't be called due to cache
		UserStore:         &mockAuthUserStore{},    // won't be called due to cache
		Logger:            &mockAuthLogger{},
		MemoryCache:       cache,
		ContextKeyUser:    testUserKey,
		ContextKeySession: testSessionKey,
		CookieName:        "session",
	})

	handler := mw.GetHandler()(next)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{Name: "session", Value: "cached-key"})
	handler.ServeHTTP(w, r)

	if !called {
		t.Fatal("next handler should have been called")
	}
	if ctxUser == nil {
		t.Fatal("user should be in context from cache")
	}
}

func TestAuthMiddleware_TypedNilSession(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	mw := AuthMiddleware(AuthMiddlewareConfig{
		SessionStore:      &typedNilAuthSessionStore{},
		UserStore:         &mockAuthUserStore{},
		Logger:            &mockAuthLogger{},
		ContextKeyUser:    testUserKey,
		ContextKeySession: testSessionKey,
		CookieName:        "session",
	})

	handler := mw.GetHandler()(next)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{Name: "session", Value: "typed-nil-session"})
	handler.ServeHTTP(w, r)

	if !called {
		t.Fatal("next handler should have been called for typed-nil session")
	}
}

func TestAuthMiddleware_TypedNilUser(t *testing.T) {
	session := &mockAuthSession{userID: "user-123", expired: false}

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	mw := AuthMiddleware(AuthMiddlewareConfig{
		SessionStore:      &mockAuthSessionStore{session: session},
		UserStore:         &typedNilAuthUserStore{},
		Logger:            &mockAuthLogger{},
		ContextKeyUser:    testUserKey,
		ContextKeySession: testSessionKey,
		CookieName:        "session",
	})

	handler := mw.GetHandler()(next)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{Name: "session", Value: "typed-nil-user"})
	handler.ServeHTTP(w, r)

	if !called {
		t.Fatal("next handler should have been called for typed-nil user")
	}
}

func TestAuthMiddleware_StoresSessionObjectInContext(t *testing.T) {
	session := &mockAuthSession{userID: "user-123", expired: false}
	user := &mockAuthUser{active: true, registered: true}

	var ctxSession any
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxSession = r.Context().Value(testSessionKey)
	})

	mw := AuthMiddleware(AuthMiddlewareConfig{
		SessionStore:      &mockAuthSessionStore{session: session},
		UserStore:         &mockAuthUserStore{user: user},
		Logger:            &mockAuthLogger{},
		ContextKeyUser:    testUserKey,
		ContextKeySession: testSessionKey,
		CookieName:        "session",
	})

	handler := mw.GetHandler()(next)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{Name: "session", Value: "valid-key"})
	handler.ServeHTTP(w, r)

	if ctxSession == nil {
		t.Fatal("session object should be in context")
	}
	if _, ok := ctxSession.(AuthSession); !ok {
		t.Fatalf("expected AuthSession in context, got %T", ctxSession)
	}
}
