package markdown

import (
	"testing"
)

func TestPreprocessListMarkers(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{
			name:   "ordered marker with 5 spaces (the reported bug)",
			input:  "1.     a",
			expect: "1. a",
		},
		{
			name:   "ordered marker with 2 spaces stays normalized",
			input:  "1.  a",
			expect: "1. a",
		},
		{
			name:   "ordered marker with 1 space is unchanged",
			input:  "1. a",
			expect: "1. a",
		},
		{
			name:   "bullet marker with many spaces",
			input:  "-      a",
			expect: "- a",
		},
		{
			name:   "plus and asterisk bullets",
			input:  "*      a\n+      b",
			expect: "* a\n+ b",
		},
		{
			name:   "two-digit ordered marker",
			input:  "10.     a",
			expect: "10. a",
		},
		{
			name:   "paren-style ordered marker",
			input:  "1)     a",
			expect: "1) a",
		},
		{
			name:   "shallow nested bullet under bullet",
			input:  "- a\n  -     b",
			expect: "- a\n  - b",
		},
		{
			name:   "leading indent up to 3 spaces preserved",
			input:  "   -     a",
			expect: "   - a",
		},
		{
			name:   "internal spacing within content untouched",
			input:  "1.     a     b",
			expect: "1. a     b",
		},
		{
			name:   "marker-only empty item untouched",
			input:  "-   \nfoo",
			expect: "-   \nfoo",
		},
		{
			name:   "content inside fenced code block untouched",
			input:  "```\n1.     a\n```",
			expect: "```\n1.     a\n```",
		},
		{
			name:   "content inside tilde fenced code block untouched",
			input:  "~~~\n1.     a\n~~~",
			expect: "~~~\n1.     a\n~~~",
		},
		{
			name:   "indented (4+ space) code block untouched",
			input:  "    1.     a",
			expect: "    1.     a",
		},
		{
			name:   "list marker inside blockquote normalized",
			input:  "> 1.     a",
			expect: "> 1. a",
		},
		{
			name:   "thematic break with wide spacing survives as a thematic break",
			input:  "-   -   -",
			expect: "- -   -",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(PreprocessListMarkers([]byte(tt.input)))
			if got != tt.expect {
				t.Errorf("PreprocessListMarkers(%q) = %q, want %q", tt.input, got, tt.expect)
			}
		})
	}
}
