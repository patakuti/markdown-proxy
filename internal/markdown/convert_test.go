package markdown

import (
	"strings"
	"testing"
)

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
