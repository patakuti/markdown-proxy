package template

import (
	"fmt"
	"html/template"
)

type PageData struct {
	Title     string
	Content   template.HTML
	Theme     string
	Themes    []string
	WatchPath string
	SourceURL string
}

type DirEntry struct {
	Name  string
	IsDir bool
	URL   string
}

type DirPageData struct {
	Title     string
	Path      string
	Entries   []DirEntry
	Theme     string
	Themes    []string
	WatchPath string
}

type ErrorPageData struct {
	Title   string
	Theme   string
	Themes  []string
	Status  int
	Message string
	Hints   []template.HTML
}

func RenderMarkdown(data *PageData) ([]byte, error) {
	return renderTemplate(markdownPageTpl, data)
}

func RenderDirectory(data *DirPageData) ([]byte, error) {
	return renderTemplate(dirPageTpl, data)
}

func RenderError(data *ErrorPageData) ([]byte, error) {
	return renderTemplate(errorPageTpl, data)
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

var markdownPageTpl = markdownPageTplHead + markdownPageTplTail

const markdownPageTplHead = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Title}} - markdown-proxy</title>
<style>` + commonCSS + `</style>
<style>` + lineAnchorCSS + `</style>
<style>` + tocCSS + `</style>
<style>` + copyButtonCSS + `</style>
<link id="theme-css" rel="stylesheet" href="/_theme/{{.Theme}}.css">
<script>var MDPROXY_THEMES=[{{range $i,$t:=.Themes}}{{if $i}},{{end}}"{{$t}}"{{end}}];</script>
`

const markdownPageTplTail = `<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/katex@0.16.11/dist/katex.min.css">
<script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.11/dist/katex.min.js"></script>
<script src="https://cdn.jsdelivr.net/npm/mermaid@10/dist/mermaid.min.js"></script>
</head>
<body class="theme-{{.Theme}}">
<div class="toolbar">
  <a href="/" class="home-link">markdown-proxy</a>
  <div class="toolbar-actions">
    {{if .SourceURL}}<a href="{{.SourceURL}}" target="_blank" rel="noopener" class="toolbar-link">Source</a>{{end}}
    <a href="javascript:void(0)" onclick="printPage()" class="toolbar-link">Print</a>
    <a href="javascript:void(0)" class="toolbar-link toc-toggle">TOC</a>
    <div class="theme-switcher">
      <label>Theme:</label>
      <select onchange="switchTheme(this.value)">
        {{range .Themes}}<option value="{{.}}"{{if eq . $.Theme}} selected{{end}}>{{.}}</option>
        {{end}}
      </select>
    </div>
  </div>
</div>
<div class="markdown-body">
{{.Content}}
</div>
<aside id="toc-panel" class="toc-panel">
  <div class="toc-header">Table of Contents</div>
  <div class="toc-body"><ul class="toc-list"></ul></div>
