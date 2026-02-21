package handler

import (
	"net/http"

	"github.com/patakuti/markdown-proxy/internal/config"
)

type TopHandler struct {
	cfg *config.Config
}

func NewTopHandler(cfg *config.Config) *TopHandler {
	return &TopHandler{cfg: cfg}
}

func (h *TopHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(topPageHTML))
}

const topPageHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>markdown-proxy</title>
<style>
  body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
    max-width: 800px;
    margin: 80px auto;
    padding: 0 20px;
    color: #24292e;
    background: #fff;
  }
  h1 { text-align: center; font-size: 2em; margin-bottom: 0.5em; }
  .subtitle { text-align: center; color: #586069; margin-bottom: 2em; }
  .input-group {
    display: flex;
    gap: 8px;
    margin-bottom: 2em;
  }
  input[type="text"] {
    flex: 1;
    padding: 10px 14px;
    font-size: 16px;
    border: 1px solid #d1d5da;
    border-radius: 6px;
    outline: none;
  }
  input[type="text"]:focus { border-color: #0366d6; box-shadow: 0 0 0 3px rgba(3,102,214,0.3); }
  button {
    padding: 10px 20px;
    font-size: 16px;
    background: #0366d6;
    color: #fff;
    border: none;
    border-radius: 6px;
    cursor: pointer;
  }
  button:hover { background: #0256b9; }
  .history { margin-top: 2em; }
  .history h2 { font-size: 1.2em; color: #586069; }
  .history ul { list-style: none; padding: 0; }
  .history li {
    padding: 8px 0;
    border-bottom: 1px solid #eaecef;
  }
  .history a { color: #0366d6; text-decoration: none; }
  .history a:hover { text-decoration: underline; }
  .history .clear-btn {
    background: none;
    border: none;
    color: #586069;
    cursor: pointer;
    font-size: 0.9em;
    padding: 4px 8px;
  }
  .history .clear-btn:hover { color: #d73a49; }
</style>
</head>
<body>
<h1>markdown-proxy</h1>
<p class="subtitle">Enter a file path or URL to view Markdown</p>
<div class="input-group">
  <input type="text" id="path-input" placeholder="/path/to/file.md or https://example.com/doc.md" autofocus>
  <button onclick="navigate()">Open</button>
</div>
<div class="history" id="history-section" style="display:none;">
  <h2>Recent files <button class="clear-btn" onclick="clearHistory()">(clear)</button></h2>
  <ul id="history-list"></ul>
</div>
<script>
function navigate() {
  var input = document.getElementById('path-input').value.trim();
  if (!input) return;
  var url = '';
  if (input.startsWith('http://')) {
    url = '/http/' + input.substring(7);
  } else if (input.startsWith('https://')) {
    url = '/https/' + input.substring(8);
  } else if (input.startsWith('/') || input.startsWith('~') || /^[A-Za-z]:\\/.test(input)) {
    url = '/local' + (input.startsWith('/') ? '' : '/') + input;
  } else {
    url = '/local/' + input;
  }
  saveHistory(input, url);
  window.location.href = url;
}

document.getElementById('path-input').addEventListener('keydown', function(e) {
  if (e.key === 'Enter') navigate();
});

function getHistory() {
  try {
    return JSON.parse(localStorage.getItem('mdproxy_history') || '[]');
  } catch(e) { return []; }
}

function saveHistory(input, url) {
  var history = getHistory().filter(function(h) { return h.input !== input; });
  history.unshift({input: input, url: url, time: Date.now()});
  if (history.length > 20) history = history.slice(0, 20);
  localStorage.setItem('mdproxy_history', JSON.stringify(history));
}

function clearHistory() {
  localStorage.removeItem('mdproxy_history');
  renderHistory();
}

function renderHistory() {
  var history = getHistory();
  var section = document.getElementById('history-section');
  var list = document.getElementById('history-list');
  if (history.length === 0) {
    section.style.display = 'none';
    return;
  }
  section.style.display = 'block';
  list.innerHTML = '';
  history.forEach(function(h) {
    var li = document.createElement('li');
    var a = document.createElement('a');
    a.href = h.url;
    a.textContent = h.input;
    li.appendChild(a);
    list.appendChild(li);
  });
}

renderHistory();
</script>
</body>
</html>`
