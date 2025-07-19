# Declarative Routing System - Multiple Configuration Sources

## Overview

This document outlines a truly declarative routing and middleware system that can load configuration from **any serializable source**: SQL databases, JSON files, YAML files, environment variables, REST APIs, etc. All configuration is stored as pure data with **handler names as strings**, not function references.

## Core Concept

Store **handler names as strings** in any configuration source, not function references. Functions are resolved at runtime from a registry.

## Configuration Structure

The same declarative structure can be represented in multiple formats:

### SQL Database

```sql
CREATE TABLE routes (
    id          INTEGER PRIMARY KEY,
    status      TEXT NOT NULL DEFAULT 'enabled',
    type        TEXT NOT NULL,                     -- 'middleware' | 'route' | 'group' | 'domain'
    parent_id   INTEGER,
    path        TEXT,
    handler     TEXT,
    method      TEXT,
    name        TEXT,
    middlewares TEXT,                              -- JSON array of middleware names
    config_json TEXT,
    sequence    INTEGER DEFAULT 0,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### JSON Configuration

```json
{
  "routes": [
    {
      "id": 1,
      "status": "enabled",
      "type": "router",
      "name": "Main API Router",
      "path": "/api",
      "middlewares": ["cors_middleware", "auth_middleware", "rate_limit_middleware"],
      "config": {
        "cors": {"origins": ["*"], "methods": ["GET", "POST"]},
        "auth": {"secret": "jwt-secret", "skip_paths": ["/health"]},
        "rate_limit": {"requests_per_minute": 100, "burst": 10}
      }
    },
    {
      "id": 2,
      "status": "enabled",
      "type": "group",
      "parent_id": 1,
      "name": "Users API",
      "path": "/users",
      "middlewares": ["logging_middleware"],
      "config": {
        "logging": {"level": "info", "format": "json"}
      }
    },
    {
      "id": 3,
      "status": "enabled",
      "type": "route",
      "parent_id": 2,
      "path": "/",
      "method": "GET",
      "handler": "list_users",
      "name": "List Users"
    }
  ]
}
```

### YAML Configuration

```yaml
routes:
  - id: 1
    status: enabled
    type: router
    name: "Main API Router"
    path: "/api"
    middlewares:
      - cors_middleware
      - auth_middleware
      - rate_limit_middleware
    config:
      cors:
        origins: ["*"]
        methods: ["GET", "POST"]
      auth:
        secret: "jwt-secret"
        skip_paths: ["/health"]
      rate_limit:
        requests_per_minute: 100
        burst: 10
        
  - id: 2
    status: enabled
    type: group
    parent_id: 1
    name: "Users API"
    path: "/users"
    middlewares:
      - logging_middleware
    config:
      logging:
        level: info
        format: json
        
  - id: 3
    status: enabled
    type: route
    parent_id: 2
    path: "/"
    method: GET
    handler: list_users
    name: "List Users"
```

## Database Schema (SQL Example)

### Main Routing Table

```sql
CREATE TABLE routes (
    id          INTEGER PRIMARY KEY,
    status      TEXT NOT NULL DEFAULT 'enabled',  -- 'enabled' | 'disabled'
    type        TEXT NOT NULL,                     -- 'middleware' | 'route' | 'group' | 'domain'
    parent_id   INTEGER,                           -- NULL for root, references routes.id
    path        TEXT,                              -- Path segment (e.g., '/api', '/users', '/:id')
    handler     TEXT,                              -- Handler name as string (e.g., 'list_users', 'cors_middleware')
    method      TEXT,                              -- HTTP method for routes ('GET', 'POST', etc.)
    name        TEXT,                              -- Human-readable name
    middlewares TEXT,                              -- JSON array of middleware names in execution order
    config_json TEXT,                              -- JSON configuration for middleware/handlers
    sequence    INTEGER DEFAULT 0,                 -- Sequence/order for processing (all types)
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (parent_id) REFERENCES routing_config(id)
);
```

## Example Data

```sql
-- Root router with global middleware chain
INSERT INTO routing_config (id, type, name, path, middlewares, config_json) VALUES 
(1, 'router', 'Main API Router', '/api', 
 '["cors_middleware", "auth_middleware", "rate_limit_middleware"]',
 '{"cors": {"origins": ["*"], "methods": ["GET", "POST"]}, "auth": {"secret": "jwt-secret", "skip_paths": ["/health"]}, "rate_limit": {"requests_per_minute": 100, "burst": 10}}');

