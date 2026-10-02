package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/patakuti/markdown-proxy/internal/config"
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

// A bearer token sent to one host must not follow a redirect to another host.
func TestDoFetch_DropsAuthorizationOnCrossHostRedirect(t *testing.T) {
	var gotAuth string
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.Write([]byte("ok"))
	}))
	defer target.Close()
	// "localhost" and "127.0.0.1" are different hosts for redirect purposes.
	targetURL := strings.Replace(target.URL, "127.0.0.1", "localhost", 1)

	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			t.Error("origin did not receive Authorization")
		}
		http.Redirect(w, r, targetURL, http.StatusFound)
	}))
	defer origin.Close()

	h := NewRemoteHandler(&config.Config{}, http.DefaultClient, nil)
	if _, _, err := h.doFetch(origin.URL, "user", "secret-token", true); err != nil {
		t.Fatal(err)
	}
	if gotAuth != "" {
		t.Errorf("Authorization leaked to redirect target: %q", gotAuth)
	}
}
