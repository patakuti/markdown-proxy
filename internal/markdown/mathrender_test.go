package markdown

import (
	"strings"
	"testing"
)

func TestInlineMathHTMLEscape(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name:     "less than",
			input:    "$a<b$",
			contains: `<span class="math inline">\(a&lt;b\)</span>`,
		},
		{
			name:     "greater than",
			input:    "$a>b$",
			contains: `<span class="math inline">\(a&gt;b\)</span>`,
		},
		{
			name:     "ampersand",
			input:    "$a&b$",
			contains: `<span class="math inline">\(a&amp;b\)</span>`,
		},
		{
			name:     "mixed special chars",
			input:    "$a<b>c$",
			contains: `<span class="math inline">\(a&lt;b&gt;c\)</span>`,
		},
		{
			name:     "no special chars",
			input:    "$x^2 + y^2 = z^2$",
			contains: `<span class="math inline">\(x^2 + y^2 = z^2\)</span>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Convert([]byte(tt.input), "")
			if err != nil {
				t.Fatalf("Convert() error: %v", err)
			}
			got := string(result)
			if !strings.Contains(got, tt.contains) {
				t.Errorf("Convert(%q)\ngot:    %s\nexpect to contain: %s", tt.input, got, tt.contains)
			}
		})
	}
}

func TestMathBlockHTMLEscape(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name:     "less than in block",
			input:    "$$\na<b\n$$",
			contains: `<span class="math display">\[a&lt;b`,
		},
		{
			name:     "greater than in block",
			input:    "$$\na>b\n$$",
			contains: `<span class="math display">\[a&gt;b`,
		},
		{
			name:     "ampersand in block",
			input:    "$$\na&b\n$$",
			contains: `<span class="math display">\[a&amp;b`,
		},
		{
			name:     "no special chars in block",
			input:    "$$\nx^2 + y^2\n$$",
			contains: `<span class="math display">\[x^2 + y^2`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Convert([]byte(tt.input), "")
			if err != nil {
				t.Fatalf("Convert() error: %v", err)
			}
			got := string(result)
			if !strings.Contains(got, tt.contains) {
				t.Errorf("Convert(%q)\ngot:    %s\nexpect to contain: %s", tt.input, got, tt.contains)
			}
		})
	}
}
