package template

const copyButtonCSS = `
.markdown-body pre {
  position: relative;
}
.copy-btn {
  position: absolute;
  top: 6px;
  right: 6px;
  padding: 3px 8px;
  font-size: 11px;
  line-height: 1.5;
  border-radius: 4px;
  border: 1px solid transparent;
  cursor: pointer;
  opacity: 0;
  transition: opacity 0.15s ease, background-color 0.15s ease;
}
.markdown-body pre:hover .copy-btn,
.copy-btn:focus {
  opacity: 1;
}
.copy-btn.copied {
  opacity: 1;
}
.theme-github .copy-btn,
.theme-simple .copy-btn {
  background: #fff;
  color: #444;
  border-color: #d0d7de;
}
.theme-github .copy-btn:hover,
.theme-simple .copy-btn:hover {
  background: #f3f4f6;
  border-color: #adb5bd;
}
.theme-github .copy-btn.copied,
.theme-simple .copy-btn.copied {
  background: #d4edda;
  color: #155724;
  border-color: #c3e6cb;
}
.theme-dark .copy-btn {
  background: #21262d;
  color: #c9d1d9;
  border-color: #30363d;
}
.theme-dark .copy-btn:hover {
  background: #2d333b;
  border-color: #6e7681;
}
.theme-dark .copy-btn.copied {
  background: #1a3d2b;
  color: #3fb950;
  border-color: #2ea043;
}
@media print {
  .copy-btn { display: none !important; }
}
`

const copyButtonJS = `<script>
(function() {
  document.addEventListener('DOMContentLoaded', function() {
    document.querySelectorAll('.markdown-body pre').forEach(function(pre) {
      var btn = document.createElement('button');
      btn.className = 'copy-btn';
      btn.textContent = 'Copy';
      btn.addEventListener('click', function() {
        var code = pre.querySelector('code');
        var text = code ? code.innerText : pre.innerText;
        navigator.clipboard.writeText(text).then(function() {
          btn.textContent = 'Copied!';
          btn.classList.add('copied');
          setTimeout(function() {
            btn.textContent = 'Copy';
            btn.classList.remove('copied');
          }, 2000);
        }).catch(function() {
          btn.textContent = 'Error';
          setTimeout(function() { btn.textContent = 'Copy'; }, 2000);
        });
      });
      pre.appendChild(btn);
    });
  });
})();
</script>`
