package template

// commonCSS contains structural layout rules only (no colors).
// Visual styling (colors, fonts, backgrounds) is provided by per-theme CSS files.
const commonCSS = `
* { box-sizing: border-box; }
body { margin: 0; padding: 0; }
.toolbar {
  position: sticky; top: 0; z-index: 100;
  display: flex; justify-content: space-between; align-items: center;
  padding: 8px 20px;
  border-bottom: 1px solid;
}
.home-link { text-decoration: none; font-weight: bold; font-size: 14px; }
.toolbar-actions { display: flex; align-items: center; gap: 12px; }
.toolbar-link { font-size: 13px; text-decoration: none; }
.toolbar-link:hover { text-decoration: underline; }
.theme-switcher { display: flex; align-items: center; gap: 6px; font-size: 13px; }
.theme-switcher select { padding: 2px 6px; font-size: 13px; }
@media print {
  .toolbar { display: none; }
  .markdown-body, .markdown-body * {
    -webkit-print-color-adjust: exact;
    print-color-adjust: exact;
  }
  table, pre, .math.display, img, blockquote, li { break-inside: avoid; }
  h1, h2, h3, h4, h5, h6 { break-after: avoid; }
}
.markdown-body pre.text-file {
  white-space: pre-wrap;
}
.plantuml-notice {
  padding: 12px 16px;
  margin: 16px 0;
  border-radius: 6px;
  font-size: 14px;
  line-height: 1.5;
  border: 1px solid;
}
.plantuml-notice code {
  padding: .2em .4em;
  border-radius: 3px;
  font-size: 85%;
}
.markdown-body {
  max-width: 980px;
  margin: 0 auto;
  padding: 40px 20px;
}
.markdown-body table {
  border-collapse: collapse;
  width: auto;
  table-layout: auto;
}
.markdown-body table th,
.markdown-body table td {
  border: 1px solid;
  padding: 6px 13px;
}
.math.display {
  font-size: 0.9em;
  overflow-x: auto;
}
`
