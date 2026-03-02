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

// PreprocessMathBlocks converts single-line display math ($$...$$ on one line)
// to multi-line format ($$\n...\n$$) so that goldmark-mathjax can parse it correctly.
// Fenced code blocks are skipped to avoid modifying code content.
func PreprocessMathBlocks(source []byte) []byte {
	lines := bytes.Split(source, []byte("\n"))
	var result [][]byte
	var inFence bool
	var fenceMarker []byte

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

		// Check for single-line $$...$$ pattern
		if sm := singleLineMathRe.FindSubmatch(line); sm != nil {
			indent := sm[1]
			content := sm[2]
			result = append(result, append(append([]byte{}, indent...), []byte("$$")...))
			result = append(result, append(append([]byte{}, indent...), content...))
			result = append(result, append(append([]byte{}, indent...), []byte("$$")...))
			continue
		}

		result = append(result, line)
	}

	return bytes.Join(result, []byte("\n"))
}
