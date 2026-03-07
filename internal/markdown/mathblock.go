package markdown

import (
	"bytes"
	"regexp"
)

// singleLineMathRe matches a line containing only $$<content>$$ (with optional surrounding whitespace).
// This does NOT match inline usage like "text $$a=b$$ text" because the line must start with $$.
var singleLineMathRe = regexp.MustCompile(`^(\s*)\$\$(.+)\$\$\s*$`)

// fenceStartRe matches the start of a fenced code block (``` or ~~~).
var fenceStartRe = regexp.MustCompile(`^(\x60{3,}|~{3,})`)

// PreprocessMathBlocks ensures that $$ delimiters are on their own lines
// so that goldmark-mathjax can parse them correctly.
// It handles three cases:
//   - Single-line: $$content$$ → $$\ncontent\n$$
//   - Opening with content: $$content...\n → $$\ncontent...\n
//   - Closing with content: ...content$$ → ...content\n$$
//
// Fenced code blocks are skipped to avoid modifying code content.
// Only exactly two consecutive $ characters are treated as delimiters
// ($$$ or more are left for goldmark-mathjax to handle directly).
func PreprocessMathBlocks(source []byte) []byte {
	lines := bytes.Split(source, []byte("\n"))
	var result [][]byte
	var inFence bool
	var fenceMarker []byte
	var inMathBlock bool

	for _, line := range lines {
		if inFence {
			// Check if this line closes the fence
			if m := fenceStartRe.FindSubmatch(line); m != nil {
				marker := m[1]
				if marker[0] == fenceMarker[0] && len(marker) >= len(fenceMarker) {
					inFence = false
					fenceMarker = nil
				}
			}
			result = append(result, line)
			continue
		}

		// Check if this line opens a fenced code block
		if m := fenceStartRe.FindSubmatch(line); m != nil {
			inFence = true
			fenceMarker = m[1]
			result = append(result, line)
			continue
		}

		if inMathBlock {
			trimmed := bytes.TrimRight(line, " \t")
			if countConsecutiveDollarsAt(trimmed, len(trimmed)-2) == 2 {
				before := trimmed[:len(trimmed)-2]
				if len(bytes.TrimSpace(before)) > 0 {
					// Content before closing $$ → split
					result = append(result, before)
					result = append(result, []byte("$$"))
				} else {
					// Just $$ → proper closing delimiter
					result = append(result, line)
				}
				inMathBlock = false
			} else {
				result = append(result, line)
			}
			continue
		}

		// Check for single-line $$...$$ pattern
		if sm := singleLineMathRe.FindSubmatch(line); sm != nil {
			indent := sm[1]
			content := sm[2]
			result = append(result, append(append([]byte{}, indent...), []byte("$$")...))
			result = append(result, append(append([]byte{}, indent...), content...))
			result = append(result, append(append([]byte{}, indent...), []byte("$$")...))
			continue
		}

		// Check for opening $$ (exactly 2 consecutive $)
		trimmed := bytes.TrimLeft(line, " \t")
		if countConsecutiveDollarsAt(trimmed, 0) == 2 {
			indent := line[:len(line)-len(trimmed)]
			content := trimmed[2:]
			if len(bytes.TrimSpace(content)) > 0 {
				// $$ followed by content → split
				result = append(result, append(append([]byte{}, indent...), []byte("$$")...))
				result = append(result, append(append([]byte{}, indent...), content...))
			} else {
				// Just $$ → opening delimiter
				result = append(result, line)
			}
			inMathBlock = true
			continue
		}

		result = append(result, line)
	}

	return bytes.Join(result, []byte("\n"))
}

// countConsecutiveDollarsAt counts how many consecutive '$' characters
// surround position pos. It returns the total length of the '$' run
// that includes the character at pos.
// If pos is out of range or the character at pos is not '$', it returns 0.
func countConsecutiveDollarsAt(b []byte, pos int) int {
	if pos < 0 || pos >= len(b) || b[pos] != '$' {
		return 0
	}
	// Find the start of the $ run
	start := pos
	for start > 0 && b[start-1] == '$' {
		start--
	}
	// Find the end of the $ run
	end := pos + 1
	for end < len(b) && b[end] == '$' {
		end++
	}
	return end - start
}
