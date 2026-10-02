package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var ok = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })

func do(h http.Handler, host string, hdr map[string]string) int {
	req := httptest.NewRequest("GET", "http://"+host+"/local/x", nil)
	req.Host = host
	for k, v := range hdr {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec.Code
}

func TestHostCheck(t *testing.T) {
	h := hostCheckMiddleware(ok)
	for host, want := range map[string]int{
		"localhost:9080": 200, "127.0.0.1:9080": 200, "[::1]:9080": 200, "LOCALHOST": 200,
		"evil.example:9080": 403, "evil.example": 403, "127.0.0.1.evil.example:9080": 403,
		"localhost.evil.example": 403, "192.168.1.5:9080": 403,
	} {
		if got := do(h, host, nil); got != want {
			t.Errorf("Host %q: got %d, want %d", host, got, want)
		}
	}
}

func TestOriginCheck(t *testing.T) {
	h := originCheckMiddleware(ok, true)
	cases := []struct {
		name string
		hdr  map[string]string
		want int
	}{
		{"plain navigation", map[string]string{"Sec-Fetch-Mode": "navigate", "Sec-Fetch-Site": "none"}, 200},
		{"cross-site navigation", map[string]string{"Sec-Fetch-Mode": "navigate", "Sec-Fetch-Site": "cross-site"}, 200},
		{"cross-site img (no-cors)", map[string]string{"Sec-Fetch-Mode": "no-cors", "Sec-Fetch-Site": "cross-site"}, 200},
		{"same-origin fetch", map[string]string{"Sec-Fetch-Mode": "cors", "Sec-Fetch-Site": "same-origin"}, 200},
		{"cross-site fetch", map[string]string{"Sec-Fetch-Mode": "cors", "Sec-Fetch-Site": "cross-site"}, 403},
		{"same-site fetch (other port)", map[string]string{"Sec-Fetch-Mode": "cors", "Sec-Fetch-Site": "same-site"}, 403},
		{"sandboxed doc fetch", map[string]string{"Sec-Fetch-Mode": "cors", "Sec-Fetch-Site": "cross-site", "Origin": "null"}, 403},
		{"foreign Origin, no fetch metadata", map[string]string{"Origin": "http://evil.example"}, 403},
		{"Origin null, no fetch metadata", map[string]string{"Origin": "null"}, 403},
		{"matching Origin", map[string]string{"Origin": "http://localhost:9080"}, 200},
		{"no headers (curl)", nil, 200},
	}
	for _, c := range cases {
		if got := do(h, "localhost:9080", c.hdr); got != c.want {
			t.Errorf("%s: got %d, want %d", c.name, got, c.want)
		}
	}
}

func TestOriginCheck_RemoteModeSkipsOriginComparison(t *testing.T) {
	h := originCheckMiddleware(ok, false)
	if got := do(h, "docs.internal:9080", map[string]string{"Origin": "https://docs.internal"}); got != 200 {
		t.Errorf("got %d, want 200", got)
	}
	if got := do(h, "docs.internal:9080", map[string]string{"Sec-Fetch-Mode": "cors", "Sec-Fetch-Site": "cross-site"}); got != 403 {
		t.Errorf("got %d, want 403", got)
	}
}

func TestSecurityHeaders(t *testing.T) {
	rec := httptest.NewRecorder()
	securityHeadersMiddleware(ok).ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Header().Get("X-Content-Type-Options") != "nosniff" || rec.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Errorf("unexpected headers: %v", rec.Header())
	}
}
