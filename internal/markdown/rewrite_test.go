package markdown

import "testing"

func TestRewriteURL(t *testing.T) {
	cases := []struct {
		name, url, scheme, server, want string
	}{
		// line references
		{"line ref", "foo.md:12", "local", "", "foo.md#L12"},
		{"line range", "foo.md:12-34", "local", "", "foo.md#L12-L34"},
		{"txt line ref", "notes/foo.txt:5", "local", "", "notes/foo.txt#L5"},
		{"markdown ext", "foo.markdown:3", "local", "", "foo.markdown#L3"},
		{"non-doc ext untouched", "foo.go:5", "local", "", "foo.go:5"},
		{"existing fragment untouched", "foo.md#sec", "local", "", "foo.md#sec"},
		{"line ref in file URL", "file:///home/u/a.md:7", "local", "", "/local/home/u/a.md#L7"},

		// file:///
		{"file URL", "file:///home/u/a.md", "local", "", "/local/home/u/a.md"},
		{"file URL windows drive", "file:///C:/Users/u/a.md", "local", "", "/local/C:/Users/u/a.md"},

		// absolute URLs
		{"https markdown", "https://github.com/o/r/blob/main/a.md", "local", "", "/https/github.com/o/r/blob/main/a.md"},
		{"https directory", "https://example.com/docs/", "local", "", "/https/example.com/docs/"},
		{"https image untouched", "https://example.com/a.png", "local", "", "https://example.com/a.png"},
		{"https page untouched", "https://example.com/page.html", "local", "", "https://example.com/page.html"},
		{"http markdown", "http://example.com/a.md", "local", "", "/http/example.com/a.md"},
		{"markdown with query", "https://example.com/a.md?x=1", "local", "", "/https/example.com/a.md?x=1"},

		// relative and fragment links are resolved by the browser
		{"relative", "sub/a.md", "local", "", "sub/a.md"},
		{"parent relative", "../a.md", "https", "github.com", "../a.md"},
		{"fragment", "#section", "local", "", "#section"},
		{"mailto", "mailto:a@b.c", "local", "", "mailto:a@b.c"},

		// absolute paths depend on where the document came from
		{"abs path local", "/docs/a.md", "local", "", "/local/docs/a.md"},
		{"abs path https", "/o/r/blob/main/a.md", "https", "github.com", "/https/github.com/o/r/blob/main/a.md"},
		{"abs path http", "/a.md", "http", "example.com", "/http/example.com/a.md"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := rewriteURL(c.url, c.scheme, c.server); got != c.want {
				t.Errorf("rewriteURL(%q, %q, %q) = %q, want %q", c.url, c.scheme, c.server, got, c.want)
			}
		})
	}
}

func TestRewriteLinks(t *testing.T) {
	in := `<p><a href="/docs/a.md:3">a</a> <a class="x" href="https://example.com/b.md">b</a> <img alt="i" src="/img/c.png"> <img src="https://example.com/d.png"></p>`
	want := `<p><a href="/local/docs/a.md#L3">a</a> <a class="x" href="/https/example.com/b.md">b</a> <img alt="i" src="/local/img/c.png"> <img src="https://example.com/d.png"></p>`
	if got := string(RewriteLinks([]byte(in), "local", "")); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}
