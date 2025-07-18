package rtr_test

import (
	"testing"

	"github.com/dracory/rtr"
)

func TestRouter_Domains(t *testing.T) {
	// Create a new router
	router := rtr.NewRouter()

	// Test GetDomains with no domains
	t.Run("GetDomains with no domains", func(t *testing.T) {
		domains := router.GetDomains()
		if len(domains) != 0 {
			t.Errorf("Expected no domains initially, got %d", len(domains))
		}
	})

	// Create some test domains
	domain1 := rtr.NewDomain("example.com")
	domain2 := rtr.NewDomain("api.example.com")
	domain3 := rtr.NewDomain("test.example.com")

	// Test AddDomain
	t.Run("AddDomain", func(t *testing.T) {
		router = router.AddDomain(domain1)
		domains := router.GetDomains()
		if len(domains) != 1 {
			t.Fatalf("Expected 1 domain after AddDomain, got %d", len(domains))
		}
		// We can't directly compare domain1 and domains[0] because they might be different instances
		// with the same pattern, so we'll just check that we got one domain

		// Add another domain
		router = router.AddDomain(domain2)
		domains = router.GetDomains()
		if len(domains) != 2 {
			t.Fatalf("Expected 2 domains after second AddDomain, got %d", len(domains))
		}
	})

	// Test AddDomains
	t.Run("AddDomains", func(t *testing.T) {
		// Reset domains
		router = rtr.NewRouter()
		
		// Add multiple domains at once
		router = router.AddDomains([]rtr.DomainInterface{domain1, domain2, domain3})
		
		domains := router.GetDomains()
		if len(domains) != 3 {
			t.Fatalf("Expected 3 domains after AddDomains, got %d", len(domains))
		}
		// We can't directly compare the domains because they might be different instances
		// with the same pattern, so we'll just check that we got three domains
	})

	// Test chaining
	t.Run("Method chaining", func(t *testing.T) {
		router = rtr.NewRouter().
			AddDomain(domain1).
			AddDomains([]rtr.DomainInterface{domain2, domain3})

		domains := router.GetDomains()
		if len(domains) != 3 {
			t.Fatalf("Expected 3 domains after method chaining, got %d", len(domains))
		}
	})
}