-- API Group with additional middleware
INSERT INTO routing_config (id, type, parent_id, name, path, middlewares, config_json) VALUES 
(2, 'group', 1, 'Users API', '/users', 
 '["logging_middleware"]',
 '{"logging": {"level": "info", "format": "json"}}');

-- Routes under the group
INSERT INTO routing_config (id, type, parent_id, path, method, handler, name) VALUES 
(3, 'route', 2, '/', 'GET', 'list_users', 'List Users'),
(4, 'route', 2, '/', 'POST', 'create_user', 'Create User'),
(5, 'route', 2, '/:id', 'GET', 'get_user', 'Get User'),
(6, 'route', 2, '/:id', 'PUT', 'update_user', 'Update User'),
(7, 'route', 2, '/:id', 'DELETE', 'delete_user', 'Delete User');

-- Routes with specific middleware
INSERT INTO routing_config (id, type, parent_id, path, method, handler, name, middlewares, config_json) VALUES 
(8, 'route', 2, '/admin', 'POST', 'create_admin_user', 'Create Admin User',
 '["admin_auth_middleware", "validation_middleware"]',
 '{"admin_auth": {"required_role": "admin"}, "validation": {"schema": "admin_user_schema"}}');

-- Admin domain with its own middleware chain
INSERT INTO routing_config (id, type, name, path, middlewares, config_json) VALUES 
(9, 'domain', 'Admin Domain', '/admin',
 '["admin_cors_middleware", "admin_auth_middleware"]', 
 '{"hosts": ["admin.example.com"], "admin_cors": {"origins": ["admin.example.com"]}, "admin_auth": {"required_role": "admin"}}');

-- Admin routes
INSERT INTO routing_config (id, type, parent_id, path, method, handler, name) VALUES 
(10, 'route', 9, '/dashboard', 'GET', 'admin_dashboard', 'Admin Dashboard'),
(11, 'route', 9, '/users', 'GET', 'admin_list_users', 'Admin User List');
```

## Simple Loading Approach

### Load All Configuration at Startup

```sql
-- Simple: load everything at application startup
SELECT * FROM routing_config WHERE status = 'enabled' ORDER BY sequence;
```

**That's it!** No complex queries needed. The application:

1. **Loads all data** with one simple query
2. **Parses in memory** to build the router structure  
3. **Fast routing** - everything is in memory
4. **Small dataset** - routing config is typically small

### Enable/Disable Routes or Middleware

```sql
-- Disable a specific route
UPDATE routing_config 
SET status = 'disabled' 
WHERE id = 7;

-- Disable all authentication middleware
UPDATE routing_config 
SET status = 'disabled' 
WHERE handler = 'auth_middleware';

-- Enable rate limiting only for user creation
UPDATE routing_config 
SET status = 'enabled' 
WHERE handler = 'rate_limit_middleware' AND parent_id = 8;
```

## Go Implementation

### 1. Database Model

```go
type RoutingConfig struct {
    ID         int64     `db:"id" json:"id"`
    Status     string    `db:"status" json:"status"`
    Type       string    `db:"type" json:"type"`
    ParentID   *int64    `db:"parent_id" json:"parent_id,omitempty"`
    Prefix     *string   `db:"prefix" json:"prefix,omitempty"`
    Endpoint   *string   `db:"endpoint" json:"endpoint,omitempty"`
    Handler    *string   `db:"handler" json:"handler,omitempty"`
    Method     *string   `db:"method" json:"method,omitempty"`
    Name       *string   `db:"name" json:"name,omitempty"`
    ConfigJSON *string   `db:"config_json" json:"config_json,omitempty"`
    OrderIndex int       `db:"order_index" json:"order_index"`
    CreatedAt  time.Time `db:"created_at" json:"created_at"`
    UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}
```

### 2. Handler Registry

```go
// Handler registry maps string names to actual functions
type HandlerRegistry struct {
    routes     map[string]http.HandlerFunc
    middleware map[string]MiddlewareFactory
}

// Add route handler
func (r *HandlerRegistry) AddRoute(name string, handler http.HandlerFunc) {
    r.routes[name] = handler
}

// Add middleware factory
func (r *HandlerRegistry) AddMiddleware(name string, factory MiddlewareFactory) {
    r.middleware[name] = factory
}

