package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserMiddleware_NilGetUser(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true })

	mw := UserMiddleware(UserMiddlewareConfig{
		GetUser: nil,
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	mw.GetHandler()(next).ServeHTTP(w, r)

	if called {
		t.Fatal("next should NOT have been called with nil GetUser")
	}
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestUserMiddleware_WithRoles_NotAuthenticated(t *testing.T) {
	called := false
	redirectCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	config := UserMiddlewareConfig{
		GetUser:      func(r *http.Request) UserMiddlewareUser { return nil },
		RequireRoles: []string{"admin", "superuser"},
		OnNotAuthenticated: func(w http.ResponseWriter, r *http.Request) {
			redirectCalled = true
		},
	}

	mw := UserMiddleware(config)
	handler := mw.GetHandler()(next)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/admin", nil)
	handler.ServeHTTP(w, r)

	if called {
		t.Fatal("next handler should NOT have been called")
	}
	if !redirectCalled {
		t.Fatal("OnNotAuthenticated should have been called")
	}
}

func TestUserMiddleware_WithRoles_NotAuthorized(t *testing.T) {
	called := false
	notAuthorizedCalled := false
	user := &mockAuthUser{active: true, admin: false, superuser: false, registered: true}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	config := UserMiddlewareConfig{
		GetUser:      func(r *http.Request) UserMiddlewareUser { return user },
		RequireRoles: []string{"admin", "superuser"},
		OnNotAuthorized: func(w http.ResponseWriter, r *http.Request) {
			notAuthorizedCalled = true
		},
	}

	mw := UserMiddleware(config)
	handler := mw.GetHandler()(next)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/admin", nil)
	handler.ServeHTTP(w, r)

	if called {
		t.Fatal("next handler should NOT have been called")
	}
	if !notAuthorizedCalled {
		t.Fatal("OnNotAuthorized should have been called")
	}
}

func TestUserMiddleware_WithRoles_AdminUser(t *testing.T) {
	called := false
	user := &mockAuthUser{active: true, admin: true, superuser: false, registered: true}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	mw := UserMiddleware(UserMiddlewareConfig{
		GetUser:      func(r *http.Request) UserMiddlewareUser { return user },
		RequireRoles: []string{"admin", "superuser"},
	})

	handler := mw.GetHandler()(next)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/admin", nil)
	handler.ServeHTTP(w, r)

	if !called {
		t.Fatal("next handler should have been called for admin user")
	}
}

func TestUserMiddleware_WithRoles_Superuser(t *testing.T) {
	called := false
	user := &mockAuthUser{active: true, admin: false, superuser: true, registered: true}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	mw := UserMiddleware(UserMiddlewareConfig{
		GetUser:      func(r *http.Request) UserMiddlewareUser { return user },
		RequireRoles: []string{"admin", "superuser"},
	})

	handler := mw.GetHandler()(next)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/admin", nil)
	handler.ServeHTTP(w, r)

	if !called {
		t.Fatal("next handler should have been called for superuser")
	}
}

func TestUserMiddleware_WithRoles_RegistrationIncomplete(t *testing.T) {
	called := false
	regIncompleteCalled := false
	user := &mockAuthUser{active: true, admin: true, superuser: false, registered: false}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	mw := UserMiddleware(UserMiddlewareConfig{
		GetUser:             func(r *http.Request) UserMiddlewareUser { return user },
		RequireRoles:        []string{"admin", "superuser"},
		RegistrationEnabled: true,
		OnRegistrationIncomplete: func(w http.ResponseWriter, r *http.Request) {
			regIncompleteCalled = true
		},
	})

	handler := mw.GetHandler()(next)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/admin", nil)
	handler.ServeHTTP(w, r)

	if called {
		t.Fatal("next handler should NOT have been called")
	}
	if !regIncompleteCalled {
		t.Fatal("OnRegistrationIncomplete should have been called")
	}
}

func TestUserMiddleware_WithRoles_NotActive(t *testing.T) {
	called := false
	notActiveCalled := false
	user := &mockAuthUser{active: false, admin: true, superuser: false, registered: true}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	mw := UserMiddleware(UserMiddlewareConfig{
		GetUser:      func(r *http.Request) UserMiddlewareUser { return user },
		RequireRoles: []string{"admin", "superuser"},
		OnNotActive: func(w http.ResponseWriter, r *http.Request) {
			notActiveCalled = true
		},
	})

	handler := mw.GetHandler()(next)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/admin", nil)
	handler.ServeHTTP(w, r)

	if called {
		t.Fatal("next handler should NOT have been called")
	}
	if !notActiveCalled {
		t.Fatal("OnNotActive should have been called")
	}
}

func TestUserMiddleware_NotAuthenticated(t *testing.T) {
	called := false
	redirectCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	mw := UserMiddleware(UserMiddlewareConfig{
		GetUser: func(r *http.Request) UserMiddlewareUser { return nil },
		OnNotAuthenticated: func(w http.ResponseWriter, r *http.Request) {
			redirectCalled = true
		},
	})

	handler := mw.GetHandler()(next)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	handler.ServeHTTP(w, r)

	if called {
		t.Fatal("next handler should NOT have been called")
	}
	if !redirectCalled {
		t.Fatal("OnNotAuthenticated should have been called")
	}
}

func TestUserMiddleware_ActiveUser(t *testing.T) {
	called := false
	user := &mockAuthUser{active: true, registered: true}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	mw := UserMiddleware(UserMiddlewareConfig{
		GetUser: func(r *http.Request) UserMiddlewareUser { return user },
	})

	handler := mw.GetHandler()(next)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	handler.ServeHTTP(w, r)

	if !called {
		t.Fatal("next handler should have been called for active registered user")
	}
}

func TestUserMiddleware_NotActive(t *testing.T) {
	called := false
	notActiveCalled := false
	user := &mockAuthUser{active: false, registered: true}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	mw := UserMiddleware(UserMiddlewareConfig{
		GetUser: func(r *http.Request) UserMiddlewareUser { return user },
		OnNotActive: func(w http.ResponseWriter, r *http.Request) {
			notActiveCalled = true
		},
	})

	handler := mw.GetHandler()(next)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	handler.ServeHTTP(w, r)

	if called {
		t.Fatal("next handler should NOT have been called")
	}
	if !notActiveCalled {
		t.Fatal("OnNotActive should have been called")
	}
}

func TestUserMiddleware_RegistrationIncomplete(t *testing.T) {
	called := false
	regIncompleteCalled := false
	user := &mockAuthUser{active: true, registered: false}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	mw := UserMiddleware(UserMiddlewareConfig{
		GetUser:             func(r *http.Request) UserMiddlewareUser { return user },
		RegistrationEnabled: true,
		RegistrationPaths:   []string{"/profile", "/register"},
		OnRegistrationIncomplete: func(w http.ResponseWriter, r *http.Request) {
			regIncompleteCalled = true
		},
	})

	handler := mw.GetHandler()(next)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	handler.ServeHTTP(w, r)

	if called {
		t.Fatal("next handler should NOT have been called")
	}
	if !regIncompleteCalled {
		t.Fatal("OnRegistrationIncomplete should have been called")
	}
}

func TestUserMiddleware_RegistrationIncompleteOnExemptPath(t *testing.T) {
	called := false
	user := &mockAuthUser{active: true, registered: false}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	mw := UserMiddleware(UserMiddlewareConfig{
		GetUser:             func(r *http.Request) UserMiddlewareUser { return user },
		RegistrationEnabled: true,
		RegistrationPaths:   []string{"/profile", "/register"},
	})

	handler := mw.GetHandler()(next)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/profile", nil)
	handler.ServeHTTP(w, r)

	if !called {
		t.Fatal("next handler should have been called on exempt path")
	}
}

func TestUserMiddleware_RegistrationIncompleteDisabled(t *testing.T) {
	called := false
	user := &mockAuthUser{active: true, registered: false}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	mw := UserMiddleware(UserMiddlewareConfig{
		GetUser:             func(r *http.Request) UserMiddlewareUser { return user },
		RegistrationEnabled: false,
		RegistrationPaths:   []string{"/profile", "/register"},
	})

	handler := mw.GetHandler()(next)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	handler.ServeHTTP(w, r)

	if !called {
		t.Fatal("next handler should be called when registration is disabled and incomplete")
	}
}

func TestIsOnRegistrationPath(t *testing.T) {
	tests := []struct {
		requestPath string
		paths       []string
		expected    bool
	}{
		{"/profile", []string{"/profile", "/register"}, true},
		{"/register", []string{"/profile", "/register"}, true},
		{"/dashboard", []string{"/profile", "/register"}, false},
		{"/profile/", []string{"/profile"}, true},
		{"/PROFILE", []string{"/profile"}, false},
		{"", []string{"/profile"}, false},
	}

	for _, tt := range tests {
		result := isOnRegistrationPath(tt.requestPath, tt.paths)
		if result != tt.expected {
			t.Errorf("isOnRegistrationPath(%q, %v) = %v, want %v", tt.requestPath, tt.paths, result, tt.expected)
		}
	}
}

func TestUserMiddleware_WithRoles_TypedNilUser(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	mw := UserMiddleware(UserMiddlewareConfig{
		GetUser: func(r *http.Request) UserMiddlewareUser {
			var user *mockAuthUser = nil
			return user
		},
		RequireRoles: []string{"admin", "superuser"},
	})

	handler := mw.GetHandler()(next)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/admin", nil)
	handler.ServeHTTP(w, r)

	if called {
		t.Fatal("next handler should NOT have been called for typed-nil user")
	}
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestUserMiddleware_WithRoles_DefaultResponseWhenCallbackNil(t *testing.T) {
	tests := []struct {
		name       string
		user       *mockAuthUser
		wantStatus int
	}{
		{"not authenticated", nil, http.StatusUnauthorized},
		{"registration incomplete", &mockAuthUser{active: true, admin: true, registered: false}, http.StatusForbidden},
		{"not active", &mockAuthUser{active: false, admin: true, registered: true}, http.StatusForbidden},
		{"not authorized", &mockAuthUser{active: true, admin: false, superuser: false, registered: true}, http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Fatal("next handler should NOT have been called")
			})

			mw := UserMiddleware(UserMiddlewareConfig{
				GetUser:             func(r *http.Request) UserMiddlewareUser { return tt.user },
				RequireRoles:        []string{"admin", "superuser"},
				RegistrationEnabled: true,
			})

			handler := mw.GetHandler()(next)
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/admin", nil)
			handler.ServeHTTP(w, r)

			if w.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestUserMiddleware_TypedNilUser(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	mw := UserMiddleware(UserMiddlewareConfig{
		GetUser: func(r *http.Request) UserMiddlewareUser {
			var user *mockAuthUser = nil
			return user
		},
	})

	handler := mw.GetHandler()(next)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	handler.ServeHTTP(w, r)

	if called {
		t.Fatal("next handler should NOT have been called for typed-nil user")
	}
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestUserMiddleware_DefaultResponseWhenCallbackNil(t *testing.T) {
	tests := []struct {
		name       string
		user       *mockAuthUser
		wantStatus int
	}{
		{"not authenticated", nil, http.StatusUnauthorized},
		{"not active", &mockAuthUser{active: false, registered: true}, http.StatusForbidden},
		{"registration incomplete", &mockAuthUser{active: true, registered: false}, http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Fatal("next handler should NOT have been called")
			})

			mw := UserMiddleware(UserMiddlewareConfig{
				GetUser:             func(r *http.Request) UserMiddlewareUser { return tt.user },
				RegistrationEnabled: true,
				RegistrationPaths:   []string{"/profile"},
			})

			handler := mw.GetHandler()(next)
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
			handler.ServeHTTP(w, r)

			if w.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestUserMiddleware_RequireRoles_UserWithoutHasRole(t *testing.T) {
	called := false
	notAuthorizedCalled := false
	user := &mockAuthUserNoRoles{active: true, registered: true}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	mw := UserMiddleware(UserMiddlewareConfig{
		GetUser:      func(r *http.Request) UserMiddlewareUser { return user },
		RequireRoles: []string{"admin"},
		OnNotAuthorized: func(w http.ResponseWriter, r *http.Request) {
			notAuthorizedCalled = true
		},
	})

	handler := mw.GetHandler()(next)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(w, r)

	if called {
		t.Fatal("next handler should NOT have been called when user lacks HasRole")
	}
	if !notAuthorizedCalled {
		t.Fatal("OnNotAuthorized should have been called")
	}
}

func TestUserMiddleware_RequireRoles(t *testing.T) {
	tests := []struct {
		name         string
		user         *mockAuthUser
		requireRoles []string
		wantStatus   int
		wantNext     bool
	}{
		{
			name:         "has required role",
			user:         &mockAuthUser{active: true, registered: true, admin: true},
			requireRoles: []string{"admin"},
			wantStatus:   http.StatusOK,
			wantNext:     true,
		},
		{
			name:         "missing required role",
			user:         &mockAuthUser{active: true, registered: true, admin: false},
			requireRoles: []string{"admin"},
			wantStatus:   http.StatusForbidden,
			wantNext:     false,
		},
		{
			name:         "any of multiple roles passes",
			user:         &mockAuthUser{active: true, registered: true, superuser: true},
			requireRoles: []string{"admin", "superuser"},
			wantStatus:   http.StatusOK,
			wantNext:     true,
		},
		{
			name:         "none of multiple roles",
			user:         &mockAuthUser{active: true, registered: true},
			requireRoles: []string{"admin", "superuser"},
			wantStatus:   http.StatusForbidden,
			wantNext:     false,
		},
		{
			name:         "no roles required",
			user:         &mockAuthUser{active: true, registered: true},
			requireRoles: nil,
			wantStatus:   http.StatusOK,
			wantNext:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
			})

			mw := UserMiddleware(UserMiddlewareConfig{
				GetUser:      func(r *http.Request) UserMiddlewareUser { return tt.user },
				RequireRoles: tt.requireRoles,
			})

			handler := mw.GetHandler()(next)
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			handler.ServeHTTP(w, r)

			if called != tt.wantNext {
				t.Fatalf("next called = %v, want %v", called, tt.wantNext)
			}
			if w.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}
