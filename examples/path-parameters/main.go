package main

import (
	"fmt"
	"log"
	"net/http"

	rtr "github.com/dracory/rtr"
)

func main() {
	// Create a new router
	r := rtr.NewRouter()

	// Define all paths with their descriptions
	paths := []struct {
		path        string
		description string
	}{
		{"/users/123", "Basic path parameter example"},
		{"/posts/456/comments/789", "Multiple parameters example"},
		{"/articles/tech/101", "Optional parameter example (with ID)"},
		{"/articles/tech", "Optional parameter example (without ID)"},
		{"/profile/john/posts/42", "Get all parameters as a map"},
	}

	// Root route that shows all available paths
	r.AddRoute(rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/").
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
    <title>Path Parameters Example</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; max-width: 800px; margin: 0 auto; padding: 20px; }
        h1 { color: #2c3e50; }
        .endpoint { 
            background: #f8f9fa; 
            border-left: 4px solid #3498db; 
            padding: 10px 15px; 
            margin: 10px 0; 
            border-radius: 0 4px 4px 0;
        }
        .endpoint a { 
            color: #3498db; 
            text-decoration: none; 
            font-weight: bold;
        }
        .endpoint a:hover { text-decoration: underline; }
        .description { color: #7f8c8d; margin: 5px 0 0 0; }
    </style>
</head>
<body>
    <h1>Path Parameters Example</h1>
    <p>This example demonstrates different ways to use path parameters in the router.</p>
    <div class="endpoints">`)

			for _, p := range paths {
				_, _ = fmt.Fprintf(w, `
        <div class="endpoint">
            <a href="%s">%s</a>
            <p class="description">%s</p>
        </div>`, p.path, p.path, p.description)
			}

			_, _ = fmt.Fprint(w, `
    </div>
</body>
</html>`)
		}))

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

	// Start the server
	port := ":8080"
	fmt.Printf("Server running on http://localhost%s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET /users/123")
	fmt.Println("  GET /posts/456/comments/789")
	fmt.Println("  GET /articles/tech/101")
	fmt.Println("  GET /articles/tech")
	fmt.Println("  GET /profile/john/posts/42")

	log.Fatal(http.ListenAndServe(port, r))
}