// Find route handler by name
func (r *HandlerRegistry) FindRoute(name string) (http.HandlerFunc, error) {
    handler, exists := r.routes[name]
    if !exists {
        return nil, fmt.Errorf("route handler not found: %s", name)
    }
    return handler, nil
}

// Find middleware by name and create with config
func (r *HandlerRegistry) FindMiddleware(name string, config map[string]any) (StdMiddleware, error) {
    factory, exists := r.middleware[name]
    if !exists {
        return nil, fmt.Errorf("middleware not found: %s", name)
    }
    return factory(config)
}

// Remove route handler
func (r *HandlerRegistry) RemoveRoute(name string) {
    delete(r.routes, name)
}

// Remove middleware factory
func (r *HandlerRegistry) RemoveMiddleware(name string) {
    delete(r.middleware, name)
}
```

### 3. Router Builder

```go
// Build router from database configuration
func BuildRouterFromDatabase(db *sql.DB, registry *HandlerRegistry) (*Router, error) {
    // Query all routing configuration
    configs, err := queryRoutingConfigs(db)
    if err != nil {
        return nil, err
    }
    
    // Build hierarchical structure
    router := NewRouter()
    
    for _, config := range configs {
        switch config.Type {
        case "middleware":
            middleware, err := resolveMiddleware(config, registry)
            if err != nil {
                return nil, err
            }
            // Apply middleware based on parent context
            applyMiddleware(router, config, middleware)
            
        case "route":
            handler, err := registry.GetRouteHandler(*config.Handler)
            if err != nil {
                return nil, err
            }
            // Add route to appropriate group/domain
            addRoute(router, config, handler)
            
        case "group":
            // Create group structure
            createGroup(router, config)
            
        case "domain":
            // Create domain structure
            createDomain(router, config)
        }
    }
    
    return router, nil
}

func resolveMiddleware(config RoutingConfig, registry *HandlerRegistry) (StdMiddleware, error) {
    var middlewareConfig map[string]any
    if config.ConfigJSON != nil {
        err := json.Unmarshal([]byte(*config.ConfigJSON), &middlewareConfig)
        if err != nil {
            return nil, err
        }
    }
    
    return registry.GetMiddleware(*config.Handler, middlewareConfig)
}
```

### 4. Usage Example

```go
func main() {
    // Connect to database
    db, err := sql.Open("sqlite3", "routing.db")
    if err != nil {
        log.Fatal(err)
    }
    
    // Create handler registry
    registry := NewHandlerRegistry()
    
    // Add route handlers
    registry.AddRoute("list_users", func(w http.ResponseWriter, r *http.Request) {
        // List users logic
        json.NewEncoder(w).Encode([]User{})
    })
    
    registry.AddRoute("create_user", func(w http.ResponseWriter, r *http.Request) {
        // Create user logic
        w.WriteHeader(http.StatusCreated)
    })
    
    // Add middleware
    registry.AddMiddleware("cors_middleware", func(config map[string]any) (StdMiddleware, error) {
        origins := config["origins"].([]string)
        return func(next http.Handler) http.Handler {
            return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.Header().Set("Access-Control-Allow-Origin", strings.Join(origins, ","))
                next.ServeHTTP(w, r)
            })
        }, nil
    })
    
    registry.AddMiddleware("auth_middleware", func(config map[string]any) (StdMiddleware, error) {
        secret := config["secret"].(string)
        return func(next http.Handler) http.Handler {
            return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                // JWT validation using secret
                next.ServeHTTP(w, r)
            })
        }, nil
    })
    
    // Build router from database
    router, err := BuildRouterFromDatabase(db, registry)
    if err != nil {
        log.Fatal(err)
    }
    
    // Start server
    log.Println("Server starting on :8080")
    http.ListenAndServe(":8080", router)
}
```

## Benefits

1. **Truly Database-Driven**: All configuration stored as pure data
2. **Dynamic Configuration**: Routes and middleware can be modified without code changes
3. **Hierarchical Structure**: Parent-child relationships for groups, domains, middleware inheritance
4. **Execution Control**: Enable/disable routes and middleware via database updates
5. **Audit Trail**: Track configuration changes with timestamps
6. **Scalable**: Can handle complex routing scenarios with proper indexing
7. **Admin Interface**: Easy to build web UI for route management

## Key Insight

The **handler names are just strings** in the database. The actual function resolution happens at runtime through the registry pattern. This makes the entire routing configuration **truly declarative and database-friendly**.
