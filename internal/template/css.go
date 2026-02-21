package template

const commonCSS = `
* { box-sizing: border-box; }
body { margin: 0; padding: 0; }
.toolbar {
  position: sticky; top: 0; z-index: 100;
  display: flex; justify-content: space-between; align-items: center;
  padding: 8px 20px;
  border-bottom: 1px solid #e1e4e8;
  background: #f6f8fa;
}
.theme-dark .toolbar { background: #1e1e1e; border-color: #444; }
.home-link { text-decoration: none; font-weight: bold; font-size: 14px; }
.theme-github .home-link, .theme-simple .home-link { color: #0366d6; }
.theme-dark .home-link { color: #58a6ff; }
.theme-switcher { display: flex; align-items: center; gap: 6px; font-size: 13px; }
.theme-switcher select { padding: 2px 6px; font-size: 13px; }
.markdown-body {
  max-width: 980px;
  margin: 0 auto;
  padding: 40px 20px;
}
.markdown-body table {
  border-collapse: collapse;
  width: 100%;
}
.markdown-body table th,
.markdown-body table td {
  border: 1px solid #dfe2e5;
  padding: 6px 13px;
}
.markdown-body table tr:nth-child(2n) {
  background-color: #f6f8fa;
}
.theme-dark .markdown-body table tr:nth-child(2n) {
  background-color: #2d2d2d;
}
.theme-dark .markdown-body table th,
.theme-dark .markdown-body table td {
  border-color: #444;
}
`

const githubCSS = `
.theme-github {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  color: #24292e;
  background: #fff;
}
.theme-github .markdown-body h1 { padding-bottom: .3em; border-bottom: 1px solid #eaecef; }
.theme-github .markdown-body h2 { padding-bottom: .3em; border-bottom: 1px solid #eaecef; }
.theme-github .markdown-body a { color: #0366d6; text-decoration: none; }
.theme-github .markdown-body a:hover { text-decoration: underline; }
.theme-github .markdown-body code {
  background: rgba(27,31,35,.05);
  padding: .2em .4em;
  border-radius: 3px;
  font-size: 85%;
}
.theme-github .markdown-body pre {
  background: #f6f8fa;
  padding: 16px;
  border-radius: 6px;
  overflow: auto;
}
.theme-github .markdown-body pre code { background: none; padding: 0; }
.theme-github .markdown-body blockquote {
  color: #6a737d;
  border-left: .25em solid #dfe2e5;
  padding: 0 1em;
  margin: 0;
}
.theme-github .markdown-body img { max-width: 100%; }
`

const simpleCSS = `
.theme-simple {
  font-family: Georgia, "Times New Roman", serif;
  color: #333;
  background: #fefefe;
  line-height: 1.8;
}
.theme-simple .markdown-body a { color: #07c; }
.theme-simple .markdown-body code {
  background: #f0f0f0;
  padding: .15em .3em;
  border-radius: 2px;
}
.theme-simple .markdown-body pre {
  background: #f0f0f0;
  padding: 14px;
  border-radius: 4px;
  overflow: auto;
}
.theme-simple .markdown-body pre code { background: none; padding: 0; }
.theme-simple .markdown-body blockquote {
  color: #666;
  border-left: 3px solid #ccc;
  padding: 0 1em;
  margin: 0;
}
.theme-simple .markdown-body img { max-width: 100%; }
`

const darkCSS = `
.theme-dark {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  color: #c9d1d9;
  background: #0d1117;
}
.theme-dark .markdown-body h1 { padding-bottom: .3em; border-bottom: 1px solid #21262d; }
.theme-dark .markdown-body h2 { padding-bottom: .3em; border-bottom: 1px solid #21262d; }
.theme-dark .markdown-body a { color: #58a6ff; text-decoration: none; }
.theme-dark .markdown-body a:hover { text-decoration: underline; }
.theme-dark .markdown-body code {
  background: rgba(110,118,129,.4);
  padding: .2em .4em;
  border-radius: 3px;
  font-size: 85%;
}
.theme-dark .markdown-body pre {
  background: #161b22;
  padding: 16px;
  border-radius: 6px;
  overflow: auto;
}
.theme-dark .markdown-body pre code { background: none; padding: 0; }
.theme-dark .markdown-body blockquote {
  color: #8b949e;
  border-left: .25em solid #30363d;
  padding: 0 1em;
  margin: 0;
}
.theme-dark .markdown-body img { max-width: 100%; }
`
