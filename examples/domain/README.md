# Domain-Based Routing Example

This example demonstrates how to use the router's domain-based routing feature to handle different domains with different routes.

## Features

- Routes requests based on the `Host` header
- Supports multiple domains with different route handlers
- Includes a default handler for unmatched domains
- Demonstrates serving both API and admin interfaces on different domains

## Setup

1. Add the following entries to your `/etc/hosts` (Unix/Mac) or `C:\\Windows\\System32\\drivers\\etc\\hosts` (Windows) file:

```
127.0.0.1 api.example.com admin.example.com
```

2. Run the example:

```bash
go run main.go
```

## Testing

Open the following URLs in your browser or use `curl`:

- API Endpoints:
  - `http://api.example.com:8080/status` - API status
  - `http://api.example.com:8080/users` - List users

- Admin Interface:
  - `http://admin.example.com:8081/` - Admin dashboard

- Unmatched domains will return a 404 with a helpful message.

## Code Overview

The example creates two domains:

1. `api.example.com` - Handles API requests
   - `/status` - Returns API status
   - `/users` - Returns a list of users

2. `admin.example.com` - Serves an admin interface
   - `/` - Admin dashboard

## How It Works

1. The router checks the `Host` header of incoming requests
2. Routes are matched based on the domain first, then the path
3. If no matching domain is found, the `NotFoundHandler` is used
