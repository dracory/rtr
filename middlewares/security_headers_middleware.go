package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dracory/rtr"
)

// SecurityHeadersConfig provides configuration for security headers middleware
type SecurityHeadersConfig struct {
	// Content Security Policy configuration
	CSP *CSPConfig
	// HSTS configuration
	HSTS *HSTSConfig
	// Frame options configuration
	FrameOptions *FrameOptionsConfig
	// Content type options
	ContentTypeNosniff bool
	// XSS protection
	XSSProtection *XSSProtectionConfig
	// Referrer policy
	ReferrerPolicy string
	// Permissions policy
	PermissionsPolicy map[string][]string
	// Custom headers allows adding custom security headers
	CustomHeaders map[string]string
}

// CSPConfig configures Content Security Policy
type CSPConfig struct {
	Enabled                 bool
	DefaultSrc              []string
	ScriptSrc               []string
	StyleSrc                []string
	FontSrc                 []string
	ImgSrc                  []string
	ConnectSrc              []string
	MediaSrc                []string
	ObjectSrc               []string
	ChildSrc                []string
	WorkerSrc               []string
	ManifestSrc             []string
	UpgradeInsecureRequests bool
}

// HSTSConfig configures HTTP Strict Transport Security
type HSTSConfig struct {
	Enabled           bool
	MaxAge            int
	IncludeSubDomains bool
	Preload           bool
}

// FrameOptionsConfig configures X-Frame-Options
type FrameOptionsConfig struct {
	Enabled bool
	Option  string // "DENY", "SAMEORIGIN", or "ALLOW-FROM uri"
}

// XSSProtectionConfig configures X-XSS-Protection
type XSSProtectionConfig struct {
	Enabled bool
	Mode    string // "block" or empty
}

// DefaultSecurityHeadersConfig returns a secure default configuration
func DefaultSecurityHeadersConfig() *SecurityHeadersConfig {
	return &SecurityHeadersConfig{
		CSP: &CSPConfig{
			Enabled:    true,
			DefaultSrc: []string{"'self'"},
			ScriptSrc: []string{
				"'self'",
				"'unsafe-inline'",
				"https://cdn.jsdelivr.net",
				"https://unpkg.com",
				"https://code.jquery.com",
				"https://cdnjs.cloudflare.com",
			},
			StyleSrc: []string{
				"'self'",
				"'unsafe-inline'",
				"https://cdn.jsdelivr.net",
				"https://maxcdn.bootstrapcdn.com",
				"https://cdnjs.cloudflare.com",
				"https://fonts.googleapis.com",
				"https://unpkg.com",
			},
			FontSrc: []string{
				"'self'",
				"https://cdn.jsdelivr.net",
				"https://fonts.googleapis.com",
				"https://fonts.gstatic.com",
				"https://cdnjs.cloudflare.com",
				"https://maxcdn.bootstrapcdn.com",
			},
			ImgSrc: []string{
				"'self'",
				"data:",
			},
			UpgradeInsecureRequests: true,
		},
		HSTS: &HSTSConfig{
			Enabled:           true,
			MaxAge:            31536000,
			IncludeSubDomains: true,
		},
		FrameOptions: &FrameOptionsConfig{
			Enabled: true,
			Option:  "DENY",
		},
		ContentTypeNosniff: true,
		XSSProtection: &XSSProtectionConfig{
			Enabled: true,
			Mode:    "block",
		},
		ReferrerPolicy: "strict-origin-when-cross-origin",
		PermissionsPolicy: map[string][]string{
			"geolocation": {},
			"microphone":  {},
			"camera":      {},
		},
		CustomHeaders: make(map[string]string),
	}
}

