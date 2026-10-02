package handler

import (
	"mime"
	"net/http"
	"strings"
)

// activeContentTypes are media types that browsers execute or render with
// script support when navigated to.
var activeContentTypes = map[string]bool{
	"text/html":             true,
	"application/xhtml+xml": true,
	"image/svg+xml":         true,
	"text/xml":              true,
	"application/xml":       true,
	"text/xsl":              true,
}

// isActiveContent reports whether contentType can run script in the page that
// loads it.
func isActiveContent(contentType string) bool {
	mt, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		mt = strings.ToLower(strings.TrimSpace(strings.SplitN(contentType, ";", 2)[0]))
	}
	return activeContentTypes[mt] || strings.HasSuffix(mt, "+xml")
}

// writePassthrough serves a non-Markdown file as-is, with headers that keep it
// from acting as a script on the proxy origin.
//
// HTML / SVG / XML get "Content-Security-Policy: sandbox allow-scripts", which
// gives the document an opaque origin: its scripts still run, but it cannot
// read other proxy URLs such as /local/... (those requests carry
// "Origin: null" and are rejected by the origin check).
func writePassthrough(w http.ResponseWriter, contentType string, data []byte) {
	h := w.Header()
	if contentType != "" {
		h.Set("Content-Type", contentType)
	}
	h.Set("X-Content-Type-Options", "nosniff")
	if isActiveContent(contentType) {
		h.Set("Content-Security-Policy", "sandbox allow-scripts")
	}
	w.Write(data)
}
