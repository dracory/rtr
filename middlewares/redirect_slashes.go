package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dracory/rtr"
)

// RedirectSlashesMiddleware returns a middleware that redirects requests with
// a trailing slash to the same path without the trailing slash using a 301
// (Moved Permanently) redirect.
//
// Behaviour mirrors chi's middleware.RedirectSlashes, with one difference:
//   - The root path "/" is never redirected.
//   - Backslashes are normalized to forward slashes to prevent protocol-relative
//     redirects such as "/\evil.com".
//   - Leading and trailing slashes are trimmed. Internal double slashes in the
//     computed path are also collapsed by http.Redirect when it sets the
//     Location header (e.g. "/api//v1/" → "/api/v1").
//   - All trailing slashes are collapsed, not just one. chi only stripped a
//     single trailing slash ("/api/v1///" → "/api/v1//"); this implementation
//     trims them all ("/api/v1///" → "/api/v1"). This is an intentional
//     improvement over chi's behaviour.
//   - The raw query string is preserved on redirect.
func RedirectSlashesMiddleware() rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Redirect Slashes Middleware").
		SetHandler(redirectSlashesHandler)
}

// redirectSlashesHandler is the stdlib-only implementation of the trailing-slash
// redirect behaviour previously provided by chi's middleware.RedirectSlashes.
func redirectSlashesHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Only redirect when the path is longer than "/" and ends with a slash.
		if len(path) <= 1 || path[len(path)-1] != '/' {
			next.ServeHTTP(w, r)
			return
		}

		// Normalize backslashes to forward slashes to prevent "/\evil.com" style
		// redirects that some clients may interpret as protocol-relative.
		path = strings.ReplaceAll(path, `\`, `/`)

		// Trim leading and trailing slashes, then force a single leading slash.
		// http.Redirect will further clean internal double slashes in the
		// Location header (e.g. "/a//b///" → "/a//b" here, then "/a/b" in Location).
		path = "/" + strings.Trim(path, "/")

		// Preserve the raw query string, if any.
		if r.URL.RawQuery != "" {
			path = fmt.Sprintf("%s?%s", path, r.URL.RawQuery)
		}

		http.Redirect(w, r, path, http.StatusMovedPermanently)
	})
}
