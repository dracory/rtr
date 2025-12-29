package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/rtr"
)

func TestPathParameters(t *testing.T) {
	// Create a new router instance with test routes
	r := setupTestRouter()

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "GET /users/123",
			method:         http.MethodGet,
			path:           "/users/123",
			expectedStatus: http.StatusOK,
			expectedBody:   "User ID: 123",
		},
		{
			name:           "GET /posts/456/comments/789",
			method:         http.MethodGet,
			path:           "/posts/456/comments/789",
			expectedStatus: http.StatusOK,
			expectedBody:   "Post ID: 456, Comment ID: 789",
		},
		{
			name:           "GET /articles/tech/101",
			method:         http.MethodGet,
			path:           "/articles/tech/101",
			expectedStatus: http.StatusOK,
			expectedBody:   "Category: tech, Article ID: 101",
		},
		{
			name:           "GET /articles/tech (optional ID)",
			method:         http.MethodGet,
			path:           "/articles/tech",
			expectedStatus: http.StatusOK,
			expectedBody:   "Category: tech (no article ID provided)",
		},
		{
			name:           "GET /profile/john/posts/42 (all params)",
			method:         http.MethodGet,
			path:           "/profile/john/posts/42",
			expectedStatus: http.StatusOK,
			expectedBody:   "All parameters: map[postID:42 username:john]",
		},
		{
			name:           "GET /nonexistent",
			method:         http.MethodGet,
			path:           "/nonexistent",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404 page not found\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
			}

			// For 404 responses, the exact body might vary by Go version
			if tc.expectedStatus != http.StatusNotFound {
				if body := rr.Body.String(); body != tc.expectedBody {
					t.Errorf("handler returned unexpected body: got %v want %v",
						body, tc.expectedBody)
				}
			}
		})
	}
}

// setupTestRouter creates a router instance with test routes
func setupTestRouter() rtr.RouterInterface {
	r := rtr.NewRouter()

	// Basic path parameter example
	r.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/users/:id").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			id := rtr.MustGetParam(r, "id")
			_, _ = fmt.Fprintf(w, "User ID: %s", id)
		}))

	// Multiple parameters example
	r.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/posts/:postID/comments/:commentID").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			postID := rtr.MustGetParam(r, "postID")
			commentID := rtr.MustGetParam(r, "commentID")
			_, _ = fmt.Fprintf(w, "Post ID: %s, Comment ID: %s", postID, commentID)
		}))

	// Optional parameter example
	r.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/articles/:category/:id?").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			category := rtr.MustGetParam(r, "category")
			if id, exists := rtr.GetParam(r, "id"); exists {
				_, _ = fmt.Fprintf(w, "Category: %s, Article ID: %s", category, id)
			} else {
				_, _ = fmt.Fprintf(w, "Category: %s (no article ID provided)", category)
			}
		}))

	// Get all parameters as a map
	r.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/profile/:username/posts/:postID").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			params := rtr.GetParams(r)
			_, _ = fmt.Fprintf(w, "All parameters: %v", params)
		}))

	return r
}
