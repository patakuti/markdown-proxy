package markdown

import (
	"strings"
	"testing"
)

func TestConvert_StripsActiveContent(t *testing.T) {
	cases := []struct {
		name, src string
		banned    []string
	}{
		{"script", "<script>alert(1)</script>\n", []string{"<script", "alert(1)"}},
		{"script in table", "<table><tr><td><script>x()</script></td></tr></table>\n", []string{"<script"}},
		{"img onerror", `<img src=x onerror="fetch('/local/etc/passwd')">` + "\n", []string{"onerror", "fetch("}},
		{"onclick", `<a href="#" onclick="x()">a</a>` + "\n", []string{"onclick"}},
		{"javascript href", `<a href="javascript:alert(1)">a</a>` + "\n", []string{"javascript:"}},
		{"javascript md link", "[a](javascript:alert(1))\n", []string{"javascript:"}},
		{"iframe", `<iframe src="/local/etc/passwd"></iframe>` + "\n", []string{"<iframe"}},
		{"object", `<object data="x"></object><embed src="x">` + "\n", []string{"<object", "<embed"}},
		{"form", `<form action="https://evil"><input name=a></form>` + "\n", []string{"<form"}},
		{"style tag", "<style>body{display:none}</style>\n", []string{"<style"}},
		{"style url", `<p style="background:url(https://evil/x)">a</p>` + "\n", []string{"url("}},
		{"base", `<base href="https://evil/">` + "\n", []string{"<base"}},
		{"meta refresh", `<meta http-equiv="refresh" content="0;url=https://evil">` + "\n", []string{"<meta"}},
		{"svg script", "```svg\n<svg><script>alert(1)</script><circle r=\"5\"/></svg>\n```\n", []string{"<script", "alert(1)"}},
		{"svg onload", "```svg\n<svg onload=\"alert(1)\"><circle r=\"5\"/></svg>\n```\n", []string{"onload"}},
		{"svg foreignObject", "```svg\n<svg><foreignObject><body onload=x()></foreignObject></svg>\n```\n", []string{"foreignobject", "onload"}},
		{"svg use href", "```svg\n<svg><use href=\"https://evil/x.svg#a\"/></svg>\n```\n", []string{"<use"}},
		{"mermaid injection", "```mermaid\n</pre><script>alert(1)</script>\n```\n", []string{"<script"}},
		{"data html link", `<a href="data:text/html,<script>alert(1)</script>">a</a>` + "\n", []string{"data:text/html"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out, err := Render([]byte(c.src), "", "local", "")
			if err != nil {
				t.Fatal(err)
			}
			lower := strings.ToLower(string(out))
			for _, b := range c.banned {
				if strings.Contains(lower, strings.ToLower(b)) {
					t.Errorf("output contains %q:\n%s", b, out)
				}
			}
		})
	}
}

func TestConvert_KeepsSafeHTML(t *testing.T) {
	cases := []struct {
		name, src string
		want      []string
	}{
		{"table", "| a | b |\n|---|:-:|\n| 1 | 2 |\n", []string{"<table>", "<th", "<td"}},
		{"raw table", "<table><tr><td align=\"center\">x</td></tr></table>\n", []string{"<table>", `align="center"`}},
		{"img width", `<img src="a.png" width="300" alt="x">` + "\n", []string{`src="a.png"`, `width="300"`, `alt="x"`}},
		{"img style", `<img src="a.png" style="max-width: 50%">` + "\n", []string{"max-width: 50%"}},
		{"md image", "![alt](https://example.com/a.png)\n", []string{`src="https://example.com/a.png"`}},
		{"details", "<details open><summary>s</summary>body</details>\n", []string{"<details open", "<summary>"}},
		{"task list", "- [x] done\n- [ ] todo\n", []string{`type="checkbox"`, "checked", "disabled"}},
		{"highlight", "```go\nfunc main() {}\n```\n", []string{`class="chroma"`, `<span class="`}},
		{"line anchor", "line1\n\nline3\n", []string{`id="L1"`, `id="L3"`}},
		{"heading id", "# Hello World\n", []string{"Hello World"}},
		{"math", "$x < y$\n\n$$\na & b\n$$\n", []string{`class="math inline"`, `class="math display"`}},
		{"mermaid", "```mermaid\ngraph TD\nA-->B\n```\n", []string{`<pre class="mermaid"`, "A--&gt;B"}},
		{"mermaid size", "```{.mermaid width=700}\ngraph TD\nA-->B\n```\n", []string{`max-width: 700px`}},
		{"svg block", "```svg\n<svg viewBox=\"0 0 10 10\" width=\"10\"><circle cx=\"5\" cy=\"5\" r=\"4\" fill=\"red\" style=\"stroke:blue\"/><linearGradient id=\"g\"><stop offset=\"0\" stop-color=\"#fff\"/></linearGradient></svg>\n```\n", []string{`class="svg-container"`, "<svg", "viewbox", "<circle", `fill="red"`, "stroke: blue", "lineargradient"}},
		{"plantuml img", "```plantuml\nA -> B\n```\n", nil},
		{"relative link", "[a](foo.md)\n[b](#sec)\n[c](mailto:a@b.c)\n", []string{`href="foo.md"`, `href="#sec"`, `href="mailto:a@b.c"`}},
		{"file link", "[a](file:///home/u/a.md)\n", []string{`href="/local/home/u/a.md"`}},
		{"line link", "[a](foo.md:12)\n", []string{`href="foo.md#L12"`}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out, err := Render([]byte(c.src), "", "local", "")
			if err != nil {
				t.Fatal(err)
			}
			for _, w := range c.want {
				if !strings.Contains(string(out), w) {
					t.Errorf("output missing %q:\n%s", w, out)
				}
			}
		})
	}
}

func TestConvert_PlantUMLImage(t *testing.T) {
	out, err := Render([]byte("```plantuml\nA -> B\n```\n"), "https://plantuml.example.com", "local", "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `<img src="https://plantuml.example.com/svg/~h`) {
		t.Errorf("plantuml image removed:\n%s", out)
	}
}
