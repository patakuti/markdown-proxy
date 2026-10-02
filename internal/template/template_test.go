package template

import (
	"regexp"
	"strings"
	"testing"
)

var (
	scriptTagRe  = regexp.MustCompile(`(?i)<script[^>]*>`)
	inlineHandRe = regexp.MustCompile(`(?i)\son[a-z]+\s*=`)
)

// Under the CSP (script-src 'nonce-...') every script needs the nonce and
// inline event handlers / javascript: URLs do not run.
func checkCSPCompatible(t *testing.T, name string, page []byte) {
	t.Helper()
	s := string(page)
	tags := scriptTagRe.FindAllString(s, -1)
	if len(tags) == 0 {
		t.Fatalf("%s: no script tags found", name)
	}
	for _, tag := range tags {
		if !strings.Contains(tag, `nonce="TESTNONCE"`) {
			t.Errorf("%s: script without nonce: %s", name, tag)
		}
	}
	if inlineHandRe.MatchString(s) {
		t.Errorf("%s: inline event handler found: %s", name, inlineHandRe.FindString(s))
	}
	if strings.Contains(strings.ToLower(s), "javascript:") {
		t.Errorf("%s: javascript: URL found", name)
	}
}

func TestPagesAreCSPCompatible(t *testing.T) {
	themes := []string{"github", "dark"}
	md, err := RenderMarkdown(&PageData{Title: "a.md", Theme: "github", Themes: themes, WatchPath: "/tmp/a.md", Nonce: "TESTNONCE"})
	if err != nil {
		t.Fatal(err)
	}
	checkCSPCompatible(t, "markdown", md)

	dir, err := RenderDirectory(&DirPageData{Title: "d", Path: "/d", Theme: "github", Themes: themes, WatchPath: "/d", Nonce: "TESTNONCE"})
	if err != nil {
		t.Fatal(err)
	}
	checkCSPCompatible(t, "directory", dir)

	e, err := RenderError(&ErrorPageData{Title: "e", Theme: "github", Themes: themes, Status: 403, Message: "m", Nonce: "TESTNONCE"})
	if err != nil {
		t.Fatal(err)
	}
	checkCSPCompatible(t, "error", e)
}
