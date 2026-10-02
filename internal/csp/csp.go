// Package csp sets a per-request Content-Security-Policy with a script nonce.
package csp

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

type ctxKey struct{}

// policy is the Content-Security-Policy template; %s is replaced by the nonce.
//
//   - Only scripts carrying the nonce run, so script injected through Markdown
//     content (which never has the nonce) is inert even if sanitizing fails.
//   - connect-src 'self' keeps scripts from sending data to other hosts
//     (same-origin requests are covered by the origin check).
//   - style-src allows inline styles for Mermaid / KaTeX / highlighting and
//     the KaTeX stylesheet from the CDN.
//   - img-src allows any host: PlantUML servers and images in remote documents.
const policyPrefix = "default-src 'none'; script-src 'nonce-"
const policySuffix = "'; style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net; " +
	"font-src https://cdn.jsdelivr.net data:; img-src * data:; connect-src 'self'; " +
	"base-uri 'none'; form-action 'self'; frame-ancestors 'none'; object-src 'none'"

func newNonce() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic("csp: cannot read random bytes: " + err.Error())
	}
	return base64.StdEncoding.EncodeToString(b)
}

// Middleware generates a nonce for each request, sets the
// Content-Security-Policy header and stores the nonce for Nonce.
// Handlers that must use a different policy (e.g. files served as-is) overwrite the header.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nonce := newNonce()
		w.Header().Set("Content-Security-Policy", policyPrefix+nonce+policySuffix)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKey{}, nonce)))
	})
}

// Nonce returns the script nonce of the request, or "" if Middleware did not run.
func Nonce(r *http.Request) string {
	n, _ := r.Context().Value(ctxKey{}).(string)
	return n
}
