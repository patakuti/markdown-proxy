package handler

import (
	"html"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/patakuti/markdown-proxy/internal/config"
	"github.com/patakuti/markdown-proxy/internal/csp"
)

func newLocal() http.Handler {
	h := NewLocalHandler(&config.Config{Theme: "github"}, []string{"github", "dark"})
	return csp.Middleware(h)
}

func get(h http.Handler, path string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", path, nil))
	return rec
}

// urlFor builds a /local/ URL for an absolute file path.
func urlFor(p string) string {
	return "/local/" + strings.TrimPrefix(filepath.ToSlash(p), "/")
}

func write(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLocal_Markdown(t *testing.T) {
	dir := t.TempDir()
	p := write(t, dir, "a.md", "# Title\n\n[next](b.md:5)\n\n<script>alert(1)</script>\n")

	rec := get(newLocal(), urlFor(p))
	body := rec.Body.String()
	if rec.Code != 200 {
		t.Fatalf("status %d: %s", rec.Code, body)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("Content-Type = %q", ct)
	}
	for _, want := range []string{`<h1 id="title">`, `href="b.md#L5"`, "/_sse?path="} {
		if !strings.Contains(body, want) {
			t.Errorf("missing %q", want)
		}
	}
	if strings.Contains(body, "alert(1)") {
		t.Error("raw script was not sanitized")
	}
	// the nonce is base64; html/template escapes '+' and '=' in attributes
	if !strings.Contains(html.UnescapeString(body), `nonce="`+nonceOf(rec)+`"`) {
		t.Error("page scripts do not carry the CSP nonce")
	}
}

// nonceOf extracts the script nonce from the CSP header.
func nonceOf(rec *httptest.ResponseRecorder) string {
	h := rec.Header().Get("Content-Security-Policy")
	const k = "script-src 'nonce-"
	i := strings.Index(h, k)
	if i < 0 {
		return ""
	}
	h = h[i+len(k):]
	return h[:strings.Index(h, "'")]
}

func TestLocal_Text(t *testing.T) {
	dir := t.TempDir()
	p := write(t, dir, "log.txt", "one\n<two>\n")

	rec := get(newLocal(), urlFor(p))
	body := rec.Body.String()
	if rec.Code != 200 || !strings.Contains(body, `id="L2" class="source-line">&lt;two&gt;`) {
		t.Errorf("status %d, body:\n%s", rec.Code, body)
	}
}

func TestLocal_NonMarkdownServedAsIs(t *testing.T) {
	dir := t.TempDir()
	png := write(t, dir, "a.png", "\x89PNG\r\n\x1a\n")
	html := write(t, dir, "a.html", "<script>1</script>")

	rec := get(newLocal(), urlFor(png))
	if rec.Header().Get("Content-Type") != "image/png" || rec.Body.String() != "\x89PNG\r\n\x1a\n" {
		t.Errorf("png: %q %q", rec.Header().Get("Content-Type"), rec.Body.String())
	}
	if rec.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Error("missing nosniff")
	}

	rec = get(newLocal(), urlFor(html))
	if csp := rec.Header().Get("Content-Security-Policy"); csp != "sandbox allow-scripts" {
		t.Errorf("html CSP = %q", csp)
	}
}

func TestLocal_Directory(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "b.md", "x")
	write(t, dir, ".hidden", "x")
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}

	rec := get(newLocal(), urlFor(dir)+"/")
	body := rec.Body.String()
	if rec.Code != 200 {
		t.Fatalf("status %d", rec.Code)
	}
	if !strings.Contains(body, ">sub</a>") || !strings.Contains(body, ">b.md</a>") || !strings.Contains(body, ">..</a>") {
		t.Errorf("listing incomplete:\n%s", body)
	}
	if strings.Contains(body, ".hidden") {
		t.Error("hidden files must not be listed")
	}
	if strings.Index(body, ">sub</a>") > strings.Index(body, ">b.md</a>") {
		t.Error("directories should be listed before files")
	}
}

func TestLocal_NotFound(t *testing.T) {
	rec := get(newLocal(), urlFor(filepath.Join(t.TempDir(), "nope.md")))
	if rec.Code != http.StatusNotFound {
		t.Errorf("status %d, want 404", rec.Code)
	}
}

func TestLocal_PathCleaning(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "a.md", "# ok")
	// "sub/.." segments are cleaned lexically before the file is opened
	rec := get(newLocal(), urlFor(dir)+"/sub/../a.md")
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), `<h1 id="ok">`) {
		t.Errorf("status %d:\n%s", rec.Code, rec.Body.String())
	}
}

func TestLocal_HomeExpansion(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home) // Windows
	write(t, home, "h.md", "# home")

	rec := get(newLocal(), "/local/~/h.md")
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), `<h1 id="home">`) {
		t.Errorf("status %d:\n%s", rec.Code, rec.Body.String())
	}
}

func TestIsLetter(t *testing.T) {
	for c, want := range map[byte]bool{'a': true, 'Z': true, '1': false, ':': false, '/': false} {
		if isLetter(c) != want {
			t.Errorf("isLetter(%q) = %v", c, !want)
		}
	}
}
