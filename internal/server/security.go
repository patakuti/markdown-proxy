package server

import (
	"net"
	"net/http"
	"net/url"
	"strings"
)

// isLoopbackHost reports whether the host part of a Host header is a loopback name.
func isLoopbackHost(hostport string) bool {
	host := hostport
	if h, _, err := net.SplitHostPort(hostport); err == nil {
		host = h
	}
	host = strings.Trim(host, "[]")
	switch strings.ToLower(host) {
	case "localhost", "127.0.0.1", "::1":
		return true
	}
	return false
}

// hostCheckMiddleware rejects requests whose Host header is not a loopback
// name. This blocks DNS rebinding: a page on attacker.example that re-resolves
// to 127.0.0.1 still sends "Host: attacker.example".
func hostCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isLoopbackHost(r.Host) {
			http.Error(w, "Forbidden: invalid Host header", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// originCheckMiddleware rejects cross-origin requests that could read
// responses from the browser (fetch / XHR / EventSource).
//
//   - Sec-Fetch-Mode: cors with Sec-Fetch-Site other than same-origin is rejected.
//   - When checkOrigin is true (local mode), a present Origin header must match
//     the Host header. This covers browsers without Fetch Metadata and the
//     "Origin: null" of sandboxed documents.
//
// Top-level navigations and no-cors subresource loads (img, script, css) are
// allowed: the requesting page cannot read those responses. No CORS headers
// are ever sent.
func originCheckMiddleware(next http.Handler, checkOrigin bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Sec-Fetch-Mode") == "cors" && r.Header.Get("Sec-Fetch-Site") != "same-origin" {
			http.Error(w, "Forbidden: cross-origin request", http.StatusForbidden)
			return
		}
		if origin := r.Header.Get("Origin"); checkOrigin && origin != "" {
			u, err := url.Parse(origin)
			if err != nil || u.Scheme != "http" || u.Host != r.Host {
				http.Error(w, "Forbidden: cross-origin request", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// securityHeadersMiddleware adds headers that apply to every response.
func securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "no-referrer")
		next.ServeHTTP(w, r)
	})
}
