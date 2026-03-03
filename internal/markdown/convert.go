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

// Convert converts Markdown source to HTML.
// plantumlServer is the PlantUML server URL for code block conversion.
func Convert(source []byte, plantumlServer string) ([]byte, error) {
	// Pre-process: expand single-line $$...$$ to multi-line for goldmark-mathjax
	source = PreprocessMathBlocks(source)

	// Pre-process: replace svg, mermaid, plantuml code blocks with raw HTML
	source = PreprocessCodeBlocks(source, plantumlServer)

	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
