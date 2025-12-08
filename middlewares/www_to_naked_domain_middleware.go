package middlewares

import (
	"net/http"
	"strings"

	"github.com/dracory/rtr"
)

// WwwToNakedDomainMiddleware redirects requests from the www subdomain to the naked domain.
func WwwToNakedDomainMiddleware() rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("WWW to Naked Domain Middleware").
		SetHandler(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				host := strings.ToLower(r.Host)

				if strings.HasPrefix(host, "www.") {
					scheme := strings.ToLower(r.URL.Scheme)
					if scheme == "" || scheme == "/" {
						scheme = "https"
					}

					targetHost := strings.TrimPrefix(r.Host, "www.")
					redirectURL := scheme + "://" + targetHost + r.URL.RequestURI()

					http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
					return
				}

				next.ServeHTTP(w, r)
			})
		})
}
