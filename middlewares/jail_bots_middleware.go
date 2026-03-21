package middlewares

import (
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/dracory/rtr"
	"github.com/jellydator/ttlcache/v3"
	"github.com/samber/lo"
)

// getIP extracts the real IP address from the request
func getIP(r *http.Request) string {
	// Check X-Forwarded-For header first
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		xff = strings.TrimSpace(xff)
		if xff != "" {
			// X-Forwarded-For can contain multiple IPs, take the first one
			if idx := strings.Index(xff, ","); idx != -1 {
				return strings.TrimSpace(xff[:idx])
			}
			return xff
		}
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		xri = strings.TrimSpace(xri)
		if xri != "" {
			return xri
		}
	}

	// Fall back to RemoteAddr
	if idx := strings.LastIndex(r.RemoteAddr, ":"); idx != -1 {
		return r.RemoteAddr[:idx]
	}
	return r.RemoteAddr
}

func JailBotsMiddleware(config JailBotsConfig) rtr.MiddlewareInterface {
	jb := new(jailBotsMiddleware)
	jb.exclude = config.Exclude
	jb.cache = ttlcache.New[string, struct{}]()
	jb.excludePaths = append([]string{}, config.ExcludePaths...)
	m := rtr.NewMiddleware().
		SetName(jb.Name()).
		SetHandler(jb.Handler)

	return m
}

// JailBotsConfig defines configuration for jail bots middleware
type JailBotsConfig struct {
	// Exclude filters items out of the internal URI blacklist lists used by
	// isJailable (e.g., if "wp" is in the blacklist but you want to allow it,
	// add "wp" here). Matches are compared literally against the blacklist
	// entries, not against request paths.
	Exclude []string

	// ExcludePaths defines request path patterns that must bypass the jail logic.
	// Supported patterns:
	//  - With a trailing '*': treated as a simple prefix match, e.g. "/blog*" matches
	//    "/blog", "/blog/", and any subpaths like "/blog/post".
	//  - Without '*': segment-aware; matches exactly the path (e.g. "/blog") or any
	//    subpath starting with that segment (e.g. "/blog/..."), but NOT lookalikes
	//    like "/blogger".
	ExcludePaths []string
}

type jailBotsMiddleware struct {
	exclude      []string
	cache        *ttlcache.Cache[string, struct{}]
	excludePaths []string
}

func (m *jailBotsMiddleware) Name() string {
	return "Jail Bots Middleware"
}

func (m *jailBotsMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		ip := getIP(r)

		// Exclude specific routes from jail logic
		if m.isExcludedPath(path) {
			next.ServeHTTP(w, r)
			return
		}

		if m.isJailed(ip) {
			w.WriteHeader(http.StatusForbidden)
			writeJailBotsHTML(w, "malicious access not allowed (jb)")
			return
		}

		jailable, reason := m.isJailable(path)

		if jailable {
			m.jail(ip)

			slog.Default().Info("Jailed bot from "+ip+" for 5 minutes",
				slog.String("reason", reason),
				slog.String("path", path),
				slog.String("ip", ip),
				slog.String("useragent", r.UserAgent()),
			)

			w.WriteHeader(http.StatusForbidden)
			writeJailBotsHTML(w, "malicious access not allowed (jb)")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func writeJailBotsHTML(w http.ResponseWriter, body string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := w.Write([]byte(body)); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (j *jailBotsMiddleware) isJailed(ip string) bool {
	if j.cache == nil {
		return false
	}
	return j.cache.Has("jail:" + ip)
}

func (j *jailBotsMiddleware) jail(ip string) {
	if j.cache == nil {
		return
	}
	j.cache.Set("jail:"+ip, struct{}{}, 5*time.Minute)
}

func (m *jailBotsMiddleware) isJailable(path string) (jailable bool, reason string) {
	startsWithList := m.startsWithBlacklistedUriList()

	for i := 0; i < len(startsWithList); i++ {
		if strings.HasPrefix(path, startsWithList[i]) {
			return true, "starts with " + startsWithList[i]
		}
	}

	containsList := m.containsBlacklistedUriList()

	for i := 0; i < len(containsList); i++ {
		if strings.Contains(path, containsList[i]) {
			return true, "contains " + containsList[i]
		}
	}

	return false, ""
}

// isExcludedPath returns true for routes that should bypass jail logic entirely.
// Supports simple wildcard '*' suffix in patterns (prefix match), e.g., '/blog*'.
// Without '*', it matches exact segment (exact path or path starting with pattern + '/').
func (m *jailBotsMiddleware) isExcludedPath(path string) bool {
	for _, pattern := range m.excludePaths {
		if pattern == "" {
			continue
		}
		if strings.HasSuffix(pattern, "*") {
			prefix := strings.TrimSuffix(pattern, "*")
			if strings.HasPrefix(path, prefix) {
				return true
			}
			continue
		}
		if path == pattern || strings.HasPrefix(path, pattern+"/") {
			return true
		}
	}
	return false
}

// containsBlacklistedUriList returns a list of strings
// which if they are found anywhere in the uri
// clearly indicate that there is a malicious bot/user
// trying to access them.
func (j *jailBotsMiddleware) containsBlacklistedUriList() []string {
	stopList := []string{
		"print(",
		"${print",
		".aws",
		".DS_Store",
		".env",
		".env.example",
		".git",
		".php",
		".vscode",
		".well-known/ALFA_DATA",
		".well-known/alfacgiapi",
		".well-known/cgialfa",
		"_ignition/health-check",
		"ALFA_DATA",
		"alfacgiapi",
		"search?folderIds=0",
		"aws/credentials",
		"backup",
		"backup/license.txt",
		"bc",
		"bk",
		"blog/license.txt",
		"bin",
		"cgialfa",
		"cloud-config.yml",
		"components/com_",
		"content/sitetree",
		"config.json",
		"cgi-bin",
		"credentials",
		"db",
		"ecp/Current/exporttool/microsoft.exchange.ediscovery.exporttool.application",
		"js/mage/cookies.js",
		"META-INF",
		"/main",
		"/new",
		"/old",
		"phpinfo",
		"server-status",
		"Telerik.Web.UI.WebResource.axd",
		"shop/license.txt",
		"sites/all/libraries/plupload/examples/upload.php",
		"simpla",
		"telescope/requests",
		"tmp/license.txt",
		"v2/_catalog",
		"wordpress",
		"wp",
		"www/license.txt",
	}

	// Check if we have any exclusion rules?

	if len(j.exclude) > 0 { // Check if exclude list is not empty
		stopList = lo.Filter(stopList, func(item string, index int) bool {
			return !slices.Contains(j.exclude, item)
		})
	}

	return stopList
}

// startsWithBlacklistedUriList returns a list of strings
// which if they are found at the start of the uri
// clearly indicate that there is a malicious bot/user
// trying to access them.
func (j jailBotsMiddleware) startsWithBlacklistedUriList() []string {
	startList := []string{
		"/content/sitetree",
		"/backup",
		"/bc",
		"/bk",
		"/main",
		"/new",
		"/old",
		"/tmp/",
		"/wordpress",
		"/wp",
		"/www",
	}

	// Apply exclusion rules to startsWith list as well
	if len(j.exclude) > 0 {
		startList = lo.Filter(startList, func(item string, index int) bool {
			// Check if any excluded item is contained in the start list item
			for _, exclude := range j.exclude {
				if strings.Contains(item, exclude) {
					return false // Exclude this item
				}
			}
			return true // Keep this item
		})
	}

	return startList
}
