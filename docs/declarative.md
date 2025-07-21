# Declarative Router Configuration

The `rtr` router supports a powerful declarative approach that allows you to define your entire routing structure using Go structs. This separates the routing *definition* from the handler *implementation*, leading to cleaner, more maintainable code.

Furthermore, because this system uses string names to reference handlers, it opens the door to loading your entire routing configuration from external data sources like JSON, YAML, or even a database.

## 1. Declarative Configuration in Go (Primary Method)

The most common way to use the declarative system is by defining your routing hierarchy directly in your Go code.

### Example: Defining a Domain in Go

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/dracory/rtr"
)

// Define your handlers first
func homeHandler(w http.ResponseWriter, r *http.Request) { w.Write([]byte("Welcome Home")) }
func statusHandler(w http.ResponseWriter, r *http.Request) string { return `{"status":"ok"}` }

func main() {
    // 1. Register handlers and middleware with string names
    registry := rtr.NewHandlerRegistry()
    registry.AddRoute(rtr.NewRoute().SetName("home").SetHandler(homeHandler))
    registry.AddRoute(rtr.NewRoute().SetName("api_status").SetJSONHandler(statusHandler))
    // ... register middleware if any

    // 2. Define the entire routing structure as a Go struct
    domainConfig := &rtr.Domain{
        Name:   "Example Domain",
        Status: rtr.StatusEnabled,
        Hosts:  []string{"localhost:8080"},
        Items: []rtr.ItemInterface{
            // A route at the root of the domain
            &rtr.Route{
                Name:    "Home",
                Method:  rtr.MethodGET,
                Path:    "/",
                Handler: "home", // Reference handler by its registered name
                Status:  rtr.StatusEnabled,
            },
            // A group of related routes
            &rtr.Group{
                Name:   "API",
                Prefix: "/api",
                Status: rtr.StatusEnabled,
                Routes: []rtr.Route{
                    {
                        Name:       "API Status",
                        Method:     rtr.MethodGET,
                        Path:       "/status",
                        JSONHandler: "api_status", // Reference handler by name
                        Status:     rtr.StatusEnabled,
                    },
                },
            },
        },
    }

    // 3. Build the router from the configuration (conceptual)
    // You would write a builder function to translate the struct into a live router.
    // router := buildRouterFromConfig(domainConfig, registry)

    fmt.Printf("Domain configuration created: %s\n", domainConfig.Name)
}
```

## 2. Loading Configuration from External Sources

The true power of this system is that the `rtr.Domain` struct and its children can be serialized. This means you can define your routes in a separate file (JSON, YAML, etc.) or a database and load them at runtime.

This is made possible by the `HandlerRegistry`, which decouples the configuration from the implementation.

### JSON Example

You can represent the exact same configuration in a `routes.json` file.

```json
{
  "name": "Example Domain",
  "status": "enabled",
  "hosts": ["localhost:8080"],
  "items": [
    {
      "type": "route",
      "name": "Home",
      "method": "GET",
      "path": "/",
      "handler": "home",
      "status": "enabled"
    },
    {
      "type": "group",
      "name": "API",
      "prefix": "/api",
      "status": "enabled",
      "routes": [
        {
          "name": "API Status",
          "method": "GET",
          "path": "/status",
          "jsonHandler": "api_status",
          "status": "enabled"
        }
      ]
    }
  ]
}
```

Your application would then read this file, unmarshal it into the `rtr.Domain` struct, and build the router.

### YAML Example

Similarly, you could use YAML.

```yaml
name: Example Domain
status: enabled
hosts:
  - localhost:8080
items:
  - type: route
    name: Home
    method: GET
    path: "/"
    handler: home
    status: enabled
  - type: group
    name: API
    prefix: /api
    status: enabled
    routes:
      - name: API Status
        method: GET
        path: /status
        jsonHandler: api_status
        status: enabled
```

### Database (SQL) Example

For ultimate flexibility, you can store your routing rules in a database. This allows you to change routing, enable/disable endpoints, and add middleware without redeploying your application.

For a detailed guide on how to structure your database tables and load the configuration, please see the in-depth guide:

**[Declarative Routing from a Database](./declarative-routing-system.md)**

```