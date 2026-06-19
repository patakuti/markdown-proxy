package template

const lineAnchorJS = `<script>
(function() {
  function highlightLines() {
    // Clear previous highlights
    document.querySelectorAll('.line-highlight').forEach(function(el) {
      el.classList.remove('line-highlight');
    });

    var hash = location.hash;
    if (!hash) return;

    // Parse #L12 or #L12-L34
    var m = hash.match(/^#L(\d+)(?:-L(\d+))?$/);
    if (!m) return;

    var startLine = parseInt(m[1], 10);
    var endLine = m[2] ? parseInt(m[2], 10) : startLine;

    // Find all line anchor elements in range
    var anchors = [];
    for (var i = startLine; i <= endLine; i++) {
      var el = document.getElementById('L' + i);
      if (el) anchors.push(el);
    }
    if (anchors.length === 0) {
      // Try to find the nearest anchor before startLine
      for (var i = startLine - 1; i >= 1; i--) {
        var el = document.getElementById('L' + i);
        if (el) {
          el.scrollIntoView({behavior: 'smooth', block: 'center'});
          return;
        }
      }
      return;
    }

    // Scroll to first anchor
    anchors[0].scrollIntoView({behavior: 'smooth', block: 'center'});

    // Highlight: for source-line spans (text files), highlight the element directly;
    // for markdown anchors, highlight the parent block element.
    var highlighted = new Set();
    anchors.forEach(function(anchor) {
      var target;
      if (anchor.classList.contains('source-line')) {
        target = anchor;
      } else {
        target = anchor.parentElement;
        while (target && target.classList.contains('markdown-body')) {
          target = null;
          break;
        }
      }
      if (target && !highlighted.has(target)) {
        highlighted.add(target);
        target.classList.add('line-highlight');
      }
    });
  }

  // Run on page load and hash change
  window.addEventListener('hashchange', highlightLines);
  window.addEventListener('DOMContentLoaded', highlightLines);
})();
</script>`

const lineAnchorCSS = `
.line-highlight {
  border-left-width: 3px;
  border-left-style: solid;
  transition: background-color 0.3s ease;
}
@media print {
  .line-highlight {
    background-color: transparent !important;
    border-left: none !important;
  }
}
`
