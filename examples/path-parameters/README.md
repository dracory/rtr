# Path Parameters Example

This example demonstrates how to use path parameters in the Dracory Router.

## Features

- Root route (`/`) with a web interface listing all available endpoints
- Basic path parameter extraction (`/users/:id`)
- Multiple parameters in a single path (`/posts/:postID/comments/:commentID`)
- Optional parameters using `?` suffix (`/articles/:category/:id?`)
- Accessing all parameters as a map (`/profile/:username/posts/:postID`)

## Running the Example

1. Make sure you have Go installed
2. Run the example:
   ```bash
   cd examples/path-parameters
   go run main.go
   ```
3. Open your browser to `http://localhost:8080` to see a list of all available endpoints with descriptions
4. Alternatively, test the endpoints directly:
   - `GET http://localhost:8080/users/123`
   - `GET http://localhost:8080/posts/456/comments/789`
   - `GET http://localhost:8080/articles/tech/101`
   - `GET http://localhost:8080/articles/tech` (optional parameter)
   - `GET http://localhost:8080/profile/john/posts/42` (all parameters)

## Code Explanation

### Basic Parameter
```go
r.Get("/users/:id", func(w http.ResponseWriter, r *http.Request) {
    id := rtr.MustGetParam(r, "id")
    fmt.Fprintf(w, "User ID: %s", id)
})
```

### Multiple Parameters
```go
r.Get("/posts/:postID/comments/:commentID", func(w http.ResponseWriter, r *http.Request) {
    postID := rtr.MustGetParam(r, "postID")
    commentID := rtr.MustGetParam(r, "commentID")
    fmt.Fprintf(w, "Post ID: %s, Comment ID: %s", postID, commentID)
})
```

### Optional Parameter
```go
r.Get("/articles/:category/:id?", func(w http.ResponseWriter, r *http.Request) {
    category := rtr.MustGetParam(r, "category")
    if id, exists := rtr.GetParam(r, "id"); exists {
        fmt.Fprintf(w, "Category: %s, Article ID: %s", category, id)
    } else {
        fmt.Fprintf(w, "Category: %s (no article ID provided)", category)
    }
})
```

### Get All Parameters
```go
r.Get("/profile/:username/posts/:postID", func(w http.ResponseWriter, r *http.Request) {
    params := rtr.GetParams(r)
    fmt.Fprintf(w, "All parameters: %v", params)
})
```
