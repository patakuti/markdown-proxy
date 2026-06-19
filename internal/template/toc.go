package template

const tocCSS = `
.markdown-body h1, .markdown-body h2, .markdown-body h3,
.markdown-body h4, .markdown-body h5, .markdown-body h6 {
  scroll-margin-top: 50px;
}
.toc-panel {
  position: fixed;
  top: 41px;
  right: 0;
  width: 280px;
  height: calc(100vh - 41px);
  overflow-y: auto;
  border-left: 1px solid;
  transform: translateX(100%);
  transition: transform 0.2s ease-out;
  z-index: 90;
  font-size: 13px;
}
body.toc-visible .toc-panel { transform: translateX(0); }
body.toc-visible .markdown-body { margin-right: 280px; }
.toc-header {
  padding: 10px 16px;
  border-bottom: 1px solid;
  font-weight: bold;
}
.toc-body { padding: 8px 0; }
.toc-list, .toc-list ul {
  list-style: none;
  padding-left: 0;
  margin: 0;
}
.toc-list ul { padding-left: 14px; display: none; }
.toc-list li.open > ul { display: block; }
.toc-list li { margin: 0; position: relative; }
.toc-list .toc-row {
  display: flex;
  align-items: flex-start;
}
.toc-list .toc-caret {
  flex: 0 0 16px;
  cursor: pointer;
  user-select: none;
  text-align: center;
  line-height: 27px;
  font-size: 10px;
}
.toc-list .toc-caret::before { content: '\25B6'; }
.toc-list li.open > .toc-row > .toc-caret::before { content: '\25BC'; }
.toc-list .toc-caret.empty { visibility: hidden; }
.toc-list a {
  flex: 1;
  display: block;
  padding: 4px 16px 4px 4px;
  text-decoration: none;
  color: inherit;
  border-left: 2px solid transparent;
  word-break: break-word;
}
.toc-list a.active { font-weight: 600; }
.toc-toggle.disabled {
  opacity: 0.4;
  pointer-events: none;
}
@media (max-width: 900px) {
  body.toc-visible .markdown-body { margin-right: 0; }
  .toc-panel { box-shadow: -2px 0 8px rgba(0,0,0,0.15); }
}
@media print {
  .toc-panel, .toc-toggle { display: none !important; }
  body.toc-visible .markdown-body { margin-right: 0 !important; }
}
`

