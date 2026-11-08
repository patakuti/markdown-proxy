package template

// wideCSS lifts the max-width limit of the article body when the "Wide"
// toolbar toggle is on, so that wide tables wrap less.
// The class is set on <html> so it can be restored before the body renders.
const wideCSS = `
html.wide .markdown-body { max-width: none; }
@media print {
  .wide-toggle { display: none !important; }
  html.wide .markdown-body { max-width: 980px; }
}
`

// wideHeadJS restores the saved state early to avoid a layout flash.
const wideHeadJS = `<script nonce="{{.Nonce}}">
try { if (localStorage.getItem('mdproxy_wide') === '1') document.documentElement.classList.add('wide'); } catch (e) {}
</script>
`

// wideJS wires up the toolbar toggle and persists its state.
const wideJS = `<script nonce="{{.Nonce}}">
(function() {
  var btn = document.querySelector('.wide-toggle');
  if (!btn) return;
  var root = document.documentElement;
  function sync() {
    var on = root.classList.contains('wide');
    btn.setAttribute('aria-pressed', on ? 'true' : 'false');
  }
  btn.addEventListener('click', function(e) {
    e.preventDefault();
    var on = root.classList.toggle('wide');
    try { localStorage.setItem('mdproxy_wide', on ? '1' : '0'); } catch (err) {}
    sync();
  });
  sync();
})();
</script>`
