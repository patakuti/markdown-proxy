package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/patakuti/markdown-proxy/internal/handler"
)

func authRequest(h http.Handler, path string, cookie *http.Cookie) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", path, nil)
	if cookie != nil {
		req.AddCookie(cookie)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestAuthMiddleware(t *testing.T) {
	h := authMiddleware(ok, "secret")
	cookie := func(v string) *http.Cookie { return &http.Cookie{Name: handler.CookieName, Value: v} }

	cases := []struct {
		name     string
		path     string
		cookie   *http.Cookie
		wantCode int
	}{
		{"no cookie redirects to login", "/https/example.com/a.md", nil, http.StatusSeeOther},
		{"wrong token redirects to login", "/", cookie("wrong"), http.StatusSeeOther},
		{"empty token redirects to login", "/", cookie(""), http.StatusSeeOther},
		{"prefix of token redirects to login", "/", cookie("secre"), http.StatusSeeOther},
		{"valid token passes", "/https/example.com/a.md", cookie("secret"), 200},
		{"login page is reachable without cookie", "/_login", nil, 200},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rec := authRequest(h, c.path, c.cookie)
			if rec.Code != c.wantCode {
				t.Errorf("status %d, want %d", rec.Code, c.wantCode)
			}
			if c.wantCode == http.StatusSeeOther && rec.Header().Get("Location") != "/_login" {
				t.Errorf("Location = %q", rec.Header().Get("Location"))
			}
		})
	}
}
