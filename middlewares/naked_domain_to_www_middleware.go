package middlewares

import (
	"net/http"
	"strings"

	"github.com/dracory/rtr"
)

// NakedDomainToWwwMiddleware redirects naked domains to the www subdomain.
// hostExcludes allows bypassing the redirect for specific hosts (e.g., localhost).
func NakedDomainToWwwMiddleware(hostExcludes []string) rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Naked Domain to WWW Middleware").
		SetHandler(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				host := strings.ToLower(r.Host)

				if host == "" {
					next.ServeHTTP(w, r)
					return
				}

				if strings.HasPrefix(host, "www.") {
					next.ServeHTTP(w, r)
					return
				}

				for _, v := range hostExcludes {
					if strings.HasPrefix(host, strings.ToLower(v)) {
						next.ServeHTTP(w, r)
						return
					}
				}

				scheme := strings.ToLower(r.URL.Scheme)
				if scheme == "" || scheme == "/" {
					scheme = "https"
				}

				targetHost := "www." + strings.TrimPrefix(r.Host, "www.")
				redirectURL := scheme + "://" + targetHost + r.URL.RequestURI()

				http.Redirect(w, r, redirectURL, http.StatusPermanentRedirect)
			})
		})
}
