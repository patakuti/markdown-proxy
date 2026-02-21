package markdown

import (
	"regexp"
	"strings"
)

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
	// External URL with .md extension → proxy via /http/ or /https/
	if strings.HasPrefix(url, "https://") {
		if isMarkdownOrDir(url) {
			return "/https/" + strings.TrimPrefix(url, "https://")
		}
		return url
	}
	if strings.HasPrefix(url, "http://") {
		if isMarkdownOrDir(url) {
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

func isMarkdownOrDir(url string) bool {
	// Remove query string and fragment
	u := url
	if idx := strings.IndexAny(u, "?#"); idx >= 0 {
		u = u[:idx]
	}
	return strings.HasSuffix(u, ".md") || strings.HasSuffix(u, "/")
}
