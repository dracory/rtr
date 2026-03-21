# Security Middleware

The security package provides configurable HTTP security middlewares for the RTR router framework.

## Features

- **HTTPS Redirect**: Automatically redirects HTTP requests to HTTPS
- **Security Headers**: Sets comprehensive security headers including CSP, HSTS, and more
- **Highly Configurable**: All aspects can be customized per project needs
- **Production Defaults**: Secure defaults out of the box

## Installation

```go
import "github.com/dracory/rtr/middleware/security"
```

## Usage

### HTTPS Redirect Middleware

```go
// Use default configuration (recommended for most cases)
router.Use(security.NewHTTPSRedirectMiddleware(nil))

// Or with custom configuration
config := &security.HTTPSRedirectConfig{
    SkipLocalhost: false,
    CustomSkipFunc: func(r *http.Request) bool {
        return r.Host == "internal.example.com"
    },
}
router.Use(security.NewHTTPSRedirectMiddleware(config))
```

### Security Headers Middleware

```go
// Use default secure configuration
router.Use(security.NewSecurityHeadersMiddleware(nil))

// Or with custom configuration
config := &security.SecurityHeadersConfig{
    CSP: &security.CSPConfig{
        Enabled:    true,
        DefaultSrc: []string{"'self'", "https://trusted.cdn.com"},
        ScriptSrc:  []string{"'self'", "https://trusted.cdn.com"},
        StyleSrc:   []string{"'self'", "https://trusted.cdn.com"},
    },
    HSTS: &security.HSTSConfig{
        Enabled:           true,
        MaxAge:            31536000,
        IncludeSubDomains: true,
        Preload:           true,
    },
    CustomHeaders: map[string]string{
        "X-Custom-Security": "enabled",
    },
}
router.Use(security.NewSecurityHeadersMiddleware(config))
```

## Configuration

### HTTPS Redirect Configuration

```go
type HTTPSRedirectConfig struct {
    SkipLocalhost   bool                    // Skip HTTPS redirect for localhost (default: true)
    TrustedProxies []string                // List of trusted proxy IPs
    CustomSkipFunc func(r *http.Request) bool // Custom logic to skip redirect
}
```

### Security Headers Configuration

```go
type SecurityHeadersConfig struct {
    CSP               *CSPConfig               // Content Security Policy
    HSTS              *HSTSConfig              // HTTP Strict Transport Security
    FrameOptions      *FrameOptionsConfig      // X-Frame-Options
    ContentTypeNosniff bool                    // X-Content-Type-Options
    XSSProtection     *XSSProtectionConfig     // X-XSS-Protection
    ReferrerPolicy    string                  // Referrer-Policy
    PermissionsPolicy map[string][]string     // Permissions-Policy
    CustomHeaders     map[string]string       // Custom security headers
}
```

### Content Security Policy Configuration

```go
type CSPConfig struct {
    Enabled                  bool
    DefaultSrc               []string // Default source for all directives
    ScriptSrc                []string // Scripts
    StyleSrc                 []string // Stylesheets
    FontSrc                  []string // Fonts
    ImgSrc                   []string // Images
    ConnectSrc               []string // AJAX, WebSockets
    MediaSrc                 []string // Video, audio
    ObjectSrc                []string // Plugins
    ChildSrc                 []string // Frames
    WorkerSrc                []string // Web Workers
    ManifestSrc              []string // Web App Manifest
    UpgradeInsecureRequests  bool     // Force HTTPS
}
```

## Examples

### E-commerce Site Configuration

```go
config := &security.SecurityHeadersConfig{
    CSP: &security.CSPConfig{
        Enabled: true,
        DefaultSrc: []string{"'self'"},
        ScriptSrc: []string{
            "'self'",
            "'unsafe-inline'", // For template literals
            "https://payment.gateway.com",
            "https://cdn.tracking.com",
        },
        StyleSrc: []string{
            "'self'",
            "'unsafe-inline'", // For CSS-in-JS
            "https://fonts.googleapis.com",
        },
        FontSrc: []string{
            "'self'",
            "https://fonts.gstatic.com",
        },
        ImgSrc: []string{
            "'self'",
            "data:",
            "https://cdn.images.com",
            "https://product.photos.com",
        },
        ConnectSrc: []string{
            "'self'",
            "https://api.ecommerce.com",
            "https://analytics.com",
        },
    },
    CustomHeaders: map[string]string{
        "X-Content-Type-Options": "nosniff",
        "X-Frame-Options": " "DENY",
    },
}

router.Use(security.NewSecurityHeadersMiddleware(config))
```

### Development Environment

```go
// More permissive settings for development
devConfig := &security.SecurityHeadersConfig{
    CSP: &security.CSPConfig{
        Enabled: true,
        DefaultSrc: []string{"'self'", "'unsafe-inline'", "'unsafe-eval'"},
        ScriptSrc: []string{"'self'", "'unsafe-inline'", "'unsafe-eval'"},
        StyleSrc:  []string{"'self'", "'unsafe-inline'"},
    },
    HSTS: &security.HSTSConfig{
        Enabled: false, // Disable HSTS in development
    },
}

router.Use(security.NewSecurityHeadersMiddleware(devConfig))
```

### API-only Configuration

```go
// For APIs that don't serve HTML
apiConfig := &security.SecurityHeadersConfig{
    CSP: &security.CSPConfig{
        Enabled: false, // No CSP needed for APIs
    },
    HSTS: &security.HSTSConfig{
        Enabled:          true,
        MaxAge:           31536000,
        IncludeSubDomains: true,
    },
    CustomHeaders: map[string]string{
        "X-Content-Type-Options": "nosniff",
        "X-Frame-Options":        "DENY",
    },
}

router.Use(security.NewSecurityHeadersMiddleware(apiConfig))
```

## Default Security Headers

When using the default configuration, the following headers are set:

- **Strict-Transport-Security**: `max-age=31536000; includeSubDomains`
- **X-Frame-Options**: `DENY`
- **X-Content-Type-Options**: `nosniff`
- **X-XSS-Protection**: `1; mode=block`
- **Referrer-Policy**: `strict-origin-when-cross-origin`
- **Permissions-Policy**: `geolocation=(), microphone=(), camera=()`
- **Content-Security-Policy**: Comprehensive CSP with common CDNs

## Testing

The package includes comprehensive tests:

```bash
go test ./middleware/security
```

## Security Considerations

1. **CSP**: Start with strict policies and relax as needed
2. **HSTS**: Use includeSubDomains only if all subdomains support HTTPS
3. **Development**: Consider disabling some headers in development environments
4. **Testing**: Always test security headers don't break functionality

## Dependencies

- Go 1.16+
- github.com/dracory/rtr
