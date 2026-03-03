package markdown

import (
	"bytes"
	"html"

	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// SafeInlineMathRenderer renders InlineMath nodes with HTML escaping.
// The upstream goldmark-mathjax renderer writes raw source bytes without
// escaping, so characters like <, >, & break browser rendering.
type SafeInlineMathRenderer struct{}

func (r *SafeInlineMathRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(mathjax.KindInlineMath, r.renderInlineMath)
}

func (r *SafeInlineMathRenderer) renderInlineMath(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString(`<span class="math inline">\(`)
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			segment := c.(*ast.Text).Segment
			value := segment.Value(source)
			if bytes.HasSuffix(value, []byte("\n")) {
				_, _ = w.WriteString(html.EscapeString(string(value[:len(value)-1])))
				if c != n.LastChild() {
					_, _ = w.Write([]byte(" "))
				}
			} else {
				_, _ = w.WriteString(html.EscapeString(string(value)))
			}
		}
		return ast.WalkSkipChildren, nil
	}
	_, _ = w.WriteString(`\)</span>`)
	return ast.WalkContinue, nil
}

// SafeMathBlockRenderer renders MathBlock nodes with HTML escaping.
type SafeMathBlockRenderer struct{}

func (r *SafeMathBlockRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(mathjax.KindMathBlock, r.renderMathBlock)
}

func (r *SafeMathBlockRenderer) renderMathBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString(`<p><span class="math display">\[`)
		l := node.Lines().Len()
		for i := 0; i < l; i++ {
			line := node.Lines().At(i)
			_, _ = w.WriteString(html.EscapeString(string(line.Value(source))))
		}
	} else {
		_, _ = w.WriteString(`\]</span></p>` + "\n")
	}
	return ast.WalkContinue, nil
}
