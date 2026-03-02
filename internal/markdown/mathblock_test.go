package markdown

import (
	"testing"
)

func TestPreprocessMathBlocks(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{
			name:   "single-line display math",
			input:  "$$a=b$$",
			expect: "$$\na=b\n$$",
		},
		{
			name:   "single-line with spaces around content",
			input:  "$$ a = b $$",
			expect: "$$\n a = b \n$$",
		},
		{
			name:   "single-line with leading whitespace",
			input:  "  $$a=b$$",
			expect: "  $$\n  a=b\n  $$",
		},
		{
			name:   "single-line with trailing whitespace",
			input:  "$$a=b$$  ",
			expect: "$$\na=b\n$$",
		},
		{
			name:   "multi-line display math unchanged",
			input:  "$$\na=b\n$$",
			expect: "$$\na=b\n$$",
		},
		{
			name:   "inline math unchanged",
			input:  "text $a=b$ text",
			expect: "text $a=b$ text",
		},
		{
			name:   "inline display math unchanged (text around $$)",
			input:  "text $$a=b$$ text",
			expect: "text $$a=b$$ text",
		},
		{
			name:   "inside fenced code block backticks",
			input:  "```\n$$a=b$$\n```",
			expect: "```\n$$a=b$$\n```",
		},
		{
			name:   "inside fenced code block tildes",
			input:  "~~~\n$$a=b$$\n~~~",
			expect: "~~~\n$$a=b$$\n~~~",
		},
		{
			name:   "inside fenced code block with language",
			input:  "```math\n$$a=b$$\n```",
			expect: "```math\n$$a=b$$\n```",
		},
		{
			name:   "after fenced code block",
			input:  "```\ncode\n```\n$$a=b$$",
			expect: "```\ncode\n```\n$$\na=b\n$$",
		},
		{
			name:   "multiple math blocks",
			input:  "$$x=1$$\ntext\n$$y=2$$",
			expect: "$$\nx=1\n$$\ntext\n$$\ny=2\n$$",
		},
		{
			name:   "empty content between $$",
			input:  "$$$$",
			expect: "$$$$",
		},
		{
			name:   "complex LaTeX expression",
			input:  `$$\frac{a}{b} + \sqrt{c}$$`,
			expect: "$$\n\\frac{a}{b} + \\sqrt{c}\n$$",
		},
		{
			name:   "mixed content",
			input:  "# Title\n\n$$E=mc^2$$\n\nSome text\n\n$$\nF=ma\n$$",
			expect: "# Title\n\n$$\nE=mc^2\n$$\n\nSome text\n\n$$\nF=ma\n$$",
		},
		{
			name:   "nested fenced code blocks (longer fence closes)",
			input:  "````\n```\n$$a=b$$\n```\n````\n$$x=1$$",
			expect: "````\n```\n$$a=b$$\n```\n````\n$$\nx=1\n$$",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(PreprocessMathBlocks([]byte(tt.input)))
			if got != tt.expect {
				t.Errorf("PreprocessMathBlocks()\ngot:    %q\nexpect: %q", got, tt.expect)
			}
		})
	}
}
