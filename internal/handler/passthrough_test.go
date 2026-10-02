package handler

import (
	"net/http/httptest"
	"testing"
)

func TestWritePassthrough(t *testing.T) {
	cases := []struct {
		ct      string
		sandbox bool
	}{
		{"text/html", true},
		{"text/html; charset=utf-8", true},
		{"TEXT/HTML", true},
		{"image/svg+xml", true},
		{"application/xhtml+xml", true},
		{"application/xml", true},
		{"application/atom+xml", true},
		{"image/png", false},
		{"application/pdf", false},
		{"text/plain", false},
		{"text/css", false},
		{"application/javascript", false},
	}
	for _, c := range cases {
		rec := httptest.NewRecorder()
		writePassthrough(rec, c.ct, []byte("x"))
		csp := rec.Header().Get("Content-Security-Policy")
		if (csp != "") != c.sandbox {
			t.Errorf("%s: CSP=%q, want sandbox=%v", c.ct, csp, c.sandbox)
		}
		if rec.Header().Get("X-Content-Type-Options") != "nosniff" {
			t.Errorf("%s: missing nosniff", c.ct)
		}
	}
}