</aside>
<script>
mermaid.initialize({startOnLoad: false, theme: document.body.className.indexOf('dark') >= 0 ? 'dark' : 'default'});
function printPage() {
  var orig = document.title;
  var name = orig.split(' - ')[0];
  name = name.replace(/\.(md|markdown)$/i, '');
  document.title = name;
  window.print();
  document.title = orig;
}
function switchTheme(theme) {
  var classes = Array.prototype.slice.call(document.body.classList);
  var kept = classes.filter(function(c) { return c.indexOf('theme-') !== 0; });
  kept.push('theme-' + theme);
  document.body.className = kept.join(' ');
  document.getElementById('theme-css').href = '/_theme/' + theme + '.css';
  localStorage.setItem('mdproxy_theme', theme);
  mermaid.initialize({startOnLoad: false, theme: theme.indexOf('dark') >= 0 ? 'dark' : 'default'});
}
(function() {
  var saved = localStorage.getItem('mdproxy_theme');
  if (saved && MDPROXY_THEMES.indexOf(saved) >= 0) {
    switchTheme(saved);
    var sel = document.querySelector('.theme-switcher select');
    if (sel) sel.value = saved;
  }
})();
</script>
<script>
document.addEventListener('DOMContentLoaded', function() {
  mermaid.run();
  document.querySelectorAll('.math.inline').forEach(function(el) {
    var tex = el.textContent.replace(/^\\\(/, '').replace(/\\\)$/, '').trim();
    katex.render(tex, el, {displayMode: false, throwOnError: false});
  });
  document.querySelectorAll('.math.display').forEach(function(el) {
    var tex = el.textContent.replace(/^\\\[/, '').replace(/\\\]$/, '').trim();
    katex.render(tex, el, {displayMode: true, throwOnError: false});
  });
});
</script>
` + lineAnchorJS + `
` + tocJS + `
` + copyButtonJS + `
{{if .WatchPath}}
<script>
(function() {
  var es = new EventSource('/_sse?path=' + encodeURIComponent('{{.WatchPath}}'));
  es.onmessage = function(e) {
    if (e.data === 'reload') {
      location.reload();
    }
  };
})();
</script>
{{end}}
</body>
</html>`

const dirPageTpl = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Title}} - markdown-proxy</title>
<style>` + commonCSS + `</style>
<link id="theme-css" rel="stylesheet" href="/_theme/{{.Theme}}.css">
<script>var MDPROXY_THEMES=[{{range $i,$t:=.Themes}}{{if $i}},{{end}}"{{$t}}"{{end}}];</script>
</head>
<body class="theme-{{.Theme}}">
<div class="toolbar">
  <a href="/" class="home-link">markdown-proxy</a>
  <div class="theme-switcher">
    <label>Theme:</label>
    <select onchange="switchTheme(this.value)">
      {{range .Themes}}<option value="{{.}}"{{if eq . $.Theme}} selected{{end}}>{{.}}</option>
      {{end}}
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
  var classes = Array.prototype.slice.call(document.body.classList);
  var kept = classes.filter(function(c) { return c.indexOf('theme-') !== 0; });
  kept.push('theme-' + theme);
  document.body.className = kept.join(' ');
  document.getElementById('theme-css').href = '/_theme/' + theme + '.css';
  localStorage.setItem('mdproxy_theme', theme);
}
(function() {
  var saved = localStorage.getItem('mdproxy_theme');
  if (saved && MDPROXY_THEMES.indexOf(saved) >= 0) {
    switchTheme(saved);
    var sel = document.querySelector('.theme-switcher select');
    if (sel) sel.value = saved;
  }
})();
</script>
{{if .WatchPath}}
<script>
(function() {
  var es = new EventSource('/_sse?path=' + encodeURIComponent('{{.WatchPath}}'));
  es.onmessage = function(e) {
    if (e.data === 'reload') {
      location.reload();
    }
  };
})();
</script>
{{end}}
</body>
</html>`

const errorPageTpl = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Title}} - markdown-proxy</title>
<style>` + commonCSS + `</style>
<link id="theme-css" rel="stylesheet" href="/_theme/{{.Theme}}.css">
<script>var MDPROXY_THEMES=[{{range $i,$t:=.Themes}}{{if $i}},{{end}}"{{$t}}"{{end}}];</script>
</head>
<body class="theme-{{.Theme}}">
<div class="toolbar">
  <a href="/" class="home-link">markdown-proxy</a>
  <div class="theme-switcher">
    <label>Theme:</label>
    <select onchange="switchTheme(this.value)">
      {{range .Themes}}<option value="{{.}}"{{if eq . $.Theme}} selected{{end}}>{{.}}</option>
      {{end}}
    </select>
  </div>
</div>
<div class="markdown-body">
<h1>{{.Status}} — {{.Title}}</h1>
<p>{{.Message}}</p>
{{if .Hints}}
<h3>What you can try</h3>
<ul>
{{range .Hints}}<li>{{.}}</li>
{{end}}
</ul>
{{end}}
</div>
<script>
function switchTheme(theme) {
  var classes = Array.prototype.slice.call(document.body.classList);
  var kept = classes.filter(function(c) { return c.indexOf('theme-') !== 0; });
  kept.push('theme-' + theme);
  document.body.className = kept.join(' ');
  document.getElementById('theme-css').href = '/_theme/' + theme + '.css';
  localStorage.setItem('mdproxy_theme', theme);
}
(function() {
  var saved = localStorage.getItem('mdproxy_theme');
  if (saved && MDPROXY_THEMES.indexOf(saved) >= 0) {
    switchTheme(saved);
    var sel = document.querySelector('.theme-switcher select');
    if (sel) sel.value = saved;
  }
})();
</script>
</body>
</html>`
