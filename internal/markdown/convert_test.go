package markdown

import (
	"strings"
	"testing"
)

func TestConvert_NonASCIIHeadingAnchor(t *testing.T) {
	// Headings with non-ASCII (e.g. Japanese) characters must produce an id
	// matching GitHub GFM behavior so that links like [text](#ヘディング) work.
	tests := []struct {
		input    string
		wantID   string
		wantLink string
	}{
		{
			// goldmark percent-encodes non-ASCII in href; browsers decode before
			// fragment matching, so "#%E3%83%98..." navigates to id="ヘディング".
			input:    "# ヘディング\n\n[link](#ヘディング)\n",
			wantID:   `id="ヘディング"`,
			wantLink: `href="#%E3%83%98%E3%83%87%E3%82%A3%E3%83%B3%E3%82%B0"`,
		},
		{
			input:  "# Hello World\n",
			wantID: `id="hello-world"`,
		},
		{
			input:  "# 見出し Test\n",
			wantID: `id="見出し-test"`,
		},
	}
	for _, tt := range tests {
		html, err := Convert([]byte(tt.input), "")
		if err != nil {
			t.Fatalf("Convert error: %v", err)
		}
		s := string(html)
		if !strings.Contains(s, tt.wantID) {
			t.Errorf("expected %q in output:\n%s", tt.wantID, s)
		}
		if tt.wantLink != "" && !strings.Contains(s, tt.wantLink) {
			t.Errorf("expected %q in output:\n%s", tt.wantLink, s)
		}
	}
}

func TestConvert_CRLFMathBlock(t *testing.T) {
	// Regression: CRLF line endings must not break $$ math block rendering.
	input := []byte("$$\r\ny = x^2\r\n$$\r\n")
	html, err := Convert(input, "")
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}
	s := string(html)
	if !strings.Contains(s, `\(y = x^2\)`) && !strings.Contains(s, `y = x^2`) {
		t.Errorf("math content missing from output:\n%s", s)
	}
	if strings.Contains(s, `\(\)`) || strings.Contains(s, `\(\r`) {
		t.Errorf("CRLF leaked into math output:\n%s", s)
	}
}
