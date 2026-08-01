package markdown

import (
	"bytes"
	"regexp"
)

// listMarkerRe matches a list item marker (bullet or ordered) at the start of
// a line's content, followed by 2 or more spaces before the item's text.
// Capture groups: (1) leading indent (0-3 spaces), (2) marker, (3) content.
// The marker itself may only be preceded by 0-3 spaces, matching CommonMark's
// own rule for what counts as a (possibly nested) list marker rather than
// indented content.
var listMarkerRe = regexp.MustCompile(`^( {0,3})([-*+]|\d{1,9}[.)]) {2,}(\S.*)$`)

// PreprocessListMarkers normalizes the whitespace between a list marker and
// its content to exactly one space.
//
// CommonMark/GFM caps the "content indent" a list item establishes at 4
// columns: if a marker is followed by 5 or more spaces, only the first space
// is treated as the separator and the rest is read as part of the item's
// text, which then renders as an indented code block (4+ columns of leading
// whitespace) instead of plain text. Authors who pad list markers for visual
// alignment (e.g. "1.     a") trigger this by accident and get a silently
// different rendering. Normalizing the spacing up front removes the
// ambiguity: a list marker is always followed by exactly one space,
// regardless of how the source was typed.
//
// This only rewrites the marker's own line, so intentional multi-line
// indented content (nested code blocks, continuation paragraphs) on
// following lines is unaffected. Fenced code blocks (``` or ~~~) are left
// untouched, and only markers preceded by 0-3 spaces are rewritten, so
// already-indented code content is never touched.
//
// Deviation from strict CommonMark/GFM: this is an intentional divergence
// from GitHub's rendering, favoring predictability over spec compliance for
// this specific, easily-mistaken construct.
func PreprocessListMarkers(source []byte) []byte {
	lines := bytes.Split(source, []byte("\n"))
	var inFence bool
	var fenceMarker []byte

	for i, line := range lines {
		bqPrefix, body := splitBlockquotePrefix(line)

		if m := fenceStartRe.FindSubmatch(body); m != nil {
			marker := m[1]
			if inFence {
				if marker[0] == fenceMarker[0] && len(marker) >= len(fenceMarker) {
					inFence = false
					fenceMarker = nil
				}
			} else {
				inFence = true
				fenceMarker = marker
			}
			continue
		}
		if inFence {
			continue
		}

		if m := listMarkerRe.FindSubmatch(body); m != nil {
			normalized := append([]byte{}, m[1]...) // leading indent
			normalized = append(normalized, m[2]...) // marker
			normalized = append(normalized, ' ')
			normalized = append(normalized, m[3]...) // content
			lines[i] = prefixed(bqPrefix, normalized)
		}
	}

	return bytes.Join(lines, []byte("\n"))
}
