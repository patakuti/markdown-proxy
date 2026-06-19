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
        var text = (code ? code.textContent : pre.textContent).replace(/\n+$/, '');
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