// NewSecurityHeadersMiddleware creates middleware that sets security headers
func NewSecurityHeadersMiddleware(config *SecurityHeadersConfig) rtr.MiddlewareInterface {
	if config == nil {
		config = DefaultSecurityHeadersConfig()
	}

	return rtr.NewMiddleware().
		SetName("Security Headers Middleware").
		SetHandler(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Set HSTS header
				if config.HSTS != nil && config.HSTS.Enabled {
					hstsValue := fmt.Sprintf("max-age=%d", config.HSTS.MaxAge)
					if config.HSTS.IncludeSubDomains {
						hstsValue += "; includeSubDomains"
					}
					if config.HSTS.Preload {
						hstsValue += "; preload"
					}
					w.Header().Set("Strict-Transport-Security", hstsValue)
				}

				// Set Frame Options header
				if config.FrameOptions != nil && config.FrameOptions.Enabled {
					w.Header().Set("X-Frame-Options", config.FrameOptions.Option)
				}

				// Set Content Type Options header
				if config.ContentTypeNosniff {
					w.Header().Set("X-Content-Type-Options", "nosniff")
				}

				// Set XSS Protection header
				if config.XSSProtection != nil && config.XSSProtection.Enabled {
					value := "1"
					if config.XSSProtection.Mode == "block" {
						value += "; mode=block"
					}
					w.Header().Set("X-XSS-Protection", value)
				}

				// Set Referrer Policy header
				if config.ReferrerPolicy != "" {
					w.Header().Set("Referrer-Policy", config.ReferrerPolicy)
				}

				// Set Permissions Policy header
				if len(config.PermissionsPolicy) > 0 {
					var policies []string
					for feature, origins := range config.PermissionsPolicy {
						if len(origins) == 0 {
							policies = append(policies, feature+"=()")
						} else {
							policies = append(policies, feature+"="+strings.Join(origins, " "))
						}
					}
					w.Header().Set("Permissions-Policy", strings.Join(policies, ", "))
				}

				// Set CSP header
				if config.CSP != nil && config.CSP.Enabled {
					cspValue := buildCSPValue(config.CSP)
					if cspValue != "" {
						w.Header().Set("Content-Security-Policy", cspValue)
					}
				}

				// Set custom headers
				for key, value := range config.CustomHeaders {
					w.Header().Set(key, value)
				}

				next.ServeHTTP(w, r)
			})
		})
}

// buildCSPValue constructs the CSP header value from configuration
func buildCSPValue(config *CSPConfig) string {
	var directives []string

	if len(config.DefaultSrc) > 0 {
		directives = append(directives, "default-src "+strings.Join(config.DefaultSrc, " "))
	}
	if len(config.ScriptSrc) > 0 {
		directives = append(directives, "script-src "+strings.Join(config.ScriptSrc, " "))
	}
	if len(config.StyleSrc) > 0 {
		directives = append(directives, "style-src "+strings.Join(config.StyleSrc, " "))
	}
	if len(config.FontSrc) > 0 {
		directives = append(directives, "font-src "+strings.Join(config.FontSrc, " "))
	}
	if len(config.ImgSrc) > 0 {
		directives = append(directives, "img-src "+strings.Join(config.ImgSrc, " "))
	}
	if len(config.ConnectSrc) > 0 {
		directives = append(directives, "connect-src "+strings.Join(config.ConnectSrc, " "))
	}
	if len(config.MediaSrc) > 0 {
		directives = append(directives, "media-src "+strings.Join(config.MediaSrc, " "))
	}
	if len(config.ObjectSrc) > 0 {
		directives = append(directives, "object-src "+strings.Join(config.ObjectSrc, " "))
	}
	if len(config.ChildSrc) > 0 {
		directives = append(directives, "child-src "+strings.Join(config.ChildSrc, " "))
	}
	if len(config.WorkerSrc) > 0 {
		directives = append(directives, "worker-src "+strings.Join(config.WorkerSrc, " "))
	}
	if len(config.ManifestSrc) > 0 {
		directives = append(directives, "manifest-src "+strings.Join(config.ManifestSrc, " "))
	}
	if config.UpgradeInsecureRequests {
		directives = append(directives, "upgrade-insecure-requests")
	}

	return strings.Join(directives, "; ")
}