const tocJS = `<script>
(function() {
  var STORAGE_KEY = 'mdproxy_toc_visible';
  var toggleBtn = document.querySelector('.toc-toggle');
  var panel = document.getElementById('toc-panel');
  if (!toggleBtn || !panel) return;

  var headings = Array.prototype.slice.call(document.querySelectorAll(
    '.markdown-body h1, .markdown-body h2, .markdown-body h3, .markdown-body h4, .markdown-body h5, .markdown-body h6'
  ));
  if (headings.length === 0) {
    toggleBtn.classList.add('disabled');
    toggleBtn.setAttribute('title', 'No headings in this document');
    return;
  }

  var listRoot = panel.querySelector('.toc-list');

  // Assign IDs and determine top level (smallest tag number).
  var topLevel = 6;
  headings.forEach(function(h, i) {
    if (!h.id) h.id = 'toc-' + i;
    var level = parseInt(h.tagName.substring(1), 10);
    if (level < topLevel) topLevel = level;
  });
  // Initially expand items whose children are at level <= topLevel + 1.
  var initialOpenThreshold = topLevel + 1;

  // Build nested tree. Each stack entry represents a heading whose li may host
  // deeper children; the ul under it is created lazily on first child.
  var linkByHeading = new Map();
  var liByHeading = new Map();
  var stack = [{ level: topLevel - 1, ul: listRoot, li: null }];

  headings.forEach(function(h) {
    var level = parseInt(h.tagName.substring(1), 10);
    while (stack.length > 1 && stack[stack.length - 1].level >= level) {
      stack.pop();
    }
    var parent = stack[stack.length - 1];
    if (!parent.ul) {
      parent.ul = document.createElement('ul');
      parent.li.appendChild(parent.ul);
    }
    var li = document.createElement('li');
    var row = document.createElement('div');
    row.className = 'toc-row';
    var caret = document.createElement('span');
    caret.className = 'toc-caret empty';
    var a = document.createElement('a');
    a.href = '#' + h.id;
    a.textContent = h.textContent.trim();
    a.addEventListener('click', function(e) {
      e.preventDefault();
      var toolbar = document.querySelector('.toolbar');
      var offset = toolbar ? toolbar.getBoundingClientRect().height : 0;
      var y = h.getBoundingClientRect().top + window.pageYOffset - offset - 8;
      window.scrollTo({top: y, behavior: 'smooth'});
      history.replaceState(null, '', '#' + h.id);
    });
    row.appendChild(caret);
    row.appendChild(a);
    li.appendChild(row);
    parent.ul.appendChild(li);
    linkByHeading.set(h, a);
    liByHeading.set(h, li);
    stack.push({ level: level, ul: null, li: li });
  });

  // Wire up carets for items that ended up with child <ul>.
  listRoot.querySelectorAll('li').forEach(function(li) {
    var childUl = li.querySelector(':scope > ul');
    if (!childUl) return;
    var caret = li.querySelector(':scope > .toc-row > .toc-caret');
    if (!caret) return;
    caret.classList.remove('empty');
    caret.addEventListener('click', function(e) {
      e.preventDefault();
      e.stopPropagation();
      li.classList.toggle('open');
    });
  });

  // Initial expansion: open ancestors of items at or above initialOpenThreshold,
  // so the first two heading levels are visible while deeper ones stay folded.
  headings.forEach(function(h) {
    var level = parseInt(h.tagName.substring(1), 10);
    if (level >= initialOpenThreshold) return;
    var li = liByHeading.get(h);
    var node = li;
    while (node && node.tagName === 'LI') {
      node.classList.add('open');
      var parentUl = node.parentElement;
      node = parentUl ? parentUl.closest('li') : null;
    }
  });

  // Toggle panel visibility.
  function setVisible(v) {
    document.body.classList.toggle('toc-visible', v);
    try { localStorage.setItem(STORAGE_KEY, v ? '1' : '0'); } catch (e) {}
  }
  toggleBtn.addEventListener('click', function() {
    setVisible(!document.body.classList.contains('toc-visible'));
  });
  var saved = null;
  try { saved = localStorage.getItem(STORAGE_KEY); } catch (e) {}
  if (saved === '1') setVisible(true);

  // Scroll following: the active heading is the last one whose top is at or
  // above a reference line just below the sticky toolbar. This mirrors the
  // reading position regardless of scroll direction.
  var active = null;
  function updateActive() {
    var toolbar = document.querySelector('.toolbar');
    var offset = toolbar ? toolbar.getBoundingClientRect().height : 0;
    var ref = offset + 10;
    var candidate = null;
    for (var i = 0; i < headings.length; i++) {
      if (headings[i].getBoundingClientRect().top <= ref) {
        candidate = headings[i];
      } else {
        break;
      }
    }
    if (!candidate) candidate = headings[0];
    if (candidate === active) return;
    if (active) {
      var prev = linkByHeading.get(active);
      if (prev) prev.classList.remove('active');
    }
    active = candidate;
    var link = linkByHeading.get(active);
    if (link) {
      link.classList.add('active');
      var li = liByHeading.get(active);
      var node = li ? li.parentElement : null;
      while (node && node !== listRoot) {
        if (node.tagName === 'LI') node.classList.add('open');
        node = node.parentElement;
      }
    }
  }
  var ticking = false;
  function onScroll() {
    if (ticking) return;
    ticking = true;
    requestAnimationFrame(function() {
      updateActive();
      ticking = false;
    });
  }
  window.addEventListener('scroll', onScroll, { passive: true });
  updateActive();
})();
</script>`
