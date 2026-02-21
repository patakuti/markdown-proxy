package template

import (
	"fmt"
	"html/template"
)

type PageData struct {
	Title   string
	Content template.HTML
	Theme   string
}

type DirEntry struct {
	Name  string
	IsDir bool
	URL   string
}

type DirPageData struct {
	Title   string
	Path    string
	Entries []DirEntry
	Theme   string
}

func RenderMarkdown(data *PageData) ([]byte, error) {
	return renderTemplate(markdownPageTpl, data)
}

func RenderDirectory(data *DirPageData) ([]byte, error) {
	return renderTemplate(dirPageTpl, data)
}

func renderTemplate(tplStr string, data interface{}) ([]byte, error) {
	tpl, err := template.New("page").Parse(tplStr)
	if err != nil {
		return nil, fmt.Errorf("template parse error: %w", err)
	}
	var buf []byte
	w := &byteWriter{buf: &buf}
	if err := tpl.Execute(w, data); err != nil {
		return nil, fmt.Errorf("template execute error: %w", err)
	}
	return buf, nil
}

type byteWriter struct {
	buf *[]byte
}

func (w *byteWriter) Write(p []byte) (int, error) {
	*w.buf = append(*w.buf, p...)
	return len(p), nil
}

const markdownPageTpl = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Title}} - markdown-proxy</title>
<style>` + githubCSS + `</style>
<style>` + simpleCSS + `</style>
<style>` + darkCSS + `</style>
<style>` + commonCSS + `</style>
<script src="https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js"></script>
</head>
<body class="theme-{{.Theme}}">
<div class="toolbar">
  <a href="/" class="home-link">markdown-proxy</a>
  <div class="theme-switcher">
    <label>Theme:</label>
    <select onchange="switchTheme(this.value)">
      <option value="github"{{if eq .Theme "github"}} selected{{end}}>GitHub</option>
      <option value="simple"{{if eq .Theme "simple"}} selected{{end}}>Simple</option>
      <option value="dark"{{if eq .Theme "dark"}} selected{{end}}>Dark</option>
    </select>
  </div>
</div>
<div class="markdown-body">
{{.Content}}
</div>
<script>
mermaid.initialize({startOnLoad: true, theme: document.body.className.includes('dark') ? 'dark' : 'default'});
function switchTheme(theme) {
  document.body.className = 'theme-' + theme;
  localStorage.setItem('mdproxy_theme', theme);
  mermaid.initialize({startOnLoad: false, theme: theme === 'dark' ? 'dark' : 'default'});
}
(function() {
  var saved = localStorage.getItem('mdproxy_theme');
  if (saved) {
    document.body.className = 'theme-' + saved;
    var sel = document.querySelector('.theme-switcher select');
    if (sel) sel.value = saved;
  }
})();
</script>
</body>
</html>`

const dirPageTpl = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Title}} - markdown-proxy</title>
<style>` + githubCSS + `</style>
<style>` + simpleCSS + `</style>
<style>` + darkCSS + `</style>
<style>` + commonCSS + `</style>
</head>
<body class="theme-{{.Theme}}">
<div class="toolbar">
  <a href="/" class="home-link">markdown-proxy</a>
  <div class="theme-switcher">
    <label>Theme:</label>
    <select onchange="switchTheme(this.value)">
      <option value="github"{{if eq .Theme "github"}} selected{{end}}>GitHub</option>
      <option value="simple"{{if eq .Theme "simple"}} selected{{end}}>Simple</option>
      <option value="dark"{{if eq .Theme "dark"}} selected{{end}}>Dark</option>
    </select>
  </div>
</div>
<div class="markdown-body">
<h1>{{.Path}}</h1>
<table>
<thead><tr><th>Name</th><th>Type</th></tr></thead>
<tbody>
{{range .Entries}}
<tr>
  <td><a href="{{.URL}}">{{.Name}}</a></td>
  <td>{{if .IsDir}}Directory{{else}}File{{end}}</td>
</tr>
{{end}}
</tbody>
</table>
</div>
<script>
function switchTheme(theme) {
  document.body.className = 'theme-' + theme;
  localStorage.setItem('mdproxy_theme', theme);
}
(function() {
  var saved = localStorage.getItem('mdproxy_theme');
  if (saved) {
    document.body.className = 'theme-' + saved;
    var sel = document.querySelector('.theme-switcher select');
    if (sel) sel.value = saved;
  }
})();
</script>
</body>
</html>`
