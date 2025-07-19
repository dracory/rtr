package rtr_test

import (
	"testing"

	"github.com/dracory/rtr"
)

func TestDeclarativeRouting(t *testing.T) {
	// Create a simple route - very readable and clean
	route := rtr.Route{
		Status:      rtr.StatusEnabled,
		Path:        "/users",
		Method:      rtr.MethodGET,
		Handler:     "list_users",
		Name:        "List Users",
		Middlewares: []string{"cors", "auth"},
	}

	// Verify route properties
	if route.Path != "/users" {
		t.Errorf("Expected path /users, got %s", route.Path)
	}

	if route.Method != rtr.MethodGET {
		t.Errorf("Expected method %s, got %s", rtr.MethodGET, route.Method)
	}

	if route.Handler != "list_users" {
		t.Errorf("Expected handler list_users, got %s", route.Handler)
	}

	if route.Status != rtr.StatusEnabled {
		t.Errorf("Expected status %s, got %s", rtr.StatusEnabled, route.Status)
	}

	if route.Name != "List Users" {
		t.Errorf("Expected name 'List Users', got %s", route.Name)
	}

	if len(route.Middlewares) != 2 {
		t.Errorf("Expected 2 middlewares, got %d", len(route.Middlewares))
	}
}

func TestHandlerRegistry(t *testing.T) {
	// Test basic registry creation
	registry := rtr.NewHandlerRegistry()
	if registry == nil {
		t.Error("Expected registry to be created")
	}

	// Test finding non-existent route
	route := registry.FindRoute("non_existent")
	if route != nil {
		t.Error("Expected nil for non-existent route")
	}

	// Test finding non-existent middleware
	middleware := registry.FindMiddleware("non_existent")
	if middleware != nil {
		t.Error("Expected nil for non-existent middleware")
	}
}

func TestDomain(t *testing.T) {
	// Create a domain configuration - clean and readable
	domain := rtr.Domain{
		Status: rtr.StatusEnabled,
		Name:   "API Domain",
		Hosts:  []string{"api.example.com", "api.test.com"},
		Items: []rtr.ItemInterface{
			rtr.Route{
				Status:  rtr.StatusEnabled,
				Path:    "/health",
				Method:  rtr.MethodGET,
				Handler: "health_check",
				Name:    "Health Check",
			},
			rtr.Group{
				Status: rtr.StatusEnabled,
				Prefix: "/api",
				Name:   "API Group",
				Routes: []rtr.Route{
					{
						Status:  rtr.StatusEnabled,
						Path:    "/users",
						Method:  rtr.MethodGET,
						Handler: "list_users",
						Name:    "List Users",
					},
				},
			},
		},
		Middlewares: []string{"cors", "rate_limit"},
	}

	// Verify domain properties
	if len(domain.Hosts) != 2 {
		t.Errorf("Expected 2 hosts, got %d", len(domain.Hosts))
	}

	if domain.Hosts[0] != "api.example.com" {
		t.Errorf("Expected first host api.example.com, got %s", domain.Hosts[0])
	}

	if len(domain.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(domain.Items))
	}

	if len(domain.Middlewares) != 2 {
		t.Errorf("Expected 2 middlewares, got %d", len(domain.Middlewares))
	}

	// Test ItemInterface methods
	if domain.GetName() != "API Domain" {
		t.Errorf("Expected name 'API Domain', got %s", domain.GetName())
	}

	if domain.GetStatus() != rtr.StatusEnabled {
		t.Errorf("Expected status '%s', got %s", rtr.StatusEnabled, domain.GetStatus())
	}

	if len(domain.GetMiddlewares()) != 2 {
		t.Errorf("Expected 2 middlewares from GetMiddlewares(), got %d", len(domain.GetMiddlewares()))
	}
}
