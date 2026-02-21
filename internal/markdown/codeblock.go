package markdown

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

// fencedBlockRe matches fenced code blocks with svg, mermaid, or plantuml language.
var fencedBlockRe = regexp.MustCompile("(?m)^```(svg|mermaid|plantuml)\\s*\n((?s:.*?))^```\\s*$")

// PreprocessCodeBlocks replaces svg, mermaid, and plantuml fenced code blocks
// in Markdown source with raw HTML before goldmark processing.
// Returns the processed source and is safe because goldmark is configured with html.WithUnsafe().
func PreprocessCodeBlocks(source []byte, plantumlServer string) []byte {
	return fencedBlockRe.ReplaceAllFunc(source, func(match []byte) []byte {
		parts := fencedBlockRe.FindSubmatch(match)
		if len(parts) < 3 {
			return match
		}
		lang := string(parts[1])
		code := string(parts[2])

		switch lang {
		case "svg":
			return []byte(fmt.Sprintf("\n<div class=\"svg-container\">\n%s</div>\n", code))
		case "mermaid":
			return []byte(fmt.Sprintf("\n<pre class=\"mermaid\">\n%s</pre>\n", code))
		case "plantuml":
			if plantumlServer != "" {
				encoded := encodePlantUML(code)
				imgURL := fmt.Sprintf("%s/svg/%s", strings.TrimRight(plantumlServer, "/"), encoded)
				return []byte(fmt.Sprintf("\n<div class=\"plantuml-container\"><img src=\"%s\" alt=\"PlantUML diagram\"></div>\n", imgURL))
			}
			return match
		}
		return match
	})
}

// encodePlantUML encodes PlantUML text for the PlantUML server URL.
// Uses the ~h (hex encoding) format: each byte is converted to its 2-digit hex representation.
func encodePlantUML(text string) string {
	return "~h" + hex.EncodeToString([]byte(strings.TrimSpace(text)))
}
