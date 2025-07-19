# Named Middleware Example

This example demonstrates the named middleware functionality of the Dracory HTTP Router package, similar to the old gouniverse/router pattern.

## Features Demonstrated

- Creating named middleware with descriptive names
- Using anonymous middleware for simple cases
- Adding named middleware to routes, groups, and router
- Middleware execution order and chaining
- Converting between regular and named middleware
- Using declarative RouteConfig with named middleware
- Better debugging and documentation through named middleware

## Named Middleware Benefits

- **Better Debugging**: Named middleware provides clear identification in logs and debugging
- **Documentation**: Middleware names serve as inline documentation
- **Reusability**: Named middleware can be easily shared across routes and groups
- **Backward Compatibility**: Works alongside existing middleware functions

## Running the Example

1. Make sure you have Go installed (1.16 or later)
2. Navigate to this directory
3. Run the example:
   ```bash
   go run main.go
   ```
4. Open your browser or use `curl` to test the endpoints:
   - `http://localhost:8080/public` - No authentication required
   - `http://localhost:8080/protected` - Requires Authorization header
   - `http://localhost:8080/admin` - Requires Authorization header and admin role
   - `http://localhost:8080/api/users` - API group with logging middleware
   - `http://localhost:8080/v1/status` - Version group with rate limiting

## Endpoints

- `GET /public` - Public endpoint with no middleware
- `GET /protected` - Protected endpoint with authentication middleware
- `GET /admin` - Admin endpoint with authentication and authorization middleware
- `GET /api/users` - API endpoint with logging middleware
- `GET /v1/status` - Versioned endpoint with rate limiting middleware

## Testing Authentication

To test protected endpoints, include an Authorization header:

```bash
# This will work
curl -H "Authorization: Bearer token123" http://localhost:8080/protected

# This will fail with 401 Unauthorized
curl http://localhost:8080/protected
```

## Middleware Examples

The example demonstrates several types of middleware:

1. **Authentication Middleware**: Checks for Authorization header
2. **Authorization Middleware**: Checks for admin role
3. **Logging Middleware**: Logs request details
4. **Rate Limiting Middleware**: Simulates rate limiting
5. **Anonymous Middleware**: Simple middleware without names

Each middleware is created with a descriptive name that helps with debugging and documentation.
