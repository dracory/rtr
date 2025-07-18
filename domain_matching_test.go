package rtr_test

import (
	"fmt"
	"testing"

	"github.com/dracory/rtr"
)

func TestDomain_SetPatterns(t *testing.T) {
	tests := []struct {
		name     string
		initial  []string
		setTo    []string
		expected []string
	}{
		{
			name:     "set single pattern",
			initial:  []string{"example.com"},
			setTo:    []string{"api.example.com"},
			expected: []string{"api.example.com"},
		},
		{
			name:     "set multiple patterns",
			initial:  []string{"example.com"},
			setTo:    []string{"api.example.com", "staging.example.com"},
			expected: []string{"api.example.com", "staging.example.com"},
		},
		{
			name:     "set empty patterns",
			initial:  []string{"example.com"},
			setTo:    []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create domain with initial patterns
			d := rtr.NewDomain(tt.initial...)
			
			// Set new patterns
			result := d.SetPatterns(tt.setTo...)
			
			// Verify the returned value is the same domain (for method chaining)
			if result != d {
				t.Error("Expected SetPatterns to return the same domain instance")
			}
			
			// Verify the patterns were set correctly
			got := d.GetPatterns()
			if len(got) != len(tt.expected) {
				t.Fatalf("Expected %d patterns, got %d", len(tt.expected), len(got))
			}
			
			for i, pattern := range tt.expected {
				if got[i] != pattern {
					t.Errorf("Pattern[%d] = %q, want %q", i, got[i], pattern)
				}
			}
		})
	}
}

func TestDomain_String(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		expected string
	}{
		{
			name:     "single pattern",
			patterns: []string{"example.com"},
			expected: "Domain(patterns=[example.com])",
		},
		{
			name:     "multiple patterns",
			patterns: []string{"example.com", "*.example.com"},
			expected: "Domain(patterns=[example.com *.example.com])",
		},
		{
			name:     "no patterns",
			patterns: []string{},
			expected: "Domain(patterns=[])",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := rtr.NewDomain(tt.patterns...)
			// Convert to string using fmt.Sprintf which will call the String() method
			got := fmt.Sprintf("%v", d)
			if got != tt.expected {
				t.Errorf("String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestDomain_MatchesPattern(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		pattern string
		want    bool
	}{
		// IPv4 with port tests
		{
			name:    "exact IPv4 with port match",
			host:    "example.com:8080",
			pattern: "example.com:8080",
			want:    true,
		},
		{
			name:    "IPv4 with port wildcard",
			host:    "example.com:8080",
			pattern: "example.com:*",
			want:    true,
		},
		{
			name:    "IPv4 with different ports",
			host:    "example.com:8080",
			pattern: "example.com:3000",
			want:    false,
		},
		{
			name:    "IPv4 pattern with port, host without port",
			host:    "example.com",
			pattern: "example.com:8080",
			want:    false,
		},

		// IPv6 tests
		{
			name:    "exact IPv6 with port match",
			host:    "[::1]:8080",
			pattern: "[::1]:8080",
			want:    true,
		},
		{
			name:    "IPv6 with port wildcard",
			host:    "[::1]:8080",
			pattern: "[::1]:*",
			want:    true,
		},
		{
			name:    "IPv6 with different ports",
			host:    "[::1]:8080",
			pattern: "[::1]:3000",
			want:    false,
		},
		{
			name:    "IPv6 pattern with port, host without port",
			host:    "[::1]",
			pattern: "[::1]:8080",
			want:    false,
		},
		{
			name:    "IPv6 without port, pattern without port",
			host:    "[2001:db8::1]",
			pattern: "[2001:db8::1]",
			want:    true,
		},

		// Wildcard subdomain tests
		{
			name:    "wildcard subdomain match",
			host:    "api.example.com",
			pattern: "*.example.com",
			want:    true,
		},
		{
			name:    "wildcard subdomain with multiple levels",
			host:    "v1.api.example.com",
			pattern: "*.example.com",
			want:    true,
		},
		{
			name:    "wildcard subdomain no match",
			host:    "example.com",
			pattern: "*.example.com",
			want:    false,
		},
		{
			name:    "wildcard subdomain with port",
			host:    "api.example.com:8080",
			pattern: "*.example.com:8080",
			want:    true,
		},

		// Edge cases
		{
			name:    "empty host and pattern",
			host:    "",
			pattern: "",
			want:    false, // Empty pattern doesn't match empty host
		},
		{
			name:    "empty host",
			host:    "",
			pattern: "example.com",
			want:    false,
		},
		{
			name:    "empty pattern",
			host:    "example.com",
			pattern: "",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := rtr.NewDomain(tt.pattern)
			got := d.Match(tt.host)
			if got != tt.want {
				t.Errorf("Match(%q) with pattern %q = %v, want %v", tt.host, tt.pattern, got, tt.want)
			}
		})
	}
}
