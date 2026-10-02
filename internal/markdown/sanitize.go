package markdown

import (
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

var sanitizer = newSanitizer()

var (
	idRe       = regexp.MustCompile(`^[\p{L}\p{N}_:.-]+$`)
	numberRe   = regexp.MustCompile(`^[0-9]+$`)
	alignRe    = regexp.MustCompile(`^(?i)(left|right|center|justify|top|middle|bottom)$`)
	checkboxRe = regexp.MustCompile(`^(?i)checkbox$`)
)

// styleValueOK accepts simple CSS values (lengths, colors, keywords) and
// rejects url(), expression() and anything with escapes or comments.
func styleValueOK(v string) bool {
	lv := strings.ToLower(v)
	if strings.Contains(lv, "url") || strings.Contains(lv, "expression") ||
		strings.ContainsAny(v, `\/;"'<>{}`) {
		return false
	}
	return true
}

// svgElements lists the SVG elements allowed inside inline SVG.
// script, foreignObject, style, use, image and animation elements are excluded.
// Names are lower-cased because bluemonday matches against tokenizer output;
// the browser restores the camelCase names when it parses the inline SVG.
var svgElements = []string{
	"svg", "g", "path", "rect", "circle", "ellipse", "line", "polyline",
	"polygon", "text", "tspan", "defs", "lineargradient", "radialgradient",
	"stop", "marker", "title", "desc", "clippath", "symbol",
}

var svgAttrs = []string{
	"viewbox", "xmlns", "version", "width", "height", "x", "y", "cx", "cy",
	"r", "rx", "ry", "x1", "y1", "x2", "y2", "dx", "dy", "d", "points",
	"transform", "fill", "stroke", "stroke-width", "stroke-linecap",
	"stroke-linejoin", "stroke-dasharray", "stroke-dashoffset",
	"stroke-opacity", "stroke-miterlimit", "fill-opacity", "fill-rule",
	"clip-rule", "opacity", "font-size", "font-family", "font-weight",
	"font-style", "text-anchor", "dominant-baseline", "offset",
	"stop-color", "stop-opacity", "gradientunits", "gradienttransform",
	"markerwidth", "markerheight", "refx", "refy", "orient",
	"markerunits", "preserveaspectratio", "marker-start",
	"marker-end", "marker-mid", "textlength", "lengthadjust",
}

func newSanitizer() *bluemonday.Policy {
	p := bluemonday.NewPolicy()

	p.AllowElements(
		"a", "abbr", "b", "blockquote", "br", "caption", "center", "cite",
		"code", "col", "colgroup", "dd", "del", "details", "div", "dl", "dt",
		"em", "figcaption", "figure", "h1", "h2", "h3", "h4", "h5", "h6",
		"hr", "i", "img", "input", "ins", "kbd", "li", "mark", "ol", "p",
		"pre", "q", "s", "section", "small", "span", "strike", "strong",
		"sub", "summary", "sup", "table", "tbody", "td", "tfoot", "th",
		"thead", "tr", "u", "ul", "var",
	)
	p.AllowElements(svgElements...)

	p.RequireNoFollowOnLinks(false)
	p.AllowStandardURLs() // http, https, mailto, relative; parseable only
	p.AllowDataURIImages()

	p.AllowAttrs("class").Globally()
	p.AllowAttrs("id").Matching(idRe).Globally()
	p.AllowAttrs("title", "lang", "dir").Globally()
	p.AllowAttrs("align", "valign").Matching(alignRe).Globally()

	p.AllowAttrs("href").OnElements("a")
	p.AllowAttrs("name").Matching(idRe).OnElements("a")
	p.AllowAttrs("src", "alt", "width", "height").OnElements("img")
	p.AllowAttrs("open").OnElements("details")
	p.AllowAttrs("colspan", "rowspan").Matching(numberRe).OnElements("td", "th")
	p.AllowAttrs("start", "value").Matching(numberRe).OnElements("ol", "li")
	p.AllowAttrs("span", "width").OnElements("col", "colgroup")
	p.AllowAttrs("width", "height").OnElements("table", "td", "th")
	p.AllowAttrs("type").Matching(checkboxRe).OnElements("input")
	p.AllowAttrs("checked", "disabled").OnElements("input")
	p.AllowAttrs("cite").OnElements("blockquote", "q")

	// Sizing styles (used by <img style="max-width:..."> and code block attributes).
	p.AllowStyles("width", "height", "max-width", "max-height").
		MatchingHandler(styleValueOK).Globally()
	// Presentation styles for inline SVG.
	p.AllowStyles("fill", "stroke", "stroke-width", "opacity", "font-size",
		"font-family", "font-weight", "text-anchor").
		MatchingHandler(styleValueOK).OnElements(svgElements...)

	p.AllowAttrs(svgAttrs...).OnElements(svgElements...)

	return p
}

// Sanitize removes active content (scripts, event handlers, dangerous URLs,
// embedded objects) from rendered HTML while keeping the HTML subset used by
// Markdown documents: tables, images, details, code highlighting classes,
// line anchors and inline SVG.
func Sanitize(htmlContent []byte) []byte {
	return sanitizer.SanitizeBytes(htmlContent)
}
