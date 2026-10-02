package markdown

import (
	"bytes"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var md goldmark.Markdown

func init() {
	md = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
				highlighting.WithFormatOptions(
					chromahtml.WithClasses(true),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithBlockParsers(
				util.Prioritized(mathjax.NewMathJaxBlockParser(), 701),
			),
			parser.WithInlineParsers(
				util.Prioritized(mathjax.NewInlineMathParser(), 501),
			),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
			renderer.WithNodeRenderers(
				util.Prioritized(&SafeMathBlockRenderer{}, 501),
				util.Prioritized(&SafeInlineMathRenderer{}, 502),
			),
		),
	)
}

// Convert converts Markdown source to HTML. The output is NOT sanitized; use Render
// for anything that is sent to the browser.
// plantumlServer is the PlantUML server URL for code block conversion.
func Convert(source []byte, plantumlServer string) ([]byte, error) {
	// Normalize CRLF to LF so all preprocessors and the parser see consistent line endings.
	source = bytes.ReplaceAll(source, []byte("\r\n"), []byte("\n"))

	// Pre-process: normalize spacing after list markers to avoid accidental
	// indented-code-block interpretation (see PreprocessListMarkers doc comment)
	source = PreprocessListMarkers(source)

	// Pre-process: expand single-line $$...$$ to multi-line for goldmark-mathjax
	source = PreprocessMathBlocks(source)

	// Pre-process: replace svg, mermaid, plantuml code blocks with raw HTML
	source = PreprocessCodeBlocks(source, plantumlServer)

	// Parse into AST using GFM-compatible heading ID generation.
	reader := text.NewReader(source)
	ctx := parser.NewContext(parser.WithIDs(newGFMIDs()))
	doc := md.Parser().Parse(reader, parser.WithContext(ctx))

	// Insert line anchors into AST
	source = insertLineAnchors(doc, source)

	// Render AST to HTML
	var buf bytes.Buffer
	if err := md.Renderer().Render(&buf, source, doc); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Render converts Markdown source to HTML that is safe to send to the browser.
// It runs Convert, rewrites links for proxy navigation (see RewriteLinks), and
// finally sanitizes the result.
//
// goldmark runs with html.WithUnsafe() so that raw HTML (tables, images,
// details, inline SVG) reaches the renderer; Sanitize is the single place that
// decides what reaches the browser. It must run last, after link rewriting,
// because rewritten URLs (e.g. file:/// -> /local/...) would otherwise be
// dropped as disallowed schemes. Handlers must use Render, not Convert.
func Render(source []byte, plantumlServer, scheme, server string) ([]byte, error) {
	htmlContent, err := Convert(source, plantumlServer)
	if err != nil {
		return nil, err
	}
	return Sanitize(RewriteLinks(htmlContent, scheme, server)), nil
}
