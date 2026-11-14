package markdown

import (
	"regexp"
	"strings"
)

// lineRefRe matches :line or :line-line at the end of a markdown file URL.
// e.g., "foo.md:12" or "foo.md:12-34"
var lineRefRe = regexp.MustCompile(`^(.*\.(?:md|markdown|txt)):(\d+)(?:-(\d+))?$`)

var (
	hrefRe = regexp.MustCompile(`(<a\s[^>]*href=")([^"]+)(")`)
	srcRe  = regexp.MustCompile(`(<img\s[^>]*src=")([^"]+)(")`)
)

// RewriteLinks rewrites links in HTML content for proxy navigation.
// scheme is "local", "http", or "https".
// server is the remote server host (empty for local).
func RewriteLinks(htmlContent []byte, scheme string, server string) []byte {
	result := hrefRe.ReplaceAllFunc(htmlContent, func(match []byte) []byte {
		parts := hrefRe.FindSubmatch(match)
		if len(parts) < 4 {
			return match
		}
		prefix := parts[1]
		url := string(parts[2])
		suffix := parts[3]

		rewritten := rewriteURL(url, scheme, server)
		return append(append(prefix, []byte(rewritten)...), suffix...)
	})

	result = srcRe.ReplaceAllFunc(result, func(match []byte) []byte {
		parts := srcRe.FindSubmatch(match)
		if len(parts) < 4 {
			return match
		}
		prefix := parts[1]
		url := string(parts[2])
		suffix := parts[3]

		rewritten := rewriteURL(url, scheme, server)
		return append(append(prefix, []byte(rewritten)...), suffix...)
	})

	return result
}

func rewriteURL(url, scheme, server string) string {
	// Convert line references: file.md:12 → file.md#L12, file.md:12-34 → file.md#L12-L34
	url = convertLineRef(url)

	// file:/// URL → proxy via /local/
	if strings.HasPrefix(url, "file:///") {
		// file:///path/to/file.md → /local/path/to/file.md
		return "/local" + strings.TrimPrefix(url, "file://")
	}

	// External URL with .md extension → proxy via /http/ or /https/
	if strings.HasPrefix(url, "https://") {
		if isMarkdown(url) {
			return "/https/" + strings.TrimPrefix(url, "https://")
		}
		return url
	}
	if strings.HasPrefix(url, "http://") {
		if isMarkdown(url) {
			return "/http/" + strings.TrimPrefix(url, "http://")
		}
		return url
	}

	// Relative links are resolved by the browser naturally - no rewrite needed
	if !strings.HasPrefix(url, "/") {
		return url
	}

	// Absolute path links
	switch scheme {
	case "local":
		return "/local" + url
	case "http":
		return "/http/" + server + url
	case "https":
		return "/https/" + server + url
	}

	return url
}

// convertLineRef converts file.md:12 to file.md#L12 and file.md:12-34 to file.md#L12-L34.
func convertLineRef(url string) string {
	// Strip existing fragment/query for matching
	base := url
	if idx := strings.IndexAny(base, "?#"); idx >= 0 {
		base = base[:idx]
	}

	m := lineRefRe.FindStringSubmatch(base)
	if m == nil {
		return url
	}

	filePath := m[1]
	startLine := m[2]
	endLine := m[3]

	if endLine != "" {
		return filePath + "#L" + startLine + "-L" + endLine
	}
	return filePath + "#L" + startLine
}

func isMarkdown(url string) bool {
	// Remove query string and fragment
	u := url
	if idx := strings.IndexAny(u, "?#"); idx >= 0 {
		u = u[:idx]
	}
	return strings.HasSuffix(u, ".md")
}
