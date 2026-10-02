package markdown

import (
	"fmt"
	"strings"
	"testing"
)

func TestLineAnchors(t *testing.T) {
	src := "# H\n\npara1\npara2\n\n```go\ncode\n```\n\n| a | b |\n|---|---|\n| 1 | 2 |\n\n> quote\n"
	out, err := Convert([]byte(src), "")
	if err != nil {
		t.Fatal(err)
	}
	html := string(out)

	// heading, each paragraph line, fence line, table header/body cells, blockquote
	for _, line := range []int{1, 3, 4, 6, 10, 12, 14} {
		if !strings.Contains(html, fmt.Sprintf(`<span id="L%d"></span>`, line)) {
			t.Errorf("missing anchor for line %d in:\n%s", line, html)
		}
	}
	// blank lines never get an anchor
	for _, line := range []int{2, 5, 9, 13} {
		if strings.Contains(html, fmt.Sprintf(`id="L%d"`, line)) {
			t.Errorf("unexpected anchor for blank line %d", line)
		}
	}
}

func TestLineAnchors_CRLF(t *testing.T) {
	out, err := Convert([]byte("# H\r\n\r\npara\r\n"), "")
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range []int{1, 3} {
		if !strings.Contains(string(out), fmt.Sprintf(`id="L%d"`, line)) {
			t.Errorf("missing anchor for line %d with CRLF input:\n%s", line, out)
		}
	}
}

func TestLineAnchors_RawHTMLBlock(t *testing.T) {
	// anchors must not break raw HTML blocks
	out, err := Convert([]byte("<div>\nhello\n</div>\n\nafter\n"), "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "<div>") || !strings.Contains(string(out), `id="L5"`) {
		t.Errorf("unexpected output:\n%s", out)
	}
}
