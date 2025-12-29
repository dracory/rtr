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

	// Add main root route with web interface listing all available endpoints
	r.AddRoute(rtr.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head>
    <title>Domain Router Example</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; max-width: 800px; margin: 0 auto; padding: 20px; }
        h1 { color: #2c3e50; }
        .domain-section {
            background: #f8f9fa;
            border: 1px solid #dee2e6;
            border-radius: 8px;
            padding: 20px;
            margin: 20px 0;
        }
        .domain-section h2 { color: #495057; margin-top: 0; }
        .endpoint { 
            background: #ffffff; 
            border-left: 4px solid #3498db; 
            padding: 10px 15px; 
            margin: 10px 0; 
            border-radius: 0 4px 4px 0;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }
        .endpoint a { 
            color: #3498db; 
            text-decoration: none; 
            font-weight: bold;
        }
        .endpoint a:hover { text-decoration: underline; }
        .description { color: #7f8c8d; margin: 5px 0 0 0; }
        .note { background: #fff3cd; border: 1px solid #ffeaa7; padding: 10px; border-radius: 4px; margin: 15px 0; }
    </style>
</head>
<body>
    <h1>Domain Router Example</h1>
    <p>This example demonstrates domain-based routing with different endpoints for different domains.</p>
    
    <div class="note">
        <strong>Note:</strong> To test domain-specific routes, add the following to your hosts file:<br>
        <code>127.0.0.1 api.example.com admin.example.com</code>
    </div>
    
    <div class="domain-section">
        <h2>API Domain (api.example.com:8080)</h2>
        <div class="endpoint">
            <a href="http://api.example.com:8080/status" target="_blank">GET /status</a>
            <p class="description">API status check endpoint</p>
        </div>
        <div class="endpoint">
            <a href="http://api.example.com:8080/users" target="_blank">GET /users</a>
            <p class="description">List of users endpoint</p>
        </div>
    </div>
    
    <div class="domain-section">
        <h2>Admin Domain (admin.example.com:8081)</h2>
        <div class="endpoint">
            <a href="http://admin.example.com:8081/" target="_blank">GET /</a>
            <p class="description">Admin panel home page</p>
        </div>
    </div>
    
    <div class="domain-section">
        <h2>Local Testing (localhost:8080)</h2>
        <p>If you haven't set up the hosts file, you can test the catch-all routes:</p>
        <div class="endpoint">
            <a href="/nonexistent">GET /nonexistent</a>
            <p class="description">Test catch-all route (404 response)</p>
        </div>
    </div>
</body>
</html>`)
	}))

	// Create domains
	apiDomain := rtr.NewDomain("api.example.com", "localhost:8080")
	adminDomain := rtr.NewDomain("admin.example.com", "localhost:8081")

	// Add routes to the API domain
	apiDomain.AddRoute(rtr.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))

	apiDomain.AddRoute(rtr.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`["user1", "user2"]`))
	}))

	// Add routes to the admin domain
	adminDomain.AddRoute(rtr.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Admin Panel</title>
		</head>
		<body>
			<h1>Welcome to Admin Panel</h1>
			<p>This is the admin interface for example.com</p>
		</body>
		</html>
		`)
	}))

	// Add catch-all route for API domain
	apiDomain.AddRoute(rtr.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "Not Found", "message": "The requested resource was not found on this server"}`))
	}))

	// Add catch-all route for Admin domain
	adminDomain.AddRoute(rtr.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>404 Not Found</title>
		</head>
		<body>
			<h1>404 Not Found</h1>
			<p>The requested page was not found on this server.</p>
			<p><a href="/">Return to Admin Panel</a></p>
		</body>
		</html>
		`)
	}))

	// Add domains to the router
	r.AddDomain(apiDomain)
	r.AddDomain(adminDomain)

	// Start the server
	fmt.Println("Starting server on :8080 (api.example.com) and :8081 (admin.example.com)")
	fmt.Println("To test, add the following to your /etc/hosts or C:\\Windows\\System32\\drivers\\etc\\hosts file:")
	fmt.Println("127.0.0.1 api.example.com admin.example.com")
	fmt.Println("Then visit:")
	fmt.Println("- http://api.example.com:8080/status")
	fmt.Println("- http://api.example.com:8080/users")
	fmt.Println("- http://admin.example.com:8081/")

	// Start the server on multiple ports
	go func() {
		log.Fatal(http.ListenAndServe(":8080", r))
	}()

	log.Fatal(http.ListenAndServe(":8081", r))
}
