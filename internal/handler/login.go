package handler

import (
	"crypto/subtle"
	"net/http"

	"github.com/patakuti/markdown-proxy/internal/config"
)

// CookieName is the name of the authentication cookie.
const CookieName = "mdproxy_token"

type LoginHandler struct {
	cfg *config.Config
}

func NewLoginHandler(cfg *config.Config) *LoginHandler {
	return &LoginHandler{cfg: cfg}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.showForm(w, "")
	case http.MethodPost:
		h.handleLogin(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *LoginHandler) showForm(w http.ResponseWriter, errMsg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	errorHTML := ""
	if errMsg != "" {
		errorHTML = `<p class="error">` + errMsg + `</p>`
	}

	w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Login - markdown-proxy</title>
<style>
  body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
    max-width: 400px;
    margin: 120px auto;
    padding: 0 20px;
    color: #24292e;
    background: #fff;
  }
  h1 { text-align: center; font-size: 1.8em; margin-bottom: 0.5em; }
  .subtitle { text-align: center; color: #586069; margin-bottom: 2em; }
  form { display: flex; flex-direction: column; gap: 12px; }
  input[type="password"] {
    padding: 10px 14px;
    font-size: 16px;
    border: 1px solid #d1d5da;
    border-radius: 6px;
    outline: none;
  }
  input[type="password"]:focus { border-color: #0366d6; box-shadow: 0 0 0 3px rgba(3,102,214,0.3); }
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
  .error { color: #d73a49; text-align: center; margin: 0; }
</style>
</head>
<body>
<h1>markdown-proxy</h1>
<p class="subtitle">Enter your access token</p>
` + errorHTML + `
<form method="POST" action="/_login">
  <input type="password" name="token" placeholder="Access token" autofocus required>
  <button type="submit">Login</button>
</form>
</body>
</html>`))
}

func (h *LoginHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")

	if subtle.ConstantTimeCompare([]byte(token), []byte(h.cfg.AuthToken)) != 1 {
		w.WriteHeader(http.StatusUnauthorized)
		h.showForm(w, "Invalid token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		MaxAge:   h.cfg.AuthCookieMaxAge * 86400, // days to seconds
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
