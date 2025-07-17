# Basic Router Example

This example demonstrates the basic usage of the Dracory HTTP Router package.

## Features Demonstrated

- Creating a new router
- Adding routes with shortcut methods
- Creating route groups
- Nested route groups
- Path parameters
- Basic request handling

## Running the Example

1. Make sure you have Go installed (1.16 or later)
2. Navigate to this directory
3. Run the example:
   ```bash
   go run main.go
   ```
4. Open your browser or use `curl` to test the endpoints:
   - `http://localhost:8080/hello`
   - `http://localhost:8080/api/status`
   - `http://localhost:8080/api/users`
   - `http://localhost:8080/api/users/123`

## Endpoints

- `GET /hello` - Simple hello world endpoint
- `GET /api/status` - API status check
- `GET /api/users` - List users
- `GET /api/users/:id` - Get user by ID

## Code Structure

The example shows:
1. Router initialization
2. Basic route definition
3. Route grouping
4. Nested routes
5. Path parameter handling
