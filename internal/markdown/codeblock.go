package markdown

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

// fencedBlockRe matches fenced code blocks with svg, mermaid, or plantuml language.
// Supports both standard syntax (```mermaid) and Pandoc-style attributes (```{.mermaid width=700}).
// Capture groups: (1) language, (2) attribute string, (3) code content.
var fencedBlockRe = regexp.MustCompile("(?m)^```(?:\\{\\.)?(svg|mermaid|plantuml)([^}\\n]*)\\}?\\s*\\n((?s:.*?))^```\\s*$")

var attrPairRe = regexp.MustCompile(`(\w+)=["']?([^"'\s}]*)["']?`)
var cssLengthRe = regexp.MustCompile(`^[\d.]+(px|%|em|rem|vh|vw|ch)?$`)

// PreprocessCodeBlocks replaces svg, mermaid, and plantuml fenced code blocks
// in Markdown source with raw HTML before goldmark processing.
// Returns the processed source and is safe because goldmark is configured with html.WithUnsafe().
func PreprocessCodeBlocks(source []byte, plantumlServer string) []byte {
	return fencedBlockRe.ReplaceAllFunc(source, func(match []byte) []byte {
		parts := fencedBlockRe.FindSubmatch(match)
		if len(parts) < 4 {
			return match
		}
		lang := string(parts[1])
		attrs := parseAttrs(string(parts[2]))
		code := string(parts[3])
		style := buildStyleAttr(attrs)

		switch lang {
		case "svg":
			// Remove blank lines from SVG content to prevent goldmark from
			// splitting the HTML block (CommonMark HTML block type 6 ends at
			// a blank line).
			code = removeBlankLines(code)
			return []byte(fmt.Sprintf("\n<div class=\"svg-container\"%s>\n%s</div>\n", style, code))
		case "mermaid":
			return []byte(fmt.Sprintf("\n<pre class=\"mermaid\"%s>\n%s</pre>\n", style, code))
		case "plantuml":
			if plantumlServer != "" {
				encoded := encodePlantUML(code)
				imgURL := fmt.Sprintf("%s/svg/%s", strings.TrimRight(plantumlServer, "/"), encoded)
				return []byte(fmt.Sprintf("\n<div class=\"plantuml-container\"%s><img src=\"%s\" alt=\"PlantUML diagram\"></div>\n", style, imgURL))
			}
			return []byte("\n<div class=\"plantuml-notice\">" +
				"<strong>PlantUML rendering is disabled.</strong> " +
				"To enable, start with <code>--plantuml-server URL</code> " +
				"or run <code>markdown-proxy --configure</code> to set up." +
				"</div>\n")
		}
		return match
	})
}

// parseAttrs parses key=value pairs from a Pandoc-style attribute string.
// e.g., " width=700 height=400" → {"width": "700", "height": "400"}
func parseAttrs(s string) map[string]string {
	attrs := make(map[string]string)
	for _, m := range attrPairRe.FindAllStringSubmatch(s, -1) {
		attrs[m[1]] = m[2]
	}
	return attrs
}

// sanitizeCSSLength validates and normalizes a CSS length value.
// Bare numbers (e.g. "700") are treated as pixels. Returns empty string if invalid.
func sanitizeCSSLength(v string) string {
	if !cssLengthRe.MatchString(v) {
		return ""
	}
	for _, unit := range []string{"px", "%", "em", "rem", "vh", "vw", "ch"} {
		if strings.HasSuffix(v, unit) {
			return v
		}
	}
	return v + "px"
}

// buildStyleAttr returns an HTML style attribute string (e.g. ` style="max-width: 700px"`)
// from attrs. Supported keys: width, height. Returns empty string if no relevant attrs.
func buildStyleAttr(attrs map[string]string) string {
	var parts []string
	if w := sanitizeCSSLength(attrs["width"]); w != "" {
		parts = append(parts, "max-width: "+w)
	}
	if h := sanitizeCSSLength(attrs["height"]); h != "" {
		parts = append(parts, "max-height: "+h)
	}
	if len(parts) == 0 {
		return ""
	}
	return fmt.Sprintf(` style="%s"`, strings.Join(parts, "; "))
}

// removeBlankLines removes blank lines (empty or whitespace-only) from the text.
func removeBlankLines(s string) string {
	var b strings.Builder
	for _, line := range strings.Split(s, "\n") {
		if strings.TrimSpace(line) != "" {
			b.WriteString(line)
			b.WriteByte('\n')
		}
	}
	return b.String()
}

// encodePlantUML encodes PlantUML text for the PlantUML server URL.
// Uses the ~h (hex encoding) format: each byte is converted to its 2-digit hex representation.
func encodePlantUML(text string) string {
	return "~h" + hex.EncodeToString([]byte(strings.TrimSpace(text)))
}
