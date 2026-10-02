package csp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMiddleware(t *testing.T) {
	var got string
	h := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { got = Nonce(r) }))

	rec1 := httptest.NewRecorder()
	h.ServeHTTP(rec1, httptest.NewRequest("GET", "/", nil))
	first := got
	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, httptest.NewRequest("GET", "/", nil))

	if first == "" || first == got {
		t.Errorf("nonce must be non-empty and differ per request: %q %q", first, got)
	}
	header := rec1.Header().Get("Content-Security-Policy")
	if !strings.Contains(header, "script-src 'nonce-"+first+"';") {
		t.Errorf("header does not carry the request nonce: %s", header)
	}
	for _, want := range []string{"default-src 'none'", "connect-src 'self'", "frame-ancestors 'none'", "object-src 'none'"} {
		if !strings.Contains(header, want) {
			t.Errorf("header missing %q: %s", want, header)
		}
	}
	if strings.Contains(header, "script-src 'unsafe-inline'") || strings.Contains(header, "unsafe-eval") {
		t.Errorf("script-src must not allow unsafe-inline/eval: %s", header)
	}
}
