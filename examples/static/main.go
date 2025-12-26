package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dracory/rtr"
)

func main() {
	// Create a new router
	router := rtr.NewRouter()

	// Add a static file server route
	// This will serve files from the "./static" directory
	router.AddRoute(rtr.GetStatic("/static/*", func(w http.ResponseWriter, r *http.Request) string {
		return "./static" // Return the path to the static directory
	}))

	// Add a root route with links to static files
	router.AddRoute(rtr.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head>
    <title>Static File Server Example</title>
    <style>
        body { 
            font-family: Arial, sans-serif; 
            line-height: 1.6; 
            max-width: 800px; 
            margin: 0 auto; 
            padding: 20px; 
        }
        h1 { color: #2c3e50; }
        .file-link { 
            display: block; 
            background: #f8f9fa; 
            border-left: 4px solid #3498db; 
            padding: 10px 15px; 
            margin: 10px 0; 
            border-radius: 0 4px 4px 0;
            text-decoration: none;
            color: #3498db;
            font-weight: bold;
        }
        .file-link:hover { 
            background: #e9ecef; 
            text-decoration: underline; 
        }
        .description { 
            color: #7f8c8d; 
            margin: 5px 0 0 0; 
            font-size: 0.9em;
        }
    </style>
</head>
<body>
    <h1>Static File Server Example</h1>
    <p>This example demonstrates how to serve static files using the rtr router.</p>
    
    <h2>Available Static Files</h2>
    
    <a href="/static/style.css" class="file-link">
        /static/style.css
        <span class="description">CSS stylesheet file</span>
    </a>
    
    <a href="/static/script.js" class="file-link">
        /static/script.js
        <span class="description">JavaScript file</span>
    </a>
    
    <a href="/static/data.json" class="file-link">
        /static/data.json
        <span class="description">JSON data file</span>
    </a>
    
    <a href="/static/image.png" class="file-link">
        /static/image.png
        <span class="description">Image file (if exists)</span>
    </a>
    
    <h2>Features</h2>
    <ul>
        <li><strong>Automatic Content-Type detection</strong> - Files are served with appropriate MIME types</li>
        <li><strong>Security</strong> - Directory traversal attacks are prevented</li>
        <li><strong>404 handling</strong> - Non-existent files return proper 404 responses</li>
        <li><strong>Integration</strong> - Works seamlessly with existing rtr routing features</li>
    </ul>
    
    <h2>Usage</h2>
    <p>Create static files in the <code>./static</code> directory and access them via the <code>/static/</code> URL prefix.</p>
    
    <h2>Code Example</h2>
    <pre><code>// Add static file server
router.AddRoute(rtr.GetStatic("/static/*", func(w http.ResponseWriter, r *http.Request) string {
    return "./static" // Path to static directory
}))</code></pre>
</body>
</html>`)
	}))

	// Start the server
	port := ":8080"
	fmt.Printf("Static file server running on http://localhost%s\n", port)
	fmt.Printf("Access static files at: http://localhost%s/static/\n", port)
	fmt.Printf("Try accessing: http://localhost%s/static/style.css\n", port)

	log.Fatal(http.ListenAndServe(port, router))
}
